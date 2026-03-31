package palm

import (
	"math"

	"fingerpaintfun/lib/stroke"
)

// ShouldDiscard returns true if the stroke looks like a palm touch:
// very few points and negligible total movement.
func ShouldDiscard(s *stroke.Stroke) bool {
	if len(s.Points) >= 3 {
		return false
	}
	return totalMovement(s.Points) < 5.0
}

// IsEdgeZone returns true if the position is within 20px of any screen edge.
func IsEdgeZone(x, y float64, width, height int) bool {
	const margin = 20.0
	return x < margin || y < margin || x > float64(width)-margin || y > float64(height)-margin
}

func totalMovement(points []stroke.Point) float64 {
	if len(points) < 2 {
		return 0
	}
	var total float64
	for i := 1; i < len(points); i++ {
		dx := points[i].X - points[i-1].X
		dy := points[i].Y - points[i-1].Y
		total += math.Sqrt(dx*dx + dy*dy)
	}
	return total
}
