package shapes

import (
	"image/color"
	"log"
	"unsafe"

	"github.com/dfirebaugh/hlg/pkg/math/matrix"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Polygon struct {
	device *wgpu.Device
	renderQueue
	*wgpu.SwapChainDescriptor
	bindGroup       *wgpu.BindGroup
	bindGroupLayout *wgpu.BindGroupLayout
	vertices        []Vertex
	center          [2]float32
	pipeline        map[RenderMode]*wgpu.RenderPipeline
	vertexBuffer    *wgpu.Buffer

	shouldeRender bool
	renderMode    RenderMode

	transform       matrix.Matrix
	transformBuffer *wgpu.Buffer

	isDisposed bool
}

func NewPolygon(device *wgpu.Device, scd *wgpu.SwapChainDescriptor, rq renderQueue, vertices []Vertex, renderMode RenderMode) *Polygon {
	p := &Polygon{
		device:              device,
		SwapChainDescriptor: scd,
		vertices:            vertices,
		renderMode:          renderMode,
		renderQueue:         rq,
	}

	p.center = p.calculateCenter()

	p.vertexBuffer = createVertexBuffer(p.device, p.vertices, float32(scd.Width), float32(scd.Height))
	p.createTransformBuffer()
	p.createBindGroupLayout(device)
	p.createBindGroup(device, p.bindGroupLayout)

	shaderModule, err := device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{Code: ShapesShaderCode},
	})
	if err != nil {
		panic(err)
	}
	defer shaderModule.Release()

	filledPipeline := p.createPipeline(device, shaderModule, scd, wgpu.PrimitiveTopology_TriangleList)
	outlinedPipeline := p.createPipeline(device, shaderModule, scd, wgpu.PrimitiveTopology_LineStrip)
	p.pipeline = map[RenderMode]*wgpu.RenderPipeline{
		RenderFilled:   filledPipeline,
		RenderOutlined: outlinedPipeline,
	}

	return p
}

func (p *Polygon) calculateCenter() [2]float32 {
	var sumX, sumY float32
	for _, vertex := range p.vertices {
		sumX += vertex.Position[0]
		sumY += vertex.Position[1]
	}
	count := float32(len(p.vertices))
	return [2]float32{sumX / count, sumY / count}
}

func (p *Polygon) createBindGroupLayout(device *wgpu.Device) {
	var err error
	p.bindGroupLayout, err = device.CreateBindGroupLayout(&wgpu.BindGroupLayoutDescriptor{
		Entries: []wgpu.BindGroupLayoutEntry{
			{
				Binding:    0,
				Visibility: wgpu.ShaderStage_Vertex,
				Buffer: wgpu.BufferBindingLayout{
					Type: wgpu.BufferBindingType_Uniform,
				},
			},
		},
		Label: "Polygon Bind Group Layout",
	})
	if err != nil {
		panic(err)
	}
}

func (p *Polygon) createBindGroup(device *wgpu.Device, layout *wgpu.BindGroupLayout) {
	var err error
	p.bindGroup, err = device.CreateBindGroup(&wgpu.BindGroupDescriptor{
		Layout: layout,
		Entries: []wgpu.BindGroupEntry{
			{
				Binding: 0,
				Buffer:  p.transformBuffer,
				Offset:  0,
				Size:    uint64(unsafe.Sizeof(p.transform)),
			},
		},
		Label: "Polygon Bind Group",
	})
	if err != nil {
		panic(err)
	}
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
	if encoder == nil || !p.shouldeRender || p.isDisposed {
		return
	}

	// Additional safety checks
	if p.bindGroup == nil {
		log.Println("Polygon RenderPass: bindGroup is nil")
		return
	}
	if p.vertexBuffer == nil {
		log.Println("Polygon RenderPass: vertexBuffer is nil")
		return
	}
	if p.pipeline == nil || p.pipeline[p.renderMode] == nil {
		log.Println("Polygon RenderPass: pipeline is nil or not set for the current render mode")
		return
	}

	encoder.SetPipeline(p.pipeline[p.renderMode])
	encoder.SetBindGroup(0, p.bindGroup, nil)
	encoder.SetVertexBuffer(0, p.vertexBuffer, 0, wgpu.WholeSize)

	vertexCount := uint32(len(p.vertices))
	encoder.Draw(vertexCount, 1, 0, 0)
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
	if p.transformBuffer != nil {
		p.transformBuffer.Release()
		p.transformBuffer = nil
	}
	if p.bindGroup != nil {
		p.bindGroup.Release()
		p.bindGroup = nil
	}
	if p.bindGroupLayout != nil {
		p.bindGroupLayout.Release()
		p.bindGroupLayout = nil
	}
	p.isDisposed = true
}

func (p *Polygon) Hide() {
	p.shouldeRender = false
}
