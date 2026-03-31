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

	dragStartX float64
	dragStartY float64
	width      int
	height     int
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

	// Mouse/single-pointer drawing via GestureDrag (idiomatic GTK4).
	drag := gtk.NewGestureDrag()
	drag.ConnectDragBegin(func(startX, startY float64) {
		s := c.active.Begin(canvas.MouseSentinel, c.currentColor, c.currentBrush, c.currentSize)
		s.Points = append(s.Points, stroke.Point{X: startX, Y: startY})
		c.dragStartX = startX
		c.dragStartY = startY
		c.area.QueueDraw()
	})
	drag.ConnectDragUpdate(func(offsetX, offsetY float64) {
		x := c.dragStartX + offsetX
		y := c.dragStartY + offsetY
		c.active.Update(canvas.MouseSentinel, x, y)
		c.area.QueueDraw()
	})
	drag.ConnectDragEnd(func(offsetX, offsetY float64) {
		s := c.active.End(canvas.MouseSentinel)
		if s != nil {
			canvas.CommitStroke(c.state, *s)
		}
		c.area.QueueDraw()
	})
	c.area.AddController(drag)

	// Touch multi-touch drawing via EventControllerLegacy.
	legacy := gtk.NewEventControllerLegacy()
	legacy.ConnectEvent(func(event gdk.Eventer) bool {
		return c.handleTouchEvent(event)
	})
	c.area.AddController(legacy)

	return c
}

func (receiver *canvasWidget) handleTouchEvent(event gdk.Eventer) bool {
	e, ok := event.(*gdk.TouchEvent)
	if !ok {
		return false
	}

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

	return false
}

// sequenceToID converts a *gdk.EventSequence to a uintptr for use as a map key.
func sequenceToID(seq *gdk.EventSequence) uintptr {
	if seq == nil {
		return canvas.MouseSentinel
	}
	return uintptr(unsafe.Pointer(seq))
}

