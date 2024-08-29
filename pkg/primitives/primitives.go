package primitives

import (
	"image/color"
	"math"

	"github.com/dfirebaugh/hlg"
)

type renderable interface {
	Render()
}

type Position struct {
	X, Y float32
}

type Size struct {
	Width, Height float32
}

func (p Position) ToFloatSlice() [2]float32 {
	return [2]float32{p.X, p.Y}
}

func (s Size) ToFloatSlice() [2]float32 {
	return [2]float32{s.Width, s.Height}
}

func colorToFloatSlice(c color.Color) [4]float32 {
	if c == nil {
		return [4]float32{1.0, 1.0, 1.0, 1.0}
	}

	r, g, b, a := c.RGBA()
	return [4]float32{
		float32(r) / float32(math.MaxUint16),
		float32(g) / float32(math.MaxUint16),
		float32(b) / float32(math.MaxUint16),
		float32(a) / float32(math.MaxUint16),
	}
}

func makeFullScreenQuad(screenWidth, screenHeight float32) []hlg.Vertex {
	v := []hlg.Vertex{
		// Bottom-left corner
		{Position: [3]float32{0, 0, 0}},
		// Bottom-right corner
		{Position: [3]float32{screenWidth, 0, 0}},
		// Top-left corner
		{Position: [3]float32{0, screenHeight, 0}},

		// Top-left corner
		{Position: [3]float32{0, screenHeight, 0}},
		// Bottom-right corner
		{Position: [3]float32{screenWidth, 0, 0}},
		// Top-right corner
		{Position: [3]float32{screenWidth, screenHeight, 0}},
	}

	return v
}
