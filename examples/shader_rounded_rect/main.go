package main

import (
	"unsafe"

	"github.com/dfirebaugh/hlg"
	"golang.org/x/image/colornames"
)

var mousePosition [2]float32

const (
	screenWidth  = 800
	screenHeight = 600
)

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetTitle("SDF Rounded Rectangle Shader")

	shaderCode := `
@group(0) @binding(0) var<uniform> mouse_pos: vec2<f32>;

@vertex
fn vs_main(@location(0) in_pos: vec3<f32>) -> @builtin(position) vec4<f32> {
    return vec4<f32>(in_pos, 1.0);
}

fn sdRoundedRect(p: vec2<f32>, size: vec2<f32>, radius: f32) -> f32 {
    let q = abs(p) - size + vec2<f32>(radius);
    return length(max(q, vec2<f32>(0.0, 0.0))) - radius;
}

@fragment
fn fs_main(@builtin(position) frag_coord: vec4<f32>) -> @location(0) vec4<f32> {
    let rect_size = vec2<f32>(150.0, 100.0); // Size of the rectangle
    let corner_radius: f32 = 20.0; // Radius for rounded corners
    let p = frag_coord.xy - mouse_pos; // Translate the rectangle to follow the mouse
    let sdf = sdRoundedRect(p, rect_size, corner_radius);

    if sdf < 0.0 {
        return vec4<f32>(1.0, 0.0, 0.0, 1.0); // Inside the rounded rectangle, color it red
    } else {
        return vec4<f32>(0.0, 0.0, 0.0, 0.0); // Outside the rectangle, make it transparent
    }
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

	shader := hlg.CompileShader(shaderCode)
	quad := hlg.CreateRenderable(shader, makeFullScreenQuad(screenWidth, screenHeight), uniforms, dataMap)

	if quad == nil {
		panic("Failed to create full-screen quad renderable")
	}

	hlg.Run(func() {
		windowWidth, windowHeight := float32(screenWidth), float32(screenHeight)
		screenWidth, screenHeight := float32(screenWidth), float32(screenHeight)

		x, y := hlg.GetCursorPosition()

		mousePosition[0] = float32(x) * (screenWidth / windowWidth)
		mousePosition[1] = float32(y) * (screenHeight / windowHeight)

		quad.UpdateUniform("mouse_pos", unsafe.Slice((*byte)(unsafe.Pointer(&mousePosition[0])), int(unsafe.Sizeof(mousePosition))))
	}, func() {
		hlg.Clear(colornames.Skyblue)
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

	return hlg.ConvertVerticesToNDC2D(v, screenWidth, screenHeight)
}
