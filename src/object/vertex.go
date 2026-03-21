package object

import (
	"fmt"
	"go_project/matrix"
	"fyne.io/fyne/v2"
)

type Vertex struct {
	x float32
	y float32
	z float32
}

func NewVertex(x float32, y float32, z float32) *Vertex {
	return &Vertex{x, y, z}
}

func (v Vertex) String() string {
	return fmt.Sprintf("(%.6f, %.6f, %.6f)", v.x, v.y, v.z)
}

func VerticeTo2DPoint(v *Vertex) *Point {
	return &Point{(v.x - min_v)/(max_v - min_v)*viewbox_w + 100, (v.y - min_v)/(max_v - min_v)*viewbox_h + 100}
}

func DrawVertex(c *fyne.Container, v *Vertex) {
	DrawPoint(c, VerticeTo2DPoint(v))
}

func DrawEdge(c *fyne.Container, v1 *Vertex, v2 *Vertex) {
	DrawLine(c, &Line{VerticeTo2DPoint(v1), VerticeTo2DPoint(v2)})
}

func VertexToMatrix(v *Vertex) (*matrix.Matrix) {
	m := matrix.NewMatrix(3,1)
	m.Buffer[0][0] = v.x
	m.Buffer[1][0] = v.y
	m.Buffer[2][0] = v.z
	return m
}

func MatrixToVertex(m *matrix.Matrix) (*Vertex) {
	var x,y,z float32
	x = m.Buffer[0][0]
	y = m.Buffer[1][0]
	z = m.Buffer[2][0]
	return &Vertex{x,y,z}
}

func RotateVertex(v *Vertex, r_Mat *matrix.Matrix) (*Vertex) {
	return MatrixToVertex(matrix.MultiplyMatrix(r_Mat, VertexToMatrix(v)))
}
 
func ScanObject(content *fyne.Container) {
	var x1, y1, z1 float32
	var objectType string
	fmt.Scanf("%s %f %f %f", &objectType, &x1, &y1, &z1)
	if objectType == "v" {
		DrawVertex(content, &Vertex{x1, y1, z1})
		content.Refresh()
	} else {
		fmt.Println("Invalid")
	}
}
