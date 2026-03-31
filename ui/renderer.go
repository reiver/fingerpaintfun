package ui

import (
	"fingerpaintfun/lib/canvas"
	"fingerpaintfun/lib/stroke"

	cairolib "github.com/diamondburned/gotk4/pkg/cairo"
)

// renderer manages the cached surface for committed strokes.
type renderer struct {
	cache      *cairolib.Surface
	cacheW     int
	cacheH     int
	cacheCount int // number of committed strokes rendered to cache
}

// Render draws the full canvas: background, cached committed strokes, and active strokes.
func (receiver *renderer) Render(cr *cairolib.Context, width, height int, state *canvas.CanvasState, active *canvas.ActiveStrokes) {
	// Background.
	bg := state.BgColor
	cr.SetSourceRGBA(bg[0], bg[1], bg[2], bg[3])
	cr.Rectangle(0, 0, float64(width), float64(height))
	cr.Fill()

	// Ensure cache surface exists and matches current dimensions.
	receiver.ensureCache(width, height, state)

	// Composite the cached committed strokes.
	if receiver.cache != nil {
		cr.SetSourceSurface(receiver.cache, 0, 0)
		cr.Paint()
	}

	// Draw active (in-progress) strokes on top.
	for _, s := range active.All() {
		RenderStroke(cr, s)
	}
}

// InvalidateCache forces a full rebuild of the cache on the next Render call.
// Call this after undo, redo, clear, or background color change.
func (receiver *renderer) InvalidateCache() {
	receiver.cacheCount = 0
	if receiver.cache != nil {
		receiver.cache.Close()
		receiver.cache = nil
	}
}

// ensureCache creates or updates the off-screen surface for committed strokes.
func (receiver *renderer) ensureCache(width, height int, state *canvas.CanvasState) {
	// Recreate if dimensions changed.
	if receiver.cache != nil && (receiver.cacheW != width || receiver.cacheH != height) {
		receiver.InvalidateCache()
	}

	committed := state.Committed
	numCommitted := len(committed)

	// Nothing to cache.
	if numCommitted == 0 {
		if receiver.cache != nil {
			receiver.InvalidateCache()
		}
		return
	}

	// Create fresh cache if needed.
	if receiver.cache == nil {
		receiver.cache = cairolib.CreateImageSurface(cairolib.FormatARGB32, width, height)
		receiver.cacheW = width
		receiver.cacheH = height
		receiver.cacheCount = 0
	}

	// If undo happened (committed count shrank), rebuild from scratch.
	if numCommitted < receiver.cacheCount {
		receiver.rebuildCache(committed)
		return
	}

	// Incrementally render only new strokes.
	if numCommitted > receiver.cacheCount {
		cr := cairolib.Create(receiver.cache)
		for i := receiver.cacheCount; i < numCommitted; i++ {
			RenderStroke(cr, &committed[i])
		}
		cr.Close()
		receiver.cacheCount = numCommitted
	}
}

// rebuildCache clears and redraws all committed strokes onto the cache.
func (receiver *renderer) rebuildCache(committed []stroke.Stroke) {
	cr := cairolib.Create(receiver.cache)
	// Clear the surface.
	cr.SetOperator(cairolib.OperatorClear)
	cr.Paint()
	cr.SetOperator(cairolib.OperatorOver)

	for i := range committed {
		RenderStroke(cr, &committed[i])
	}
	cr.Close()
	receiver.cacheCount = len(committed)
}
