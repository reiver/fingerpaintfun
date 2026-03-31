package stroke

// Point represents an (x, y) coordinate on the canvas.
type Point struct {
	X, Y float64
}

// Stroke represents a single drawing stroke.
type Stroke struct {
	Color     [4]float64 // RGBA, each in [0, 1]
	BrushType int
	BrushSize float64
	Points    []Point
}
