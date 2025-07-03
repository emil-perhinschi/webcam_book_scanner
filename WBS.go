package WBS

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Hello Fyne")

	hello := widget.NewLabel("Hello, Fyne!")
	button := widget.NewButton("Click me!", func() {
		hello.SetText("You clicked the button!")
	})

	myWindow.SetContent(container.NewVBox(hello, button))

	myWindow.ShowAndRun()
}
