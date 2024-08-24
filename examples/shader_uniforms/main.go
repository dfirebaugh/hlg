package main

import (
	"unsafe"

	"github.com/dfirebaugh/hlg"
	"golang.org/x/image/colornames"
)

var (
	colorFactor   [4]float32
	mousePosition [2]float32
)

const (
	screenWidth  = 720
	screenHeight = 480
)

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetTitle("Moving Circle with Uniforms")

	shaderCode := `
@group(0) @binding(0) var<uniform> colorFactor: vec4<f32>;
@group(0) @binding(1) var<uniform> mouse_pos: vec2<f32>;

@vertex
fn vs_main(@location(0) in_pos: vec3<f32>) -> @builtin(position) vec4<f32> {
    return vec4<f32>(in_pos, 1.0);
}

@fragment
fn fs_main(@builtin(position) frag_coord: vec4<f32>) -> @location(0) vec4<f32> {
    let circle_radius: f32 = 50.0;
    let dist_to_center = distance(frag_coord.xy, mouse_pos);
    if dist_to_center < circle_radius {
        return colorFactor; // Inside the circle, apply the color factor
    } else {
        return vec4<f32>(0.0, 0.0, 0.0, 0.0); // Outside the circle, make it transparent
    }
}
	`

	colorFactor = [4]float32{1.0, 0.0, 0.0, 1.0}

	uniforms := map[string]hlg.Uniform{
		"colorFactor": {
			Binding: 0,
			Size:    uint64(unsafe.Sizeof(colorFactor)),
		},
		"mouse_pos": {
			Binding: 1,
			Size:    uint64(unsafe.Sizeof(mousePosition)),
		},
	}

	dataMap := map[string][]byte{
		"colorFactor":  unsafe.Slice((*byte)(unsafe.Pointer(&colorFactor[0])), int(unsafe.Sizeof(colorFactor))),
		"mouse_pos":    unsafe.Slice((*byte)(unsafe.Pointer(&mousePosition[0])), int(unsafe.Sizeof(mousePosition))),
	}

	quad := hlg.CreateRenderable(shaderCode, makeFullScreenQuad(screenWidth, screenHeight), uniforms, dataMap)

	if quad == nil {
		panic("Failed to create quad renderable")
	}

	hlg.Run(func() {
		x, y := hlg.GetCursorPosition()

		colorFactor[0] = float32(x) / screenWidth
		colorFactor[2] = float32(y) / screenHeight

		mousePosition[0] = float32(x)
		mousePosition[1] = float32(y)

		quad.UpdateUniform("colorFactor", unsafe.Slice((*byte)(unsafe.Pointer(&colorFactor[0])), int(unsafe.Sizeof(colorFactor))))
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
