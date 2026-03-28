package object

import (
	"math"
	// "runtime"
	//"sync"
	"image/color"
)

func isPointInSet(p Point, s []Point) bool {
	for _, p_s := range s {
		if (math.Round(float64(p_s.x)) == math.Round(float64(p.x))) && (math.Round(float64(p_s.y)) == math.Round(float64(p.y))) {
			
			return true
		}
	}
	return false
}

func GetRender2DVertices(OR *ObjectRender) []Point {
	p := []Point{}
	for _, vert := range OR.Vertices {
		if !isPointInSet(*VertexTo2DPoint(vert), p) {
			p = append(p, *VertexTo2DPoint(vert))
		}
	}
	return MergeSortPointsX(p)
}

func PixelZValue(x, y int, OR *ObjectRender, p_idx int) float32 {
	p := OR.Polygons[p_idx]
	x_, y_ := float32(x), float32(y)
	x1, y1 := VertexToFynePos(OR.Vertices[p.face.v1_idx-1]).X, VertexToFynePos(OR.Vertices[p.face.v1_idx-1]).Y
	z1 := OR.Vertices[p.face.v1_idx-1].z
	x2, y2 := VertexToFynePos(OR.Vertices[p.face.v2_idx-1]).X, VertexToFynePos(OR.Vertices[p.face.v2_idx-1]).Y
	z2 := OR.Vertices[p.face.v2_idx-1].z
	x3, y3 := VertexToFynePos(OR.Vertices[p.face.v3_idx-1]).X, VertexToFynePos(OR.Vertices[p.face.v3_idx-1]).Y
	z3 := OR.Vertices[p.face.v3_idx-1].z
	A := y1*(z2-z3) + y2*(z3-z1) + y3*(z1-z2)
	B := z1*(x2-x3) + z2*(x3-x1) + z3*(x1-x2)
	C := x1*(y2-y3) + x2*(y3-y1) + x3*(y1-y2)
	D := -x1*(y2*z3-y3*z2) - x2*(y3*z1-y1*z3) - x3*(y1*z2-y2*z1)
	return (-D - A*x_ - B*y_) / C
}

func isPixelInConvexHull(x, y int, ch []Point) bool {
	n_vert := len(ch)
	x_calc, y_calc := float32(x), float32(y)
	inside := false
	p1 := ch[0]
	for i := 1; i < n_vert+1; i++ {
		p2 := ch[i%n_vert]
		if y_calc > min(p1.y, p2.y) {
			if y_calc < max(p1.y, p2.y) {
				if x_calc < max(p1.x, p2.x) {
					x_intersection := (y_calc-p1.y)*(p2.x-p1.x)/(p2.y-p1.y) + p1.x
					if p1.x == p2.x || x_calc <= x_intersection {
						inside = !inside
					}
				}
			}
		}
		p1 = p2
	}
	return inside
}

func isPixelInPolygon(x, y int, OR *ObjectRender, p_idx int) bool {
	p := OR.Polygons[p_idx]
	x_, y_ := float32(x), float32(y)
	v_a := VertexToFynePos(OR.Vertices[p.face.v1_idx-1])
	v_b := VertexToFynePos(OR.Vertices[p.face.v2_idx-1])
	v_c := VertexToFynePos(OR.Vertices[p.face.v3_idx-1])
	denom := ((v_b.Y-v_c.Y)*(v_a.X-v_c.X) + (v_c.X-v_b.X)*(v_a.Y-v_c.Y))
	a := ((v_b.Y-v_c.Y)*(x_-v_c.X) + (v_c.X-v_b.X)*(y_-v_c.Y)) / denom
	b := ((v_c.Y-v_a.Y)*(x_-v_c.X) + (v_a.X-v_c.X)*(y_-v_c.Y)) / denom
	c := 1 - a - b
	if a >= 0 && b >= 0 && c >= 0 {
		return true
	}
	return false
}

func GetPixelColor(x, y int, OR *ObjectRender) color.Color {
	if !isPixelInConvexHull(x, y, OR.ConvexHull) {
		return color.Transparent
	}
	for i, _ := range OR.Polygons {
		r_i := len(OR.Polygons) - 1 - i
		if isPixelInPolygon(x, y, OR, r_i) {
			higher_exist_below := false
			j := 1
			n := 0
			for {
				if (r_i-j >= 0) && (j < len(OR.Polygons)/10 || len(OR.Polygons) < 20) {
					if isPixelInPolygon(x, y, OR, r_i-j) {
						if PixelZValue(x, y, OR, r_i) < PixelZValue(x, y, OR, r_i-j) {
							higher_exist_below = true
							break
						}
						n++
					}
				} else {
					break
				}
				j++
				if n >= 10 {
					break
				}
			}
			if !higher_exist_below {
				return *OR.Polygons[r_i].color
			}
		}
	}
	return color.Transparent
}

func InitializeRasterBuffer(OR *ObjectRender) {
	for y, _ := range OR.rasterBuffer {
		for x, _ := range OR.rasterBuffer[y] {
			OR.rasterBuffer[y][x] = new(color.Color)
			*(OR.rasterBuffer[y][x]) = color.Transparent
		}
	}
}

func UpdateRasterBuffer(OR *ObjectRender) {
	start_y := (len(OR.rasterBuffer) - viewbox_h) / 2
	end_y := len(OR.rasterBuffer) - start_y
	start_x := (len(OR.rasterBuffer[0]) - viewbox_w) / 2
	end_x := len(OR.rasterBuffer[0]) - start_x
	for y := start_y; y < end_y; y++ {
		for x := start_x; x < end_x; x++ {
			go func() {
				*(OR.rasterBuffer[y][x]) = GetPixelColor(x, y, OR)
			}()
		}
	}
}
