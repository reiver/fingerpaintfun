package ui

import (
	"math"
	"math/rand"

	"fingerpaintfun/cfg"
	"fingerpaintfun/lib/brush"
	"fingerpaintfun/lib/stroke"

	"github.com/diamondburned/gotk4/pkg/cairo"
)

// RenderStroke draws a stroke using the appropriate brush renderer.
func RenderStroke(cr *cairo.Context, s *stroke.Stroke) {
	if len(s.Points) == 0 {
		return
	}

	// Interpolate points to fill gaps from fast swipes.
	points := brush.Interpolate(s.Points, s.BrushSize/2)

	switch s.BrushType {
	case cfg.BrushCrayon:
		renderCrayon(cr, points, s.Color, s.BrushSize)
	case cfg.BrushMarker:
		renderMarker(cr, points, s.Color, s.BrushSize)
	case cfg.BrushPencil:
		renderPencil(cr, points, s.Color, s.BrushSize)
	case cfg.BrushNeon:
		renderNeon(cr, points, s.Color, s.BrushSize)
	case cfg.BrushRainbow:
		renderRainbow(cr, points, s.BrushSize)
	case cfg.BrushSpray:
		renderSpray(cr, points, s.Color, s.BrushSize)
	case cfg.BrushChalk:
		renderChalk(cr, points, s.Color, s.BrushSize)
	case cfg.BrushWatercolor:
		renderWatercolor(cr, points, s.Color, s.BrushSize)
	default:
		renderRound(cr, points, s.Color, s.BrushSize)
	}
}

func renderRound(cr *cairo.Context, points []stroke.Point, color [4]float64, size float64) {
	cr.SetSourceRGBA(color[0], color[1], color[2], color[3])
	radius := size / 2
	for _, p := range points {
		cr.Arc(p.X, p.Y, radius, 0, 2*math.Pi)
		cr.Fill()
	}
}

func renderCrayon(cr *cairo.Context, points []stroke.Point, color [4]float64, size float64) {
	params := brush.GetParams(cfg.BrushCrayon)
	radius := size / 2
	// Use a deterministic seed based on the first point so redraws look the same.
	seed := int64(0)
	if len(points) > 0 {
		seed = int64(points[0].X*1000) + int64(points[0].Y*7)
	}
	rng := rand.New(rand.NewSource(seed))

	for _, p := range points {
		for sub := 0; sub < params.SubStrokes; sub++ {
			jx := (rng.Float64()*2 - 1) * params.Jitter * size
			jy := (rng.Float64()*2 - 1) * params.Jitter * size
			alpha := params.Alpha + rng.Float64()*0.15
			cr.SetSourceRGBA(color[0], color[1], color[2], alpha)
			cr.Arc(p.X+jx, p.Y+jy, radius*0.8, 0, 2*math.Pi)
			cr.Fill()
		}
	}
}

func renderMarker(cr *cairo.Context, points []stroke.Point, color [4]float64, size float64) {
	params := brush.GetParams(cfg.BrushMarker)
	cr.SetOperator(cairo.OperatorOver)
	cr.SetSourceRGBA(color[0], color[1], color[2], params.Alpha)
	radius := size / 2
	for _, p := range points {
		cr.Arc(p.X, p.Y, radius, 0, 2*math.Pi)
		cr.Fill()
	}
}

func renderPencil(cr *cairo.Context, points []stroke.Point, color [4]float64, size float64) {
	cr.SetSourceRGBA(color[0], color[1], color[2], color[3])
	cr.SetLineWidth(math.Max(size/4, 1))
	if len(points) > 0 {
		cr.MoveTo(points[0].X, points[0].Y)
		for _, p := range points[1:] {
			cr.LineTo(p.X, p.Y)
		}
		cr.Stroke()
	}
}

func renderNeon(cr *cairo.Context, points []stroke.Point, color [4]float64, size float64) {
	// Glow layer: wider, translucent.
	cr.SetSourceRGBA(color[0], color[1], color[2], 0.15)
	glowRadius := size * 1.5
	for _, p := range points {
		cr.Arc(p.X, p.Y, glowRadius, 0, 2*math.Pi)
		cr.Fill()
	}
	// Sharp core.
	cr.SetSourceRGBA(color[0], color[1], color[2], 1.0)
	radius := size / 2
	for _, p := range points {
		cr.Arc(p.X, p.Y, radius, 0, 2*math.Pi)
		cr.Fill()
	}
}

func renderRainbow(cr *cairo.Context, points []stroke.Point, size float64) {
	radius := size / 2
	var cumDist float64
	for i, p := range points {
		if i > 0 {
			dx := p.X - points[i-1].X
			dy := p.Y - points[i-1].Y
			cumDist += math.Sqrt(dx*dx + dy*dy)
		}
		// Cycle hue based on cumulative distance.
		hue := math.Mod(cumDist/50.0, 1.0)
		r, g, b := hsvToRGB(hue, 1.0, 1.0)
		cr.SetSourceRGBA(r, g, b, 1.0)
		cr.Arc(p.X, p.Y, radius, 0, 2*math.Pi)
		cr.Fill()
	}
}

func renderSpray(cr *cairo.Context, points []stroke.Point, color [4]float64, size float64) {
	seed := int64(0)
	if len(points) > 0 {
		seed = int64(points[0].X*1000) + int64(points[0].Y*7)
	}
	rng := rand.New(rand.NewSource(seed))

	cr.SetSourceRGBA(color[0], color[1], color[2], 0.6)
	for _, p := range points {
		dots := int(size)
		for d := 0; d < dots; d++ {
			angle := rng.Float64() * 2 * math.Pi
			dist := rng.Float64() * size
			dx := math.Cos(angle) * dist
			dy := math.Sin(angle) * dist
			cr.Arc(p.X+dx, p.Y+dy, 1.0, 0, 2*math.Pi)
			cr.Fill()
		}
	}
}

func renderChalk(cr *cairo.Context, points []stroke.Point, color [4]float64, size float64) {
	seed := int64(0)
	if len(points) > 0 {
		seed = int64(points[0].X*1000) + int64(points[0].Y*7)
	}
	rng := rand.New(rand.NewSource(seed))
	radius := size / 2

	for _, p := range points {
		for d := 0; d < 6; d++ {
			if rng.Float64() < 0.3 {
				continue // gaps for chalky texture
			}
			jx := (rng.Float64()*2 - 1) * radius
			jy := (rng.Float64()*2 - 1) * radius
			alpha := 0.2 + rng.Float64()*0.4
			cr.SetSourceRGBA(color[0], color[1], color[2], alpha)
			cr.Arc(p.X+jx, p.Y+jy, 1.5, 0, 2*math.Pi)
			cr.Fill()
		}
	}
}

func renderWatercolor(cr *cairo.Context, points []stroke.Point, color [4]float64, size float64) {
	seed := int64(0)
	if len(points) > 0 {
		seed = int64(points[0].X*1000) + int64(points[0].Y*7)
	}
	rng := rand.New(rand.NewSource(seed))

	cr.SetOperator(cairo.OperatorOver)
	for _, p := range points {
		alpha := 0.08 + rng.Float64()*0.12
		radius := (size / 2) + rng.Float64()*(size/4)
		cr.SetSourceRGBA(color[0], color[1], color[2], alpha)
		cr.Arc(p.X, p.Y, radius, 0, 2*math.Pi)
		cr.Fill()
	}
}

// hsvToRGB converts HSV (h in [0,1], s in [0,1], v in [0,1]) to RGB.
func hsvToRGB(h, s, v float64) (r, g, b float64) {
	i := int(h * 6)
	f := h*6 - float64(i)
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)
	switch i % 6 {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}
	return
}
