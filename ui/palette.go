package ui

import (
	"math"

	"fingerpaintfun/cfg"

	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// paletteWidget is a 2-row × 6-column grid of color swatches.
type paletteWidget struct {
	grid     *gtk.Grid
	selected int
	swatches []*gtk.DrawingArea
	onChange func(color [4]float64)
}

func newPaletteWidget(onChange func(color [4]float64)) *paletteWidget {
	p := &paletteWidget{
		grid:     gtk.NewGrid(),
		selected: cfg.DefaultColor,
		onChange: onChange,
	}

	p.grid.SetRowSpacing(8)
	p.grid.SetColumnSpacing(8)
	p.grid.SetHAlign(gtk.AlignCenter)
	p.grid.SetVAlign(gtk.AlignCenter)

	for i, color := range cfg.Palette {
		idx := i
		c := color

		swatch := gtk.NewDrawingArea()
		swatch.SetSizeRequest(48, 48)
		swatch.SetDrawFunc(func(area *gtk.DrawingArea, cr *cairo.Context, w, h int) {
			drawSwatch(cr, w, h, c, idx == p.selected)
		})

		click := gtk.NewGestureClick()
		click.ConnectPressed(func(nPress int, x, y float64) {
			p.selected = idx
			if p.onChange != nil {
				p.onChange([4]float64{c.R, c.G, c.B, c.A})
			}
			PlayPop()
			// Bounce animation on selected swatch.
			swatch.AddCSSClass("swatch-bounce")
			glib.TimeoutAdd(150, func() bool {
				swatch.RemoveCSSClass("swatch-bounce")
				return false // don't repeat
			})
			// Redraw all swatches to update highlight.
			for _, sw := range p.swatches {
				sw.QueueDraw()
			}
		})
		swatch.AddController(click)

		p.swatches = append(p.swatches, swatch)
		row := idx / 6
		col := idx % 6
		p.grid.Attach(swatch, col, row, 1, 1)
	}

	return p
}

func drawSwatch(cr *cairo.Context, w, h int, c cfg.Color, selected bool) {
	cx := float64(w) / 2
	cy := float64(h) / 2
	radius := math.Min(float64(w), float64(h))/2 - 2

	// Draw filled circle.
	cr.SetSourceRGBA(c.R, c.G, c.B, c.A)
	cr.Arc(cx, cy, radius, 0, 2*math.Pi)
	cr.Fill()

	// Draw border (white for dark colors, dark for light colors).
	if selected {
		cr.SetLineWidth(4)
		cr.SetSourceRGBA(1, 1, 1, 1)
		cr.Arc(cx, cy, radius, 0, 2*math.Pi)
		cr.Stroke()
		cr.SetLineWidth(2)
		cr.SetSourceRGBA(0, 0, 0, 1)
		cr.Arc(cx, cy, radius+1, 0, 2*math.Pi)
		cr.Stroke()
	} else {
		cr.SetLineWidth(1)
		cr.SetSourceRGBA(0.3, 0.3, 0.3, 0.5)
		cr.Arc(cx, cy, radius, 0, 2*math.Pi)
		cr.Stroke()
	}
}
