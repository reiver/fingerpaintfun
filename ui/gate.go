package ui

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// ShowParentalGate displays a math problem that a 3-5 year old cannot solve
// but an adult can. Calls onSuccess only if the adult answers correctly.
func ShowParentalGate(parent *adw.ApplicationWindow, onSuccess func()) {
	a := rand.Intn(8) + 2 // 2-9
	b := rand.Intn(8) + 2 // 2-9
	answer := a + b

	dialog := adw.NewMessageDialog(&parent.Window, "Parental Gate", fmt.Sprintf("What is %d + %d?", a, b))
	dialog.AddResponse("cancel", "Cancel")
	dialog.AddResponse("ok", "OK")
	dialog.SetDefaultResponse("ok")
	dialog.SetCloseResponse("cancel")

	entry := gtk.NewEntry()
	entry.SetInputPurpose(gtk.InputPurposeDigits)
	entry.SetHAlign(gtk.AlignCenter)
	entry.SetSizeRequest(120, -1)
	dialog.SetExtraChild(entry)

	dialog.ConnectResponse(func(response string) {
		if response == "ok" {
			text := entry.Text()
			if val, err := strconv.Atoi(text); err == nil && val == answer {
				onSuccess()
			}
		}
	})

	dialog.Present()
}
