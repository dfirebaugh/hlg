package gui

import "image/color"

type OpCode float32

const (
	DrawTriangle OpCode = iota
)

type Vertex struct {
	Position  [3]float32
	FillColor [4]float32
}

type VertexAttribute struct {
	ShaderLocation uint32
	Offset         uint64
	Format         string
}

type VertexBufferLayout struct {
	ArrayStride uint64
	Attributes  []VertexAttribute
}

func colorToFloat32(c color.Color) [4]float32 {
	r, g, b, a := c.RGBA()
	return [4]float32{
		float32(r) / 0xffff,
		float32(g) / 0xffff,
		float32(b) / 0xffff,
		float32(a) / 0xffff,
	}
}

// screenToNDC transforms screen space coordinates to NDC.
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
