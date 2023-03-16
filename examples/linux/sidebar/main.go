package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)

	// Source: https://github.com/gotk3/gotk3
	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Sidebar")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	// End Quote
	// Create a new box
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal("Unable to create box:", err)
	}

	// Create a new horizontal box for the label and icon
	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		log.Fatal("Unable to create horizontal box:", err)
	}

	// Add the horizontal box to the vertical box
	box.Add(hbox)

	// Create a new label
	label, err := gtk.LabelNew("Sidebar's are nice")
	if err != nil {
		log.Fatal("Error during creation of the hello world label. The error is: ", err)
	}

	// Create a new icon
	icon_button, err := gtk.ButtonNewFromIconName("format-justify-fill", gtk.ICON_SIZE_BUTTON)
	if err != nil {
		log.Fatal("Error during the creation of the hamburger menu button: ", err)
	}

	// Add the icon button to the horizontal box
	hbox.PackStart(icon_button, false, false, 0)

	// Create a new separator
	separator, err := gtk.SeparatorNew(gtk.ORIENTATION_VERTICAL)
	if err != nil {
		log.Fatal("Unable to create separator:", err)
	}

	// Add the separator to the horizontal box
	hbox.PackStart(separator, false, false, 10)

	// Add the label to the horizontal box
	hbox.PackStart(label, true, true, 0)

	// Add the box widget to the window
	win.Add(box)

	// Set a default window size
	win.SetDefaultSize(500, 400)

	// Make all elements appear on the window
	win.ShowAll()

	// Begin executing GTKS's main loop.
	gtk.Main()
}
