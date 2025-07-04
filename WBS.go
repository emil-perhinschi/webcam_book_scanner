package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"gocv.io/x/gocv"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Hello Fyne")
	myWindow.Resize(fyne.NewSize(640, 480))

	webcam, errWebcam := gocv.VideoCaptureDevice(0)

	if errWebcam != nil {
		fmt.Println("ERORR WEBCAM INITIALIZATION: ", errWebcam)
		return
	}

	fyneImage := canvas.NewImageFromResource(nil)
	myWindow.SetContent(fyneImage)
	go func() {
		imgFrame := gocv.NewMat()
		defer imgFrame.Close()

		for {
			fyne.Do(func() {
				ok := webcam.Read(&imgFrame)
				if !ok {
					fmt.Println("WEBCAM READ error: Cannot read webcam frame")
					return
				}

				if imgFrame.Empty() {
					fmt.Println("imgFrane is empty")
					return
				}

				// rgbaImg := gocv.NewMat()
				// gocv.CvtColor(imgFrame, &rgbaImg, gocv.ColorBGR555ToRGBA)

				image, err := imgFrame.ToImage()
				if err != nil {
					fmt.Println("ERROR: ", err)
					return
				}
				fyneImage.Image = image
				fyneImage.Refresh()

				time.Sleep(500 * time.Millisecond)
			})
		}
	}()
	myWindow.ShowAndRun()
}
