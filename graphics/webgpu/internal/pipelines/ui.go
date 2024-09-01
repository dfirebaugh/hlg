package pipelines

import (
	"log"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/context"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/shader"
	"github.com/dfirebaugh/hlg/gui"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

// func translateVertexBufferLayout(layout graphics.VertexBufferLayout) wgpu.VertexBufferLayout {
// 	var translatedAttributes []wgpu.VertexAttribute
// 	for _, attr := range layout.Attributes {
// 		var format wgpu.VertexFormat
// 		switch attr.Format {
// 		case "float32x3":
// 			format = wgpu.VertexFormat_Float32x3
// 		case "float32x4":
// 			format = wgpu.VertexFormat_Float32x4
// 		default:
// 			log.Fatalf("Unknown vertex format: %s", attr.Format)
// 		}
// 		translatedAttributes = append(translatedAttributes, wgpu.VertexAttribute{
// 			ShaderLocation: attr.ShaderLocation,
// 			Offset:         attr.Offset,
// 			Format:         format,
// 		})
// 	}
//
// 	return wgpu.VertexBufferLayout{
// 		ArrayStride: layout.ArrayStride,
// 		Attributes:  translatedAttributes,
// 	}
// }

type PrimitiveBuffer struct {
	context.RenderContext

	pipeline *wgpu.RenderPipeline

	vertexBuffer *wgpu.Buffer
	vertices     []gui.Vertex

	isDisposed bool
}

func NewPrimitiveBuffer(ctx context.RenderContext, vertices []gui.Vertex, layout graphics.VertexBufferLayout) *PrimitiveBuffer {
	if ctx == nil {
		log.Fatal("RenderContext is nil")
	}
	if ctx.GetPipelineManager() == nil {
		log.Fatal("PipelineManager is not initialized")
	}
	if ctx.GetDevice() == nil {
		log.Fatal("Device is not initialized")
	}
	if ctx.GetSwapChainDescriptor() == nil {
		log.Fatal("SwapChainDescriptor is not initialized")
	}

	p := &PrimitiveBuffer{
		RenderContext: ctx,
		vertices:      vertices,
	}

	p.createVertexBuffer()
	if p.vertexBuffer == nil {
		log.Fatal("Vertex buffer is nil")
	}

	translatedLayout := translateVertexBufferLayout(layout)

	p.pipeline = ctx.GetPipelineManager().GetPipeline("primitve buffer",
		&wgpu.PipelineLayoutDescriptor{
			Label: "Render Pipeline Layout",
		},
		p.GetShader(shader.PrimitiveBufferShader),
		p.GetSwapChainDescriptor(), wgpu.PrimitiveTopology_TriangleList, []wgpu.VertexBufferLayout{translatedLayout})

	if p.pipeline == nil {
		log.Fatal("Pipeline creation failed")
	}

	return p
}

func (p *PrimitiveBuffer) RenderPass(encoder *wgpu.RenderPassEncoder) {
	if encoder == nil || p.isDisposed {
		return
	}

	if p.pipeline == nil {
		log.Println("PrimitiveBuffer RenderPass: pipeline is nil or not set for the current render mode")
		return
	}

	if len(p.vertices) == 0 {
		return
	}

	encoder.SetPipeline(p.pipeline)
	encoder.SetVertexBuffer(0, p.vertexBuffer, 0, wgpu.WholeSize)

	vertexCount := uint32(len(p.vertices))
	encoder.Draw(vertexCount, 1, 0, 0)
}

func (p *PrimitiveBuffer) Render() {
	// if p.isDisposed {
	// 	return
	// }
	// p.AddToRenderQueue(p)
}

func (p *PrimitiveBuffer) Dispose() {
	if p.vertexBuffer != nil {
		p.vertexBuffer.Release()
		p.vertexBuffer = nil
	}
	p.isDisposed = true
}

func (p *PrimitiveBuffer) IsDisposed() bool {
	return p.isDisposed
}

func (p *PrimitiveBuffer) UpdateVertexBuffer(vertices []gui.Vertex) {
	if p.vertexBuffer != nil {
		p.vertexBuffer.Destroy()
	}
	var err error
	p.vertices = vertices
	p.vertexBuffer, err = p.RenderContext.GetDevice().CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Primitive Vertex Buffer",
		Contents: wgpu.ToBytes(vertices[:]),
		Usage:    wgpu.BufferUsage_Vertex,
	})
	if err != nil {
		panic(err)
	}
}

func (p *PrimitiveBuffer) createVertexBuffer() {
	var err error
	p.vertexBuffer, err = p.RenderContext.GetDevice().CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Primitive Vertex Buffer",
		Contents: wgpu.ToBytes(p.vertices[:]),
		Usage:    wgpu.BufferUsage_Vertex,
	})
	if err != nil {
		panic(err)
	}
}
