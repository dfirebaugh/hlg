package shapes

import (
	"image/color"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Polygon struct {
	device *wgpu.Device
	*wgpu.SwapChainDescriptor
	renderQueue
	vertices     []Vertex
	pipeline     map[RenderMode]*wgpu.RenderPipeline
	vertexBuffer *wgpu.Buffer

	shouldeRender bool
	renderMode    RenderMode
}

func NewPolygon(device *wgpu.Device, scd *wgpu.SwapChainDescriptor, rq renderQueue, vertices []Vertex, renderMode RenderMode) *Polygon {
	p := &Polygon{
		device:              device,
		SwapChainDescriptor: scd,
		vertices:            vertices,
		renderMode:          renderMode,
		renderQueue:         rq,
	}

	p.vertexBuffer = createVertexBuffer(p.device, p.vertices, float32(scd.Width), float32(scd.Height))

	shaderModule, err := device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{Code: ShapesShaderCode},
	})
	if err != nil {
		panic(err)
	}
	filledPipeline := createPipeline(device, shaderModule, scd, wgpu.PrimitiveTopology_TriangleList)
	outlinedPipeline := createPipeline(device, shaderModule, scd, wgpu.PrimitiveTopology_LineStrip)
	p.pipeline = map[RenderMode]*wgpu.RenderPipeline{
		RenderFilled:   filledPipeline,
		RenderOutlined: outlinedPipeline,
	}

	return p
}

func (p *Polygon) SetColor(c color.Color) {
	r, g, b, a := c.RGBA()
	for i := range p.vertices {
		p.vertices[i].Color = [4]float32{
			float32(r) / 0xffff, // RGBA values are 0-65535 so divide by 0xffff
			float32(g) / 0xffff,
			float32(b) / 0xffff,
			float32(a) / 0xffff,
		}
	}

	p.vertexBuffer = createVertexBuffer(p.device, p.vertices, float32(p.SwapChainDescriptor.Width), float32(p.SwapChainDescriptor.Height))
}

func (p *Polygon) RenderPass(encoder *wgpu.RenderPassEncoder) {
	if !p.shouldeRender {
		return
	}

	currentPipeline := p.pipeline[p.renderMode]
	encoder.SetPipeline(currentPipeline)
	encoder.SetVertexBuffer(0, p.vertexBuffer, 0, wgpu.WholeSize)
	if p.renderMode == RenderFilled {
		encoder.Draw(uint32(len(p.vertices)), 1, 0, 0)
	} else if p.renderMode == RenderOutlined {
		encoder.Draw(uint32(len(p.vertices)+1), 1, 0, 0) // +1 for loop back to the first vertex
	}
}

func (p *Polygon) Render() {
	p.shouldeRender = true
	p.renderQueue.AddToRenderQueue(p)
}

// Dispose releases any resources used by the Triangle.
func (p *Polygon) Dispose() {
	if p.pipeline != nil {
		for _, pipeline := range p.pipeline {
			pipeline.Release()
			pipeline = nil
		}
	}
	if p.vertexBuffer != nil {
		p.vertexBuffer.Release()
		p.vertexBuffer = nil
	}
}
func (p *Polygon) Hide() {
	p.shouldeRender = false
}
