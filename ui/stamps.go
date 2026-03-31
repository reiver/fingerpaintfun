package ui

import (
	"math"

	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// Placeholder stamps drawn with Cairo. Replace with PNG/SVG assets later.

// StampDef defines a stamp with a name and a draw function.
type StampDef struct {
	Name string
	Draw func(cr *cairo.Context, cx, cy, size float64)
}

// Stamps is the list of available stamps.
var Stamps = []StampDef{
	{"Circle", drawStampCircle},
	{"Square", drawStampSquare},
	{"Triangle", drawStampTriangle},
	{"Star", drawStampStar},
	{"Heart", drawStampHeart},
}

// RenderStamp draws a stamp at (cx, cy) with the given size and color.
func RenderStamp(cr *cairo.Context, stampID int, cx, cy, size float64, color [4]float64) {
	if stampID < 0 || stampID >= len(Stamps) {
		return
	}
	cr.Save()
	cr.SetSourceRGBA(color[0], color[1], color[2], color[3])
	Stamps[stampID].Draw(cr, cx, cy, size)
	cr.Restore()
}

func drawStampCircle(cr *cairo.Context, cx, cy, size float64) {
	cr.Arc(cx, cy, size/2, 0, 2*math.Pi)
	cr.Fill()
}

func drawStampSquare(cr *cairo.Context, cx, cy, size float64) {
	half := size / 2
	cr.Rectangle(cx-half, cy-half, size, size)
	cr.Fill()
}

func drawStampTriangle(cr *cairo.Context, cx, cy, size float64) {
	half := size / 2
	cr.MoveTo(cx, cy-half)
	cr.LineTo(cx+half, cy+half)
	cr.LineTo(cx-half, cy+half)
	cr.ClosePath()
	cr.Fill()
}

func drawStampStar(cr *cairo.Context, cx, cy, size float64) {
	outerR := size / 2
	innerR := outerR * 0.4
	points := 5
	for i := 0; i < points*2; i++ {
		angle := float64(i)*math.Pi/float64(points) - math.Pi/2
		r := outerR
		if i%2 == 1 {
			r = innerR
		}
		x := cx + r*math.Cos(angle)
		y := cy + r*math.Sin(angle)
		if i == 0 {
			cr.MoveTo(x, y)
		} else {
			cr.LineTo(x, y)
		}
	}
	cr.ClosePath()
	cr.Fill()
}

func drawStampHeart(cr *cairo.Context, cx, cy, size float64) {
	s := size / 2
	cr.MoveTo(cx, cy+s*0.6)
	cr.CurveTo(cx-s*1.2, cy-s*0.2, cx-s*0.5, cy-s, cx, cy-s*0.4)
	cr.CurveTo(cx+s*0.5, cy-s, cx+s*1.2, cy-s*0.2, cx, cy+s*0.6)
	cr.Fill()
}

// stampPickerWidget lets the user select a stamp.
type stampPickerWidget struct {
	grid     *gtk.Grid
	selected int
	buttons  []*gtk.DrawingArea
	onChange func(stampID int)
}

func newStampPickerWidget(onChange func(int)) *stampPickerWidget {
	sp := &stampPickerWidget{
		grid:     gtk.NewGrid(),
		selected: 0,
		onChange: onChange,
	}

	sp.grid.SetColumnSpacing(8)
	sp.grid.SetRowSpacing(8)
	sp.grid.SetHAlign(gtk.AlignCenter)
	sp.grid.SetVAlign(gtk.AlignCenter)

	for i, stamp := range Stamps {
		idx := i
		drawFn := stamp.Draw

		btn := gtk.NewDrawingArea()
		btn.SetSizeRequest(48, 48)
		btn.SetDrawFunc(func(area *gtk.DrawingArea, cr *cairo.Context, w, h int) {
			cxp := float64(w) / 2
			cyp := float64(h) / 2

			if idx == sp.selected {
				cr.SetSourceRGBA(0.3, 0.5, 1.0, 0.3)
				cr.Rectangle(0, 0, float64(w), float64(h))
				cr.Fill()
			}

			cr.SetSourceRGBA(0, 0, 0, 1)
			drawFn(cr, cxp, cyp, 24)
		})

		click := gtk.NewGestureClick()
		click.ConnectPressed(func(nPress int, x, y float64) {
			sp.selected = idx
			if sp.onChange != nil {
				sp.onChange(idx)
			}
			PlayPop()
			for _, b := range sp.buttons {
				b.QueueDraw()
			}
		})
		btn.AddController(click)

		sp.buttons = append(sp.buttons, btn)
		sp.grid.Attach(btn, i, 0, 1, 1)
	}

	return sp
}
