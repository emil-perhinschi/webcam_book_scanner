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

	mainContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	app.panicIfErr(err, "Unable to create main container:")
	win.Add(mainContainer)

	menubar := app.makeMenuBar()
	mainContainer.PackStart(menubar, false, false, 0)

	devicesContainer, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	app.panicIfErr(err, "Unable to create the devices container:")

	refreshDevicesButton, err := gtk.ButtonNewWithLabel("Refresh devices")
	app.panicIfErr(err, "Unable to create the refresh button:")
	refreshDevicesButton.Connect("clicked", app.RefreshDevicesList)

	app.captureImageButton, err = gtk.ButtonNewWithLabel("Capture image")
	app.panicIfErr(err, "Unable to create the refresh button:")
	app.captureImageButton.SetSensitive(false)
	app.captureImageButton.Connect("clicked", app.CaptureImage)

	combo, err := app.makeDevicesComboBox()
	app.panicIfErr(err, "failed to make devices combo: ")
	devicesContainer.PackStart(combo, true, true, 1)
	devicesContainer.PackEnd(refreshDevicesButton, false, true, 0)
	devicesContainer.PackEnd(app.captureImageButton, false, true, 0)

	mainContainer.PackStart(devicesContainer, false, true, 0)

	// Create a drawing area
	app.viewport, err = gtk.DrawingAreaNew()
	app.panicIfErr(err, "Unable to create drawing area:")
	app.viewport.SetSizeRequest(400, 200)
	mainContainer.PackStart(app.viewport, true, true, 0)

	// Connect the draw signal to draw on the drawing area
	app.viewport.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		// Simple drawing: fill with light gray
		cr.SetSourceRGB(0.8, 0.8, 0.8)
		cr.Paint()
		cr.Fill()
	})

	app.logArea, err = gtk.LabelNew("Will add messages here")
	app.panicIfErr(err, "failed to initialize log area")
	mainContainer.PackStart(app.logArea, false, false, 0)

	// Show all widgets
	win.ShowAll()

	// Start the GTK main loop
	gtk.Main()
}
