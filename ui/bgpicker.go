package ui

import (
	"math"

	"fingerpaintfun/cfg"

	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// bgPickerWidget lets the user select a canvas background color.
type bgPickerWidget struct {
	grid     *gtk.Grid
	selected int
	swatches []*gtk.DrawingArea
	onChange func(color [4]float64)
}

func newBgPickerWidget(onChange func(color [4]float64)) *bgPickerWidget {
	bp := &bgPickerWidget{
		grid:     gtk.NewGrid(),
		selected: 9, // White (index 9 in palette)
		onChange: onChange,
	}

	bp.grid.SetRowSpacing(8)
	bp.grid.SetColumnSpacing(8)
	bp.grid.SetHAlign(gtk.AlignCenter)
	bp.grid.SetVAlign(gtk.AlignCenter)

	for i, color := range cfg.Palette {
		idx := i
		c := color

		swatch := gtk.NewDrawingArea()
		swatch.SetSizeRequest(48, 48)
		swatch.SetDrawFunc(func(area *gtk.DrawingArea, cr *cairo.Context, w, h int) {
			drawBgSwatch(cr, w, h, c, idx == bp.selected)
		})

		click := gtk.NewGestureClick()
		click.ConnectPressed(func(nPress int, x, y float64) {
			bp.selected = idx
			if bp.onChange != nil {
				bp.onChange([4]float64{c.R, c.G, c.B, c.A})
			}
			PlayPop()
			for _, sw := range bp.swatches {
				sw.QueueDraw()
			}
		})
		swatch.AddController(click)

		bp.swatches = append(bp.swatches, swatch)
		row := idx / 6
		col := idx % 6
		bp.grid.Attach(swatch, col, row, 1, 1)
	}

	return bp
}

func drawBgSwatch(cr *cairo.Context, w, h int, c cfg.Color, selected bool) {
	cx := float64(w) / 2
	cy := float64(h) / 2
	radius := math.Min(float64(w), float64(h))/2 - 2

	// Square with rounded corners to distinguish from the drawing palette circles.
	cr.SetSourceRGBA(c.R, c.G, c.B, c.A)
	cr.Rectangle(cx-radius, cy-radius, radius*2, radius*2)
	cr.Fill()

	if selected {
		cr.SetLineWidth(3)
		cr.SetSourceRGBA(1, 0.4, 0, 1) // orange highlight
		cr.Rectangle(cx-radius, cy-radius, radius*2, radius*2)
		cr.Stroke()
	} else {
		cr.SetLineWidth(1)
		cr.SetSourceRGBA(0.3, 0.3, 0.3, 0.5)
		cr.Rectangle(cx-radius, cy-radius, radius*2, radius*2)
		cr.Stroke()
	}
}
