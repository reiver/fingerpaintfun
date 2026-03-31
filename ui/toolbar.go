package ui

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// toolbarWidget manages the bottom toolbar with mode switching
// between color palette, brush picker, and background picker.
type toolbarWidget struct {
	container *gtk.Box
	stack     *gtk.Stack

	palette     *paletteWidget
	brushPicker *brushPickerWidget
	bgPicker    *bgPickerWidget
	actions     *actionsWidget

	currentPage string // "palette", "brushes", "bg"
}

func newToolbarWidget(canvas *canvasWidget) *toolbarWidget {
	t := &toolbarWidget{currentPage: "palette"}

	t.container = gtk.NewBox(gtk.OrientationHorizontal, 4)
	t.container.AddCSSClass("toolbar")
	t.container.SetHAlign(gtk.AlignFill)
	t.container.SetHExpand(true)

	// Actions (undo/redo) on the left.
	t.actions = newActionsWidget(canvas)
	t.container.Append(t.actions.box)

	// Center: stack that toggles between palette, brush picker, and bg picker.
	t.stack = gtk.NewStack()
	t.stack.SetTransitionType(gtk.StackTransitionTypeCrossfade)
	t.stack.SetTransitionDuration(150)
	t.stack.SetHExpand(true)

	// Color palette.
	t.palette = newPaletteWidget(func(color [4]float64) {
		canvas.currentColor = color
	})
	t.stack.AddNamed(t.palette.grid, "palette")

	// Brush picker.
	t.brushPicker = newBrushPickerWidget(
		func(brushType int) { canvas.currentBrush = brushType },
		func(size float64) { canvas.currentSize = size },
	)
	t.stack.AddNamed(t.brushPicker.grid, "brushes")

	// Stamp picker.
	stampPicker := newStampPickerWidget(func(stampID int) {
		canvas.currentStamp = stampID
	})
	t.stack.AddNamed(stampPicker.grid, "stamps")

	// Template picker.
	templatePicker := newTemplatePickerWidget(func(templateID int) {
		canvas.currentTemplate = templateID
		canvas.renderer.templateID = templateID
		canvas.area.QueueDraw()
	})
	t.stack.AddNamed(templatePicker.grid, "templates")

	// Background color picker.
	t.bgPicker = newBgPickerWidget(func(color [4]float64) {
		canvas.state.BgColor = color
		canvas.renderer.InvalidateCache()
		canvas.area.QueueDraw()
	})
	t.stack.AddNamed(t.bgPicker.grid, "bg")

	t.stack.SetVisibleChildName("palette")
	t.container.Append(t.stack)

	// Right side buttons.
	btnBox := gtk.NewBox(gtk.OrientationVertical, 2)

	// Toggle button: cycles palette → brushes → bg → palette.
	toggleBtn := gtk.NewButton()
	toggleBtn.SetIconName("applications-graphics-symbolic")
	toggleBtn.SetSizeRequest(48, 48)
	toggleBtn.AddCSSClass("flat")
	toggleBtn.ConnectClicked(func() {
		switch t.currentPage {
		case "palette":
			t.currentPage = "brushes"
		case "brushes":
			t.currentPage = "stamps"
		case "stamps":
			t.currentPage = "templates"
		case "templates":
			t.currentPage = "bg"
		default:
			t.currentPage = "palette"
		}
		t.stack.SetVisibleChildName(t.currentPage)
	})
	btnBox.Append(toggleBtn)

	// Mirror mode toggle.
	mirrorBtn := gtk.NewToggleButton()
	mirrorBtn.SetIconName("object-flip-horizontal-symbolic")
	mirrorBtn.SetSizeRequest(48, 48)
	mirrorBtn.AddCSSClass("flat")
	mirrorBtn.ConnectToggled(func() {
		canvas.state.MirrorMode = mirrorBtn.Active()
		canvas.area.QueueDraw()
	})
	btnBox.Append(mirrorBtn)

	// Clear button.
	btnBox.Append(t.actions.clearBtn)

	t.container.Append(btnBox)

	return t
}
