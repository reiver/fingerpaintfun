package brush

import (
	"fingerpaintfun/cfg"
)

// Params holds per-brush-type rendering parameters.
type Params struct {
	Alpha      float64 // opacity (1.0 for round, 0.4 for marker, etc.)
	Jitter     float64 // random offset as a fraction of size (0 for round, 0.25 for crayon)
	SubStrokes int     // number of sub-circles per point (1 for round, 4 for crayon)
}

// GetParams returns the rendering parameters for a brush type.
func GetParams(brushType int) Params {
	switch brushType {
	case cfg.BrushCrayon:
		return Params{Alpha: 0.4, Jitter: 0.25, SubStrokes: 4}
	case cfg.BrushMarker:
		return Params{Alpha: 0.4, Jitter: 0, SubStrokes: 1}
	case cfg.BrushPencil:
		return Params{Alpha: 1.0, Jitter: 0, SubStrokes: 1}
	default: // BrushRound
		return Params{Alpha: 1.0, Jitter: 0, SubStrokes: 1}
	}
}
