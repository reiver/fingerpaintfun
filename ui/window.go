package ui

import (
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func newWindow(app *adw.Application) *adw.ApplicationWindow {
	win := adw.NewApplicationWindow(&app.Application)
	win.SetDefaultSize(360, 648)
	win.SetTitle("Finger Paint Fun")

	// Navigation view for canvas ↔ gallery transitions.
	nav := adw.NewNavigationView()

	// Main drawing page: canvas on top, toolbar on bottom.
	vbox := gtk.NewBox(gtk.OrientationVertical, 0)

	canvas := newCanvasWidget()
	canvas.area.SetVExpand(true)
	canvas.area.SetHExpand(true)
	vbox.Append(canvas.area)

	// Bottom toolbar with palette, brush picker, bg picker, and action buttons.
	toolbar := newToolbarWidget(canvas)

	// Gallery button in the toolbar.
	galleryBtn := gtk.NewButton()
	galleryBtn.SetIconName("view-grid-symbolic")
	galleryBtn.SetSizeRequest(48, 48)
	galleryBtn.AddCSSClass("flat")
	galleryBtn.ConnectClicked(func() {
		// Auto-save before entering gallery.
		if len(canvas.state.Committed) > 0 {
			SaveCanvas(canvas)
		}
		gallery := newGalleryPage(canvas, win, nav)
		nav.Push(gallery.page)
	})
	toolbar.actions.box.Append(galleryBtn)

	vbox.Append(toolbar.container)

	drawPage := adw.NewNavigationPage(vbox, "Draw")
	nav.Add(drawPage)

	win.SetContent(nav)

	// Auto-save on close.
	win.ConnectCloseRequest(func() bool {
		if len(canvas.state.Committed) > 0 {
			SaveCanvas(canvas)
		}
		return false // allow close to proceed
	})

	loadCSS()

	return win
}
