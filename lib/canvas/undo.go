package canvas

import (
	"fingerpaintfun/lib/stroke"
)

// CommitStroke appends a stroke to the committed list and clears the redo stack.
func CommitStroke(state *CanvasState, s stroke.Stroke) {
	state.Committed = append(state.Committed, s)
	state.RedoStack = state.RedoStack[:0]
}

// Undo removes the last committed stroke and pushes it onto the redo stack.
// Returns true if a stroke was undone.
func Undo(state *CanvasState) bool {
	n := len(state.Committed)
	if n == 0 {
		return false
	}
	last := state.Committed[n-1]
	state.Committed = state.Committed[:n-1]
	state.RedoStack = append(state.RedoStack, last)
	return true
}

// Redo restores the last undone stroke back to the committed list.
// Returns true if a stroke was redone.
func Redo(state *CanvasState) bool {
	n := len(state.RedoStack)
	if n == 0 {
		return false
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
