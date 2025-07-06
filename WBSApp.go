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
	webcam        *WebcamDevice
}

func (app *WBSApp) OpenWebcamDevice() {
	var err error
	app.webcam.Open()
	app.panicIfErr(err, "Could not open video device 0")
	app.webcam.isOpen = true
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
			app.closeCamera()
			return
		}
		// log.Println("Selected:", str)
		app.logArea.SetText("Selected device " + str)
		app.displayVideoInViewport(str)
	})

	return combo, nil
}

func (app *WBSApp) closeCamera() {
	app.webcam.Close()
}

func (app *WBSApp) displayVideoInViewport(device_path string) {

	app.viewport.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		if app.currentPixbuf != nil {
			gtk.GdkCairoSetSourcePixBuf(cr, app.currentPixbuf, 0, 0)
			cr.Paint()
		}
	})

	go func() {

		app.OpenWebcamDevice()
		defer app.webcam.Close()

		frame := gocv.NewMat()
		defer frame.Close()

		for {
			// fmt.Println("Reading another frame")
			if !app.webcam.isOpen {
				return
			}

			if app.webcam == nil || !app.webcam.IsOpened() {
				time.Sleep(1 * time.Second)
				break
			}

			ok := app.webcam.captureDevice.Read(&frame)
			if !ok {
				log.Println("Cannot read from webcam")
				continue
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
