package main

import (
	"fmt"

	"github.com/conformal/gotk3/glib"
	"github.com/conformal/gotk3/gtk"
)

type LoginWindow struct {
	Builder   *gtk.Builder
	Window    *gtk.Window
	CancelBtn chan bool
	LoginBtn  chan bool
}

func NewLoginWindow() *LoginWindow {
	lw := &LoginWindow{Builder: getBuilder()}
	lw.CancelBtn = make(chan bool)
	lw.LoginBtn = make(chan bool)

	lw.Window = setupWindow(lw.Builder)
	lw.connectSignals()
	return lw
}

func GtkInit() {
	gtk.Main()
}

func getBuilder() *gtk.Builder {
	gtk.Init(nil)

	b, err := gtk.BuilderNew()
	if err != nil {
		panic(err)
	}

	err = b.AddFromFile("/home/andrew/login.glade")
	if err != nil {
		panic(err)
	}

	return b
}

func setupWindow(b *gtk.Builder) *gtk.Window {
	obj, err := b.GetObject("login_window")
	if err != nil {
		panic(err)
	}

	if lw, ok := obj.(*gtk.Window); ok {
		lw.ShowAll()
		return lw
	} else {
		panic("could not get login_window")
	}
}

func (lw *LoginWindow) connectSignals() {
	co, err := lw.Builder.GetObject("cancel_button")
	if err != nil {
		panic(err)
	}

	if cb, ok := co.(*gtk.Button); ok {
		cb.Connect("clicked", func(o *glib.Object) {
			fmt.Println("Cancel button clicked")

			if lw.CancelBtn != nil {
				lw.CancelBtn <- true
			}
		})
	}

	lo, err := lw.Builder.GetObject("login_button")
	if err != nil {
		panic(err)
	}

	if lb, ok := lo.(*gtk.Button); ok {
		lb.Connect("clicked", func(o glib.IObject) {
			fmt.Println("Login button clicked")

			if lw.LoginBtn != nil {
				lw.LoginBtn <- true
			}
		})
	}

	lw.Window.Connect("activate-default", func(o glib.IObject) {
		fmt.Println("Window Default Action")
		lw.LoginBtn <- true
	})

	lw.Window.Connect("destroy", func(o glib.IObject) {
		fmt.Println("Window Closed")
		if lw.CancelBtn != nil {
			lw.CancelBtn <- true
		}
	})
}

func (lw *LoginWindow) SetUsername(username string) {
	uo, err := lw.Builder.GetObject("username_entry")
	if err != nil {
		panic(err)
	}

	if ue, ok := uo.(*gtk.Entry); ok {
		ue.SetText(username)
	}
}

func (lw *LoginWindow) GetUserData() {
	uo, err := lw.Builder.GetObject("username_entry")
	if err != nil {
		panic(err)
	}

	if ue, ok := uo.(*gtk.Entry); ok {
		Username, _ = ue.GetText()
	}

	po, err := lw.Builder.GetObject("password_entry")
	if err != nil {
		panic(err)
	}

	if pe, ok := po.(*gtk.Entry); ok {
		Password, _ = pe.GetText()
	}

}

func (lw *LoginWindow) SetError(e error) {
	eo, err := lw.Builder.GetObject("error_label")
	if err != nil {
		panic(err)
	}

	if el, ok := eo.(*gtk.Label); ok {
		el.SetText(e.Error())
	}

}
