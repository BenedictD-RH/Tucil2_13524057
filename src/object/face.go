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

func (f *Face) Increment(inc int) *Face {
	return NewFace(f.v1_idx + inc, f.v2_idx + inc, f.v3_idx + inc)
}

func (f *Face) ListOfIdx() []int {
	return []int{f.v1_idx, f.v2_idx, f.v3_idx}
}

func DrawFace(content *fyne.Container, O *Object, idx int) {
	DrawEdgeTo2DSpace(content, O.VertexArray[O.FaceArray[idx].v1_idx-1], O.VertexArray[O.FaceArray[idx].v2_idx-1])
	DrawEdgeTo2DSpace(content, O.VertexArray[O.FaceArray[idx].v2_idx-1], O.VertexArray[O.FaceArray[idx].v3_idx-1])
	DrawEdgeTo2DSpace(content, O.VertexArray[O.FaceArray[idx].v1_idx-1], O.VertexArray[O.FaceArray[idx].v3_idx-1])

}
