package canvas

import (
	"fingerpaintfun/lib/stroke"
)

// CanvasState holds all drawing state for the canvas.
type CanvasState struct {
	Committed []stroke.Stroke // finished strokes (undo operates on this)
	RedoStack []stroke.Stroke // strokes removed by undo
	BgColor   [4]float64      // background RGBA, default white
}

// NewCanvasState creates a new canvas state with a white background.
func NewCanvasState() *CanvasState {
	return &CanvasState{
		BgColor: [4]float64{1, 1, 1, 1},
	}
}
