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

func (v *Vertex) GetX() float32 {
	return v.x
}

func (v *Vertex) GetY() float32 {
	return v.y
}

func (v *Vertex) GetZ() float32 {
	return v.z
}

func (v Vertex) String() string {
	return fmt.Sprintf("(%.6f, %.6f, %.6f)", v.x, v.y, v.z)
}

func VertexTo2DPoint(v *Vertex) *Point {
	return &Point{(v.x-Min_v)/(Max_v-Min_v)*float32(Viewbox_w) + float32(Viewbox_start_x),
		(v.y-Min_v)/(Max_v-Min_v)*float32(Viewbox_h) + float32(Viewbox_start_y)}
}

func VertexToFynePos(v *Vertex) fyne.Position {
	v2d := VertexTo2DPoint(v)
	return fyne.NewPos(v2d.x, v2d.y)
}

func DrawVertex(c *fyne.Container, v *Vertex) {
	DrawPoint(c, VertexTo2DPoint(v))
}

func VertexToMatrix(v *Vertex) *matrix.Matrix {
	m := matrix.NewMatrix(3, 1)
	m.Buffer[0][0] = v.x
	m.Buffer[1][0] = v.y
	m.Buffer[2][0] = v.z
	return m
}

func MatrixToVertex(m *matrix.Matrix) *Vertex {
	var x, y, z float32
	x = m.Buffer[0][0]
	y = m.Buffer[1][0]
	z = m.Buffer[2][0]
	return &Vertex{x, y, z}
}

func RotateVertex(v *Vertex, r_Mat *matrix.Matrix) *Vertex {
	return MatrixToVertex(matrix.MultiplyMatrix(r_Mat, VertexToMatrix(v)))
}

func SumVertex(v1, v2 *Vertex) *Vertex {
	return NewVertex(v1.x+v2.x, v1.y+v2.y, v1.z+v2.z)
}

func SubtractVertex(v1, v2 *Vertex) *Vertex {
	return NewVertex(v1.x-v2.x, v1.y-v2.y, v1.z-v2.z)
}

func ScaleVertex(v *Vertex, n float32) *Vertex {
	return NewVertex(v.x*n, v.y*n, v.z*n)
}

func MidVertex(v1, v2 *Vertex) *Vertex {
	return NewVertex((v1.x+v2.x)/2, (v1.y+v2.y)/2, (v1.z+v2.z)/2)
}
