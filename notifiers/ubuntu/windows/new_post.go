package windows

import "github.com/conformal/gotk3/gtk"

func PostWindow() string {
	gtk.Init(nil)

	d, _ := gtk.DialogNew()
	d.SetTitle("Post New Done")
	d.SetIconFromFile("/opt/idonethis/indicate.png")

	d.AddButton("Post", gtk.RESPONSE_OK)

	te, _ := gtk.TextViewNew()
	te.SetEditable(true)
	te.Show()
	te.SetSizeRequest(300, 48)
	te.SetWrapMode(gtk.WRAP_WORD)

	tb, _ := gtk.TextBufferNew(nil)
	te.SetBuffer(tb)

	bx, _ := d.GetContentArea()
	bx.PackStart(te, true, true, 4)

	if gtk.ResponseType(d.Run()) == gtk.RESPONSE_OK {
		s, e := tb.GetBounds()
		text, _ := tb.GetText(s, e, false)
		d.Destroy()
		return text
	}
	d.Destroy()

	return ""
}
