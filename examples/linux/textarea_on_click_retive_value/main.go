package main

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"log"
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
	win.SetTitle("Simple Example")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	// End Quote
	// Create a new box
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal("Unable to create box:", err)
	}

	// Create a new label
	label, err := gtk.LabelNew("Hello, World!")
	if err != nil {
		log.Fatal("Error during creation of the hello world label. The error is: ", err)
	}

	// Add the label to the window
	box.Add(label)

	// Create a new text input
	textArea, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Error during creation of the hello world label. The error is: ", err)
	}

	// Add the textArea to the window
	box.Add(textArea)

	// Add a button to print the content
	button, err := gtk.ButtonNewWithLabel("Print Text")
	if err != nil {
		log.Fatal("Error during creation of the button. The error is: ", err)
	}

	button.Connect("clicked", func() {
		// Retirive the content of the textArea
		textAreaContent, err := textArea.GetText()
		if err != nil {
			log.Fatal("Error during Retiriving the text from the textarea:", err)
		}
		// Retrieve the content of the text input and print it to the console
		fmt.Println("Text input content:", textAreaContent)
	})

	// Add the button to the window
	box.Add(button)

	// Add the box widget to the window
	win.Add(box)

	// Set a default window size
	win.SetDefaultSize(500, 400)

	// Make all elements appear on the window
	win.ShowAll()

	// Begin executing GTKS's main loop.
	gtk.Main()
}
