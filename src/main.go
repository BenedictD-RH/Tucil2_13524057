package main

import (
	"fmt"
	"go_project/matrix"
	"go_project/object"
	"go_project/voxel"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var r_Mat *matrix.Matrix = matrix.RotationMatrix_X(matrix.DegreeToRad(180))
var last_r_Mat = matrix.RotationMatrix_X(matrix.DegreeToRad(180))
var O *object.Object = object.NewObject()
var unvoxelized_O *object.Object = object.NewObject()
var last_x, last_y float64 = 0, 0
var view_3d *fyne.Container = container.NewStack()
var curr_depth = 1
var message_interface *fyne.Container = container.NewWithoutLayout()
var curr_file_name = ""

type draggableRect struct {
	*widget.Label
}

func (d *draggableRect) Dragged(ev *fyne.DragEvent) {
	last_x += float64(ev.Dragged.DX)
	last_y += float64(ev.Dragged.DY)
	deg := float32(math.Sqrt(float64(last_x*last_x+last_y*last_y))) * 360 / 800
	if !O.IsEmpty() {
		last_r_Mat = matrix.RotateMatrixByDrag2D(last_y*-1, last_x, deg, r_Mat)
		object.DrawObjectPolygons(object.RotateObject(O, last_r_Mat))
		view_3d.Refresh()
	}
	d.Refresh()
}

func (d *draggableRect) DragEnd() {
	d.Move(fyne.NewPos(100, 100))
	last_x, last_y = 0, 0
	if !O.IsEmpty() {
		r_Mat = last_r_Mat
		object.DrawObjectPolygons(object.RotateObject(O, last_r_Mat))
		view_3d.Refresh()
	}
	d.Refresh()
}

func SendNewMessage(s string) {
	msg := widget.NewLabel(s)
	msg.Resize(fyne.NewSize(200, 40))
	msg.Move(fyne.NewPos(0, 750))
	for _, m := range message_interface.Objects {
		m.Move(m.Position().Add(fyne.NewDelta(0, -15)))
	}
	fyne.Do(func() { message_interface.Refresh() })
	message_interface.Add(msg)
	go func() {
		time.Sleep(15 * time.Second)
		fyne.Do(func() {
			message_interface.Remove(msg)
			message_interface.Refresh()
		})
	}()
}

func VoxelizeObject(O *object.Object, depth int) *object.Object {
	O_N := object.NormalizeObject(O, voxel.MaxEdgeLength)
	Oct := voxel.NewOctreeRoot(O_N)
	for i := 0; i < depth; i++ {
		voxel.IncreaseDepth(Oct)
	}
	SendNewMessage(fmt.Sprintf("Amount of Voxels : %d", voxel.OctreeVoxelAmount(Oct)))
	O_N2 := voxel.EraseDuplicates(voxel.OctreeToVoxelObject(Oct))
	SendNewMessage(fmt.Sprintf("Amount of Vertex : %d", len(O_N2.VertexArray)))
	SendNewMessage(fmt.Sprintf("Amount of Faces : %d", len(O_N2.FaceArray)))
	SendNewMessage("Amount of Octree Nodes at each depth : ")
	for i := 1; i <= depth; i++ {
		SendNewMessage(fmt.Sprintf("%d : %d", i, voxel.GetAmountOfNodeAtDepth(Oct, i)))
	}
	SendNewMessage("Amount of Empty Octree Nodes at each depth : ")
	for i := 1; i <= depth; i++ {
		SendNewMessage(fmt.Sprintf("%d : %d", i, voxel.GetAmountOfEmptyNodesAtDepth(Oct, i)))
	}
	SendNewMessage("Translating...")
	ONew := O_N2.Translate()
	go func() {
		O_New2 := voxel.EraseInnerFaces(O_N2).Translate()

		O = O_New2
		fyne.Do(func() {
			view_3d.RemoveAll()
			object.DrawObjectInitial(view_3d, object.RotateObject(O, r_Mat))
			view_3d.Refresh()
		})
	}()
	SendNewMessage("Saving...")
	save_path := object.SaveObject(curr_file_name+"_Voxel", ONew)
	SendNewMessage("Object sucessfully saved at" + save_path)
	return ONew
}

func main() {
	//r_Mat = matrix.MultiplyMatrix(matrix.RotationMatrix_Y(matrix.DegreeToRad(90)), r_Mat)
	// r_Mat = matrix.MultiplyMatrix(matrix.RotationMatrix_Z(matrix.DegreeToRad(30)), r_Mat)
	// r_Mat = matrix.MultiplyMatrix(matrix.RotationMatrix_X(matrix.DegreeToRad(-30)), r_Mat)
	last_r_Mat = r_Mat
	a := app.NewWithID("1")
	w := a.NewWindow("3D Model Viewer")
	w.Resize(fyne.NewSize(800, 800))
	content := container.NewWithoutLayout()
	watermark := widget.NewLabel("Made by Benedict Darrel Setiawan")
	watermark.Resize(fyne.NewSize(800, 800))
	label_empty := widget.NewLabel("")
	d_rect := &draggableRect{label_empty}
	d_rect.Resize(fyne.NewSize(600, 600))
	d_rect.Move(fyne.NewPos(100, 100))
	reset_button := widget.NewButton("Reset", func() {
		if !O.IsEmpty() {
			r_Mat = matrix.RotationMatrix_X(matrix.DegreeToRad(180))
			object.InitializeRasterBuffer(object.RenderedObject)
			object.DrawObjectPolygons(object.RotateObject(O, r_Mat))
			view_3d.Refresh()
		}
	})
	reset_button.Resize(fyne.NewSize(100, 40))
	reset_button.Move(fyne.NewPos(680, 480))

	zoom_in_button := widget.NewButton("+", func() {
		if !O.IsEmpty() {
			object.Viewbox_h = int(float32(object.Viewbox_h) * 1.2)
			object.Viewbox_w = int(float32(object.Viewbox_w) * 1.2)
			object.Viewbox_start_x = (object.Screen_w - object.Viewbox_w) / 2
			object.Viewbox_start_y = (object.Screen_h - object.Viewbox_h) / 2
			object.InitializeRasterBuffer(object.RenderedObject)
			object.DrawObjectPolygons(object.RotateObject(O, r_Mat))
			view_3d.Refresh()
		}
	})
	zoom_in_button.Resize(fyne.NewSize(50, 50))
	zoom_in_button.Move(fyne.NewPos(730, 680))

	zoom_out_button := widget.NewButton("-", func() {
		if !O.IsEmpty() {
			object.Viewbox_h = int(float32(object.Viewbox_h) / 1.2)
			object.Viewbox_w = int(float32(object.Viewbox_w) / 1.2)
			object.Viewbox_start_x = (object.Screen_w - object.Viewbox_w) / 2
			object.Viewbox_start_y = (object.Screen_h - object.Viewbox_h) / 2
			object.InitializeRasterBuffer(object.RenderedObject)
			object.DrawObjectPolygons(object.RotateObject(O, r_Mat))
			view_3d.Refresh()
		}
	})
	zoom_out_button.Resize(fyne.NewSize(50, 50))
	zoom_out_button.Move(fyne.NewPos(730, 735))

	file_entry := widget.NewEntry()
	file_entry.SetPlaceHolder("Input File Path")
	file_entry.SetText("../data/")
	file_entry.Resize(fyne.NewSize(200, 40))
	file_entry.Move(fyne.NewPos(550, 20))
	file_entry.OnSubmitted = func(s string) {
		_, err := os.ReadFile(s)
		if err != nil {
			SendNewMessage("Incorrect file input")
		} else {
			view_3d.RemoveAll()
			O = object.ParseObject(s)
			substrings := strings.Split(s, "/")
			fmt.Println(substrings[len(substrings)-1])
			curr_file_name = strings.Split(substrings[len(substrings)-1], ".obj")[0]
			unvoxelized_O = O
			O = object.NormalizeObject(O, voxel.MaxEdgeLength)
			object.DrawObjectInitial(view_3d, object.RotateObject(O, r_Mat))
			SendNewMessage(s + " successfully parsed!")
		}
	}

	pitch_entry := widget.NewEntry()
	pitch_entry.SetPlaceHolder("Pitch")
	pitch_entry.Resize(fyne.NewSize(100, 40))
	pitch_entry.Move(fyne.NewPos(680, 530))
	pitch_entry.OnSubmitted = func(s string) {
		if !O.IsEmpty() {
			deg, err := strconv.ParseFloat(s, 32)
			if err != nil {
				SendNewMessage("Invalid input for Pitch!")
			} else {
				r_Mat = matrix.MultiplyMatrix(matrix.RotationMatrix_X(matrix.DegreeToRad(float32(deg))), r_Mat)
				object.DrawObjectPolygons(object.RotateObject(O, r_Mat))
				view_3d.Refresh()
			}
		}
		pitch_entry.SetText("")
	}

	yaw_entry := widget.NewEntry()
	yaw_entry.SetPlaceHolder("Yaw")
	yaw_entry.Resize(fyne.NewSize(100, 40))
	yaw_entry.Move(fyne.NewPos(680, 580))
	yaw_entry.OnSubmitted = func(s string) {
		if !O.IsEmpty() {
			deg, err := strconv.ParseFloat(s, 32)
			if err != nil {
				SendNewMessage("Invalid input for Yaw!")
			} else {
				r_Mat = matrix.MultiplyMatrix(matrix.RotationMatrix_Y(matrix.DegreeToRad(float32(deg))), r_Mat)
				object.DrawObjectPolygons(object.RotateObject(O, r_Mat))
				view_3d.Refresh()
			}
		}
		yaw_entry.SetText("")
	}

	roll_entry := widget.NewEntry()
	roll_entry.SetPlaceHolder("Roll")
	roll_entry.Resize(fyne.NewSize(100, 40))
	roll_entry.Move(fyne.NewPos(680, 630))
	roll_entry.OnSubmitted = func(s string) {
		if !O.IsEmpty() {
			deg, err := strconv.ParseFloat(s, 32)
			if err != nil {
				SendNewMessage("Invalid input for Roll!")
			} else {
				r_Mat = matrix.MultiplyMatrix(matrix.RotationMatrix_Z(matrix.DegreeToRad(float32(deg))), r_Mat)
				object.DrawObjectPolygons(object.RotateObject(O, r_Mat))
				view_3d.Refresh()
			}
		}
		roll_entry.SetText("")
	}

	voxelize_button := widget.NewButton(fmt.Sprintf("Voxelize Object at Depth %d", curr_depth), func() {
		if !O.IsEmpty() {
			start := time.Now()
			O = VoxelizeObject(unvoxelized_O, curr_depth)
			SendNewMessage(fmt.Sprintf("Time Elapsed : %s", time.Since(start)))
			fyne.Do(func() {
				view_3d.RemoveAll()
				object.DrawObjectInitial(view_3d, object.RotateObject(O, r_Mat))
				view_3d.Refresh()
			})
		}
	})
	voxelize_button.Resize(fyne.NewSize(200, 60))
	voxelize_button.Move(fyne.NewPos(300, 10))

	inc_depth_button := widget.NewButton(">", func() {
		curr_depth++
		voxelize_button.SetText(fmt.Sprintf("Voxelize Object at Depth %d", curr_depth))
	})
	inc_depth_button.Resize(fyne.NewSize(30, 60))
	inc_depth_button.Move(fyne.NewPos(510, 10))

	dcr_depth_button := widget.NewButton("<", func() {
		if curr_depth > 1 {
			curr_depth--
		}
		voxelize_button.SetText(fmt.Sprintf("Voxelize Object at Depth %d", curr_depth))
	})
	dcr_depth_button.Resize(fyne.NewSize(30, 60))
	dcr_depth_button.Move(fyne.NewPos(260, 10))

	content.Add(d_rect)
	content.Add(reset_button)
	content.Add(zoom_in_button)
	content.Add(zoom_out_button)
	content.Add(file_entry)
	content.Add(pitch_entry)
	content.Add(yaw_entry)
	content.Add(roll_entry)
	content.Add(voxelize_button)
	content.Add(inc_depth_button)
	content.Add(dcr_depth_button)
	content.Add(message_interface)
	content.Add(watermark)
	stack := container.NewStack(view_3d)
	stack.Add(content)
	w.SetContent(stack)
	w.ShowAndRun()
}
