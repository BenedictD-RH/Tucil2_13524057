package object

import (
	"fmt"
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Point struct {
	x float32
	y float32
}

const point_rad = 5

func (p Point) String() string {
	return fmt.Sprintf("(%.0f, %.0f)", p.x, p.y)
}

func DrawPoint(c *fyne.Container, p *Point) {
	point := canvas.NewCircle(color.White)
	point.StrokeWidth = 1
	point.Position1 = fyne.NewPos(p.x - point_rad, p.y - point_rad)
    point.Position2 = fyne.NewPos(p.x + point_rad, p.y + point_rad)
	c.Add(point)
}