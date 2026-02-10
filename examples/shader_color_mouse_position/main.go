package main

import (
	"unsafe"

	"github.com/dfirebaugh/hlg"
)

var mousePosition [2]float32

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

// getShaderCode returns the appropriate shader code based on the backend.
// For OpenGL/WebGL: GLSL with #vertex and #fragment markers
// For WebGPU: WGSL
func getShaderCode() string {
	// OpenGL and WebGL both use GLSL; only WebGPU uses WGSL
	if hlg.GetBackend() != hlg.BackendWebGPU {
		// GLSL shader for OpenGL backend
		return `
#vertex
#version 410 core

layout(location = 0) in vec3 in_pos;

void main() {
    gl_Position = vec4(in_pos, 1.0);
}
#fragment
#version 410 core

uniform vec2 mouse_pos;
out vec4 fragColor;

void main() {
    float x_factor = mouse_pos.x / 720.0;
    float y_factor = mouse_pos.y / 480.0;
    fragColor = vec4(x_factor, y_factor, 0.0, 1.0);
}
`
	}

	// WGSL shader for WebGPU backend
	return `
@group(0) @binding(0) var<uniform> mouse_pos: vec2<f32>;

@vertex
fn vs_main(@location(0) in_pos: vec3<f32>) -> @builtin(position) vec4<f32> {
    return vec4<f32>(in_pos, 1.0);
}

@fragment
fn fs_main() -> @location(0) vec4<f32> {
    let x_factor = mouse_pos.x / 720.0;
    let y_factor = mouse_pos.y / 480.0;
    return vec4<f32>(x_factor, y_factor, 0.0, 1.0);
}
`
}

func main() {
	// Use OpenGL backend (default). Change to BackendWebGPU to test WGSL shaders.
	// hlg.SetBackend(hlg.BackendWebGPU)

	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetTitle("Color by Mouse Position Shader")
	hlg.SetVSync(true)

	shaderCode := getShaderCode()

	uniforms := map[string]hlg.Uniform{
		"mouse_pos": {
			Binding: 0,
			Size:    uint64(unsafe.Sizeof(mousePosition)),
		},
	}

	dataMap := map[string][]byte{
		"mouse_pos": unsafe.Slice((*byte)(unsafe.Pointer(&mousePosition[0])), int(unsafe.Sizeof(mousePosition))),
	}

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

	quad := hlg.CreateRenderable(shader, quadVertexData, vertexLayout, uniforms, dataMap)

	if quad == nil {
		panic("Failed to create full-screen quad renderable")
	}

	hlg.Run(func() {
		x, y := hlg.GetCursorPosition()
		mousePosition[0] = float32(x)
		mousePosition[1] = float32(y)
		quad.UpdateUniform("mouse_pos", unsafe.Slice((*byte)(unsafe.Pointer(&mousePosition[0])), int(unsafe.Sizeof(mousePosition))))
	}, func() {
		quad.Render()
	})
}

func makeFullScreenQuad() []Vertex {
	v := []Vertex{
		// Bottom-left corner
		{Position: [3]float32{-1.0, -1.0, 0.0}},
		// Bottom-right corner
		{Position: [3]float32{1.0, -1.0, 0.0}},
		// Top-left corner
		{Position: [3]float32{-1.0, 1.0, 0.0}},

		// Top-left corner
		{Position: [3]float32{-1.0, 1.0, 0.0}},
		// Bottom-right corner
		{Position: [3]float32{1.0, -1.0, 0.0}},
		// Top-right corner
		{Position: [3]float32{1.0, 1.0, 0.0}},
	}

	return v
}
