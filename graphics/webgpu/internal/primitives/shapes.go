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
	var vertices []Vertex
	center := Vertex{
		Position: [3]float32{float32(cx), float32(cy), 0},
		Color:    colorArray,
	}

	for i := 0; i <= sides; i++ {
		angle := float32(i) * 2 * float32(math.Pi) / float32(sides)
		x := float32(cx) + (width/2)*float32(math.Cos(float64(angle)))
		y := float32(cy) + (width/2)*float32(math.Sin(float64(angle)))

		vertex := Vertex{
			Position: [3]float32{x, y, 0},
			Color:    colorArray,
		}

		vertices = append(vertices, center, vertex)

		if i < sides {
			nextAngle := float32(i+1) * 2 * float32(math.Pi) / float32(sides)
			nextX := float32(cx) + (width/2)*float32(math.Cos(float64(nextAngle)))
			nextY := float32(cy) + (width/2)*float32(math.Sin(float64(nextAngle)))

			nextVertex := Vertex{
				Position: [3]float32{nextX, nextY, 0},
				Color:    colorArray,
			}

			vertices = append(vertices, nextVertex)
		}
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
