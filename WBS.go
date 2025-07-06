package main

import (
	"runtime"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	runtime.GOMAXPROCS(1)
	gtk.Init(nil)

	var app WBSApp
	app.webcam = &WebcamDevice{}

	win := app.makeMainWindow("Webcam book scanner")

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	app.panicIfErr(err, "Unable to create main container box:")
	win.Add(vbox)

	menubar := app.makeMenuBar()
	vbox.PackStart(menubar, false, false, 0)

	combo, err := app.makeDevicesComboBox()
	app.panicIfErr(err, "failed to make devices combo: ")
	vbox.PackStart(combo, false, false, 5)

	// Create a drawing area
	app.viewport, err = gtk.DrawingAreaNew()
	app.panicIfErr(err, "Unable to create drawing area:")
	app.viewport.SetSizeRequest(400, 200)
	vbox.PackStart(app.viewport, true, true, 0)

	// Connect the draw signal to draw on the drawing area
	app.viewport.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		// Simple drawing: fill with light gray and draw a red rectangle
		cr.SetSourceRGB(0.8, 0.8, 0.8)
		cr.Paint()
		cr.SetSourceRGB(1, 0, 0)
		cr.Rectangle(50, 50, 100, 100)
		cr.Fill()
	})

	app.logArea, err = gtk.LabelNew("Will add messages here")
	app.panicIfErr(err, "failed to initialize log area")
	vbox.PackStart(app.logArea, true, true, 0)

	// Show all widgets
	win.ShowAll()

	// Start the GTK main loop
	gtk.Main()
}
