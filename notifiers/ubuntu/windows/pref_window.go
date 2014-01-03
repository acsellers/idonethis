package windows

import (
	"fmt"
	"strconv"

	"github.com/conformal/gotk3/glib"
	"github.com/conformal/gotk3/gtk"
)

type PrefResult struct {
	CheckInterval int
	WipeUserData  bool
}

func PrefWindow(interval int, user string) PrefResult {
	gtk.Init(nil)

	d, _ := gtk.DialogNew()
	d.SetTitle("iDoneThis Settings")
	d.SetIconFromFile("/opt/idonethis/indicate.png")

	cl, _ := gtk.LabelNew("Interval to Check for Dones (minutes)")
	ci, _ := gtk.EntryNew()
	ci.SetText(fmt.Sprint(interval))
	cb, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 4)
	cb.PackStart(cl, true, true, 0)
	cb.PackEnd(ci, true, true, 0)

	bx, _ := d.GetContentArea()
	bx.PackStart(cb, true, true, 4)

	var clearData bool
	if user != "" {
		dl, _ := gtk.LabelNew("Currently Saving Password for " + user)
		dl.Set("xalign", 0.0)
		di, _ := gtk.ButtonNewWithLabel("Clear")
		di.Connect("clicked", func(o glib.IObject) {
			clearData = true
		})
		db, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 4)
		db.PackStart(dl, true, true, 0)
		db.PackEnd(di, true, true, 0)
		bx.PackStart(db, true, true, 4)
	}

	d.AddButton("Apply", gtk.RESPONSE_OK)
	d.ShowAll()

	if gtk.ResponseType(d.Run()) == gtk.RESPONSE_OK {
		ct, _ := ci.GetText()
		ni, e := strconv.Atoi(ct)
		if e == nil && ni > 0 {
			interval = ni
		}
		d.Destroy()
		return PrefResult{interval, clearData}
	}
	d.Destroy()
	return PrefResult{interval, clearData}
}
