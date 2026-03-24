package object

import (
	//"fmt"
	"math"
	"sync"
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/PerformLine/go-stockutil/colorutil"
)

type PolygonRender struct {
	face    *Face
	z_value float32
	render  *canvas.ArbitraryPolygon
}

type ObjectRender struct {
	polygons []*PolygonRender
	vertices []*Vertex
	stack *fyne.Container
}

func NewPolygonRender(O *Object, f_idx int) *PolygonRender {
	v1 := O.vertexArray[O.faceArray[f_idx].v1_idx-1]
	v2 := O.vertexArray[O.faceArray[f_idx].v2_idx-1]
	v3 := O.vertexArray[O.faceArray[f_idx].v3_idx-1]
	avg_z := (v1.z + v2.z + v3.z) / 3
	render := canvas.NewArbitraryPolygon([]fyne.Position{VertexToFynePos(v1),
														 VertexToFynePos(v2),
														 VertexToFynePos(v3)}, color.White)
	render.Hide()
	return &PolygonRender{O.faceArray[f_idx], avg_z, render}
}

func MergePolygons(p1, p2 []*PolygonRender) []*PolygonRender {
	i_1 := 0
	i_2 := 0
	merged := []*PolygonRender{}
	for {
		if i_1 == len(p1) && i_2 == len(p2) {
			break
		} else if i_1 == len(p1) {
			merged = append(merged, p2[i_2])
			i_2++
		} else if i_2 == len(p2) {
			merged = append(merged, p1[i_1])
			i_1++
		} else {
			if p1[i_1].z_value < p2[i_2].z_value {
				merged = append(merged, p1[i_1])
				i_1++
			} else {
				merged = append(merged, p2[i_2])
				i_2++
			}
		}
	}
	return merged
}

func MergeSortPolygons(p []*PolygonRender) []*PolygonRender {
	if len(p) <= 1 {
		return p
	} else {
		half_point := int(len(p) / 2)
		f_h := p[0:half_point]
		s_h := p[half_point:]
		sorted_f_h := MergeSortPolygons(f_h)
		sorted_s_h := MergeSortPolygons(s_h)
		return MergePolygons(sorted_f_h, sorted_s_h)
	}
}

func NewObjectRender(O *Object, content *fyne.Container) *ObjectRender {
	OR := &ObjectRender{[]*PolygonRender{}, O.vertexArray, content}
	for i, _ := range O.faceArray {
		OR.polygons = append(OR.polygons, NewPolygonRender(O, i))
	}
	OR.polygons = MergeSortPolygons(OR.polygons)
	for _, p := range OR.polygons {
		OR.stack.Add(p.render)
	}
	return OR
}

func isVertexInPolygon(OR *ObjectRender, p_idx int, v *Vertex) bool {
	p := OR.polygons[p_idx]
	v_a := OR.vertices[p.face.v1_idx-1]
	v_b := OR.vertices[p.face.v2_idx-1]
	v_c := OR.vertices[p.face.v3_idx-1]
	denom := ((v_b.y-v_c.y)*(v_a.x-v_c.x) + (v_c.x-v_b.x)*(v_a.y-v_c.y))
	a := ((v_b.y-v_c.y)*(v.x-v_c.x) + (v_c.x-v_b.x)*(v.y-v_c.y)) / denom
	b := ((v_c.y-v_a.y)*(v.x-v_c.x) + (v_a.x-v_c.x)*(v.y-v_c.y)) / denom
	c := 1 - a - b
	if a >= 0 && b >= 0 && c >= 0 {
		if a == 1 || b == 1 || c == 1 {
			return false
		}
		return true
	}
	return false
}

func shouldPolygonBeRendered(OR *ObjectRender, p_idx int) bool {
	p := OR.polygons[p_idx]
	v_a := OR.vertices[p.face.v1_idx-1]
	v_b := OR.vertices[p.face.v2_idx-1]
	v_c := OR.vertices[p.face.v3_idx-1]
	o_a, o_b, o_c := false, false, false
	for i, _ := range OR.polygons {
		if len(OR.polygons)-1-i > p_idx {
			if !o_a {
				o_a = isVertexInPolygon(OR, len(OR.polygons)-1-i, v_a)
			}
			if !o_b {
				o_b = isVertexInPolygon(OR, len(OR.polygons)-1-i, v_b)
			}
			if !o_c {
				o_c = isVertexInPolygon(OR, len(OR.polygons)-1-i, v_c)
			}
			if o_a && o_b && o_c {
				break
			}
		} else {
			return true
		}

	}
	if o_a && o_b && o_c {
		return false
	}

	return true
}

func UpdatePolygonList(OR *ObjectRender) {
	OR.stack.RemoveAll()
	var wg sync.WaitGroup
	for _, p := range OR.polygons {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.render.Hide()
			v1:= OR.vertices[p.face.v1_idx-1]
			v2:= OR.vertices[p.face.v2_idx-1]
			v3:= OR.vertices[p.face.v3_idx-1]
			p.render.Points = []fyne.Position{VertexToFynePos(v1), VertexToFynePos(v2), VertexToFynePos(v3)}
			p.z_value = (v1.z + v2.z + v3.z)/3
		}()
	}
	wg.Wait()
	OR.polygons = MergeSortPolygons(OR.polygons)
	for _, p := range OR.polygons {
		OR.stack.Add(p.render)
	}
}

func UpdatePolygon(OR *ObjectRender, idx int) {
	p := OR.polygons[idx]
	p.render.FillColor = PolygonColor(OR, idx)
	if shouldPolygonBeRendered(OR, idx) {
		p.render.Show()
	}
}

func PolygonColor(OR *ObjectRender, p_idx int) (color.Color) {
	p := OR.polygons[p_idx]
	v1:= OR.vertices[p.face.v1_idx-1]
	v2:= OR.vertices[p.face.v2_idx-1]
	v3:= OR.vertices[p.face.v3_idx-1]
	e1_x := v2.x - v1.x
	e1_y := v2.y - v1.y
	e1_z := v2.z - v1.z
	e2_x := v3.x - v1.x
	e2_y := v3.y - v1.y
	e2_z := v3.z - v1.z
	c_product_x := e1_y*e2_z - e1_z*e2_y
	c_product_y := - (e1_x*e2_z - e1_z*e2_x)
	c_product_z := e1_x*e2_y - e1_y*e2_x
	magnitude := math.Sqrt(float64(c_product_x*c_product_x + c_product_y*c_product_y + c_product_z*c_product_z))
	intensity := float64(math.Abs(float64(c_product_z)))/magnitude
	r, g, b := colorutil.HslToRgb(0, 0, max(0,intensity))
	return color.RGBA{r,g,b,255}
}
