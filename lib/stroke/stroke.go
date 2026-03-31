package stroke

// Point represents an (x, y) coordinate on the canvas.
type Point struct {
	X, Y float64
}

// Stroke types.
const (
	TypeDraw  = iota // normal brush stroke
	TypeStamp        // stamp placement
	TypeFill         // flood-fill operation
)

// Stroke represents a single drawing action on the canvas.
type Stroke struct {
	Type      int        // TypeDraw, TypeStamp, or TypeFill
	Color     [4]float64 // RGBA, each in [0, 1]
	BrushType int
	BrushSize float64
	Points    []Point
	Mirrored  bool // true if this is a mirror-mode copy (undo removes both original + copy)

	// Stamp-specific fields.
	StampID int   // index into the stamp list
	StampX  float64
	StampY  float64

	// Fill-specific fields.
	FillX int // pixel coordinate of fill origin
	FillY int
}
