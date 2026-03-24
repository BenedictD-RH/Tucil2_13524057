package object

import (
	"fyne.io/fyne/v2"
)

type Face struct {
	v1_idx int
	v2_idx int
	v3_idx int
}

func NewFace(i1, i2, i3 int) *Face {
	return &Face{i1, i2, i3}
}

func DrawFace(content *fyne.Container, O *Object, idx int) {
	// v1 := O.vertexArray[O.faceArray[idx].v1_idx - 1]
	// v2 := O.vertexArray[O.faceArray[idx].v2_idx - 1]
	// v3 := O.vertexArray[O.faceArray[idx].v3_idx - 1]
	DrawEdgeTo2DSpace(content, O.vertexArray[O.faceArray[idx].v1_idx-1], O.vertexArray[O.faceArray[idx].v2_idx-1])
	DrawEdgeTo2DSpace(content, O.vertexArray[O.faceArray[idx].v2_idx-1], O.vertexArray[O.faceArray[idx].v3_idx-1])
	DrawEdgeTo2DSpace(content, O.vertexArray[O.faceArray[idx].v1_idx-1], O.vertexArray[O.faceArray[idx].v3_idx-1])
	// face := canvas.NewArbitraryPolygon([]fyne.Position{VertexToFynePos(v1), VertexToFynePos(v2), VertexToFynePos(v3)}, color.White)
	// content.Add(face)
}