package main

import (
	"github.com/dfirebaugh/hlg"
)

const (
	screenWidth  = 720
	screenHeight = 480
)

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
	quad := hlg.CreateRenderable(shader, makeFullScreenQuad(screenWidth, screenHeight), nil, nil)

	if quad == nil {
		panic("Failed to create full-screen quad renderable")
	}

	hlg.Run(func() {
	}, func() {
		quad.Render()
	})
}

func makeFullScreenQuad(screenWidth, screenHeight float32) []hlg.Vertex {
	v := []hlg.Vertex{
		{Position: [3]float32{0, 0, 0}},            // Bottom-left
		{Position: [3]float32{screenWidth, 0, 0}},  // Bottom-right
		{Position: [3]float32{0, screenHeight, 0}}, // Top-left

		{Position: [3]float32{0, screenHeight, 0}},           // Top-left
		{Position: [3]float32{screenWidth, 0, 0}},            // Bottom-right
		{Position: [3]float32{screenWidth, screenHeight, 0}}, // Top-right
	}

	return hlg.ConvertVerticesToNDC2D(v, screenWidth, screenHeight)
}
