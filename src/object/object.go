package object

import (
	"fmt"
	"fyne.io/fyne/v2"
	"os"
	"bufio"
	"log"
	"go_project/matrix"
)

const viewbox_w = 600
const viewbox_h = 600
var max_v float32 = 5
var min_v float32 = -5


type Face struct {
	v1_idx int
	v2_idx int
	v3_idx int
}

func NewFace(i1, i2, i3 int) *Face {
	return &Face{i1, i2, i3}
}

func (f Face) String() string {
	return fmt.Sprintf("(%d, %d, %d)", f.v1_idx, f.v2_idx, f.v3_idx)
}

type Object struct {
	vertexArray []*Vertex
	faceArray []*Face
}

func NewObject() *Object {
	return &Object{[]*Vertex{}, []*Face{}}
}

func AddVertex(O *Object, v *Vertex) {
	if (max(v.x, v.y) > max_v) {
		max_v = max(v.x, v.y)
	}
	if (min(v.x, v.y) < min_v) {
		min_v = min(v.x, v.y)
	}
	(*O).vertexArray = append((*O).vertexArray, v)
}

func AddFace(O *Object, f *Face) {
	(*O).faceArray = append((*O).faceArray, f)
}

func PrintObject(O *Object) {
	for _, vert := range O.vertexArray {
		fmt.Println("v", vert)
	}

	for _, face := range O.faceArray {
		fmt.Println("f", face)
	}
}

func DrawFace(content *fyne.Container, O* Object, idx int) {
	DrawEdge(content, O.vertexArray[O.faceArray[idx].v1_idx - 1], O.vertexArray[O.faceArray[idx].v2_idx - 1])
	DrawEdge(content, O.vertexArray[O.faceArray[idx].v2_idx - 1], O.vertexArray[O.faceArray[idx].v3_idx - 1])
	DrawEdge(content, O.vertexArray[O.faceArray[idx].v1_idx - 1], O.vertexArray[O.faceArray[idx].v3_idx - 1])
}

func DrawObject(content *fyne.Container, O *Object) {
	for _, vert := range O.vertexArray {
		DrawVertex(content, vert)
	}

	for i, _ := range O.faceArray {
		DrawFace(content, O, i)
	}
}

func RotateObject(O *Object, r_Mat *matrix.Matrix) (*Object) {
	r_O := NewObject()
	r_O.faceArray = O.faceArray
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
			if (identifier == "v ") {
				var x, y, z float32
				fmt.Sscanf(values, "%f %f %f", &x, &y, &z)
				AddVertex(O, NewVertex(x,y,z))
			} else if (identifier == "f ") {
				var i1, i2, i3 int
				fmt.Sscanf(values, "%d %d %d", &i1, &i2, &i3)
				AddFace(O, NewFace(i1,i2,i3))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error during scanning: %s", err)
	}
	return O
}