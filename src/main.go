package main

import (
	//"fmt"
	"math"
	"go_project/matrix"
	"go_project/object"
	//"image/color"

	// "image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	//"fyne.io/fyne/v2/canvas"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var	r_Mat *matrix.Matrix = matrix.RotationMatrix_X(matrix.DegreeToRad(180))
var O *object.Object = object.ParseObject("../data/cow.obj")
var view_3d *fyne.Container = container.NewWithoutLayout()
type draggableRect struct {
	*widget.Label
}

func (d *draggableRect) Dragged(ev *fyne.DragEvent) {
    var dx, dy = float64(ev.Dragged.DX/600)*-1, float64(ev.Dragged.DY/600)
	deg := float32(math.Sqrt(float64(dx*dx + dy*dy))*360)
	r_Mat = matrix.RotateMatrixByDrag2D(dy, dx, deg, r_Mat)
	view_3d.RemoveAll()
	object.DrawObject(view_3d, object.RotateObject(O, r_Mat))
	view_3d.Refresh()
	d.Refresh()
}

func (d *draggableRect) DragEnd() {
	d.Move(fyne.NewPos(0,0))
	d.Refresh()
}

func main() {

	a := app.New()
	w := a.NewWindow("3D Model Viewer")
	w.Resize(fyne.NewSize(800, 800))
	content := container.NewWithoutLayout(view_3d)
	label := widget.NewLabel("HEY")
	label.Resize(fyne.NewSize(800,800))
	d_rect := &draggableRect{label}
	reset_button := widget.NewButton("Reset", func() {
		r_Mat = matrix.RotationMatrix_X(matrix.DegreeToRad(180))
		view_3d.RemoveAll()
		object.DrawObject(view_3d, object.RotateObject(O, r_Mat))
		view_3d.Refresh()
	})
	reset_button.Resize(fyne.NewSize(100,30))
	content.Add(d_rect)
	content.Add(reset_button)
	w.SetContent(content)

    
    object.DrawObject(view_3d, object.RotateObject(O, r_Mat))

	w.ShowAndRun()
}
