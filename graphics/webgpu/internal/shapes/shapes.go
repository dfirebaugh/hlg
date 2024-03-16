package shapes

import (
	"unsafe"

	_ "embed"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/common"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

//go:embed shapes.wgsl
var ShapesShaderCode string

type renderQueue interface {
	AddToRenderQueue(r graphics.Renderable)
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
				ArrayStride: uint64(unsafe.Sizeof(common.Vertex{})),
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
			Topology:         topology,
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
