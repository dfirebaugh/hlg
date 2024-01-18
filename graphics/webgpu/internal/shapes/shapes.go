package shapes

import (
	"unsafe"

	_ "embed"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

//go:embed shapes.wgsl
var ShapesShaderCode string

type renderQueue interface {
	AddToRenderQueue(r graphics.Renderable)
}

type RenderMode int

const (
	RenderFilled RenderMode = iota
	RenderOutlined
)

// Vertex represents a single vertex in the shape.
type Vertex struct {
	Position [3]float32 // x, y, z coordinates
	Color    [4]float32 // RGBA color
}

// ScreenToNDC transforms screen space coordinates to NDC.
// screenWidth and screenHeight are the dimensions of the screen.
func ScreenToNDC(x, y, screenWidth, screenHeight float32) [3]float32 {
	// Normalize coordinates to [0, 1]
	normalizedX := x / screenWidth
	normalizedY := y / screenHeight

	// Map to NDC [-1, 1]
	ndcX := normalizedX*2 - 1
	ndcY := 1 - normalizedY*2 // Y is inverted in NDC

	return [3]float32{ndcX, ndcY, 0} // Assuming Z coordinate to be 0 for 2D
}

// convertVerticesToNDC converts an array of vertices from screen space to NDC.
func convertVerticesToNDC(vertices []Vertex, screenWidth, screenHeight float32) []Vertex {
	ndcVertices := make([]Vertex, len(vertices))
	for i, v := range vertices {
		ndcPosition := ScreenToNDC(v.Position[0], v.Position[1], screenWidth, screenHeight)
		ndcVertices[i] = Vertex{
			Position: ndcPosition,
			Color:    v.Color,
		}
	}
	return ndcVertices
}

func createVertexBuffer(device *wgpu.Device, vertices []Vertex, width float32, height float32) *wgpu.Buffer {
	ndcVertices := convertVerticesToNDC(vertices, width, height)
	vertexBuffer, err := device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Vertex Buffer",
		Contents: wgpu.ToBytes(ndcVertices[:]),
		Usage:    wgpu.BufferUsage_Vertex,
	})
	if err != nil {
		panic(err)
	}

	return vertexBuffer
}

func (p *Polygon) createPipeline(device *wgpu.Device, shaderModule *wgpu.ShaderModule, scd *wgpu.SwapChainDescriptor, topology wgpu.PrimitiveTopology) *wgpu.RenderPipeline {
	renderPipelineLayout, err := device.CreatePipelineLayout(&wgpu.PipelineLayoutDescriptor{
		Label: "Render Pipeline Layout",
		BindGroupLayouts: []*wgpu.BindGroupLayout{
			p.bindGroupLayout,
		},
	})
	if err != nil {
		panic(err)
	}
	defer renderPipelineLayout.Release()
	pipeline, err := device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Layout: renderPipelineLayout,
		Vertex: wgpu.VertexState{
			Module:     shaderModule,
			EntryPoint: "vs_main",
			Buffers: []wgpu.VertexBufferLayout{{
				ArrayStride: uint64(unsafe.Sizeof(Vertex{})),
				Attributes: []wgpu.VertexAttribute{
					{
						ShaderLocation: 0,
						Offset:         0,
						Format:         wgpu.VertexFormat_Float32x3, // Position
					},
					{
						ShaderLocation: 1,
						Offset:         uint64(unsafe.Sizeof([3]float32{})),
						Format:         wgpu.VertexFormat_Float32x4, // Color
					},
				},
			}},
		},
		Primitive: wgpu.PrimitiveState{
			Topology:         wgpu.PrimitiveTopology_TriangleList,
			StripIndexFormat: wgpu.IndexFormat_Undefined,
			FrontFace:        wgpu.FrontFace_CCW,
			CullMode:         wgpu.CullMode_None,
		},
		Multisample: wgpu.MultisampleState{
			Count:                  1,
			Mask:                   0xFFFFFFFF,
			AlphaToCoverageEnabled: false,
		},
		Fragment: &wgpu.FragmentState{
			Module:     shaderModule,
			EntryPoint: "fs_main",
			Targets: []wgpu.ColorTargetState{{
				Format:    scd.Format,
				Blend:     &wgpu.BlendState_Replace,
				WriteMask: wgpu.ColorWriteMask_All,
			}},
		},
	})
	if err != nil {
		panic(err)
	}

	return pipeline
}
