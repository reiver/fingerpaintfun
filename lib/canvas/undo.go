package canvas

import (
	"fingerpaintfun/lib/stroke"
)

// CommitStroke appends a stroke to the committed list and clears the redo stack.
// If mirror mode is active, also commits a horizontally mirrored copy.
// Mirrored strokes are marked so undo removes both as a pair.
func CommitStroke(state *CanvasState, s stroke.Stroke, canvasWidth float64) {
	s.Mirrored = false
	state.Committed = append(state.Committed, s)

	if state.MirrorMode && canvasWidth > 0 {
		ms := mirrorStroke(s, canvasWidth)
		ms.Mirrored = true
		state.Committed = append(state.Committed, ms)
	}

	state.RedoStack = state.RedoStack[:0]
}

// Undo removes the last committed stroke and pushes it onto the redo stack.
// If the last stroke is a mirrored copy, also removes its original (and vice versa).
// Returns true if a stroke was undone.
func Undo(state *CanvasState) bool {
	n := len(state.Committed)
	if n == 0 {
		return false
	}

	// If last stroke is mirrored, remove both mirror and original.
	if state.Committed[n-1].Mirrored && n >= 2 {
		pair := state.Committed[n-2:]
		state.Committed = state.Committed[:n-2]
		state.RedoStack = append(state.RedoStack, pair...)
		return true
	}

	last := state.Committed[n-1]
	state.Committed = state.Committed[:n-1]
	state.RedoStack = append(state.RedoStack, last)
	return true
}

// Redo restores the last undone stroke back to the committed list.
// Handles mirrored pairs.
// Returns true if a stroke was redone.
func Redo(state *CanvasState) bool {
	n := len(state.RedoStack)
	if n == 0 {
		return false
	}

	// If there's a mirrored pair at the end, restore both.
	if state.RedoStack[n-1].Mirrored && n >= 2 {
		pair := state.RedoStack[n-2:]
		state.RedoStack = state.RedoStack[:n-2]
		state.Committed = append(state.Committed, pair...)
		return true
	}

	last := state.RedoStack[n-1]
	state.RedoStack = state.RedoStack[:n-1]
	state.Committed = append(state.Committed, last)
	return true
}

// ClearAll removes all committed strokes and the redo stack.
func ClearAll(state *CanvasState) {
	state.Committed = state.Committed[:0]
	state.RedoStack = state.RedoStack[:0]
}

// mirrorStroke creates a horizontally mirrored copy of a stroke.
func mirrorStroke(s stroke.Stroke, canvasWidth float64) stroke.Stroke {
	ms := stroke.Stroke{
		Color:     s.Color,
		BrushType: s.BrushType,
		BrushSize: s.BrushSize,
		Points:    make([]stroke.Point, len(s.Points)),
	}
	for i, p := range s.Points {
		ms.Points[i] = stroke.Point{
			X: canvasWidth - p.X,
			Y: p.Y,
		}
	}
	return ms
}
