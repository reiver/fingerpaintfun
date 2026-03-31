package brush

import (
	"math"

	"fingerpaintfun/lib/stroke"
)

// Interpolate fills gaps between consecutive points so that fast swipes
// produce continuous strokes. Returns a new slice with intermediate points
// inserted wherever the distance between consecutive points exceeds maxGap.
func Interpolate(points []stroke.Point, maxGap float64) []stroke.Point {
	if len(points) < 2 || maxGap <= 0 {
		return points
	}

	result := make([]stroke.Point, 0, len(points)*2)
	result = append(result, points[0])

	for i := 1; i < len(points); i++ {
		prev := points[i-1]
		curr := points[i]
		dx := curr.X - prev.X
		dy := curr.Y - prev.Y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist > maxGap {
			steps := int(math.Ceil(dist / maxGap))
			for s := 1; s < steps; s++ {
				t := float64(s) / float64(steps)
				result = append(result, stroke.Point{
					X: prev.X + dx*t,
					Y: prev.Y + dy*t,
				})
			}
		}
		result = append(result, curr)
	}

	return result
}
