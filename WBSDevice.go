package main

import (
	"fmt"
	"reflect"
	"sync"

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
}

func (webcam *WebcamDevice) IsOpened() bool {
	if webcam.captureDevice == nil {
		return false
	}
	fmt.Println(";;;; checking if camera is opened ;;;;", webcam)
	webcam.mtx.Lock()
	defer webcam.mtx.Unlock()
	isOpened := webcam.isOpen
	fmt.Println(";;;; DONE checking if camera is opened ;;;;", webcam)
	return isOpened
}
