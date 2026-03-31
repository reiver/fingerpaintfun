package canvas

import (
	"fingerpaintfun/lib/stroke"
)

// MouseSentinel is the reserved key for mouse input in the active strokes map.
// Real *gdk.EventSequence pointers are heap-allocated and will never be zero.
const MouseSentinel uintptr = 0

// ActiveStrokes tracks in-progress strokes keyed by an opaque touch sequence ID.
type ActiveStrokes struct {
	strokes map[uintptr]*stroke.Stroke
}

// NewActiveStrokes creates a new ActiveStrokes tracker.
func NewActiveStrokes() *ActiveStrokes {
	return &ActiveStrokes{
		strokes: make(map[uintptr]*stroke.Stroke),
	}
}

// Begin starts a new stroke for the given sequence ID.
func (receiver *ActiveStrokes) Begin(id uintptr, color [4]float64, brushType int, brushSize float64) *stroke.Stroke {
	s := &stroke.Stroke{
		Color:     color,
		BrushType: brushType,
		BrushSize: brushSize,
	}
	receiver.strokes[id] = s
	return s
}

// Update appends a point to the active stroke for the given sequence ID.
func (receiver *ActiveStrokes) Update(id uintptr, x, y float64) {
	s, ok := receiver.strokes[id]
	if !ok {
		return
	}
	s.Points = append(s.Points, stroke.Point{X: x, Y: y})
}

// End removes and returns the completed stroke for the given sequence ID.
// Returns nil if no stroke exists for that ID.
func (receiver *ActiveStrokes) End(id uintptr) *stroke.Stroke {
	s, ok := receiver.strokes[id]
	if !ok {
		return nil
	}
	delete(receiver.strokes, id)
	return s
}

// All returns all in-progress strokes for rendering.
func (receiver *ActiveStrokes) All() []*stroke.Stroke {
	result := make([]*stroke.Stroke, 0, len(receiver.strokes))
	for _, s := range receiver.strokes {
		result = append(result, s)
	}
	return result
}
