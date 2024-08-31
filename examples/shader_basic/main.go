package main

import (
	"unsafe"

	"github.com/dfirebaugh/hlg"
)

const (
	screenWidth  = 720
	screenHeight = 480
)

type Vertex struct {
	Position [3]float32
}

func (v Vertex) ToBytes() []byte {
	return *(*[]byte)(unsafe.Pointer(&v))
}

func VerticesToBytes(vertices []Vertex) []byte {
	size := len(vertices) * int(unsafe.Sizeof(Vertex{}))
	data := make([]byte, size)
	copy(data, unsafe.Slice((*byte)(unsafe.Pointer(&vertices[0])), size))
	return data
}

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetTitle("Full Screen Red Shader")

	shaderCode := `
@vertex
fn vs_main(@location(0) in_pos: vec3<f32>) -> @builtin(position) vec4<f32> {
    return vec4<f32>(in_pos, 1.0);
}

@fragment
fn fs_main() -> @location(0) vec4<f32> {
    return vec4<f32>(1.0, 0.0, 0.0, 1.0); // Red color
}
	`

	shader := hlg.CompileShader(shaderCode)

	vertexLayout := hlg.VertexBufferLayout{
		ArrayStride: 3 * 4,
		Attributes: []hlg.VertexAttributeLayout{
			{
				ShaderLocation: 0,
				Offset:         0,
				Format:         "float32x3",
			},
		},
	}

	quadVertices := makeFullScreenQuad()
	quadVertexData := VerticesToBytes(quadVertices)

	quad := hlg.CreateRenderable(shader, quadVertexData, vertexLayout, nil, nil)

	if quad == nil {
		panic("Failed to create full-screen quad renderable")
	}

	hlg.Run(func() {
	}, func() {
		quad.Render()
	})
}

func makeFullScreenQuad() []Vertex {
	v := []Vertex{
		{Position: [3]float32{-1.0, -1.0, 0.0}}, // Bottom-left
		{Position: [3]float32{1.0, -1.0, 0.0}},  // Bottom-right
		{Position: [3]float32{-1.0, 1.0, 0.0}},  // Top-left

		{Position: [3]float32{-1.0, 1.0, 0.0}}, // Top-left
		{Position: [3]float32{1.0, -1.0, 0.0}}, // Bottom-right
		{Position: [3]float32{1.0, 1.0, 0.0}},  // Top-right
	}

	return v
}
