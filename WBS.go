package main

import (
	"fmt"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"gocv.io/x/gocv"
)

func main() {
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
	refreshDevicesButton.Connect("clicked", func() {
		app.RefreshDevicesList()
	})

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
	app.makeViewport()
	mainContainer.PackStart(app.viewport, true, true, 0)

	app.logArea, err = gtk.LabelNew("Will add messages here")
	app.panicIfErr(err, "failed to initialize log area")
	mainContainer.PackStart(app.logArea, false, false, 0)

	// Show all widgets
	win.ShowAll()

	go func(myapp *WBSApp) {
		frame := gocv.NewMat()
		defer frame.Close()
		// counter := 0
		for {
			if myapp.cameraShouldBeOpen {
				myapp.OpenWebcamDevice()
				// if counter == 0 {
				// running this only once stops the leak
				myapp.captureFrameFromDevice(&frame)
				// counter++
				// }
				myapp.viewport.SetFromPixbuf(app.currentPixbuf)
				app.currentPixbuf = nil

				// this does not seem to indicate Mat-s are leaking
				// go run -tags matprofile .
				// var b bytes.Buffer
				// gocv.MatProfile.WriteTo(&b, 1)
				// fmt.Print(b.String())
				// time.Sleep(33 * time.Millisecond)
			} else {
				fmt.Println("camera not open, sleeping")
				time.Sleep(1 * time.Second)
			}
		}
	}(&app)

	// Start the GTK main loop
	gtk.Main()

}
