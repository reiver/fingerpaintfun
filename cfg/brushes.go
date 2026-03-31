package cfg

// Brush type constants.
const (
	BrushRound  = iota
	BrushCrayon
	BrushMarker
	BrushPencil
	BrushNeon
	BrushRainbow
	BrushSpray
	BrushChalk
	BrushWatercolor
)

// Brush size constants (in pixels).
const (
	SizeSmall  float64 = 4
	SizeMedium float64 = 12
	SizeLarge  float64 = 28
)

// DefaultBrush is the initial brush type.
const DefaultBrush = BrushRound

// DefaultSize is the initial brush size.
const DefaultSize = SizeLarge
