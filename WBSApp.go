package main

import (
	"WBS/internal/v4l2"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/korandiz/v4l"
	"gocv.io/x/gocv"
)

type WBSApp struct {
	viewport           *gtk.Image
	logArea            *gtk.Label
	currentPixbuf      *gdk.Pixbuf
	webcam             *V4l2Device
	captureImageButton *gtk.Button
	cameraShouldBeOpen bool
}

func (app *WBSApp) makeViewport() {
	var err error
	app.viewport, err = gtk.ImageNew()
	app.panicIfErr(err, "Unable to create drawing area:")
	app.viewport.SetSizeRequest(1024, 768)
}

func (app *WBSApp) makeMainWindow(title string) *gtk.Window {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	app.panicIfErr(err, "Unable to create window:")

	win.SetTitle(title)
	// win.SetDefaultSize(1024, 768)
	win.Maximize()
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	return win
}

func (app *WBSApp) panicIfErr(err error, message string) {
	if err != nil {
		log.Panic(message, err)
	}
}

func (app *WBSApp) makeDevicesComboBox() (*gtk.ComboBox, error) {
	store, err := gtk.ListStoreNew(glib.TYPE_STRING)
	app.panicIfErr(err, "Unable to create list store:")

	defaultStoreItemValue := "Choose a device"
	defaultStoreItem := store.Append()
	store.SetValue(defaultStoreItem, 0, defaultStoreItemValue)

	items, err := v4l2.ListDevices()
	if err != nil {
		errorString := err.Error()
		if !strings.Contains(errorString, "ls: cannot access '/dev/v4l/by-id/usb*'") {
			return nil, errors.New("could not find any useable device under '/dev/v4l/by-id/usb*'")
		}
	}

	for _, item := range items {
		iter := store.Append()
		store.SetValue(iter, 0, item)
	}

	combo, err := gtk.ComboBoxNewWithModel(store)
	app.panicIfErr(err, "Unable to create combo box: ")

	renderer, err := gtk.CellRendererTextNew()
	app.panicIfErr(err, "Unable to create cell renderer: ")
	combo.PackStart(renderer, true)

	combo.AddAttribute(renderer, "text", 0)
	combo.SetActive(0) // Set default selection

	combo.Connect("changed", func(thisCombo *gtk.ComboBox) {
		iter, err := thisCombo.GetActiveIter()
		app.panicIfErr(err, "Error getting active iter:")

		value, err := store.GetValue(iter, 0)
		app.panicIfErr(err, "Error getting value: ")

		str, err := value.GetString()
		app.panicIfErr(err, "Error converting to string")
		if str == defaultStoreItemValue {
			app.cameraShouldBeOpen = false
			app.closeCamera()
		} else {
			// log.Println("Selected:", str)
			app.logArea.SetText("Selected device " + str)
			app.cameraShouldBeOpen = true
			// go app.webcam.CaptureFrames(str, nil)

		}
	})

	return combo, nil
}

func (app *WBSApp) OpenWebcamDevice(devicePath string, config v4l.DeviceConfig) {
	// var err error // why do I need this ?
	app.webcam.CaptureFrames(devicePath, config)
	// app.panicIfErr(err, "Could not open video device 0")
	app.webcam.isOpen = true
	app.captureImageButton.SetSensitive(true)

}

func (app *WBSApp) closeCamera() {
	app.webcam.Close()
	app.captureImageButton.SetSensitive(false)
}

// func (app *WBSApp) captureFrameFromDevice(frame *gocv.Mat) {
// 	if app.webcam == nil || !app.webcam.isOpen {
// 		fmt.Println("webcam is nil or not open")
// 		// time.Sleep(1 * time.Second)
// 		return
// 	}

// 	app.webcam.mtx.Lock()
// 	ok := app.webcam.captureDevice.Read(frame)
// 	app.webcam.mtx.Unlock()
// 	if !ok {
// 		log.Println("Cannot read from webcam, waiting for a second")
// 		time.Sleep(1000 * time.Millisecond)
// 		return
// 	}

// 	if frame.Empty() {
// 		log.Println("Frame is empty")
// 		return
// 	}
// 	// something happens in the next line related to the memory leak
// 	pixbuf, err := app.MatToPixbuf(frame)
// 	if err != nil {
// 		log.Println("Failed to convert video frame to pixbuf: ", err)
// 		pixbuf = nil
// 		return
// 	}
// 	app.currentPixbuf = pixbuf
// 	pixbuf = nil
// }

func (app *WBSApp) makeMenuBar() *gtk.MenuBar {
	// Create the menubar
	menubar, err := gtk.MenuBarNew()
	app.panicIfErr(err, "Unable to create menubar:")

	// Create menu items
	menu, err := gtk.MenuNew()
	app.panicIfErr(err, "Unable to create menu:")

	itemOne, err := gtk.MenuItemNewWithLabel("One")
	app.panicIfErr(err, "Unable to create menu item:")
	itemOne.Connect("activate", func() {
		log.Println("Menu item 'One' clicked")
	})
	menu.Append(itemOne)

	itemTwo, err := gtk.MenuItemNewWithLabel("Two")
	app.panicIfErr(err, "Unable to create menu item: ")

	itemTwo.Connect("activate", func() {
		log.Println("Menu item 'Two' clicked")
	})
	menu.Append(itemTwo)

	// Create a File menu to hold the items
	fileMenu, err := gtk.MenuItemNewWithLabel("File")
	app.panicIfErr(err, "Unable to creater file menu:")

	fileMenu.SetSubmenu(menu)
	menubar.Append(fileMenu)

	return menubar
}

func (app *WBSApp) RefreshDevicesList() {
	fmt.Println("RefreshDevicesList .............")
}

func (app *WBSApp) CaptureImage() {
	// get the next file ID in project folder using the project
	fmt.Println("SaveImage ......................")
}

func (app *WBSApp) MatToPixbuf(mat *gocv.Mat) (*gdk.Pixbuf, error) {

	// Convert Mat to bytes (assuming RGB image)
	data, err := gocv.IMEncode(".png", *mat)
	if err != nil {
		return nil, fmt.Errorf("failed to encode Mat to PNG: %v", err)
	}
	defer data.Close()

	// Create a PixbufLoader to load the image data
	loader, err := gdk.PixbufLoaderNew()
	if err != nil {
		return nil, fmt.Errorf("failed to create PixbufLoader: %v", err)
	}

	// Write the image data to the loader
	pixbuf, err := loader.WriteAndReturnPixbuf(data.GetBytes())
	if err != nil {
		return nil, fmt.Errorf("failed to write to PixbufLoader: %v", err)
	}

	// // Get the Pixbuf from the loader
	// pixbuf, err := loader.GetPixbuf()
	err = loader.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close loader: %v", err)
	}

	return pixbuf, nil
}

/*
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

				time.Sleep(33 * time.Millisecond)
			} else {
				fmt.Println("camera not open, sleeping")
				time.Sleep(1 * time.Second)
			}
		}
	}(&app)
*/
