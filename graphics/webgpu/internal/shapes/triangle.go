package shapes

import (
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Triangle struct {
	*wgpu.SwapChainDescriptor
	renderQueue
	vertices     [3]Vertex
	pipeline     map[RenderMode]*wgpu.RenderPipeline
	vertexBuffer *wgpu.Buffer

	shouldBeRendered bool
	renderMode       RenderMode
}

func NewTriangle(device *wgpu.Device, scd *wgpu.SwapChainDescriptor, rq renderQueue, vertices []Vertex, renderMode RenderMode) *Triangle {
	ndcVertices := ConvertVerticesToNDC(vertices, float32(scd.Width), float32(scd.Height))
	vertexBuffer, err := device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Vertex Buffer",
		Contents: wgpu.ToBytes(ndcVertices),
		Usage:    wgpu.BufferUsage_Vertex,
	})
	if err != nil {
		panic(err)
	}

	shaderModule, err := device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{Code: ShapesShaderCode},
	})
	if err != nil {
		panic(err)
	}

	filledPipeline := createPipeline(device, shaderModule, scd, wgpu.PrimitiveTopology_TriangleList)
	outlinedPipeline := createPipeline(device, shaderModule, scd, wgpu.PrimitiveTopology_LineStrip)
	return &Triangle{
		SwapChainDescriptor: scd,
		pipeline: map[RenderMode]*wgpu.RenderPipeline{
			RenderFilled:   filledPipeline,
			RenderOutlined: outlinedPipeline,
		},
		vertexBuffer: vertexBuffer,
		renderMode:   renderMode,
		renderQueue:  rq,
	}
}

func (t *Triangle) SetShouldBeRendered(shouldRender bool) {
	t.shouldBeRendered = shouldRender
}

func (t *Triangle) RenderPass(encoder *wgpu.RenderPassEncoder) {
	if !t.shouldBeRendered {
		return
	}

	currentPipeline := t.pipeline[t.renderMode]
	encoder.SetPipeline(currentPipeline)
	encoder.SetVertexBuffer(0, t.vertexBuffer, 0, wgpu.WholeSize)
	if t.renderMode == RenderFilled {
		encoder.Draw(uint32(len(t.vertices)), 1, 0, 0)
	} else if t.renderMode == RenderOutlined {
		encoder.Draw(uint32(len(t.vertices)+1), 1, 0, 0) // +1 for loop back to the first vertex
	}
}

func (t *Triangle) Render() {
	t.SetShouldBeRendered(true)
	t.renderQueue.AddToRenderQueue(t)
}

// Dispose releases any resources used by the Triangle.
func (t *Triangle) Dispose() {
	if t.pipeline != nil {
		for _, pipeline := range t.pipeline {
			pipeline.Release()
			pipeline = nil
		}
	}
	if t.vertexBuffer != nil {
		t.vertexBuffer.Release()
		t.vertexBuffer = nil
	}
}
