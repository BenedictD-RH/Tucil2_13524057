package main

import (
	"go_project/object"

	//"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	//"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)


func main() {
	a := app.New()
	w := a.NewWindow("3D Model Viewer")
	w.Resize(fyne.NewSize(800, 800))
	content := container.NewWithoutLayout()
	w.SetContent(content)

    O := object.ParseObject("../data/cow.obj")
    object.DrawObject(content, O)

	w.ShowAndRun()
}
