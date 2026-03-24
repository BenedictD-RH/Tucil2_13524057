package object

import (
	"bufio"
	"fmt"
	"go_project/matrix"
	"log"
	"math"
	"os"
	

	//"image/color"
	"fyne.io/fyne/v2"
	//"fyne.io/fyne/v2/container"
)

const viewbox_w = 600
const viewbox_h = 600

var Max_v float32 = -math.MaxFloat32
var Min_v float32 = math.MaxFloat32

var RenderedObject *ObjectRender

func (f Face) String() string {
	return fmt.Sprintf("(%d, %d, %d)", f.v1_idx, f.v2_idx, f.v3_idx)
}

type Object struct {
	vertexArray []*Vertex
	faceArray   []*Face
	edgeArray   []*Edge
}

func NewObject() *Object {
	return &Object{[]*Vertex{}, []*Face{}, []*Edge{}}
}

func AddVertex(O *Object, v *Vertex) {
	if max(v.x, v.y) > Max_v {
		Max_v = max(v.x, v.y)
	}
	if min(v.x, v.y) < Min_v {
		Min_v = min(v.x, v.y)
	}
	(*O).vertexArray = append((*O).vertexArray, v)
}

func AddFace(O *Object, f *Face) {
	(*O).faceArray = append((*O).faceArray, f)
	for _, edge := range getEdgeList(f) {
		if !isEdgeInObject(O, edge) {
			(*O).edgeArray = append((*O).edgeArray, edge)
		}
	}
}

func PrintObject(O *Object) {
	for _, vert := range O.vertexArray {
		fmt.Println("v", vert)
	}

	for _, face := range O.faceArray {
		fmt.Println("f", face)
	}
}

func DrawObject(content *fyne.Container, O *Object) {
	for i, _ := range O.edgeArray {
		DrawEdge(content, O, i)
	}
}

func DrawObjectInitial(content *fyne.Container, O *Object) {
	RenderedObject = NewObjectRender(O, content)
	for i, _ := range RenderedObject.polygons {
		UpdatePolygon(RenderedObject, i)
	}
}

func DrawObjectPolygons(O *Object) {
	RenderedObject.vertices = O.vertexArray
	UpdatePolygonList(RenderedObject)
	for i, _ := range RenderedObject.polygons {
		go UpdatePolygon(RenderedObject, i)
	}
}

func RotateObject(O *Object, r_Mat *matrix.Matrix) *Object {
	r_O := NewObject()
	r_O.faceArray = O.faceArray
	r_O.edgeArray = O.edgeArray
	Max_v = -math.MaxFloat32
	Min_v = math.MaxFloat32
	for _, vert := range O.vertexArray {
		AddVertex(r_O, RotateVertex(vert, r_Mat))
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

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error during scanning: %s", err)
	}
	return O
}
