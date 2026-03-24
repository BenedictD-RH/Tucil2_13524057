package main

import (
	"go_project/matrix"
	"go_project/object"
	"math"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var	r_Mat *matrix.Matrix = matrix.RotationMatrix_X(matrix.DegreeToRad(180))
var last_r_Mat = matrix.RotationMatrix_X(matrix.DegreeToRad(180))
var O *object.Object = object.ParseObject("../data/cube.obj")
var last_x, last_y float64 = 0,0
var view_3d *fyne.Container = container.NewStack()
type draggableRect struct {
	*widget.Label
}

func (d *draggableRect) Dragged(ev *fyne.DragEvent) {
    last_x += float64(ev.Dragged.DX)
	last_y += float64(ev.Dragged.DY)
	deg := float32(math.Sqrt(float64(last_x*last_x + last_y*last_y)))*360/800
	last_r_Mat = matrix.RotateMatrixByDrag2D(last_y*-1, last_x, deg, r_Mat)
	object.DrawObjectPolygons(object.RotateObject(O, last_r_Mat))
	view_3d.Refresh()
	d.Refresh()
}

func (d *draggableRect) DragEnd() {
	d.Move(fyne.NewPos(0,0))
	last_x, last_y = 0,0
	r_Mat = last_r_Mat
	object.DrawObjectPolygons(object.RotateObject(O, last_r_Mat))
	view_3d.Refresh()
	d.Refresh()
	fmt.Println("Max_v : ", object.Max_v, "Min_v : ", object.Min_v)
}

func main() {
	a := app.NewWithID("1")
	w := a.NewWindow("3D Model Viewer")
	w.Resize(fyne.NewSize(800, 800))
	content := container.NewWithoutLayout()
	label := widget.NewLabel("Hey")
	label.Resize(fyne.NewSize(800,800))
	d_rect := &draggableRect{label}
	reset_button := widget.NewButton("Reset", func() {
		r_Mat = matrix.RotationMatrix_X(matrix.DegreeToRad(180))
		object.DrawObjectPolygons(object.RotateObject(O, r_Mat))
		view_3d.Refresh()
	})
	reset_button.Resize(fyne.NewSize(100,30))
	content.Add(d_rect)
	content.Add(reset_button)
	stack := container.NewStack(view_3d)
	stack.Add(content)
	w.SetContent(stack)
    object.DrawObjectInitial(view_3d, object.RotateObject(O, r_Mat))
	fmt.Println("Max_v : ", object.Max_v, "Min_v : ", object.Min_v)
	w.ShowAndRun()
}
