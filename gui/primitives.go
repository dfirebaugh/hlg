package gui

import (
	"math"
	// "github.com/sirupsen/logrus"
)

func (d *Draw) drawTriangle(x1, y1, x2, y2, x3, y3 int, op *DrawOptions) {
	// if op.OutlineSize > 0 {
	// logrus.Warn("triangle outlines aren't supported")
	// }

	ndc1 := screenToNDC(float32(x1), float32(y1), float32(d.ScreenWidth), float32(d.ScreenHeight))
	ndc2 := screenToNDC(float32(x2), float32(y2), float32(d.ScreenWidth), float32(d.ScreenHeight))
	ndc3 := screenToNDC(float32(x3), float32(y3), float32(d.ScreenWidth), float32(d.ScreenHeight))

	fillColor := colorToFloat32(op.Style.FillColor)
	vertex1 := Vertex{
		Position: ndc1, FillColor: fillColor,
	}
	vertex2 := Vertex{
		Position: ndc2, FillColor: fillColor,
	}
	vertex3 := Vertex{
		Position: ndc3, FillColor: fillColor,
	}

	d.vertices = append(d.vertices, vertex1, vertex2, vertex3)
}

func (d *Draw) drawRectangle(x, y, width, height int, op *DrawOptions) {
	ndcTopLeft := screenToNDC(float32(x), float32(y), float32(d.ScreenWidth), float32(d.ScreenHeight))
	ndcTopRight := screenToNDC(float32(x+width), float32(y), float32(d.ScreenWidth), float32(d.ScreenHeight))
	ndcBottomLeft := screenToNDC(float32(x), float32(y+height), float32(d.ScreenWidth), float32(d.ScreenHeight))
	ndcBottomRight := screenToNDC(float32(x+width), float32(y+height), float32(d.ScreenWidth), float32(d.ScreenHeight))

	fillColor := colorToFloat32(op.FillColor)

	vertex1 := Vertex{
		Position: ndcTopLeft, FillColor: fillColor,
	}
	vertex2 := Vertex{
		Position: ndcBottomLeft, FillColor: fillColor,
	}
	vertex3 := Vertex{
		Position: ndcBottomRight, FillColor: fillColor,
	}
	vertex4 := Vertex{
		Position: ndcTopLeft, FillColor: fillColor,
	}
	vertex5 := Vertex{
		Position: ndcBottomRight, FillColor: fillColor,
	}
	vertex6 := Vertex{
		Position: ndcTopRight, FillColor: fillColor,
	}

	d.vertices = append(d.vertices, vertex1, vertex2, vertex3, vertex4, vertex5, vertex6)
}

func (d *Draw) drawRoundedRectangle(x, y, width, height, radius int, op *DrawOptions) {
	if radius > width/2 || radius > height/2 {
		radius = min(width/2, height/2)
	}

	// Draw the corner circles
	d.drawCircle(x+radius, y+radius, radius, op)
	d.drawCircle(x+width-radius, y+radius, radius, op)        // Top-right corner
	d.drawCircle(x+radius, y+height-radius, radius, op)       // Bottom-left corner
	d.drawCircle(x+width-radius, y+height-radius, radius, op) // Bottom-right corner

	// Draw the four side rectangles
	d.drawRectangle(x+radius, y, width-2*radius, radius, op)               // Top side
	d.drawRectangle(x+radius, y+height-radius, width-2*radius, radius, op) // Bottom side
	d.drawRectangle(x, y+radius, radius, height-2*radius, op)              // Left side
	d.drawRectangle(x+width-radius, y+radius, radius, height-2*radius, op) // Right side

	// Draw the center rectangle
	d.drawRectangle(x+radius, y+radius, width-2*radius, height-2*radius, op)
}

func (d *Draw) drawRoundedRectangleWithOutline(x, y, width, height, radius, outlineWidth int, op *DrawOptions) {
	// Draw the filled rounded rectangle
	d.drawRoundedRectangle(x, y, width, height, radius, op)

	outerX := x - outlineWidth
	outerY := y - outlineWidth
	outerWidth := width + 2*outlineWidth
	outerHeight := height + 2*outlineWidth
	outerRadius := radius + outlineWidth

	outlineOptions := *op // Copy the original options
	outlineOptions.FillColor = op.OutlineColor

	d.drawRoundedRectangle(outerX, outerY, outerWidth, outerHeight, outerRadius, &outlineOptions)

	innerX := x + outlineWidth
	innerY := y + outlineWidth
	innerWidth := width - 2*outlineWidth
	innerHeight := height - 2*outlineWidth
	innerRadius := radius

	innerFillOptions := *op // Copy the original options
	innerFillOptions.FillColor = op.FillColor

	d.drawRoundedRectangle(innerX, innerY, innerWidth, innerHeight, innerRadius, &innerFillOptions)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (d *Draw) drawCircleWithOutline(x, y, radius, outlineWidth int, op *DrawOptions) {
	// Draw the outer circle (the outline)
	outlineRadius := radius + outlineWidth
	outlineOptions := *op // Copy the original options
	outlineOptions.FillColor = op.OutlineColor

	d.drawCircle(x, y, outlineRadius, &outlineOptions)

	// Draw the inner circle (the filled area)
	innerRadius := radius
	innerFillOptions := *op // Copy the original options
	innerFillOptions.FillColor = op.FillColor

	d.drawCircle(x, y, innerRadius, &innerFillOptions)
}

func (d *Draw) drawCircle(x, y, radius int, op *DrawOptions) {
	segments := 36 // Number of segments to approximate the circle
	angleStep := 2 * math.Pi / float64(segments)

	centerNDC := screenToNDC(float32(x), float32(y), float32(d.ScreenWidth), float32(d.ScreenHeight))
	fillColor := colorToFloat32(op.FillColor)

	for i := 0; i < segments; i++ {
		theta1 := float64(i) * angleStep
		theta2 := float64(i+1) * angleStep

		x1 := float32(x) + float32(radius)*float32(math.Cos(theta1))
		y1 := float32(y) + float32(radius)*float32(math.Sin(theta1))

		x2 := float32(x) + float32(radius)*float32(math.Cos(theta2))
		y2 := float32(y) + float32(radius)*float32(math.Sin(theta2))

		ndc1 := screenToNDC(x1, y1, float32(d.ScreenWidth), float32(d.ScreenHeight))
		ndc2 := screenToNDC(x2, y2, float32(d.ScreenWidth), float32(d.ScreenHeight))

		vertex1 := Vertex{
			Position: centerNDC, FillColor: fillColor,
		}
		vertex2 := Vertex{
			Position: ndc1, FillColor: fillColor,
		}
		vertex3 := Vertex{
			Position: ndc2, FillColor: fillColor,
		}
		d.vertices = append(d.vertices, vertex1, vertex2, vertex3)
	}
}
