package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Hello Fyne")
	myWindow.Resize(fyne.NewSize(640, 480))

	myWindow.ShowAndRun()
}
