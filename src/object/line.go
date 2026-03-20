package object

import (
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

const line_w = 1

type Line struct {
	p1 *Point
	p2 *Point
}

func DrawLine(c *fyne.Container, l *Line) {
	line := canvas.NewLine(color.White)
	line.StrokeWidth = line_w
	line.Position1 = fyne.NewPos(l.p1.x, l.p1.y)
    line.Position2 = fyne.NewPos(l.p2.x, l.p2.y)
	c.Add(line)
}

