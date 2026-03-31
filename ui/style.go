package ui

import (
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func loadCSS() {
	css := gtk.NewCSSProvider()
	css.LoadFromData(`
		.toolbar {
			background: @headerbar_bg_color;
			min-height: 60px;
			padding: 4px;
		}
	`)
	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		css,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)
}
