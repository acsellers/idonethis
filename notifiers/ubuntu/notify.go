package main

import "github.com/acsellers/ubuntu/notify"

var ms = 1000

func Notify(email, name, message string) {
	var icon string
	if email != "" {
		icon, _ = LocalAvatarPath(email)
	}
	notify.Notify(name, message, icon, 7*ms)
}
