package pipeline

import (
	"sync"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type PipelineManager struct {
	pipelines map[string]*wgpu.RenderPipeline
	mu        sync.Mutex
	device    *wgpu.Device
}

func NewPipelineManager(device *wgpu.Device) *PipelineManager {
	return &PipelineManager{
		pipelines: make(map[string]*wgpu.RenderPipeline),
		device:    device,
	}
}

func (pm *PipelineManager) GetPipeline(
	key string,
	layout *wgpu.PipelineLayoutDescriptor,
	shaderModule *wgpu.ShaderModule,
	scd *wgpu.SwapChainDescriptor,
	topology wgpu.PrimitiveTopology,
	vertexBufferLayout []wgpu.VertexBufferLayout,
) *wgpu.RenderPipeline {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if pipeline, exists := pm.pipelines[key]; exists {
		return pipeline
	}

	pipeline := createPipeline(pm.device, layout, shaderModule, scd, topology, vertexBufferLayout)
	pm.pipelines[key] = pipeline

	return pipeline
}

func createPipeline(
	device *wgpu.Device,
	layout *wgpu.PipelineLayoutDescriptor,
	shaderModule *wgpu.ShaderModule,
	scd *wgpu.SwapChainDescriptor,
	topology wgpu.PrimitiveTopology,
	vertexBufferLayout []wgpu.VertexBufferLayout,
) *wgpu.RenderPipeline {
	renderPipelineLayout, err := device.CreatePipelineLayout(layout)
	if err != nil {
		panic(err)
	}

	pipeline, err := device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Layout: renderPipelineLayout,
		Vertex: wgpu.VertexState{
			Module:     shaderModule,
			EntryPoint: "vs_main",
			Buffers:    vertexBufferLayout,
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
				Format: scd.Format,
				Blend: &wgpu.BlendState{
					Color: wgpu.BlendComponent{
						SrcFactor: wgpu.BlendFactor_SrcAlpha,
						DstFactor: wgpu.BlendFactor_OneMinusSrcAlpha,
						Operation: wgpu.BlendOperation_Add,
					},
					// Alpha: wgpu.BlendComponent{
					// 	SrcFactor: wgpu.BlendFactor_One,
					// 	DstFactor: wgpu.BlendFactor_Zero,
					// 	Operation: wgpu.BlendOperation_Add,
					// },
				},
				WriteMask: wgpu.ColorWriteMask_All,
			}},
		},
	})
	if err != nil {
		panic(err)
	}

	return pipeline
}
