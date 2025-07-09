package main

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"gocv.io/x/gocv"
)

type WebcamDevice struct {
	mtx           sync.Mutex
	captureDevice *gocv.VideoCapture
	isOpen        bool
}

func (webcam *WebcamDevice) Open() {
	if webcam.isOpen {
		return
	}
	var err error
	webcam.captureDevice, err = gocv.VideoCaptureDevice(0) //TODO get the device id from string
	if err != nil {
		fmt.Println("ERROR OPENING VIDEO CAPTGURE DEVICE:", err)
		return
	}
	webcam.captureDevice.Set(gocv.VideoCaptureAutoFocus, 1)
	webcam.captureDevice.Set(gocv.VideoCaptureFrameWidth, 1920)
	webcam.captureDevice.Set(gocv.VideoCaptureFrameHeight, 1080)
	webcam.isOpen = true
}

func (webcam *WebcamDevice) Close() {
	webcam.mtx.Lock()
	defer webcam.mtx.Unlock()
	if !webcam.isOpen {
		return
	}

	if webcam.captureDevice != nil && webcam.captureDevice.IsOpened() {
		fmt.Println(":::::: camera is open, getting ready to close ", reflect.TypeOf(webcam), " ", webcam)
		webcam.captureDevice.Close()
		fmt.Println(":::::: camera was closed type is ", reflect.TypeOf(webcam), " ", webcam)
	}
	webcam.isOpen = false
	time.Sleep(1000 * time.Millisecond)
}

// func (webcam *WebcamDevice) IsOpened() bool {
// 	if webcam.captureDevice == nil {
// 		return false
// 	}
// 	isOpened := webcam.isOpen
// 	return isOpened
// }

func (webcam *WebcamDevice) CaptureImage(imageFilePath string) bool {
	if !webcam.isOpen {
		return false
	}

	img := gocv.NewMat()
	defer img.Close()
	if ok := webcam.captureDevice.Read(&img); !ok {
		fmt.Println("cannot read video device")
		return false
	}

	gocv.IMWrite(imageFilePath, img)

	return true
}
