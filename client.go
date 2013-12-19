package idonethis

import (
	"bytes"
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
		if len(results) >= 2 {
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

func (c *Client) PostDone(text string) (Done, error) {
	current, err := c.CurrentUser()
	if err != nil {
		return Done{}, err
	}
	nd := newDone{
		Calendar:  c.TeamName,
		DoneDate:  time.Now().Format(ymd),
		Owner:     current.Email,
		ShortName: c.TeamName,
		Text:      text,
		User:      current,
	}
	body, err := json.Marshal(nd)
	if err != nil {
		return Done{}, err
	}

	result, err := c.postJson("dones/", body)
	if err != nil {
		return Done{}, err
	}

	var d Done
	err = json.Unmarshal(result, &d)
	if err != nil {
		fmt.Println(string(result))
		return Done{}, err
	}

	return d, nil
}

func (c *Client) LikeDone(d Done) error {
	a := likeDone{
		Action:    "like",
		DoneId:    d.Id,
		ShortName: c.TeamName,
		Username:  c.UserName,
	}

	body, err := json.Marshal(a)
	if err != nil {
		return err
	}

	lu := fmt.Sprintf("feedback/done/%s/%d/", c.UserName, d.Id)
	_, err = c.postJson(lu, body)
	return err
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

func (c *Client) postJson(page string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", c.teamUrl()+page, bytes.NewReader(body))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", c.teamUrl())
	token, err := c.getCSRF()
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("X-CSRFToken", token)

	resp, err := c.session.Do(req)
	if err != nil {
		return []byte{}, err
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return result, nil
}

func (c *Client) getCSRF() (string, error) {
	// load the login page to get a csrf token for the login page
	homeUrl := fmt.Sprintf("https://idonethis.com/cal/%s/", c.TeamName)
	home, err := c.session.Get(homeUrl)
	if err != nil {
		return "", err
	}

	homePage, err := ioutil.ReadAll(home.Body)
	if err != nil {
		return "", err
	}

	r := regexp.MustCompile(`name='csrfmiddlewaretoken' value='(.*)' `)
	results := r.FindSubmatch(homePage)
	if len(results) >= 2 {
		return string(results[1]), nil
	}

	return "", fmt.Errorf("Could not retrieve csrf token")
}
