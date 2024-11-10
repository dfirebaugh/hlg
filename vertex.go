package hlg

import (
	"image/color"

	"github.com/dfirebaugh/hlg/graphics"
)

type Vertex struct {
	Position [3]float32
	Color    color.Color
}

// ScreenToNDC transforms screen space coordinates to NDC.
// screenWidth and screenHeight are the dimensions of the screen.
func screenToNDC(x, y, screenWidth, screenHeight float32) [3]float32 {
	// Normalize coordinates to [0, 1]
	normalizedX := x / screenWidth
	normalizedY := y / screenHeight

	// Map to NDC [-1, 1]
	ndcX := normalizedX*2 - 1
	ndcY := 1 - normalizedY*2 // Y is inverted in NDC

	return [3]float32{ndcX, ndcY, 0} // Assuming Z coordinate to be 0 for 2D
}

func ConvertVerticesToNDC2D(vertices []Vertex, screenWidth, screenHeight float32) []Vertex {
	ndcVertices := make([]Vertex, len(vertices))
	for i, v := range vertices {
		ndcPosition := screenToNDC(v.Position[0], v.Position[1], screenWidth, screenHeight)
		ndcVertices[i] = Vertex{
			Position: ndcPosition,
			Color:    v.Color,
		}
	}
	return ndcVertices
}

func convertToGraphicsVertex(hv Vertex) graphics.Vertex {
	gv := graphics.Vertex{
		Position: hv.Position,
		Color:    toRGBA(hv.Color),
	}
	return gv
}

func convertVerticesToGraphics(vertices []Vertex) []graphics.Vertex {
	gvs := make([]graphics.Vertex, len(vertices))
	for i, v := range vertices {
		gvs[i] = convertToGraphicsVertex(v)
	}
	return gvs
}

func toRGBA(c color.Color) [4]float32 {
	r, g, b, a := c.RGBA()
	return [4]float32{
		float32(r) / 0xffff,
		float32(g) / 0xffff,
		float32(b) / 0xffff,
		float32(a) / 0xffff,
	}
}
