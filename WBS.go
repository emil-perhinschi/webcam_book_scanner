package main

import (
	"errors"
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/korandiz/v4l"
	"github.com/korandiz/v4l/fmt/yuyv"
)

func main() {
	gtk.Init(nil)

	var app WBSApp
	app.webcam = NewW4l2Device(&app)

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
	fillConfigList()
	// counter := 0
	// glib.IdleAdd(func() bool {
	// 	fmt.Println(counter)
	// 	counter++
	// 	if counter > 10 {
	// 		return false
	// 	}
	// 	time.Sleep(1000 * time.Millisecond) // this runs in the main thread and it blocks the gui
	// 	return true
	// })
	// Start the GTK main loop
	gtk.Main()

}

func fillConfigList() error {
	dev, err := v4l.Open(`/dev/v4l/by-id/usb-046d_Logitech_BRIO_50316219-video-index0`)
	if err != nil {
		return errors.New("Open: " + err.Error())
	}
	defer dev.Close()
	cfgs, err := dev.ListConfigs()
	if err != nil {
		return errors.New("ListConfigs: " + err.Error())
	}
	fmt.Println(cfgs)
	for i := range cfgs {
		if cfgs[i].Format != yuyv.FourCC {
			continue
		}
		fmt.Println(cfg2str(cfgs[i]), cfgs[i])
	}

	return nil
}

func cfg2str(cfg v4l.DeviceConfig) string {
	w := cfg.Width
	h := cfg.Height
	f := cfg.FPS
	return fmt.Sprintf("%dx%d @ %.4g FPS", w, h, float64(f.N)/float64(f.D))
}
