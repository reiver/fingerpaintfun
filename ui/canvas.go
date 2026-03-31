package ui

import (
	"unsafe"

	"fingerpaintfun/cfg"
	"fingerpaintfun/lib/canvas"
	"fingerpaintfun/lib/palm"
	"fingerpaintfun/lib/stroke"

	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// canvasWidget wraps a gtk.DrawingArea with canvas state and input handling.
type canvasWidget struct {
	area     *gtk.DrawingArea
	state    *canvas.CanvasState
	active   *canvas.ActiveStrokes
	renderer *renderer
	suspect  map[uintptr]bool // sequences flagged by palm rejection

	currentColor [4]float64
	currentBrush int
	currentSize  float64

	mouseDown bool
	width     int
	height    int
}

func newCanvasWidget() *canvasWidget {
	c := &canvasWidget{
		area:     gtk.NewDrawingArea(),
		state:    canvas.NewCanvasState(),
		active:   canvas.NewActiveStrokes(),
		renderer: &renderer{},
		suspect:  make(map[uintptr]bool),

		currentColor: [4]float64{
			cfg.Palette[cfg.DefaultColor].R,
			cfg.Palette[cfg.DefaultColor].G,
			cfg.Palette[cfg.DefaultColor].B,
			cfg.Palette[cfg.DefaultColor].A,
		},
		currentBrush: cfg.DefaultBrush,
		currentSize:  cfg.DefaultSize,
	}

	c.area.SetDrawFunc(func(area *gtk.DrawingArea, cr *cairo.Context, width, height int) {
		c.width = width
		c.height = height
		c.renderer.Render(cr, width, height, c.state, c.active)
	})

	// Single event controller for both touch and mouse.
	legacy := gtk.NewEventControllerLegacy()
	legacy.ConnectEvent(func(event gdk.Eventer) bool {
		return c.handleEvent(event)
	})
	c.area.AddController(legacy)

	return c
}

func (receiver *canvasWidget) handleEvent(event gdk.Eventer) bool {
	switch e := event.(type) {

	// --- Mouse input ---
	case *gdk.ButtonEvent:
		if e.EventType() == gdk.ButtonPress {
			x, y, _ := e.Position()
			receiver.mouseDown = true
			s := receiver.active.Begin(canvas.MouseSentinel, receiver.currentColor, receiver.currentBrush, receiver.currentSize)
			s.Points = append(s.Points, stroke.Point{X: x, Y: y})
			receiver.area.QueueDraw()
			return true
		}
		if e.EventType() == gdk.ButtonRelease {
			receiver.mouseDown = false
			s := receiver.active.End(canvas.MouseSentinel)
			if s != nil {
				canvas.CommitStroke(receiver.state, *s)
			}
			receiver.area.QueueDraw()
			return true
		}

	case *gdk.MotionEvent:
		if receiver.mouseDown {
			x, y, _ := e.Position()
			receiver.active.Update(canvas.MouseSentinel, x, y)
			receiver.area.QueueDraw()
			return true
		}

	// --- Touch input ---
	case *gdk.TouchEvent:
		seq := e.EventSequence()
		id := sequenceToID(seq)
		x, y, _ := e.Position()

		switch e.EventType() {
		case gdk.TouchBegin:
			if palm.IsEdgeZone(x, y, receiver.width, receiver.height) {
				receiver.suspect[id] = true
			}
			s := receiver.active.Begin(id, receiver.currentColor, receiver.currentBrush, receiver.currentSize)
			s.Points = append(s.Points, stroke.Point{X: x, Y: y})
			receiver.area.QueueDraw()
			return true

		case gdk.TouchUpdate:
			receiver.active.Update(id, x, y)
			receiver.area.QueueDraw()
			return true

		case gdk.TouchEnd:
			s := receiver.active.End(id)
			if s != nil {
				discard := receiver.suspect[id] || palm.ShouldDiscard(s)
				if !discard {
					canvas.CommitStroke(receiver.state, *s)
				}
			}
			delete(receiver.suspect, id)
			receiver.area.QueueDraw()
			return true
		}
	}

	return false
}

// sequenceToID converts a *gdk.EventSequence to a uintptr for use as a map key.
func sequenceToID(seq *gdk.EventSequence) uintptr {
	if seq == nil {
		return canvas.MouseSentinel
	}
	return uintptr(unsafe.Pointer(seq))
}

