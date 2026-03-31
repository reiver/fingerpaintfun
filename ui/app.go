package ui

import (
	"os"

	"fingerpaintfun/cfg"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
)

// Run starts the application main loop.
func Run() {
	app := adw.NewApplication(cfg.AppID, gio.ApplicationFlagsNone)
	app.ConnectActivate(func() {
		initSound()
		win := newWindow(app)
		win.Show()
	})
	app.Run(os.Args)
}
