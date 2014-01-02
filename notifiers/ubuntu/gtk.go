package main

import (
	"fmt"

	"github.com/conformal/gotk3/glib"
	"github.com/conformal/gotk3/gtk"
)

var Builder *gtk.Builder

type LoginWindow struct {
	Window    *gtk.Window
	CancelBtn chan bool
	LoginBtn  chan bool
}

func init() {
	Builder = getBuilder()
}
func NewLoginWindow() *LoginWindow {
	lw := &LoginWindow{}
	lw.CancelBtn = make(chan bool)
	lw.LoginBtn = make(chan bool)

	lw.Window = setupWindow(Builder)
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
	co, err := Builder.GetObject("cancel_button")
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

	lo, err := Builder.GetObject("login_button")
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
	uo, err := Builder.GetObject("username_entry")
	if err != nil {
		panic(err)
	}

	if ue, ok := uo.(*gtk.Entry); ok {
		ue.SetText(username)
	}
}

func (lw *LoginWindow) GetUserData() {
	uo, err := Builder.GetObject("username_entry")
	if err != nil {
		panic(err)
	}

	if ue, ok := uo.(*gtk.Entry); ok {
		Username, _ = ue.GetText()
	}

	po, err := Builder.GetObject("password_entry")
	if err != nil {
		panic(err)
	}

	if pe, ok := po.(*gtk.Entry); ok {
		Password, _ = pe.GetText()
	}

}

func (lw *LoginWindow) SetError(e error) {
	eo, err := Builder.GetObject("error_label")
	if err != nil {
		panic(err)
	}

	if el, ok := eo.(*gtk.Label); ok {
		el.SetText(e.Error())
	}

}

func NewPostWindow() {
	obj, err := Builder.GetObject("post_window")
	if err != nil {
		panic(err)
	}

	if pw, ok := obj.(*gtk.Window); ok {
		pw.ShowAll()
		pbo, err := Builder.GetObject("post_button")
		if err != nil {
			panic(err)
		}

		if pb, ok := pbo.(*gtk.Button); ok {
			pb.Connect("clicked", func(o glib.IObject) {
				tbo, err := Builder.GetObject("post_buffer")
				if err != nil {
					panic(err)
				}

				if tb, ok := tbo.(*gtk.TextBuffer); ok {
					s, e := tb.GetBounds()
					text, err := tb.GetText(s, e, false)
					if err != nil {
						panic(err)
					}
					fmt.Println(text)
					d, err := Client.PostDone(text)
					fmt.Println("Done response", e)
					fmt.Println(d.Id)
				}
				pw.Destroy()
			})
		}

	} else {
		panic("could not get post_window")
	}

}

func PrefWindow() {

}
