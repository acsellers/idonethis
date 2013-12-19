package idonethis

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

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

type newDone struct {
	Calendar  string `json:"calendar"`
	DoneDate  string `json:"done_date"`
	Owner     string `json:"owner"`
	ShortName string `json:"short_name"`
	Text      string `json:"text"`
	User      User   `json:"user"`
}

type likeDone struct {
	Action    string `json:"action"`
	DoneId    int    `json:"done_id"`
	ShortName string `json:"short_name"`
	Username  string `json:"username"`
}
