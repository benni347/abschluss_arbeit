package main

import (
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"
	utils "github.com/benni347/messengerutils"
)

func main() {
	m := &utils.MessengerUtils{Verbose: true}
	a := app.New()
	a.Settings().SetTheme(&myTheme{})
	w := a.NewWindow("Alarm")
	input := newNumericalEntry()
	input.SetPlaceHolder("Enter a time (e.g. 1330)")
	previewLabel := widget.NewLabel("")
	content := container.NewVBox(
		input,
		widget.NewButton(
			"Retrive time",
			func() {
				timeValue := input.Text
				m.PrintInfo("The value entered in time is: ", timeValue)
				previewLabel.SetText(timeValue)
			},
		),
		previewLabel,
	)
	w.SetContent(content)
	w.ShowAndRun()
}

type numericalEntry struct {
	widget.Entry
}

func newNumericalEntry() *numericalEntry {
	entry := &numericalEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *numericalEntry) TypedRune(r rune) {
	if (r >= '0' && r <= '9') || r == '.' || r == ':' {
		e.Entry.TypedRune(r)
	}
}

func (e *numericalEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}
	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}

func (e *numericalEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}
