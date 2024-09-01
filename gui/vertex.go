package gui

import (
	"image/color"
)

type OpCode float32

const (
	DrawTriangle OpCode = iota
)

type Vertex struct {
	Position      [3]float32
	LocalPosition [2]float32
	OpCode        float32
	Radius        float32
	Color         [4]float32
}

type VertexBufferLayout struct {
	ArrayStride uint64
	Attributes  []VertexAttributeLayout
}

type VertexAttributeLayout struct {
	ShaderLocation uint32
	Offset         uint64
	Format         string
}

	var (
vertexLayout = VertexBufferLayout{
		ArrayStride: 4*3 + 4*2 + 4 + 4 + 4*4,
		Attributes: []VertexAttributeLayout{
			{
				ShaderLocation: 0,
				Offset:         0,
				Format:         "float32x3", // Position (Clip space)
			},
			{
				ShaderLocation: 1,
				Offset:         3 * 4,
				Format:         "float32x2", // Local Position
			},
			{
				ShaderLocation: 2,
				Offset:         (3 + 2) * 4,
				Format:         "float32", // OpCode
			},
			{
				ShaderLocation: 3,
				Offset:         (3 + 2 + 1) * 4,
				Format:         "float32", // Radius
			},
			{
				ShaderLocation: 4,
				Offset:         (3 + 2 + 1 + 1) * 4,
				Format:         "float32x4", // Color
			},
		},
	}
)

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
