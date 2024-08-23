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

func main() {
	hlg.SetWindowSize(720, 480)
	hlg.SetTitle("Color by Mouse Position Shader")

	shaderCode := `
@group(0) @binding(0) var<uniform> mouse_pos: vec2<f32>;

@vertex
fn vs_main(@location(0) in_pos: vec3<f32>) -> @builtin(position) vec4<f32> {
    return vec4<f32>(in_pos, 1.0);
}

@fragment
fn fs_main() -> @location(0) vec4<f32> {
    let x_factor = mouse_pos.x / 720.0;
    let y_factor = mouse_pos.y / 480.0;
    return vec4<f32>(x_factor, y_factor, 0.0, 1.0); // Color based on mouse position
}
	`

	uniforms := map[string]hlg.Uniform{
		"mouse_pos": {
			Binding: 0,
			Size:    uint64(unsafe.Sizeof(mousePosition)),
		},
	}

	dataMap := map[string][]byte{
		"mouse_pos": unsafe.Slice((*byte)(unsafe.Pointer(&mousePosition[0])), int(unsafe.Sizeof(mousePosition))),
	}

	quad := hlg.CreateRenderable(shaderCode, makeFullScreenQuad(screenWidth, screenHeight), uniforms, dataMap)

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

func makeFullScreenQuad(screenWidth, screenHeight float32) []hlg.Vertex {
	// Define a full-screen quad using screen space coordinates
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

	return hlg.ConvertVerticesToNDC2D(v, screenWidth, screenHeight)
}
