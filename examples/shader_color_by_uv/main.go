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
	hlg.SetTitle("Color by Pixel Position Shader")

	shaderCode := `
@vertex
fn vs_main(@location(0) in_pos: vec3<f32>) -> @builtin(position) vec4<f32> {
    return vec4<f32>(in_pos, 1.0);
}

@fragment
fn fs_main(@builtin(position) frag_coord: vec4<f32>) -> @location(0) vec4<f32> {
    let x_factor = frag_coord.x / 720.0;
    let y_factor = frag_coord.y / 480.0;
    return vec4<f32>(x_factor, y_factor, 0.0, 1.0); // Color based on pixel position
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
