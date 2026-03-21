package main

import (
	// "fmt"
	"go_project/object"
	"go_project/matrix"
	// "image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	// "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)


func main() {
	r_Mat := matrix.RotationMatrix_X(matrix.DegreeToRad(180))
	r_Mat = matrix.MultiplyMatrix(matrix.RotationMatrix_Y(matrix.DegreeToRad(240)), r_Mat)
	r_Mat = matrix.MultiplyMatrix(matrix.RotationMatrix_X(matrix.DegreeToRad(-10)), r_Mat)
	a := app.New()
	w := a.NewWindow("3D Model Viewer")
	w.Resize(fyne.NewSize(800, 800))
	content := container.NewWithoutLayout()
	w.SetContent(content)

    O := object.ParseObject("../data/cow.obj")
    object.DrawObject(content, object.RotateObject(O, r_Mat))

	w.ShowAndRun()
}
