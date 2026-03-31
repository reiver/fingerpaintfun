package ui

import (
	"fingerpaintfun/lib/canvas"
	"fingerpaintfun/lib/stroke"

	cairolib "github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
)

// floodFill performs a BFS flood fill on the canvas at the given pixel coordinate.
// It runs in a goroutine and updates the UI via glib.IdleAdd when done.
func floodFill(c *canvasWidget, pixelX, pixelY int, fillColor [4]float64) {
	w := c.width
	h := c.height
	if w <= 0 || h <= 0 || pixelX < 0 || pixelY < 0 || pixelX >= w || pixelY >= h {
		return
	}

	// Render current canvas state to a temporary surface to read pixels.
	surface := cairolib.CreateImageSurface(cairolib.FormatARGB32, w, h)
	cr := cairolib.Create(surface)

	// Background.
	bg := c.state.BgColor
	cr.SetSourceRGBA(bg[0], bg[1], bg[2], bg[3])
	cr.Rectangle(0, 0, float64(w), float64(h))
	cr.Fill()

	// All committed strokes.
	for i := range c.state.Committed {
		RenderStroke(cr, &c.state.Committed[i])
	}
	cr.Close()
	surface.Flush()

	data := surface.Data()
	stride := surface.Stride()
	if len(data) == 0 {
		surface.Close()
		return
	}

	// Get target color at the click point.
	targetR, targetG, targetB, targetA := getPixel(data, stride, pixelX, pixelY)

	// Convert fill color to 0-255.
	fR := uint8(fillColor[0] * 255)
	fG := uint8(fillColor[1] * 255)
	fB := uint8(fillColor[2] * 255)
	fA := uint8(fillColor[3] * 255)

	// Don't fill if target color matches fill color.
	if colorMatch(targetR, targetG, targetB, targetA, fR, fG, fB, fA, 5) {
		surface.Close()
		return
	}

	// BFS flood fill.
	visited := make([]bool, w*h)
	queue := []point{{pixelX, pixelY}}
	visited[pixelY*w+pixelX] = true
	const tolerance = 10

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]

		// Set pixel to fill color (ARGB32 format: B, G, R, A in memory on little-endian).
		offset := p.y*stride + p.x*4
		if offset+3 < len(data) {
			data[offset+0] = fB
			data[offset+1] = fG
			data[offset+2] = fR
			data[offset+3] = fA
		}

		// Check 4 neighbors.
		for _, d := range [4]point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			nx, ny := p.x+d.x, p.y+d.y
			if nx < 0 || ny < 0 || nx >= w || ny >= h {
				continue
			}
			idx := ny*w + nx
			if visited[idx] {
				continue
			}
			nr, ng, nb, na := getPixel(data, stride, nx, ny)
			if colorMatch(nr, ng, nb, na, targetR, targetG, targetB, targetA, tolerance) {
				visited[idx] = true
				queue = append(queue, point{nx, ny})
			}
		}
	}

	surface.MarkDirty()

	// Update UI on the main thread.
	glib.IdleAdd(func() bool {
		// Store the filled surface as the new cache and commit a fill marker stroke.
		c.renderer.cache = surface
		c.renderer.cacheW = w
		c.renderer.cacheH = h
		c.renderer.cacheCount = len(c.state.Committed) + 1

		fillStroke := stroke.Stroke{
			Type:  stroke.TypeFill,
			Color: fillColor,
			FillX: pixelX,
			FillY: pixelY,
		}
		canvas.CommitStroke(c.state, fillStroke, float64(w))
		c.area.QueueDraw()
		return false
	})
}

type point struct {
	x, y int
}

func getPixel(data []byte, stride, x, y int) (r, g, b, a uint8) {
	offset := y*stride + x*4
	if offset+3 >= len(data) {
		return 0, 0, 0, 0
	}
	// ARGB32 on little-endian: B, G, R, A
	return data[offset+2], data[offset+1], data[offset+0], data[offset+3]
}

func colorMatch(r1, g1, b1, a1, r2, g2, b2, a2 uint8, tolerance int) bool {
	return abs(int(r1)-int(r2)) <= tolerance &&
		abs(int(g1)-int(g2)) <= tolerance &&
		abs(int(b1)-int(b2)) <= tolerance &&
		abs(int(a1)-int(a2)) <= tolerance
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
