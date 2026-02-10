package main

import (
	"unsafe"

	"github.com/dfirebaugh/hlg"
	"golang.org/x/image/colornames"
)

// mousePosition stores the mouse position scaled to framebuffer coordinates
var mousePosition [2]float32

// screenParams stores framebuffer width and height for the shader
var screenParams [2]float32

const (
	screenWidth  = 800
	screenHeight = 600
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

uniform vec2 mouse_pos;
uniform vec2 screen_size;
out vec4 fragColor;

float sdRoundedRect(vec2 p, vec2 size, float radius) {
    vec2 q = abs(p) - size + vec2(radius);
    return length(max(q, vec2(0.0, 0.0))) - radius;
}

void main() {
    vec2 rect_size = vec2(150.0, 100.0);
    float corner_radius = 20.0;
    // Flip Y: OpenGL gl_FragCoord has origin at bottom-left
    // mouse_pos is already scaled to framebuffer coordinates from Go code
    vec2 frag_flipped = vec2(gl_FragCoord.x, screen_size.y - gl_FragCoord.y);
    vec2 p = frag_flipped - mouse_pos;
    float sdf = sdRoundedRect(p, rect_size, corner_radius);

    if (sdf < 0.0) {
        fragColor = vec4(1.0, 0.0, 0.0, 1.0);
    } else {
        fragColor = vec4(0.0, 0.0, 0.0, 0.0);
    }
}
`
	}

	return `
@group(0) @binding(0) var<uniform> mouse_pos: vec2<f32>;
@group(0) @binding(1) var<uniform> screen_size: vec2<f32>;

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
    let rect_size = vec2<f32>(150.0, 100.0);
    let corner_radius: f32 = 20.0;
    let p = frag_coord.xy - mouse_pos;
    let sdf = sdRoundedRect(p, rect_size, corner_radius);

    if sdf < 0.0 {
        return vec4<f32>(1.0, 0.0, 0.0, 1.0);
    } else {
        return vec4<f32>(0.0, 0.0, 0.0, 0.0);
    }
}
`
}

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetTitle("SDF Rounded Rectangle Shader")
	hlg.SetVSync(true)

	shaderCode := getShaderCode()

	uniforms := map[string]hlg.Uniform{
		"mouse_pos": {
			Binding: 0,
			Size:    uint64(unsafe.Sizeof(mousePosition)),
		},
		"screen_size": {
			Binding: 1,
			Size:    uint64(unsafe.Sizeof(screenParams)),
		},
	}

	dataMap := map[string][]byte{
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
		panic("Failed to create full-screen quad renderable")
	}

	hlg.Run(func() {
		x, y := hlg.GetCursorPosition()

		// Scale mouse coordinates to framebuffer space for HiDPI displays
		scale := hlg.GetPixelScale()
		mousePosition[0] = float32(x) * scale
		mousePosition[1] = float32(y) * scale

		// Update screen params with actual framebuffer size
		fbW, fbH := hlg.GetFramebufferSize()
		screenParams[0] = float32(fbW)
		screenParams[1] = float32(fbH)

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
