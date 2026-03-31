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

	currentColor    [4]float64
	currentBrush    int
	currentSize     float64
	currentStamp    int // selected stamp ID for BrushStamp mode
	currentTemplate int // selected template ID (0 = blank)
	mascot          *mascot

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
		mascot:       newMascot(),
	}

	c.area.SetDrawFunc(func(area *gtk.DrawingArea, cr *cairo.Context, width, height int) {
		c.width = width
		c.height = height
		c.renderer.Render(cr, width, height, c.state, c.active, c.mascot)
	})

	// Mouse/single-pointer drawing via GestureDrag (idiomatic GTK4).
	drag := gtk.NewGestureDrag()
	drag.ConnectDragBegin(func(startX, startY float64) {
		// Fill mode: tap to flood-fill.
		if c.currentBrush == cfg.BrushFill {
			go floodFill(c, int(startX), int(startY), c.currentColor)
			return
		}
		// Stamp mode: place stamp on tap, don't start a drag stroke.
		if c.currentBrush == cfg.BrushStamp {
			st := stroke.Stroke{
				Type:    stroke.TypeStamp,
				Color:   c.currentColor,
				StampID: c.currentStamp,
				StampX:  startX,
				StampY:  startY,
			}
			canvas.CommitStroke(c.state, st, float64(c.width))
			c.renderer.InvalidateCache()
			c.area.QueueDraw()
			PlayPop()
			return
		}
		s := c.active.Begin(canvas.MouseSentinel, c.strokeColor(), c.currentBrush, c.currentSize)
		s.Points = append(s.Points, stroke.Point{X: startX, Y: startY})
		c.dragStartX = startX
		c.dragStartY = startY
		c.mascot.SetState(mascotDrawing)
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
			canvas.CommitStroke(c.state, *s, float64(c.width))
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

// strokeColor returns the background color for eraser, or the current color for other brushes.
func (receiver *canvasWidget) strokeColor() [4]float64 {
	if receiver.currentBrush == cfg.BrushEraser {
		return receiver.state.BgColor
	}
	return receiver.currentColor
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
		s := receiver.active.Begin(id, receiver.strokeColor(), receiver.currentBrush, receiver.currentSize)
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
				canvas.CommitStroke(receiver.state, *s, float64(receiver.width))
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

