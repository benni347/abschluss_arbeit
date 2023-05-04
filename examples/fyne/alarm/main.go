package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(&myTheme{})
	w := a.NewWindow("Alarm")
	w.SetContent(widget.NewLabel("Hello Fyne!"))
	w.ShowAndRun()
}
