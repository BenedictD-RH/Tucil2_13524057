package object

import (
	"fmt"
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type RoundedPoint struct {
	x int
	y int
}

type Point struct {
	x float32
	y float32
}

const point_rad = 5

func (p Point) String() string {
	return fmt.Sprintf("(%.0f, %.0f)", p.x, p.y)
}

func (p Point) Round() RoundedPoint {
	return RoundedPoint{int(p.x), int(p.y)}
}

func NewPoint(x,y float32) (*Point) {
	return &Point{x,y}
}

func Orientation(a,b,c Point) int {
	res := (b.y - a.y) * (c.x - b.x) - (c.y - b.y) * (b.x - a.x)
	if res == 0 {
		return 0
	} else if res > 0 {
		return 1
	}
	return -1
}

func IsPointInTriangle(p Point, tri [3]Point) bool {
	p_a := tri[0]
	p_b := tri[1]
	p_c := tri[2]
	denom := ((p_b.y-p_c.y)*(p_a.x-p_c.x) + (p_c.x-p_b.x)*(p_a.y-p_c.y))
	a := ((p_b.y-p_c.y)*(p.x-p_c.x) + (p_c.x-p_b.x)*(p.y-p_c.y)) / denom
	b := ((p_c.y-p_a.y)*(p.x-p_c.x) + (p_a.x-p_c.x)*(p.y-p_c.y)) / denom
	c := 1 - a - b
	if a >= 0 && b >= 0 && c >= 0 {
		return true
	}
	return false
}

func IsPolygonInTriangle(p1 []Point, tri [3]Point) bool {
	for _,p := range p1 {
		if (IsPointInTriangle(p, tri)) { return true }
	}
	return false
}

func IsPointInPolygon(p Point, poly []Point) bool {
	n_vert := len(poly)
	x_calc, y_calc := p.x, p.y
	inside := false
	p1 := poly[0]
	for i := 1; i < n_vert+1; i++ {
		p2 := poly[i%n_vert]
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

func IsPolygonInPolygon(p1, p2 []Point) (bool) {
	for _, p := range p1 {
		if (IsPointInPolygon(p, p2)) { return true }
	}
	return false
}


func OnSegment(p, q, r Point) bool {
	return (r.x <= max(p.x, q.x) && r.x >= min(p.x, q.x)) && (r.y <= max(p.y, q.y) && r.y >= min(p.y, q.y))
}


func DrawPoint(c *fyne.Container, p *Point) {
	point := canvas.NewCircle(color.White)
	point.StrokeWidth = 1
	point.Position1 = fyne.NewPos(p.x - point_rad, p.y - point_rad)
    point.Position2 = fyne.NewPos(p.x + point_rad, p.y + point_rad)
	c.Add(point)
}

func sortX(p1, p2 []Point) ([]Point) {
	i_1 := 0
	i_2 := 0
	merged := []Point{}
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
			if p1[i_1].x < p2[i_2].x {
				merged = append(merged, p1[i_1])
				i_1++
			} else if p1[i_1].x == p2[i_2].x {
				if p1[i_1].y < p2[i_2].y {
					merged = append(merged, p1[i_1])
					i_1++
				} else {
					merged = append(merged, p2[i_2])
					i_2++
				}
			} else {
				merged = append(merged, p2[i_2])
				i_2++
			}
		}
	}
	return merged
}

func MergeSortPointsX(p []Point) []Point {
	if len(p) <= 1 {
		return p
	} else {
		half_point := int(len(p) / 2)
		f_h := p[:half_point]
		s_h := p[half_point:]
		var sorted_f_h, sorted_s_h []Point
		sorted_f_h = MergeSortPointsX(f_h)
		sorted_s_h = MergeSortPointsX(s_h)
		return sortX(sorted_f_h, sorted_s_h)
	}
}
