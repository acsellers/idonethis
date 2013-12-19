package idonethis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const ymd = "2006-01-02"

type Client struct {
	UserName string
	password string
	TeamName string

	session *http.Client
}

func NewClient(username, password string) (*Client, error) {
	c := &Client{UserName: username, password: password}

	jar, _ := cookiejar.New(&cookiejar.Options{})
	c.session = &http.Client{Jar: jar}

	// load the login page to get a csrf token for the login page
	login, err := c.session.Get("https://idonethis.com/accounts/login/")
	if err != nil {
		return nil, err
	}

	loginPage, err := ioutil.ReadAll(login.Body)
	if err != nil {
		return nil, err
	}

	r := regexp.MustCompile("value=(.*) ")
	results := r.FindSubmatch(loginPage)
	var result []byte
	if len(results) == 2 {
		result = results[1]
		result = result[1 : len(result)-1]
	} else {
		return nil, fmt.Errorf("Could not retrieve csrf token")
	}

	// create a login request to setup our http.Client's info
	// we would use client.PostForm, but idonethis wants the referer
	// header set, so we have to construct the request manually
	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)
	form.Add("csrfmiddlewaretoken", string(result))

	formRequest, err := http.NewRequest("POST",
		"https://idonethis.com/accounts/login/",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, err
	}
	formRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	formRequest.Header.Set("Referer", "https://idonethis.com/accounts/login/")

	second, err := c.session.Do(formRequest)
	if err != nil {
		return nil, err
	}

	secondPage, e := ioutil.ReadAll(second.Body)
	if e != nil {
		fmt.Println(e)
	} else {
		rt := regexp.MustCompile("/cal/(.*)/\\\" title")
		results = rt.FindSubmatch(secondPage)
		if len(results) == 2 {
			c.TeamName = string(results[1])
		}
	}

	return c, nil
}

func (c *Client) Team() (Team, error) {
	body, err := c.getJson("")
	if err != nil {
		return Team{}, err
	}

	var t Team
	err = json.Unmarshal(body, &t)
	if err != nil {
		return Team{}, err
	}

	return t, nil
}

func (c *Client) CurrentUser() (User, error) {
	body, err := c.getJson("members/current/")
	if err != nil {
		return User{}, err
	}

	var u User
	err = json.Unmarshal(body, &u)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func (c *Client) AllDones() ([]Done, error) {
	body, err := c.getJson("dones/")
	if err != nil {
		return []Done{}, err
	}

	var d []Done
	err = json.Unmarshal(body, &d)
	if err != nil {
		return []Done{}, err
	}

	return d, nil
}

func (c *Client) FilteredDones(df DoneFilter) ([]Done, error) {
	body, err := c.getJson("dones/" + df.String())
	if err != nil {
		return []Done{}, err
	}

	var d []Done
	err = json.Unmarshal(body, &d)
	if err != nil {
		return []Done{}, err
	}

	return d, nil
}

func (c *Client) FollowData() ([]Follow, error) {
	body, err := c.getJson("follow/")
	if err != nil {
		return []Follow{}, err
	}

	var fd []Follow
	err = json.Unmarshal(body, &fd)
	if err != nil {
		return []Follow{}, err
	}

	return fd, nil

}
func (c *Client) TeamMembers() ([]User, error) {
	body, err := c.getJson("members/")
	if err != nil {
		return []User{}, err
	}

	var u []User
	err = json.Unmarshal(body, &u)
	if err != nil {
		return []User{}, err
	}

	return u, nil

}

func (c *Client) Tags() ([]Tag, error) {
	body, err := c.getJson("tags/")
	if err != nil {
		return []Tag{}, err
	}

	var t []Tag
	err = json.Unmarshal(body, &t)
	if err != nil {
		return []Tag{}, err
	}

	return t, nil
}

func (c *Client) teamUrl() string {
	return fmt.Sprintf("https://idonethis.com/api/v3/team/%s/", c.TeamName)
}

func (c *Client) getJson(page string) ([]byte, error) {
	resp, err := c.session.Get(c.teamUrl() + page)
	if err != nil {
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

type User struct {
	Id            int `json:"id"`
	TeamProfileId int `json:"team_profile_id"`

	Username   string `json:"username"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	NicestName string `json:"nicest_name"`

	HasJoined bool `json:"has_joined"`
	IsStaff   bool `json:"is_staff"`
	IsAdmin   bool `json:"is_admin"`

	Email           string   `json:"email"`
	AvatarUrl       string   `json:"avatar_url"`
	AlternateEmails []string `json:"alternate_emails"`
}

// Follower and followee strings are the Usernames of the follower and followee
type Follow struct {
	Id       int    `json:"id"`
	Follower string `json:"follower"`
	Followee string `json:"followee"`
}

type Team struct {
	Name        string `json:"name"`
	ShortName   string `json:"short_name"`
	Active      bool   `json:"active"`
	MemberCount int    `json:"member_count"`
	Question    string `json:"question"`
	/*
	  Unmarshaled fields

	  "fixed_discount"
	  "is_free_forever"
	  "monthly_price_per_user"
	  "total_discount"
	  "referrer_reward"
	  "cc_last4"
	  "referral_code"
	  "trial_expiration_date"
	  "coupon_code"
	  "referred_reward"
	  "has_discount"
	  "is_paying"
	  "on_trial"
	  "percent_discount"
	  "monthly_charge"
	  "cancellation_requested"
	*/

}

type Tag struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	DoneIds []int  `json:"done_ids"`
}

type Done struct {
	Id            int    `json:"id"`
	TeamShortName string `json:"team_short_name"`
	DoneDate      string `json:"done_date"`
	Owner         string `json:"owner"`

	Text         string `json:"text"`
	MarkedupText string `json:"markedup_text"`

	Tags     []Tag  `json:"tags"`
	Comments []Done `json:"comments"`
	Likes    []Like `json:"likes"`

	/*
	  Unknown attributes

	  "origin": null
	  "rawintegrationdata": {}
	*/
}

type Like struct {
	Email string `json:"user"`
}

type DoneFilter struct {
	Tags       []string
	Start, End time.Time
	Page       int
	PerPage    int
}

func (df DoneFilter) String() string {
	vals := url.Values{}

	if len(df.Tags) > 0 {
		vals.Add("tags", strings.Join(df.Tags, ","))
	}

	if df.Page != 0 {
		vals.Add("page", fmt.Sprint(df.Page))
	}

	if df.PerPage != 0 {
		vals.Add("per_page", fmt.Sprint(df.Page))
	}

	if !df.Start.IsZero() {
		vals.Add("start", df.Start.Format(ymd))
	}
	if !df.End.IsZero() {
		vals.Add("end", df.End.Format(ymd))
	}
	vs := vals.Encode()

	if vs != "" {
		return "?" + vals.Encode()
	}

	return ""
}
