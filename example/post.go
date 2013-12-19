package main

import (
	"flag"
	"fmt"
	"strings"

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
		panic(e)
	}
	d, e := c.PostDone(strings.Join(flag.Args(), " "))
	if e != nil {
		panic(e)
	}
	fmt.Println(d.Text)
}
