package texture

import (
	_ "embed"
	"unsafe"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

//go:embed texture.wgsl
var TextureShaderCode string

type Vertex struct {
	position  [3]float32
	texCoords [2]float32
}

var VertexBufferLayout = wgpu.VertexBufferLayout{
	ArrayStride: uint64(unsafe.Sizeof(Vertex{})),
	StepMode:    wgpu.VertexStepMode_Vertex,
	Attributes: []wgpu.VertexAttribute{
		{
			Offset:         0,
			ShaderLocation: 0,
			Format:         wgpu.VertexFormat_Float32x3,
		},
		{
			Offset:         uint64(unsafe.Sizeof([3]float32{})),
			ShaderLocation: 1,
			Format:         wgpu.VertexFormat_Float32x2,
		},
	},
}

var VERTICES = [...]Vertex{
	{
		position:  [3]float32{-1.0, -1.0, 0.0}, // bottom left
		texCoords: [2]float32{0.0, 1},
	},
	{
		position:  [3]float32{1.0, -1.0, 0.0}, // bottom right
		texCoords: [2]float32{1.0, 1},
	},
	{
		position:  [3]float32{-1.0, 1.0, 0.0}, // top left
		texCoords: [2]float32{0.0, 0},
	},
	{
		position:  [3]float32{1.0, 1.0, 0.0}, // top right
		texCoords: [2]float32{1.0, 0},
	},
}

var INDICES = [...]uint16{
	0, 1, 2, // first triangle
	2, 1, 3, // second triangle
}
