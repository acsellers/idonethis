package windows

import "github.com/conformal/gotk3/gtk"

func LoginWindow(username string, err error) (string, string) {
	gtk.Init(nil)

	d, _ := gtk.DialogNew()
	d.SetTitle("Login to iDoneThis")
	d.SetIconFromFile("/opt/idonethis/indicate.png")

	d.AddButton("Cancel", gtk.RESPONSE_CANCEL)
	d.AddButton("Login", gtk.RESPONSE_OK)

	ul, _ := gtk.LabelNew("Username")
	ui, _ := gtk.EntryNew()
	ui.SetText(username)
	ui.Set("margin-left", 12)
	ub, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 4)
	ub.PackStart(ul, true, true, 0)
	ub.PackEnd(ui, true, true, 0)

	pl, _ := gtk.LabelNew("Password")
	pi, _ := gtk.EntryNew()
	pi.SetVisibility(false)
	pi.Set("margin-left", 12)
	pb, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 4)
	pb.PackStart(pl, true, true, 0)
	pb.PackEnd(pi, true, true, 0)

	tl, _ := gtk.LabelNew("Client")
	tl.SetMarkup("<big>iDoneThis Client</big>")

	bx, _ := d.GetContentArea()
	bx.PackStart(tl, true, true, 4)

	if err != nil {
		el, _ := gtk.LabelNew("error")
		el.SetMarkup("<span color=\"#ff0000\">" + err.Error() + "</span>")
		bx.PackStart(el, true, true, 4)
	}

	bx.PackEnd(pb, true, true, 4)
	bx.PackEnd(ub, true, true, 4)
	d.ShowAll()

	if gtk.ResponseType(d.Run()) == gtk.RESPONSE_OK {
		ut, _ := ui.GetText()
		pt, _ := pi.GetText()
		d.Destroy()
		return ut, pt
	}
	d.Destroy()
	return "", ""
}
