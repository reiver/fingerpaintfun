package ui

import (
	"math"

	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// Placeholder templates drawn with Cairo. Replace with PNG assets later.

// TemplateDef defines a coloring page template.
type TemplateDef struct {
	Name string
	Draw func(cr *cairo.Context, width, height float64)
}

// Templates is the list of available coloring page templates.
// Index 0 is "Blank" (no template).
var Templates = []TemplateDef{
	{"Blank", nil},
	{"Circle", drawTemplateCircle},
	{"House", drawTemplateHouse},
	{"Fish", drawTemplateFish},
	{"Star", drawTemplateStar},
}

// RenderTemplate draws the template outline onto the canvas.
func RenderTemplate(cr *cairo.Context, templateID int, width, height float64) {
	if templateID <= 0 || templateID >= len(Templates) || Templates[templateID].Draw == nil {
		return
	}
	cr.Save()
	cr.SetSourceRGBA(0, 0, 0, 0.6)
	cr.SetLineWidth(3)
	Templates[templateID].Draw(cr, width, height)
	cr.Restore()
}

func drawTemplateCircle(cr *cairo.Context, w, h float64) {
	cx, cy := w/2, h/2
	r := math.Min(w, h)*0.35
	cr.Arc(cx, cy, r, 0, 2*math.Pi)
	cr.Stroke()
}

func drawTemplateHouse(cr *cairo.Context, w, h float64) {
	// Body.
	bx := w * 0.25
	by := h * 0.45
	bw := w * 0.5
	bh := h * 0.4
	cr.Rectangle(bx, by, bw, bh)
	cr.Stroke()
	// Roof.
	cr.MoveTo(bx-w*0.05, by)
	cr.LineTo(w/2, h*0.2)
	cr.LineTo(bx+bw+w*0.05, by)
	cr.ClosePath()
	cr.Stroke()
	// Door.
	dw := bw * 0.25
	dh := bh * 0.5
	cr.Rectangle(w/2-dw/2, by+bh-dh, dw, dh)
	cr.Stroke()
	// Window.
	ww := bw * 0.2
	cr.Rectangle(bx+bw*0.15, by+bh*0.2, ww, ww)
	cr.Stroke()
}

func drawTemplateFish(cr *cairo.Context, w, h float64) {
	cx, cy := w/2, h/2
	// Body (ellipse).
	cr.Save()
	cr.Translate(cx, cy)
	cr.Scale(1, 0.5)
	cr.Arc(0, 0, w*0.25, 0, 2*math.Pi)
	cr.Restore()
	cr.Stroke()
	// Tail.
	tx := cx + w*0.25
	cr.MoveTo(tx, cy)
	cr.LineTo(tx+w*0.1, cy-h*0.1)
	cr.LineTo(tx+w*0.1, cy+h*0.1)
	cr.ClosePath()
	cr.Stroke()
	// Eye.
	cr.Arc(cx-w*0.1, cy-h*0.03, w*0.02, 0, 2*math.Pi)
	cr.Stroke()
}

func drawTemplateStar(cr *cairo.Context, w, h float64) {
	cx, cy := w/2, h/2
	outerR := math.Min(w, h) * 0.35
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
	cr.Stroke()
}

// templatePickerWidget lets the user select a coloring page template.
type templatePickerWidget struct {
	grid     *gtk.Grid
	selected int
	buttons  []*gtk.DrawingArea
	onChange func(templateID int)
}

func newTemplatePickerWidget(onChange func(int)) *templatePickerWidget {
	tp := &templatePickerWidget{
		grid:     gtk.NewGrid(),
		selected: 0, // Blank
		onChange: onChange,
	}

	tp.grid.SetColumnSpacing(8)
	tp.grid.SetRowSpacing(8)
	tp.grid.SetHAlign(gtk.AlignCenter)
	tp.grid.SetVAlign(gtk.AlignCenter)

	for i, tmpl := range Templates {
		idx := i
		drawFn := tmpl.Draw

		btn := gtk.NewDrawingArea()
		btn.SetSizeRequest(48, 48)
		btn.SetDrawFunc(func(area *gtk.DrawingArea, cr *cairo.Context, w, h int) {
			if idx == tp.selected {
				cr.SetSourceRGBA(0.3, 0.5, 1.0, 0.3)
				cr.Rectangle(0, 0, float64(w), float64(h))
				cr.Fill()
			}

			if drawFn != nil {
				cr.SetSourceRGBA(0, 0, 0, 0.8)
				cr.SetLineWidth(1.5)
				drawFn(cr, float64(w), float64(h))
			} else {
				// "Blank" label — just a border.
				cr.SetSourceRGBA(0.5, 0.5, 0.5, 0.5)
				cr.SetLineWidth(1)
				cr.Rectangle(4, 4, float64(w)-8, float64(h)-8)
				cr.Stroke()
			}
		})

		click := gtk.NewGestureClick()
		click.ConnectPressed(func(nPress int, x, y float64) {
			tp.selected = idx
			if tp.onChange != nil {
				tp.onChange(idx)
			}
			PlayPop()
			for _, b := range tp.buttons {
				b.QueueDraw()
			}
		})
		btn.AddController(click)

		tp.buttons = append(tp.buttons, btn)
		tp.grid.Attach(btn, i, 0, 1, 1)
	}

	return tp
}
