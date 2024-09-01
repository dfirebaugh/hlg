package gui

import (
	"image/color"
)

const (
	drawOpCircle = iota
	drawOpRoundedRect
	drawOpTriangle
)

func (d *Draw) drawTriangle(x1, y1, x2, y2, x3, y3 int, op *DrawOptions) {
	vertices := d.makeTriangleVertices(float32(x1), float32(y1), float32(x2), float32(y2), float32(x3), float32(y3), op.FillColor)
	d.vertices = append(d.vertices, vertices...)
}

func (d *Draw) drawRectangle(x, y, width, height int, op *DrawOptions) {
	centerX := float32(x + width/2)
	centerY := float32(y + height/2)
	vertices := d.makeRoundedRectVertices(centerX, centerY, float32(width), float32(height), op.FillColor)
	d.vertices = append(d.vertices, vertices...)
}

func (d *Draw) drawRoundedRectangle(x, y, width, height, radius int, op *DrawOptions) {
	centerX := float32(x + width/2)
	centerY := float32(y + height/2)
	vertices := d.makeRoundedRectVertices(centerX, centerY, float32(width), float32(height), op.FillColor)
	d.vertices = append(d.vertices, vertices...)
}

func (d *Draw) drawCircle(x, y, radius int, op *DrawOptions) {
	centerX := float32(x)
	centerY := float32(y)
	vertices := d.makeCircleVertices(centerX, centerY, float32(radius), op.FillColor)
	d.vertices = append(d.vertices, vertices...)
}

func (d *Draw) drawRoundedRectangleWithOutline(x, y, width, height, radius, outlineWidth int, op *DrawOptions) {
	// Draw the outer rounded rectangle (the outline)
	outerX := x - outlineWidth
	outerY := y - outlineWidth
	outerWidth := width + 2*outlineWidth
	outerHeight := height + 2*outlineWidth
	// outerRadius := radius + outlineWidth

	outlineOptions := *op // Copy the original options
	outlineOptions.FillColor = op.OutlineColor

	outerVertices := d.makeRoundedRectVertices(float32(outerX+outerWidth/2), float32(outerY+outerHeight/2), float32(outerWidth), float32(outerHeight), outlineOptions.FillColor)
	d.vertices = append(d.vertices, outerVertices...)

	// Draw the inner rounded rectangle (the filled area)
	innerX := x + outlineWidth
	innerY := y + outlineWidth
	innerWidth := width - 2*outlineWidth
	innerHeight := height - 2*outlineWidth

	innerVertices := d.makeRoundedRectVertices(float32(innerX+innerWidth/2), float32(innerY+innerHeight/2), float32(innerWidth), float32(innerHeight), op.FillColor)
	d.vertices = append(d.vertices, innerVertices...)
}

func (d *Draw) drawCircleWithOutline(x, y, radius, outlineWidth int, op *DrawOptions) {
	// Draw the outer circle (the outline)
	outlineRadius := radius + outlineWidth
	outlineOptions := *op // Copy the original options
	outlineOptions.FillColor = op.OutlineColor

	outlineVertices := d.makeCircleVertices(float32(x), float32(y), float32(outlineRadius), outlineOptions.FillColor)
	d.vertices = append(d.vertices, outlineVertices...)

	// Draw the inner circle (the filled area)
	innerVertices := d.makeCircleVertices(float32(x), float32(y), float32(radius), op.FillColor)
	d.vertices = append(d.vertices, innerVertices...)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (d *Draw) makeCircleVertices(centerX, centerY, radius float32, color color.Color) []Vertex {
	screenWidth := float32(d.ScreenWidth)
	screenHeight := float32(d.ScreenHeight)

	ndcCenterX := (centerX/screenWidth)*2.0 - 1.0
	ndcCenterY := 1.0 - (centerY/screenHeight)*2.0

	left := ndcCenterX - radius/screenWidth*2.0
	right := ndcCenterX + radius/screenWidth*2.0
	bottom := ndcCenterY - radius/screenHeight*2.0
	top := ndcCenterY + radius/screenHeight*2.0

	colorVec := colorToFloat32(color)

	return []Vertex{
		{Position: [3]float32{left, bottom, 0.0}, LocalPosition: [2]float32{-1.0, -1.0}, OpCode: drawOpCircle, Radius: radius, Color: colorVec},
		{Position: [3]float32{right, bottom, 0.0}, LocalPosition: [2]float32{1.0, -1.0}, OpCode: drawOpCircle, Radius: radius, Color: colorVec},
		{Position: [3]float32{left, top, 0.0}, LocalPosition: [2]float32{-1.0, 1.0}, OpCode: drawOpCircle, Radius: radius, Color: colorVec},

		{Position: [3]float32{left, top, 0.0}, LocalPosition: [2]float32{-1.0, 1.0}, OpCode: drawOpCircle, Radius: radius, Color: colorVec},
		{Position: [3]float32{right, bottom, 0.0}, LocalPosition: [2]float32{1.0, -1.0}, OpCode: drawOpCircle, Radius: radius, Color: colorVec},
		{Position: [3]float32{right, top, 0.0}, LocalPosition: [2]float32{1.0, 1.0}, OpCode: drawOpCircle, Radius: radius, Color: colorVec},
	}
}

func (d *Draw) makeRoundedRectVertices(centerX, centerY, width, height float32, color color.Color) []Vertex {
	screenWidth := float32(d.ScreenWidth)
	screenHeight := float32(d.ScreenHeight)

	ndcCenterX := (centerX/screenWidth)*2.0 - 1.0
	ndcCenterY := 1.0 - (centerY/screenHeight)*2.0
	ndcWidth := (width / screenWidth) * 2.0
	ndcHeight := (height / screenHeight) * 2.0

	left := ndcCenterX - ndcWidth/2
	right := ndcCenterX + ndcWidth/2
	bottom := ndcCenterY - ndcHeight/2
	top := ndcCenterY + ndcHeight/2

	colorVec := colorToFloat32(color)

	return []Vertex{
		{Position: [3]float32{left, bottom, 0.0}, LocalPosition: [2]float32{-1.0, -1.0}, OpCode: drawOpRoundedRect, Radius: 0.1, Color: colorVec},
		{Position: [3]float32{right, bottom, 0.0}, LocalPosition: [2]float32{1.0, -1.0}, OpCode: drawOpRoundedRect, Radius: 0.1, Color: colorVec},
		{Position: [3]float32{left, top, 0.0}, LocalPosition: [2]float32{-1.0, 1.0}, OpCode: drawOpRoundedRect, Radius: 0.1, Color: colorVec},

		{Position: [3]float32{left, top, 0.0}, LocalPosition: [2]float32{-1.0, 1.0}, OpCode: drawOpRoundedRect, Radius: 0.1, Color: colorVec},
		{Position: [3]float32{right, bottom, 0.0}, LocalPosition: [2]float32{1.0, -1.0}, OpCode: drawOpRoundedRect, Radius: 0.1, Color: colorVec},
		{Position: [3]float32{right, top, 0.0}, LocalPosition: [2]float32{1.0, 1.0}, OpCode: drawOpRoundedRect, Radius: 0.1, Color: colorVec},
	}
}

func (d *Draw) makeTriangleVertices(x1, y1, x2, y2, x3, y3 float32, color color.Color) []Vertex {
	screenWidth := float32(d.ScreenWidth)
	screenHeight := float32(d.ScreenHeight)

	// Convert each vertex position to NDC
	ndcX1 := (x1/screenWidth)*2.0 - 1.0
	ndcY1 := 1.0 - (y1/screenHeight)*2.0

	ndcX2 := (x2/screenWidth)*2.0 - 1.0
	ndcY2 := 1.0 - (y2/screenHeight)*2.0

	ndcX3 := (x3/screenWidth)*2.0 - 1.0
	ndcY3 := 1.0 - (y3/screenHeight)*2.0

	colorVec := colorToFloat32(color)

	return []Vertex{
		{Position: [3]float32{ndcX1, ndcY1, 0.0}, LocalPosition: [2]float32{0.0, 1.0}, OpCode: drawOpTriangle, Radius: 0.0, Color: colorVec},
		{Position: [3]float32{ndcX2, ndcY2, 0.0}, LocalPosition: [2]float32{-1.0, -1.0}, OpCode: drawOpTriangle, Radius: 0.0, Color: colorVec},
		{Position: [3]float32{ndcX3, ndcY3, 0.0}, LocalPosition: [2]float32{1.0, -1.0}, OpCode: drawOpTriangle, Radius: 0.0, Color: colorVec},
	}
}
