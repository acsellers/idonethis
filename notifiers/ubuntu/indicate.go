package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/conformal/gotk3/glib"
	"github.com/conformal/gotk3/gtk"
	"github.com/doxxan/appindicator"
)

var Indicator *appindicator.AppIndicator

func StartIndicator() {
	items := []string{"New Done", "Preferences", "", "Quit"}
	handlers := []ItemHandler{NewHandler, PrefHandler, nil, QuitHandler}
	m, _ := gtk.MenuNew()

	for i, _ := range items {
		if items[i] != "" {
			mi, _ := gtk.MenuItemNewWithLabel(items[i])
			mi.Connect("activate", handlers[i])
			m.Append(mi)
		} else {
			si, _ := gtk.SeparatorMenuItemNew()
			m.Append(si)
		}
	}
	m.ShowAll()

	Indicator = appindicator.NewAppIndicatorWithPath(
		"idonethis",
		"/opt/idonethis/indicate.png",
		"/opt/idonethis/indicate.png",
		int(appindicator.CategoryApplicationStatus),
	)
	fmt.Println("Created indicator")
	Indicator.SetStatus(appindicator.StatusActive)
	Indicator.SetTitle("iDoneThis")

	fmt.Println("setting menu")
	Indicator.C_SetMenu(unsafe.Pointer(m.Native()))

	go gtk.Main()
}

type ItemHandler func(o *glib.Object)

func NewHandler(o *glib.Object) {
	NewPostWindow()
}

func PrefHandler(o *glib.Object) {
	PrefWindow()
}

func QuitHandler(o *glib.Object) {
	os.Exit(0)
}
