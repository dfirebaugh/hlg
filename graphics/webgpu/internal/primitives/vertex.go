//go:build !js

package primitives

import (
	"image/color"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

// Vertex represents a single vertex in the shape.
type Vertex struct {
	Position  [3]float32 // x, y, z coordinates
	Color     [4]float32 // RGBA color
	TexCoords [2]float32
}

func (v *Vertex) SetColor(c color.Color) {
	r, g, b, a := c.RGBA()
	v.Color = [4]float32{
		float32(r) / 0xffff,
		float32(g) / 0xffff,
		float32(b) / 0xffff,
		float32(a) / 0xffff,
	}
}

// ScreenToNDC transforms screen space coordinates to NDC.
// screenWidth and screenHeight are the dimensions of the screen.
func ScreenToNDC(x, y, screenWidth, screenHeight float32) [3]float32 {
	// Normalize coordinates to [0, 1]
	normalizedX := x / screenWidth
	normalizedY := y / screenHeight

	// Map to NDC [-1, 1]
	ndcX := normalizedX*2 - 1
	ndcY := 1 - normalizedY*2 // Y is inverted in NDC

	return [3]float32{ndcX, ndcY, 0} // Assuming Z coordinate to be 0 for 2D
}

// convertVerticesToNDC converts an array of vertices from screen space to NDC.
func ConvertVerticesToNDC(vertices []Vertex, screenWidth, screenHeight float32) []Vertex {
	ndcVertices := make([]Vertex, len(vertices))
	for i, v := range vertices {
		ndcPosition := ScreenToNDC(v.Position[0], v.Position[1], screenWidth, screenHeight)
		ndcVertices[i] = Vertex{
			Position:  ndcPosition,
			Color:     v.Color,
			TexCoords: v.TexCoords,
		}
	}
	return ndcVertices
}

func CreateVertexBuffer(device *wgpu.Device, vertices []Vertex, width float32, height float32) *wgpu.Buffer {
	ndcVertices := ConvertVerticesToNDC(vertices, width, height)
	vertexBuffer, err := device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Vertex Buffer",
		Contents: wgpu.ToBytes(ndcVertices[:]),
		Usage:    wgpu.BufferUsage_Vertex | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		panic(err)
	}

	return vertexBuffer
}

// UpdateVertexBuffer writes new vertex data to an existing buffer without recreating it
func UpdateVertexBuffer(device *wgpu.Device, buffer *wgpu.Buffer, vertices []Vertex, width float32, height float32) {
	ndcVertices := ConvertVerticesToNDC(vertices, width, height)
	_ = device.GetQueue().WriteBuffer(buffer, 0, wgpu.ToBytes(ndcVertices[:]))
}

func CalculateCenter(v []Vertex) [2]float32 {
	var sumX, sumY float32
	for _, vertex := range v {
		sumX += vertex.Position[0]
		sumY += vertex.Position[1]
	}
	count := float32(len(v))
	return [2]float32{sumX / count, sumY / count}
}
