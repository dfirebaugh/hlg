package main

import (
	"unsafe"

	"github.com/dfirebaugh/hlg"
	"golang.org/x/image/colornames"
)

var (
	colorFactor   [4]float32
	mousePosition [2]float32
	screenParams  [2]float32
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

// getShaderCode returns the appropriate shader code based on the backend.
func getShaderCode() string {
	// OpenGL and WebGL both use GLSL; only WebGPU uses WGSL
	if hlg.GetBackend() != hlg.BackendWebGPU {
		return `
#vertex
#version 410 core

layout(location = 0) in vec3 in_pos;

void main() {
    gl_Position = vec4(in_pos, 1.0);
}
#fragment
#version 410 core

uniform vec4 colorFactor;
uniform vec2 mouse_pos;
uniform vec2 screen_size;
out vec4 fragColor;

void main() {
    float circle_radius = 50.0;
    // Flip Y: OpenGL gl_FragCoord has origin at bottom-left
    // mouse_pos is already scaled to framebuffer coordinates from Go code
    vec2 frag_flipped = vec2(gl_FragCoord.x, screen_size.y - gl_FragCoord.y);
    float dist_to_center = distance(frag_flipped, mouse_pos);
    if (dist_to_center < circle_radius) {
        fragColor = colorFactor;
    } else {
        fragColor = vec4(0.0, 0.0, 0.0, 0.0);
    }
}
`
	}

	return `
@group(0) @binding(0) var<uniform> colorFactor: vec4<f32>;
@group(0) @binding(1) var<uniform> mouse_pos: vec2<f32>;
@group(0) @binding(2) var<uniform> screen_size: vec2<f32>;

@vertex
fn vs_main(@location(0) in_pos: vec3<f32>) -> @builtin(position) vec4<f32> {
    return vec4<f32>(in_pos, 1.0);
}

@fragment
fn fs_main(@builtin(position) frag_coord: vec4<f32>) -> @location(0) vec4<f32> {
    let circle_radius: f32 = 50.0;
    let dist_to_center = distance(frag_coord.xy, mouse_pos);
    if dist_to_center < circle_radius {
        return colorFactor;
    } else {
        return vec4<f32>(0.0, 0.0, 0.0, 0.0);
    }
}
`
}

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetTitle("Moving Circle with Uniforms")
	hlg.SetVSync(true)

	shaderCode := getShaderCode()

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
		"screen_size": {
			Binding: 2,
			Size:    uint64(unsafe.Sizeof(screenParams)),
		},
	}

	dataMap := map[string][]byte{
		"colorFactor": unsafe.Slice((*byte)(unsafe.Pointer(&colorFactor[0])), int(unsafe.Sizeof(colorFactor))),
		"mouse_pos":   unsafe.Slice((*byte)(unsafe.Pointer(&mousePosition[0])), int(unsafe.Sizeof(mousePosition))),
		"screen_size": unsafe.Slice((*byte)(unsafe.Pointer(&screenParams[0])), int(unsafe.Sizeof(screenParams))),
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
		panic("Failed to create quad renderable")
	}

	hlg.Run(func() {
		x, y := hlg.GetCursorPosition()

		colorFactor[0] = float32(x) / screenWidth
		colorFactor[2] = float32(y) / screenHeight

		// Scale mouse coordinates to framebuffer space for HiDPI displays
		scale := hlg.GetPixelScale()
		mousePosition[0] = float32(x) * scale
		mousePosition[1] = float32(y) * scale

		// Update screen params with actual framebuffer size
		fbW, fbH := hlg.GetFramebufferSize()
		screenParams[0] = float32(fbW)
		screenParams[1] = float32(fbH)

		quad.UpdateUniform("colorFactor", unsafe.Slice((*byte)(unsafe.Pointer(&colorFactor[0])), int(unsafe.Sizeof(colorFactor))))
		quad.UpdateUniform("mouse_pos", unsafe.Slice((*byte)(unsafe.Pointer(&mousePosition[0])), int(unsafe.Sizeof(mousePosition))))
		quad.UpdateUniform("screen_size", unsafe.Slice((*byte)(unsafe.Pointer(&screenParams[0])), int(unsafe.Sizeof(screenParams))))
	}, func() {
		hlg.Clear(colornames.Skyblue)
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
