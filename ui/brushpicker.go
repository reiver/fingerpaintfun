package ui

import (
	"math"

	"fingerpaintfun/cfg"

	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type brushPickerWidget struct {
	grid          *gtk.Grid
	selectedBrush int
	selectedSize  int // 0=small, 1=medium, 2=large
	brushButtons  []*gtk.DrawingArea
	sizeButtons   []*gtk.DrawingArea
	onBrushChange func(brushType int)
	onSizeChange  func(size float64)
}

var brushTypes = []struct {
	Type int
	Name string
}{
	{cfg.BrushRound, "Round"},
	{cfg.BrushCrayon, "Crayon"},
	{cfg.BrushMarker, "Marker"},
	{cfg.BrushNeon, "Neon"},
	{cfg.BrushRainbow, "Rainbow"},
	{cfg.BrushSpray, "Spray"},
	{cfg.BrushChalk, "Chalk"},
	{cfg.BrushWatercolor, "Water"},
	{cfg.BrushEraser, "Eraser"},
}

var brushSizes = []struct {
	Size  float64
	Label string
}{
	{cfg.SizeSmall, "S"},
	{cfg.SizeMedium, "M"},
	{cfg.SizeLarge, "L"},
}

func newBrushPickerWidget(onBrush func(int), onSize func(float64)) *brushPickerWidget {
	bp := &brushPickerWidget{
		grid:          gtk.NewGrid(),
		selectedBrush: 0,
		selectedSize:  2, // Large by default
		onBrushChange: onBrush,
		onSizeChange:  onSize,
	}

	bp.grid.SetRowSpacing(8)
	bp.grid.SetColumnSpacing(12)
	bp.grid.SetHAlign(gtk.AlignCenter)
	bp.grid.SetVAlign(gtk.AlignCenter)

	// Brush types in rows of 4.
	for i, bt := range brushTypes {
		idx := i
		brushType := bt.Type

		btn := gtk.NewDrawingArea()
		btn.SetSizeRequest(40, 40)
		btn.SetDrawFunc(func(area *gtk.DrawingArea, cr *cairo.Context, w, h int) {
			drawBrushPreview(cr, w, h, brushType, idx == bp.selectedBrush)
		})

		click := gtk.NewGestureClick()
		click.ConnectPressed(func(nPress int, x, y float64) {
			bp.selectedBrush = idx
			if bp.onBrushChange != nil {
				bp.onBrushChange(brushType)
			}
			bp.redrawAll()
			PlayPop()
		})
		btn.AddController(click)

		bp.brushButtons = append(bp.brushButtons, btn)
		row := idx / 4
		col := idx % 4
		bp.grid.Attach(btn, col, row, 1, 1)
	}

	// Sizes in the row after brushes.
	sizeRow := (len(brushTypes) + 3) / 4 // next row after brush rows
	for i, bs := range brushSizes {
		idx := i
		size := bs.Size

		btn := gtk.NewDrawingArea()
		btn.SetSizeRequest(40, 40)
		btn.SetDrawFunc(func(area *gtk.DrawingArea, cr *cairo.Context, w, h int) {
			drawSizePreview(cr, w, h, size, idx == bp.selectedSize)
		})

		click := gtk.NewGestureClick()
		click.ConnectPressed(func(nPress int, x, y float64) {
			bp.selectedSize = idx
			if bp.onSizeChange != nil {
				bp.onSizeChange(size)
			}
			bp.redrawAll()
			PlayPop()
		})
		btn.AddController(click)

		bp.sizeButtons = append(bp.sizeButtons, btn)
		bp.grid.Attach(btn, i, sizeRow, 1, 1)
	}

	return bp
}

func (receiver *brushPickerWidget) redrawAll() {
	for _, b := range receiver.brushButtons {
		b.QueueDraw()
	}
	for _, b := range receiver.sizeButtons {
		b.QueueDraw()
	}
}

func drawBrushPreview(cr *cairo.Context, w, h int, brushType int, selected bool) {
	cx := float64(w) / 2
	cy := float64(h) / 2

	if selected {
		cr.SetSourceRGBA(0.3, 0.5, 1.0, 0.3)
		cr.Rectangle(0, 0, float64(w), float64(h))
		cr.Fill()
	}

	cr.SetSourceRGBA(0, 0, 0, 1)
	switch brushType {
	case cfg.BrushCrayon:
		for i := 0; i < 5; i++ {
			ox := float64(i-2) * 5
			cr.Arc(cx+ox, cy, 2.5, 0, 2*math.Pi)
			cr.Fill()
		}
	case cfg.BrushMarker:
		cr.SetSourceRGBA(0, 0, 0, 0.4)
		cr.Rectangle(cx-12, cy-5, 24, 10)
		cr.Fill()
	case cfg.BrushNeon:
		cr.SetSourceRGBA(0.2, 1.0, 0.2, 0.3)
		cr.Arc(cx, cy, 12, 0, 2*math.Pi)
		cr.Fill()
		cr.SetSourceRGBA(0.2, 1.0, 0.2, 1.0)
		cr.Arc(cx, cy, 5, 0, 2*math.Pi)
		cr.Fill()
	case cfg.BrushRainbow:
		for i := 0; i < 5; i++ {
			hue := float64(i) / 5.0
			r, g, b := hsvToRGB(hue, 1, 1)
			cr.SetSourceRGBA(r, g, b, 1)
			cr.Arc(cx+float64(i-2)*6, cy, 3, 0, 2*math.Pi)
			cr.Fill()
		}
	case cfg.BrushSpray:
		cr.SetSourceRGBA(0, 0, 0, 0.5)
		for i := 0; i < 12; i++ {
			ox := float64(i%4-2) * 4
			oy := float64(i/4-1) * 4
			cr.Arc(cx+ox, cy+oy, 1, 0, 2*math.Pi)
			cr.Fill()
		}
	case cfg.BrushChalk:
		cr.SetSourceRGBA(0, 0, 0, 0.3)
		for i := 0; i < 7; i++ {
			ox := float64(i-3) * 4
			if i%2 == 0 {
				cr.Arc(cx+ox, cy, 2, 0, 2*math.Pi)
				cr.Fill()
			}
		}
	case cfg.BrushWatercolor:
		cr.SetSourceRGBA(0.2, 0.4, 1.0, 0.15)
		cr.Arc(cx, cy, 12, 0, 2*math.Pi)
		cr.Fill()
		cr.SetSourceRGBA(0.2, 0.4, 1.0, 0.3)
		cr.Arc(cx-2, cy+1, 8, 0, 2*math.Pi)
		cr.Fill()
	case cfg.BrushEraser:
		// Pink eraser rectangle.
		cr.SetSourceRGBA(1.0, 0.7, 0.7, 1.0)
		cr.Rectangle(cx-10, cy-7, 20, 14)
		cr.Fill()
		cr.SetSourceRGBA(0.5, 0.3, 0.3, 1.0)
		cr.SetLineWidth(1.5)
		cr.Rectangle(cx-10, cy-7, 20, 14)
		cr.Stroke()
	default:
		cr.Arc(cx, cy, 6, 0, 2*math.Pi)
		cr.Fill()
	}
}

func drawSizePreview(cr *cairo.Context, w, h int, size float64, selected bool) {
	cx := float64(w) / 2
	cy := float64(h) / 2

	if selected {
		cr.SetSourceRGBA(0.3, 0.5, 1.0, 0.3)
		cr.Rectangle(0, 0, float64(w), float64(h))
		cr.Fill()
	}

	cr.SetSourceRGBA(0, 0, 0, 1)
	cr.Arc(cx, cy, size/2, 0, 2*math.Pi)
	cr.Fill()
}
