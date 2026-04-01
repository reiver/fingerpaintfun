# Finger Paint Fun — User Guide

**Finger Paint Fun** is a finger painting app for kids ages 3 to 5.
Your child can draw, stamp, and color with their fingers — no instructions needed.
This guide is for parents.

## Quick Start

1. Open **Finger Paint Fun**
2. Hand the device to your child
3. They draw by touching the screen (or clicking and dragging with a mouse)
4. That's it — everything is designed to be discovered by poking

## What Your Child Can Do

Your child can do all of these without help:

* **Draw** with fingers or a mouse — just touch the canvas and move
* **Pick colors** — tap any color circle at the bottom of the screen
* **Switch brushes** — tap the paint palette button on the right to cycle through: colors, brushes, stamps, templates, and background colors
* **Undo a mistake** — tap the undo arrow (left arrow button)
* **Redo** — tap the redo arrow (right arrow button)
* **Clear the canvas** — **long-press** (press and hold) the trash button. A single tap won't work — this prevents accidental clearing

## Tools and Brushes

The app includes several brush types. Your child switches between them by tapping the palette toggle button on the right side of the toolbar, then tapping a brush:

| Brush          | What It Does                                    |
|----------------|-------------------------------------------------|
| **Round**      | Smooth, solid lines                             |
| **Crayon**     | Rough, textured strokes                         |
| **Marker**     | Semi-transparent — overlapping areas get darker |
| **Pencil**     | Thin, precise lines                             |
| **Neon**       | Glowing lines with a soft halo                  |
| **Rainbow**    | Color shifts automatically along the stroke     |
| **Spray**      | Scattered dots like spray paint                 |
| **Chalk**      | Dusty, textured strokes with gaps               |
| **Watercolor** | Very soft and translucent, like wet paint       |
| **Eraser**     | Removes paint (draws with the background color) |

Three brush sizes are available: small, medium, and large.

## Stamps

Tap the palette toggle button until you see the stamp picker.
Tap a stamp shape (circle, square, triangle, star, or heart), then tap anywhere on the canvas to place it.

## Coloring Pages

Tap the palette toggle button until you see the template picker.
Choose a template (house, fish, star, or circle) and it will appear as an outline on the canvas.
Your child paints freely over the template — there are no enforced boundaries.
Select "Blank" to return to a plain canvas.

## Background Colors

Tap the palette toggle button until you see the background color picker (square swatches instead of circles).
Tap a color to change the canvas background.
Try black for a neon-on-dark effect.

## Mirror Mode

Tap the mirror button (on the right side of the toolbar) to turn on symmetry mode.
Everything drawn on one side of the canvas is automatically mirrored on the other side.
A faint dashed line shows the mirror axis.
Tap the button again to turn it off.

## Saving and Gallery

Paintings are saved automatically when your child opens the gallery or when the app closes.

* Tap the **grid button** in the toolbar to open the gallery
* Saved paintings appear as thumbnails
* Tap a thumbnail to view it fullscreen
* Tap the **+** button to start a new blank painting

## Fill Tool

When the fill (paint bucket) brush is selected, tapping on the canvas floods that area with the current color.
Good for quickly coloring large regions.

## Mascot

A small yellow character appears in the top-right corner of the canvas.
It reacts to what your child does — wiggling when they draw, looking surprised when they switch colors, and celebrating their creativity.
It can be turned off (see Parental Controls below).

## For Parents

### What Your Child Can Do vs. What Needs a Parent

| Action                                            | Who                    |
|---------------------------------------------------|------------------------|
| Drawing, colors, brushes, undo, stamps, templates | Child                  |
| Deleting saved paintings                          | Parent (parental gate) |
| Muting sound effects                              | Parent (parental gate) |
| Toggling the mascot                               | Parent (parental gate) |

### Parental Gate

Some actions are protected by a simple math question (e.g., "What is 5 + 7?") that a 3-to-5-year-old cannot answer.
This prevents accidental deletions and setting changes.

### Privacy and Safety

* **No internet required** — the app works entirely offline
* **No data collection** — no analytics, no tracking, no accounts
* **No ads** — ever
* **No external links** — your child cannot navigate away to a website
* **No microphone or camera** — not used
* **Paintings are stored locally** on the device only

### Where Are Paintings Saved?

Paintings are saved as PNG files in:
```
~/.local/share/fingerpaintfun/paintings/
```

### Closing the App

Use your system's standard window close method:
* **Desktop**: click the window close button or press `Alt+F4`
* **Phosh/mobile**: swipe the app away from the task switcher

The current painting is automatically saved when the app closes.

## Troubleshooting

**Nothing happens when I tap the clear button.**
The clear button requires a **long-press** (press and hold for about half a second). This is intentional — it prevents your child from accidentally erasing their work.

**I can't hear any sounds.**
Check that your device volume is up.
If sound was muted via the parental gate, unmute it by tapping the speaker icon and answering the math question.

**The app won't build / won't start.**
See [README.md](README.md) for build requirements and instructions.
