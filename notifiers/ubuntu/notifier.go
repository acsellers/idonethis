package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/acsellers/idonethis"
	"github.com/conformal/gotk3/gtk"
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
	CheckForDones()

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

func GetInfoFromLoginWindow() {
	lw := NewLoginWindow()
	lw.SetUsername(Username)
	go func() {
		gtk.Main()
	}()

	select {
	case <-lw.CancelBtn:
		close(lw.LoginBtn)
		close(lw.CancelBtn)
		lw.LoginBtn = nil
		lw.CancelBtn = nil

		lw.Window.Destroy()
		fmt.Println("Closing because Login was Cancelled")
		os.Exit(0)
	case <-lw.LoginBtn:
		close(lw.LoginBtn)
		close(lw.CancelBtn)
		lw.LoginBtn = nil
		lw.CancelBtn = nil

		lw.GetUserData()
		lw.Window.Destroy()
	}
}

func GetLoginInfoWithError(e error) {
	lw := NewLoginWindow()
	lw.SetUsername(Username)
	lw.SetError(e)
	go func() {
		gtk.Main()
	}()

	select {
	case <-lw.CancelBtn:
		close(lw.LoginBtn)
		close(lw.CancelBtn)
		lw.LoginBtn = nil
		lw.CancelBtn = nil

		lw.Window.Destroy()
		fmt.Println("Closing because Login was Cancelled")
		os.Exit(0)
	case <-lw.LoginBtn:
		close(lw.LoginBtn)
		close(lw.CancelBtn)
		lw.LoginBtn = nil
		lw.CancelBtn = nil

		lw.GetUserData()
		lw.Window.Destroy()
	}
}
