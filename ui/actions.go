package ui

import (
	"fingerpaintfun/lib/canvas"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// actionsWidget holds undo, redo, and clear buttons.
type actionsWidget struct {
	box      *gtk.Box
	undoBtn  *gtk.Button
	redoBtn  *gtk.Button
	clearBtn *gtk.Button
	canvas   *canvasWidget
}

func newActionsWidget(c *canvasWidget) *actionsWidget {
	a := &actionsWidget{canvas: c}

	a.box = gtk.NewBox(gtk.OrientationHorizontal, 4)

	// Undo button.
	a.undoBtn = gtk.NewButton()
	a.undoBtn.SetIconName("edit-undo-symbolic")
	a.undoBtn.SetSizeRequest(48, 48)
	a.undoBtn.AddCSSClass("flat")
	a.undoBtn.ConnectClicked(func() {
		canvas.Undo(c.state)
		c.renderer.InvalidateCache()
		c.area.QueueDraw()
		a.updateSensitivity()
		PlayWhoosh()
	})
	a.box.Append(a.undoBtn)

	// Redo button.
	a.redoBtn = gtk.NewButton()
	a.redoBtn.SetIconName("edit-redo-symbolic")
	a.redoBtn.SetSizeRequest(48, 48)
	a.redoBtn.AddCSSClass("flat")
	a.redoBtn.ConnectClicked(func() {
		canvas.Redo(c.state)
		c.renderer.InvalidateCache()
		c.area.QueueDraw()
		a.updateSensitivity()
		PlayWhoosh()
	})
	a.box.Append(a.redoBtn)

	// Clear button (long-press only).
	a.clearBtn = gtk.NewButton()
	a.clearBtn.SetIconName("edit-clear-all-symbolic")
	a.clearBtn.SetSizeRequest(48, 48)
	a.clearBtn.AddCSSClass("flat")

	longPress := gtk.NewGestureLongPress()
	longPress.ConnectPressed(func(x, y float64) {
		canvas.ClearAll(c.state)
		c.renderer.InvalidateCache()
		c.area.QueueDraw()
		a.updateSensitivity()
		PlaySplash()
	})
	a.clearBtn.AddController(longPress)

	a.updateSensitivity()

	return a
}

func (receiver *actionsWidget) updateSensitivity() {
	receiver.undoBtn.SetSensitive(len(receiver.canvas.state.Committed) > 0)
	receiver.redoBtn.SetSensitive(len(receiver.canvas.state.RedoStack) > 0)
	receiver.clearBtn.SetSensitive(len(receiver.canvas.state.Committed) > 0 || len(receiver.canvas.state.RedoStack) > 0)
}
