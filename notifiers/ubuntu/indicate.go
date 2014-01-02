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
	items := []string{"new_item", "pref_item", "quit_item"}
	handlers := []ItemHandler{NewHandler, PrefHandler, QuitHandler}

	for i, _ := range items {
		io, err := Builder.GetObject(items[i])
		if err != nil {
			fmt.Println("error", items[i])
			continue
		}
		if ii, ok := io.(*gtk.MenuItem); ok {
			fmt.Println(items[i])
			ii.Connect("activate", handlers[i])
		} else {
			fmt.Println("Not widget:", items[i])
		}
	}

	Indicator = appindicator.NewAppIndicatorWithPath(
		"idonethis",
		"/opt/idonethis/indicate.png",
		"/opt/idonethis/indicate.png",
		int(appindicator.CategoryApplicationStatus),
	)
	fmt.Println("Created indicator")
	Indicator.SetStatus(appindicator.StatusActive)
	Indicator.SetTitle("iDoneThis")

	mo, e := Builder.GetObject("menu")
	if e != nil {
		panic(e)
	}
	if m, ok := mo.(*gtk.Menu); ok {
		fmt.Println("setting menu")
		Indicator.C_SetMenu(unsafe.Pointer(m.Native()))
	} else {
		fmt.Println("Menu not available")
	}
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
