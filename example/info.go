package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/acsellers/idonethis"
)

var (
	Username, Password string
)

func init() {
	flag.StringVar(&Username, "user", "", "Username that you use to log into idonethis")
	flag.StringVar(&Password, "pass", "", "Password that you use to log into idonethis")
	flag.Parse()
}
func main() {
	if Username == "" || Password == "" {
		fmt.Println("Must add user and pass arguments to log into idonethis")
		return
	}

	c, e := idonethis.NewClient(Username, Password)
	if e != nil {
		log.Fatal(e)
	}

	cu, e := c.CurrentUser()
	if e != nil {
		log.Println(e)
	} else {
		fmt.Println("Your UserName is", cu.Username, "and your display name is", cu.NicestName)
	}

	fmt.Println("\nTeam Information:")
	t, e := c.Team()
	if e != nil {
		log.Println(e)
	} else {
		fmt.Println("You are part of", t.Name, "with", t.MemberCount-1, "other members")
	}

	fmt.Println("\nTeam Members:")
	team, e := c.TeamMembers()
	if e != nil {
		log.Println(e)
	} else {
		for i, member := range team {
			fmt.Println(i+1, "-", member.NicestName, "-", member.Email)
		}
	}

	fmt.Println("\nHashtags used by your team:")
	tags, e := c.Tags()
	if e != nil {
		log.Println(e)
	} else {
		for _, tag := range tags {
			fmt.Println(tag.Name, "used", len(tag.DoneIds), "times")
		}
	}

	fmt.Println("\nYou are following:")
	follows, e := c.FollowData()
	if e != nil {
		log.Println(e)
	} else {
		for _, follow := range follows {
			fmt.Println(follow.Followee)
		}
	}

	fmt.Println("\nDones by the team:")
	dones, e := c.AllDones()
	if e != nil {
		log.Println(e)
	} else {
		for _, done := range dones {
			fmt.Println(done.Owner, done.DoneDate, done.Text)
		}
	}

	fmt.Println("\nDones written today:")
	dones, e := c.FilteredDones(idonethis.DoneFilter{Start: time.Now(), End: time.Now()})
	if e != nil {
		log.Println(e)
	} else {
		for _, done := range dones {
			fmt.Println(done.Owner, done.DoneDate, done.Text)
		}
	}

}
