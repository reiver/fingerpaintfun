package cfg

// Color represents an RGBA color with components in the range [0, 1].
type Color struct {
	R, G, B, A float64
}

// Palette is the fixed 12-color palette for the app.
var Palette = []Color{
	{1.0, 0.0, 0.0, 1.0},       // Red
	{1.0, 0.6, 0.0, 1.0},       // Orange
	{1.0, 1.0, 0.0, 1.0},       // Yellow
	{0.0, 0.8, 0.0, 1.0},       // Green
	{0.0, 0.8, 0.8, 1.0},       // Cyan
	{0.0, 0.0, 1.0, 1.0},       // Blue
	{0.6, 0.0, 0.8, 1.0},       // Purple
	{1.0, 0.4, 0.7, 1.0},       // Pink
	{0.6, 0.3, 0.0, 1.0},       // Brown
	{1.0, 1.0, 1.0, 1.0},       // White
	{0.6, 0.6, 0.6, 1.0},       // Grey
	{0.0, 0.0, 0.0, 1.0},       // Black
}

// DefaultColor is the initial drawing color index (Black).
const DefaultColor = 11
