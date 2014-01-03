package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/acsellers/idonethis"
	"github.com/acsellers/idonethis/notifiers/ubuntu/windows"
)

var (
	PasswordNotFound = fmt.Errorf("No password stored")
	ErrorLog         *log.Logger
	CheckInterval    <-chan time.Time
	CheckMinutes     = 5
)

var (
	Username, Password string
	SeenDones          map[int]bool
	Client             *idonethis.Client
)

func init() {
	SeenDones = make(map[int]bool)
}

func main() {
	GetLoginInfo()
	for Username == "" || Password == "" {
		GetLoginInfo()
	}

	var e error
	Client, e = idonethis.NewClient(Username, Password)
	for e != nil {
		fmt.Println("Could not authenticate, asking for new user/pass")
		GetLoginInfoWithError(e)
		Client, e = idonethis.NewClient(Username, Password)
	}

	SaveSimple()
	SetPassword(Username, Password)

	StartIndicator()
	fmt.Println("checking for dones")
	CheckForPreviousDones()

	fmt.Println("starting to tick")
	if CheckInterval == nil {
		CheckInterval = time.Tick(time.Duration(CheckMinutes) * time.Minute)
	}
	for _ = range CheckInterval {
		fmt.Println("checking for new dones")
		CheckForDones()
	}
}

func CheckForDones() {
	f := idonethis.DoneFilter{Start: time.Now(), End: time.Now()}
	dones, e := Client.FilteredDones(f)
	if e != nil {
		fmt.Println(e)
		os.Exit(0)
	}

	for _, done := range dones {
		if _, ok := SeenDones[done.Id]; !ok {
			SeenDones[done.Id] = true
			Notify("", done.Owner, done.Text)
		}
	}
}

func CheckForPreviousDones() {
	f := idonethis.DoneFilter{Start: time.Now(), End: time.Now()}
	dones, e := Client.FilteredDones(f)
	if e != nil {
		fmt.Println(e)
		os.Exit(0)
	}

	for _, done := range dones {
		if _, ok := SeenDones[done.Id]; !ok {
			SeenDones[done.Id] = true
		}
	}
}
func GetLoginInfo() {
	GetSimple()
	if Username == "" {
		GetInfoFromLoginWindow()
		return
	}

	var e error
	Password, e = GetPassword(Username)
	if e != nil || Password == "" {
		GetInfoFromLoginWindow()
	}
}

func NewPostWindow() {
	text := windows.PostWindow()
	if text == "" {
		return
	}

	fmt.Println(text)
	d, err := Client.PostDone(text)
	if err != nil {
		return
	}
	fmt.Println(d.Id)
}

func GetInfoFromLoginWindow() {
	Username, Password = windows.LoginWindow(Username, nil)
	if Username == "" {
		fmt.Println("User did not wish to login, terminating")
		os.Exit(0)
	}
}

func GetLoginInfoWithError(e error) {
	Username, Password = windows.LoginWindow(Username, e)
	if Username == "" {
		fmt.Println("User did not wish to login, terminating")
		os.Exit(0)
	}
}

func PrefWindow() {
	r := windows.PrefWindow(CheckMinutes, Username)
	CheckMinutes = r.CheckInterval
	CheckInterval = time.Tick(time.Duration(CheckMinutes) * time.Minute)

	if r.WipeUserData {
		SetPassword(Username, "")
	}
}
