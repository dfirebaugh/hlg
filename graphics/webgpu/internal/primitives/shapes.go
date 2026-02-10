//go:build !js

package primitives

import (
	"image/color"
	"math"

	"github.com/dfirebaugh/hlg/graphics"
)

// MakeTriangle creates a new Triangle renderable and adds it to the RenderQueue.
// It returns a reference to the created Triangle.
func MakeTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) []Vertex {
	r, g, b, a := c.RGBA()
	triangle := []Vertex{
		{
			Position: [3]float32{float32(x1), float32(y1), 0},
			Color:    [4]float32{float32(r) / 0xffff, float32(g) / 0xffff, float32(b) / 0xffff, float32(a) / 0xffff},
		},
		{
			Position: [3]float32{float32(x2), float32(y2), 0},
			Color:    [4]float32{float32(r) / 0xffff, float32(g) / 0xffff, float32(b) / 0xffff, float32(a) / 0xffff},
		},
		{
			Position: [3]float32{float32(x3), float32(y3), 0},
			Color:    [4]float32{float32(r) / 0xffff, float32(g) / 0xffff, float32(b) / 0xffff, float32(a) / 0xffff},
		},
	}

	return triangle
}

// MakeRectangle creates a new Rectangle renderable and adds it to the RenderQueue.
// It returns a reference to the created Rectangle.
func MakeRectangle(x, y, width, height int, c color.Color) []Vertex {
	r, g, b, a := c.RGBA()
	colorArray := [4]float32{float32(r) / 0xffff, float32(g) / 0xffff, float32(b) / 0xffff, float32(a) / 0xffff}

	topLeft := Vertex{
		Position: [3]float32{float32(x), float32(y), 0},
		Color:    colorArray,
	}
	topRight := Vertex{
		Position: [3]float32{float32(x + width), float32(y), 0},
		Color:    colorArray,
	}
	bottomLeft := Vertex{
		Position: [3]float32{float32(x), float32(y + height), 0},
		Color:    colorArray,
	}
	bottomRight := Vertex{
		Position: [3]float32{float32(x + width), float32(y + height), 0},
		Color:    colorArray,
	}

	rectangleVertices := []Vertex{topLeft, bottomLeft, topRight, bottomLeft, bottomRight, topRight}

	return rectangleVertices
}

// MakeCircle creates a new Circle renderable and adds it to the RenderQueue.
// It returns a reference to the created Circle.
func MakeCircle(cx, cy int, radius float32, c color.Color, segments int) []Vertex {
	// note: we could probably more efficiently draw circles with a custom shader -- but this is a good start
	circle := MakePolygon(cx, cy, radius*2, c, segments)
	return circle
}

func MakePolygonFromVertices(cx, cy int, width float32, vertices []graphics.Vertex) []Vertex {
	v := make([]Vertex, len(vertices))
	for i, vertex := range vertices {
		v[i] = Vertex{
			Position: [3]float32{vertex.Position[0], vertex.Position[1], vertex.Position[2]},
			Color:    [4]float32{vertex.Color[0], vertex.Color[1], vertex.Color[2], vertex.Color[3]},
		}
	}

	return v
}

// MakePolygon creates a new Polygon renderable and adds it to the RenderQueue.
// It returns a reference to the created Polygon.
func MakePolygon(cx, cy int, width float32, c color.Color, sides int) []Vertex {
	r, g, b, a := c.RGBA()
	colorArray := [4]float32{float32(r) / 0xffff, float32(g) / 0xffff, float32(b) / 0xffff, float32(a) / 0xffff}

	// Pre-allocate vertices: each side needs 3 vertices (center + 2 edge points)
	vertices := make([]Vertex, 0, sides*3)

	center := Vertex{
		Position: [3]float32{float32(cx), float32(cy), 0},
		Color:    colorArray,
	}

	// Pre-calculate angle step and radius
	angleStep := 2 * float32(math.Pi) / float32(sides)
	radius := width / 2
	fcx, fcy := float32(cx), float32(cy)

	// Calculate first vertex position
	x := fcx + radius*float32(math.Cos(0))
	y := fcy + radius*float32(math.Sin(0))

	for i := 0; i < sides; i++ {
		// Current vertex (reuse from previous iteration or initial calculation)
		currentVertex := Vertex{
			Position: [3]float32{x, y, 0},
			Color:    colorArray,
		}

		// Calculate next vertex position
		nextAngle := float32(i+1) * angleStep
		nextX := fcx + radius*float32(math.Cos(float64(nextAngle)))
		nextY := fcy + radius*float32(math.Sin(float64(nextAngle)))

		nextVertex := Vertex{
			Position: [3]float32{nextX, nextY, 0},
			Color:    colorArray,
		}

		// Add triangle: center, current, next
		vertices = append(vertices, center, currentVertex, nextVertex)

		// Next vertex becomes current for the next iteration
		x, y = nextX, nextY
	}

	return vertices
}

// MakeLine creates a new Line renderable and adds it to the RenderQueue.
// It returns a reference to the created Line.
func MakeLine(x1, y1, x2, y2 int, width float32, c color.Color) []Vertex {
	r, g, b, a := c.RGBA()
	colorArray := [4]float32{float32(r) / 0xffff, float32(g) / 0xffff, float32(b) / 0xffff, float32(a) / 0xffff}

	dx := float32(x2 - x1)
	dy := float32(y2 - y1)
	len := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	sin := dy / len
	cos := dx / len

	// Calculate the four corners of the line (as a very thin rectangle)
	halfWidth := width / 2
	vertices := []Vertex{
		{Position: [3]float32{float32(x1) - sin*halfWidth, float32(y1) + cos*halfWidth, 0}, Color: colorArray},
		{Position: [3]float32{float32(x2) - sin*halfWidth, float32(y2) + cos*halfWidth, 0}, Color: colorArray},
		{Position: [3]float32{float32(x2) + sin*halfWidth, float32(y2) - cos*halfWidth, 0}, Color: colorArray},
		{Position: [3]float32{float32(x1) + sin*halfWidth, float32(y1) - cos*halfWidth, 0}, Color: colorArray},
	}

	// Creating two triangles to form the line
	lineVertices := []Vertex{vertices[0], vertices[1], vertices[2], vertices[0], vertices[2], vertices[3]}

	return lineVertices
}
