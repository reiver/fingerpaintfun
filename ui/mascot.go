package ui

import (
	"math"

	cairolib "github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
)

// Mascot states.
const (
	mascotIdle = iota
	mascotDrawing
	mascotSurprised
	mascotCelebrating
	mascotWhoops
	mascotWave
)

// mascot manages the encouraging character overlay.
type mascot struct {
	state   int
	frame   int
	enabled bool
	timerID glib.SourceHandle // glib timeout source ID
}

func newMascot() *mascot {
	m := &mascot{
		state:   mascotIdle,
		enabled: true,
	}
	// Start idle animation loop.
	m.timerID = glib.TimeoutAdd(500, func() bool {
		m.frame = (m.frame + 1) % 3
		return true // keep repeating
	})
	return m
}

// SetState transitions the mascot to a new state for a short duration,
// then returns to idle.
func (receiver *mascot) SetState(state int) {
	receiver.state = state
	receiver.frame = 0
	// Return to idle after 1 second.
	glib.TimeoutAdd(1000, func() bool {
		receiver.state = mascotIdle
		return false
	})
}

// Render draws the mascot in the top-right corner of the canvas.
func (receiver *mascot) Render(cr *cairolib.Context, canvasWidth, canvasHeight int) {
	if !receiver.enabled {
		return
	}

	size := 48.0
	margin := 8.0
	cx := float64(canvasWidth) - size/2 - margin
	cy := size/2 + margin

	// Draw body (circle).
	cr.Save()

	// Subtle bounce for idle animation.
	bounce := 0.0
	if receiver.state == mascotIdle {
		bounce = math.Sin(float64(receiver.frame)*2*math.Pi/3) * 2
	}
	cy += bounce

	// Body.
	cr.SetSourceRGBA(1.0, 0.85, 0.2, 1.0) // yellow
	cr.Arc(cx, cy, size/2-2, 0, 2*math.Pi)
	cr.Fill()

	// Outline.
	cr.SetSourceRGBA(0.8, 0.6, 0.0, 1.0)
	cr.SetLineWidth(2)
	cr.Arc(cx, cy, size/2-2, 0, 2*math.Pi)
	cr.Stroke()

	// Eyes.
	eyeY := cy - 4
	leftEyeX := cx - 8
	rightEyeX := cx + 8
	eyeSize := 4.0

	switch receiver.state {
	case mascotSurprised:
		// Wide eyes.
		cr.SetSourceRGBA(0, 0, 0, 1)
		cr.Arc(leftEyeX, eyeY, eyeSize+1, 0, 2*math.Pi)
		cr.Fill()
		cr.Arc(rightEyeX, eyeY, eyeSize+1, 0, 2*math.Pi)
		cr.Fill()
	case mascotWhoops:
		// Squinting eyes.
		cr.SetSourceRGBA(0, 0, 0, 1)
		cr.SetLineWidth(2)
		cr.MoveTo(leftEyeX-4, eyeY)
		cr.LineTo(leftEyeX+4, eyeY)
		cr.Stroke()
		cr.MoveTo(rightEyeX-4, eyeY)
		cr.LineTo(rightEyeX+4, eyeY)
		cr.Stroke()
	case mascotCelebrating:
		// Star eyes.
		cr.SetSourceRGBA(0, 0, 0, 1)
		drawMiniStar(cr, leftEyeX, eyeY, eyeSize)
		drawMiniStar(cr, rightEyeX, eyeY, eyeSize)
	default:
		// Normal dot eyes.
		cr.SetSourceRGBA(0, 0, 0, 1)
		cr.Arc(leftEyeX, eyeY, eyeSize, 0, 2*math.Pi)
		cr.Fill()
		cr.Arc(rightEyeX, eyeY, eyeSize, 0, 2*math.Pi)
		cr.Fill()
	}

	// Mouth.
	mouthY := cy + 6
	cr.SetSourceRGBA(0, 0, 0, 1)
	cr.SetLineWidth(2)
	switch receiver.state {
	case mascotSurprised:
		// O mouth.
		cr.Arc(cx, mouthY+2, 5, 0, 2*math.Pi)
		cr.Stroke()
	case mascotWhoops:
		// Wavy mouth.
		cr.MoveTo(cx-8, mouthY)
		cr.CurveTo(cx-4, mouthY+4, cx+4, mouthY-4, cx+8, mouthY)
		cr.Stroke()
	case mascotCelebrating:
		// Big smile.
		cr.Arc(cx, mouthY-2, 10, 0.2, math.Pi-0.2)
		cr.Stroke()
	case mascotWave:
		// Small smile + raised arm.
		cr.Arc(cx, mouthY-2, 7, 0.2, math.Pi-0.2)
		cr.Stroke()
		// Arm waving.
		cr.MoveTo(cx+size/2-4, cy)
		cr.LineTo(cx+size/2+8, cy-12)
		cr.Stroke()
	default:
		// Gentle smile.
		cr.Arc(cx, mouthY-2, 7, 0.2, math.Pi-0.2)
		cr.Stroke()
	}

	cr.Restore()
}

func drawMiniStar(cr *cairolib.Context, cx, cy, size float64) {
	for i := 0; i < 10; i++ {
		angle := float64(i)*math.Pi/5 - math.Pi/2
		r := size
		if i%2 == 1 {
			r = size * 0.4
		}
		x := cx + r*math.Cos(angle)
		y := cy + r*math.Sin(angle)
		if i == 0 {
			cr.MoveTo(x, y)
		} else {
			cr.LineTo(x, y)
		}
	}
	cr.ClosePath()
	cr.Fill()
}
