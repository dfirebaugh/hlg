package gui

import (
	"image/color"
	"math"
)

type (
	DrawContext interface {
		DrawCircle(x int, y int, radius int, op *DrawOptions)
		DrawLine(points []Position, op *DrawOptions)
		DrawRectangle(x int, y int, width int, height int, op *DrawOptions)
		DrawSegment(x1 int, y1 int, x2 int, y2 int, op *DrawOptions)
		DrawTriangle(x1 int, y1 int, x2 int, y2 int, x3 int, y3 int, op *DrawOptions)
		Encode() []Vertex
		GetTheme() *Theme
		SetTheme(t *Theme)
	}

	Draw struct {
		ScreenWidth, ScreenHeight int
		*Theme

		vertices []Vertex
	}
	DrawOptions struct {
		Style
	}
	Position struct {
		X, Y, Z int
	}
	Style struct {
		FillColor    color.Color
		OutlineColor color.Color
		OutlineSize  int
		CornerRadius int
	}
)

func NewDrawContext(width, height int) DrawContext {
	return &Draw{
		ScreenWidth:  width,
		ScreenHeight: height,
		Theme:        DefaultTheme(),
	}
}

func (d *Draw) DrawTriangle(x1, y1, x2, y2, x3, y3 int, op *DrawOptions) {
	if op.OutlineSize > 0 {
		outlineOp := *op
		outlineOp.FillColor = op.OutlineColor
	}
	d.drawTriangle(x1, y1, x2, y2, x3, y3, op)
}

func (d *Draw) DrawRectangle(x, y, width, height int, op *DrawOptions) {
	if op.CornerRadius > 0 {
		if op.OutlineSize > 0 && op.CornerRadius > 0 {
			d.drawRoundedRectangleWithOutline(x, y, width, height, op.CornerRadius, op.OutlineSize, op)
		} else {
			d.drawRoundedRectangle(x, y, width, height, op.CornerRadius, op)
		}
	} else {
		if op.OutlineSize > 0 {
			outlineOp := *op
			outlineOp.FillColor = op.OutlineColor
		}
		d.drawRectangle(x, y, width, height, op)
	}
}

func (d *Draw) DrawCircle(x, y, radius int, op *DrawOptions) {
	if op.OutlineSize > 0 {
		d.drawCircleWithOutline(x, y, radius, op.OutlineSize, op)
	} else {
		d.drawCircle(x, y, radius, op)
	}
}

func (d *Draw) DrawSegment(x1, y1, x2, y2 int, op *DrawOptions) {
	dx := x2 - x1
	dy := y2 - y1
	length := int(math.Sqrt(float64(dx*dx + dy*dy)))
	_ = length
	length = 1

	d.DrawRectangle(x1, y1, length, op.OutlineSize, op)
}

func (d *Draw) DrawLine(points []Position, op *DrawOptions) {
	for i := 0; i < len(points)-1; i++ {
		d.DrawSegment(points[i].X, points[i].Y, points[i+1].X, points[i+1].Y, op)
	}
}

func (d *Draw) Encode() []Vertex {
	return d.vertices
}
