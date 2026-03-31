package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"fingerpaintfun/lib/canvas"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	cairolib "github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// paintingsDir returns the directory where paintings are saved.
func paintingsDir() string {
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		home, _ := os.UserHomeDir()
		dataDir = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dataDir, "fingerpaintfun", "paintings")
}

// SaveCanvas renders the current canvas to a PNG file.
func SaveCanvas(c *canvasWidget) error {
	dir := paintingsDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	w := c.width
	h := c.height
	if w <= 0 || h <= 0 {
		w, h = 360, 540
	}

	// Render to off-screen surface.
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

	filename := fmt.Sprintf("%s.png", time.Now().Format("20060102_150405"))
	path := filepath.Join(dir, filename)
	return surface.WriteToPNG(path)
}

// galleryPage shows saved paintings as a thumbnail grid.
type galleryPage struct {
	page   *adw.NavigationPage
	grid   *gtk.FlowBox
	canvas *canvasWidget
	win    *adw.ApplicationWindow
}

func newGalleryPage(c *canvasWidget, win *adw.ApplicationWindow, nav *adw.NavigationView) *galleryPage {
	g := &galleryPage{canvas: c, win: win}

	content := gtk.NewBox(gtk.OrientationVertical, 0)

	// Header with back and new-painting buttons.
	header := adw.NewHeaderBar()
	newBtn := gtk.NewButton()
	newBtn.SetIconName("document-new-symbolic")
	newBtn.ConnectClicked(func() {
		canvas.ClearAll(c.state)
		c.renderer.InvalidateCache()
		c.area.QueueDraw()
		nav.Pop()
	})
	header.PackEnd(newBtn)
	content.Append(header)

	// Scrollable thumbnail grid.
	scrolled := gtk.NewScrolledWindow()
	scrolled.SetPolicy(gtk.PolicyNever, gtk.PolicyAutomatic)
	scrolled.SetVExpand(true)

	g.grid = gtk.NewFlowBox()
	g.grid.SetMaxChildrenPerLine(3)
	g.grid.SetMinChildrenPerLine(3)
	g.grid.SetSelectionMode(gtk.SelectionNone)
	g.grid.SetColumnSpacing(4)
	g.grid.SetRowSpacing(4)
	g.grid.SetHomogeneous(true)

	g.loadThumbnails(nav)

	scrolled.SetChild(g.grid)
	content.Append(scrolled)

	g.page = adw.NewNavigationPage(content, "Gallery")

	return g
}

func (receiver *galleryPage) loadThumbnails(nav *adw.NavigationView) {
	dir := paintingsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	// Sort newest first.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() > entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".png") {
			continue
		}

		path := filepath.Join(dir, entry.Name())

		pic := gtk.NewPicture()
		pic.SetFilename(path)
		pic.SetContentFit(gtk.ContentFitCover)
		pic.SetSizeRequest(100, 100)

		// Long-press to delete (behind parental gate).
		longPress := gtk.NewGestureLongPress()
		filePath := path
		longPress.ConnectPressed(func(x, y float64) {
			ShowParentalGate(receiver.win, func() {
				os.Remove(filePath)
				// Refresh the grid by removing all children and reloading.
				for {
					child := receiver.grid.ChildAtIndex(0)
					if child == nil {
						break
					}
					receiver.grid.Remove(child)
				}
				receiver.loadThumbnails(nav)
			})
		})
		pic.AddController(longPress)

		receiver.grid.Append(pic)
	}
}
