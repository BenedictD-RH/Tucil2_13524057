package object

import (
	"bufio"
	"fmt"
	"go_project/matrix"
	"log"
	"math"
	"os"

	"fyne.io/fyne/v2"
	//"fyne.io/fyne/v2/container"
)

const screen_w = 800
const screen_h = 800
const viewbox_w = 600
const viewbox_h = 600
const viewbox_start_x = (screen_w - viewbox_w) / 2
const viewbox_start_y = (screen_h - viewbox_h) / 2

var Max_v float32 = -math.MaxFloat32
var Min_v float32 = math.MaxFloat32
var Center_v Vertex = Vertex{0, 0, 0}

var RenderedObject *ObjectRender

func (f Face) String() string {
	return fmt.Sprintf("(%d, %d, %d)", f.v1_idx, f.v2_idx, f.v3_idx)
}

type Object struct {
	VertexArray []*Vertex
	FaceArray   []*Face
	EdgeArray   []*Edge
	Dimension   *ObjectDim
}

type ObjectDim struct {
	min_x        float32
	max_x        float32
	min_y        float32
	max_y        float32
	min_z        float32
	max_z        float32
	centerVertex *Vertex
}

func GetSizeX(O *Object) float32 {
	return O.Dimension.max_x - O.Dimension.min_x
}

func GetSizeY(O *Object) float32 {
	return O.Dimension.max_y - O.Dimension.min_y
}

func GetSizeZ(O *Object) float32 {
	return O.Dimension.max_z - O.Dimension.min_z
}

func GetCenterVertex(O *Object) *Vertex {
	return O.Dimension.centerVertex
}

func FindObjectDimension(v []*Vertex) *ObjectDim {
	var max_x, min_x float32 = -math.MaxFloat32, math.MaxFloat32
	var max_y, min_y float32 = -math.MaxFloat32, math.MaxFloat32
	var max_z, min_z float32 = -math.MaxFloat32, math.MaxFloat32
	for _, vert := range v {
		if max_x < vert.x {
			max_x = vert.x
		}
		if min_x > vert.x {
			min_x = vert.x
		}
		if max_y < vert.y {
			max_y = vert.y
		}
		if min_y > vert.y {
			min_y = vert.y
		}
		if max_z < vert.z {
			max_z = vert.z
		}
		if min_z > vert.z {
			min_z = vert.z
		}
	}
	centerVertex := NewVertex((max_x+min_x)/2, (max_y+min_y)/2, (max_z+min_z)/2)
	return &ObjectDim{min_x, max_x, min_y, max_y, min_z, max_z, centerVertex}
}

func NewObject() *Object {
	return &Object{[]*Vertex{},
		[]*Face{},
		[]*Edge{},
		&ObjectDim{math.MaxFloat32, -math.MaxFloat32,
			math.MaxFloat32, -math.MaxFloat32,
			math.MaxFloat32, -math.MaxFloat32,
			nil}}
}

func GetFaceVertices(O *Object, f *Face) []*Vertex {
	vert := [3]*Vertex{}
	vert[0] = O.VertexArray[f.v1_idx - 1]
	vert[1] = O.VertexArray[f.v2_idx - 1]
	vert[2] = O.VertexArray[f.v3_idx - 1]
	return vert[:]
}

func AddVertex(O *Object, v *Vertex) {
	if max(v.x, v.y) > Max_v {
		Max_v = max(v.x, v.y)
	}
	if min(v.x, v.y) < Min_v {
		Min_v = min(v.x, v.y)
	}
	(*O).VertexArray = append((*O).VertexArray, v)
}

func AddFace(O *Object, f *Face) {
	(*O).FaceArray = append((*O).FaceArray, f)
	for _, edge := range getEdgeList(f) {
		if !isEdgeInObject(O, edge) {
			(*O).EdgeArray = append((*O).EdgeArray, edge)
		}
	}
}

func PrintObject(O *Object) {
	for _, vert := range O.VertexArray {
		fmt.Println("v", vert)
	}

	for _, face := range O.FaceArray {
		fmt.Println("f", face)
	}
}

func DrawObject(content *fyne.Container, O *Object) {
	for i, _ := range O.EdgeArray {
		DrawEdge(content, O, i)
	}
}

func DrawObjectInitial(content *fyne.Container, O *Object) {
	RenderedObject = NewObjectRender(O, content)
	// for i, _ := range RenderedObject.polygons {
	// 	UpdatePolygon(RenderedObject, i)
	// }
}

func DrawObjectPolygons(O *Object) {
	RenderedObject.Vertices = O.VertexArray
	UpdatePolygonList(RenderedObject)
	// for i, _ := range RenderedObject.polygons {
	// 	go UpdatePolygon(RenderedObject, i)
	// }
}

func RotateObject(O *Object, r_Mat *matrix.Matrix) *Object {
	r_O := NewObject()
	r_O.FaceArray = O.FaceArray
	r_O.EdgeArray = O.EdgeArray
	r_O.Dimension = O.Dimension
	Max_v = -math.MaxFloat32
	Min_v = math.MaxFloat32
	for _, vert := range O.VertexArray {
		rotated_vert := RotateVertex(SubtractVertex(vert, r_O.Dimension.centerVertex), r_Mat)
		AddVertex(r_O, SumVertex(rotated_vert, r_O.Dimension.centerVertex))
	}
	return r_O
}

func ParseObject(filename string) *Object {
	O := NewObject()
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) >= 2 {
			identifier := line[0:2]
			values := line[2:]
			if identifier == "v " {
				var x, y, z float32
				fmt.Sscanf(values, "%f %f %f", &x, &y, &z)
				AddVertex(O, NewVertex(x, y, z))
			} else if identifier == "f " {
				var i1, i2, i3 int
				fmt.Sscanf(values, "%d %d %d", &i1, &i2, &i3)
				AddFace(O, NewFace(i1, i2, i3))
			}
		}
	}
	O.Dimension = FindObjectDimension(O.VertexArray)
	Center_v = *O.Dimension.centerVertex
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error during scanning: %s", err)
	}
	return O
}
