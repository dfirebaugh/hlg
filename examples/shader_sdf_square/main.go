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

type Vertex struct {
	Position [3]float32
}

func VerticesToBytes(vertices []Vertex) []byte {
	size := len(vertices) * int(unsafe.Sizeof(Vertex{}))
	data := make([]byte, size)
	copy(data, unsafe.Slice((*byte)(unsafe.Pointer(&vertices[0])), size))
	return data
}

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetTitle("SDF Square Shader")

	shaderCode := `
@group(0) @binding(0) var<uniform> mouse_pos: vec2<f32>;

@vertex
fn vs_main(@location(0) in_pos: vec3<f32>) -> @builtin(position) vec4<f32> {
    return vec4<f32>(in_pos, 1.0);
}

@fragment
fn fs_main(@builtin(position) frag_coord: vec4<f32>) -> @location(0) vec4<f32> {
    let square_size: vec2<f32> = vec2<f32>(50.0, 50.0);
    let dist_to_edge = abs(frag_coord.xy - mouse_pos) - square_size;
    let sdf = length(max(dist_to_edge, vec2<f32>(0.0, 0.0))) + min(max(dist_to_edge.x, dist_to_edge.y), 0.0);

    if sdf < 0.0 {
        return vec4<f32>(1.0, 0.0, 0.0, 1.0); // Inside the square, color it red
    } else {
        return vec4<f32>(0.0, 0.0, 0.0, 0.0); // Outside the square, make it transparent
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
