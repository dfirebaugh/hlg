package hlg

import (
	"image/color"

	"github.com/dfirebaugh/hlg/graphics"
)

type Shape interface {
	graphics.Shape
}

// PolygonFromVertices creates a polygon shape using a specified array of vertices.
// The vertices should be defined with their positions and colors. The function converts
// the input vertices from the local Vertex type to the graphics.Vertex type required by
// the graphics backend. The created polygon is then added to the render queue.
//
// Parameters:
//   - x, y: The x and y coordinates of the polygon's position.
//   - width: The width of the polygon (not used directly in this function, included for interface consistency).
//   - vertices: A slice of Vertex that defines the positions and colors of the polygon's vertices.
//
// Returns:
//   - A graphics.Shape that represents the created polygon.
func PolygonFromVertices(x, y int, width float32, vertices []Vertex) graphics.Shape {
	ensureSetupCompletion()

	graphicsVertices := make([]graphics.Vertex, len(vertices))
	for i, v := range vertices {
		graphicsVertices[i] = graphics.Vertex(v)
	}

	return hlg.graphicsBackend.AddPolygonFromVertices(x, y, width, graphicsVertices)
}

// Polygon creates a polygon shape with a specified number of sides, position, width, and color.
// x, y define the center of the polygon.
// width defines the diameter of the circumcircle of the polygon.
// sides specify the number of sides (vertices) of the polygon.
// c specifies the color of the polygon.
func Polygon(x, y int, width float32, sides int, c color.Color) graphics.Shape {
	ensureSetupCompletion()
	return hlg.graphicsBackend.AddPolygon(x, y, width, c, sides)
}

// Triangle creates a triangle shape with specified vertices and color.
// x1, y1, x2, y2, x3, y3 define the coordinates of the three vertices of the triangle.
// c specifies the color of the triangle.
func Triangle(x1, y1, x2, y2, x3, y3 int, c color.Color) graphics.Shape {
	ensureSetupCompletion()
	return hlg.graphicsBackend.AddTriangle(x1, y1, x2, y2, x3, y3, c)
}

// Rectangle creates a rectangle shape with specified position, dimensions, and color.
// x, y define the top-left corner of the rectangle.
// width, height define the dimensions of the rectangle.
// c specifies the color of the rectangle.
func Rectangle(x, y, width, height int, c color.Color) graphics.Shape {
	ensureSetupCompletion()
	return hlg.graphicsBackend.AddRectangle(x, y, width, height, c)
}

// Circle creates a circle shape with specified center, radius, and color.
// x, y define the center of the circle.
// radius defines the radius of the circle.
// c specifies the color of the circle.
func Circle(x, y int, radius float32, c color.Color) graphics.Shape {
	ensureSetupCompletion()
	return hlg.graphicsBackend.AddCircle(x, y, radius, c, 32)
}

// Line creates a line with specified start and end points, width, and color.
// x1, y1 define the start point of the line.
// x2, y2 define the end point of the line.
// width defines the thickness of the line.
// c specifies the color of the line.
func Line(x1, y1, x2, y2 int, width float32, c color.Color) graphics.Shape {
	ensureSetupCompletion()
	return hlg.graphicsBackend.AddLine(x1, y1, x2, y2, width, c)
}

// PrintAt renders text at a specified position with a specified color.
// s is the string to be rendered.
// x, y define the position where the text will be rendered.
// c specifies the color of the text.
//
// Deprecated: the original implementation of this was inefficient...
// TODO: fix this
func PrintAt(s string, x int, y int, c color.Color) {
}
