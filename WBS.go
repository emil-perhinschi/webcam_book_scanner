package main

import (
	"WBS/internal/device"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"gocv.io/x/gocv"
)

type WBSApp struct {
	viewport      *gtk.DrawingArea
	logArea       *gtk.Label
	currentPixbuf *gdk.Pixbuf
	webcam        *gocv.VideoCapture
}

func main() {
	gtk.Init(nil)

	var app WBSApp

	win := app.makeMainWindow("Webcam book scanner")

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	app.panicIfErr(err, "Unable to create main container box:")
	win.Add(vbox)

	menubar := app.makeMenuBar()
	vbox.PackStart(menubar, false, false, 0)

	combo, err := app.makeDevicesComboBox()
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

	app.connectoToCamera(app.viewport)

	app.logArea, err = gtk.LabelNew("Will add messages here")
	app.panicIfErr(err, "failed to initialize log area")
	vbox.PackStart(app.logArea, true, true, 0)

	// Show all widgets
	win.ShowAll()

	// Start the GTK main loop
	gtk.Main()
}

func (app *WBSApp) connectoToCamera(viewport *gtk.DrawingArea) {

}

func (app *WBSApp) makeMainWindow(title string) *gtk.Window {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	app.panicIfErr(err, "Unable to create window:")

	win.SetTitle(title)
	win.SetDefaultSize(400, 300)
	win.Connect("destroy", gtk.MainQuit)

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

	items, err := device.ListDevices()
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
	combo.Connect("changed", func() {
		iter, err := combo.GetActiveIter()
		app.panicIfErr(err, "Error getting active iter:")

		value, err := store.GetValue(iter, 0)
		app.panicIfErr(err, "Error getting value: ")

		str, err := value.GetString()
		app.panicIfErr(err, "Error converting to string")
		if str == defaultStoreItemValue {
			return
		}
		// log.Println("Selected:", str)
		app.logArea.SetText("Selected device " + str)
		go func() { app.displayVideoInViewport(str) }()
	})

	return combo, nil
}

func (app *WBSApp) displayVideoInViewport(device_path string) {

	app.viewport.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		if app.currentPixbuf != nil {
			gtk.GdkCairoSetSourcePixBuf(cr, app.currentPixbuf, 0, 0)
			cr.Paint()
		}
	})

	go func() {
		var err error
		app.webcam, err = gocv.VideoCaptureDevice(0) //TODO get the device id from string
		app.panicIfErr(err, "Could not open video device 0")
		defer app.webcam.Close()

		frame := gocv.NewMat()
		defer frame.Close()

		for {
			// fmt.Println("Reading another frame")
			ok := app.webcam.Read(&frame)
			if !ok {
				log.Println("Cannot read from webcam")
			}

			if frame.Empty() {
				log.Println("Frame is empty")
				continue
			}

			pixbuf, err := device.MatToPixbuf(frame)
			if err != nil {
				log.Println("Failed to convert video frame to pixbuf: ", err)
				continue
			}
			app.currentPixbuf = pixbuf
			app.viewport.QueueDraw()
			time.Sleep(100 * time.Millisecond)
		}
	}()

}

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
