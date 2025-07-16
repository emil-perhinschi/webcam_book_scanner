package main

import (
	"errors"
	"fmt"
	"image"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/korandiz/v4l"
	"github.com/korandiz/v4l/fmt/yuyv"
)

type V4l2Device struct {
	shouldCapture bool
	isOpen        bool
	app           *WBSApp
	yuyv          *yuyv.Image
}

func NewW4l2Device(app *WBSApp) *V4l2Device {

	webcam := &V4l2Device{app: app, shouldCapture: false}

	return webcam
}

// call in go func
func (webcam *V4l2Device) CaptureFrames(devicePath string, config v4l.DeviceConfig) {
	if !webcam.shouldCapture {
		return
	}

	device, err := v4l.Open(devicePath)
	if err != nil {
		fmt.Println("Open", err)
		return
	}
	defer device.Close()

	if err := device.SetConfig(config); err != nil {
		fmt.Println("SetConfig", err)
		return
	}
	if err := device.TurnOn(); err != nil {
		fmt.Println("TurnOn", err)
		return
	}

	config, err = device.GetConfig()
	if err != nil {
		fmt.Println("GetConfig", err)
		return
	}
	if config.Format != yuyv.FourCC {
		fmt.Println("SetConfig", errors.New("failed to set YUYV format"))
		return
	}
	binfo, err := device.BufferInfo()
	if err != nil {
		fmt.Println("BufferInfo", err)
		return
	}
	webcam.yuyv = &yuyv.Image{
		Pix:    make([]byte, binfo.BufferSize),
		Stride: binfo.ImageStride,
		Rect:   image.Rect(0, 0, config.Width, config.Height),
	}
	for {
		if !webcam.shouldCapture {
			time.Sleep(1 * time.Second)
			continue
		}
		buffer, err := device.Capture()
		if err != nil {
			fmt.Println("Capture", err)
			break
		}
		buffer.ReadAt(webcam.yuyv.Pix, 0)
		ok := webcam.updateImage()
		if !ok {
			break
		}
	}
}

func (webcam *V4l2Device) updateImage() bool {
	ch := make(chan bool)
	glib.IdleAdd(func() {
		if !webcam.shouldCapture {
			ch <- false
			return
		}
		rgba := &image.RGBA{
			Pix:    webcam.app.currentPixbuf.GetPixels(),
			Stride: webcam.app.currentPixbuf.GetRowstride(),
			Rect:   image.Rect(0, 0, webcam.app.currentPixbuf.GetWidth(), webcam.app.currentPixbuf.GetHeight()),
		}
		yuyv.ToRGBA(rgba, rgba.Rect, webcam.yuyv, webcam.yuyv.Rect.Min)
		webcam.app.viewport.SetFromPixbuf(webcam.app.currentPixbuf)
		ch <- true
	})
	return <-ch
}

func (webcam *V4l2Device) Close() {
	webcam.shouldCapture = false
}

type V4l2DeviceConfig struct {
	width  int
	height int
	fps    int
}

func (webcam *V4l2Device) ReadDeviceCapabilities(devicePath string) ([]V4l2DeviceConfig, error) {
	var result []V4l2DeviceConfig

	dev, err := v4l.Open(devicePath)
	if err != nil {
		return nil, errors.New("Open: " + err.Error())
	}
	defer dev.Close()
	cfgs, err := dev.ListConfigs()
	if err != nil {
		return nil, errors.New("ListConfigs: " + err.Error())
	}
	for i := range cfgs {
		if cfgs[i].Format != yuyv.FourCC {
			continue
		}
		var item V4l2DeviceConfig
		item.width = cfgs[i].Width
		item.height = cfgs[i].Height
		item.fps = int(cfgs[i].FPS.N / cfgs[i].FPS.D)
		result = append(result, item)
	}
	return result, nil
}
