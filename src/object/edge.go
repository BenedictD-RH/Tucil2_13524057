package object

import (
	"fyne.io/fyne/v2"
)

type Edge struct {
	v1_idx int
	v2_idx int
}

func getEdgeList(f *Face) []*Edge {
	return []*Edge{{f.v1_idx, f.v2_idx},
		{f.v2_idx, f.v3_idx},
		{f.v1_idx, f.v3_idx}}
}

func isEdgeEqual(e1, e2 *Edge) bool {
	if e1.v1_idx == e2.v1_idx && e1.v2_idx == e2.v2_idx {
		return true
	} else if e1.v2_idx == e2.v1_idx && e1.v1_idx == e2.v2_idx {
		return true
	}
	return false
}

func isEdgeInObject(O *Object, e *Edge) bool {
	for _, edge := range O.edgeArray {
		if isEdgeEqual(edge, e) {
			return true
		}
	}
	return false
}

func isVertexInAnEdge(O *Object, v *Vertex) bool {
	for _, edge := range O.edgeArray {
		if (v == O.vertexArray[edge.v1_idx - 1] || v == O.vertexArray[edge.v2_idx - 1]) {
			return true
		}
	}
	return false
}

func DrawEdge(content *fyne.Container, O *Object, idx int) {
	DrawEdgeTo2DSpace(content, O.vertexArray[O.edgeArray[idx].v1_idx-1], O.vertexArray[O.edgeArray[idx].v2_idx-1])
}

func DrawEdgeTo2DSpace(c *fyne.Container, v1 *Vertex, v2 *Vertex) {
	DrawLine(c, &Line{VertexTo2DPoint(v1), VertexTo2DPoint(v2)})
}