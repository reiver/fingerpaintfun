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

		.swatch-bounce {
			transition: transform 150ms ease-out;
			transform: scale(1.3);
		}

		.swatch-normal {
			transition: transform 150ms ease-in;
			transform: scale(1.0);
		}
	`)
	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		css,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)
}
