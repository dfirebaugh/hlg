package shapes

import (
	"image/color"
	"log"
	"unsafe"

	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/common"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Polygon struct {
	surface common.Surface
	device  *wgpu.Device
	renderQueue
	*wgpu.SwapChainDescriptor
	bindGroup       *wgpu.BindGroup
	bindGroupLayout *wgpu.BindGroupLayout
	pipeline        *wgpu.RenderPipeline

	*common.Transform
	vertexBuffer *wgpu.Buffer
	vertices     []common.Vertex

	shouldeRender bool
	isDisposed    bool
}

func NewPolygon(surface common.Surface, device *wgpu.Device, scd *wgpu.SwapChainDescriptor, rq renderQueue, vertices []common.Vertex) *Polygon {
	p := &Polygon{
		surface:             surface,
		device:              device,
		SwapChainDescriptor: scd,
		vertices:            vertices,
		renderQueue:         rq,
	}

	sw, sh := surface.GetSurfaceSize()

	p.Transform = common.NewTransform(surface, device, scd, "Polygon Transform Buffer", float32(sw), float32(sh))
	p.vertexBuffer = common.CreateVertexBuffer(p.device, p.vertices, float32(sw), float32(sh))

	p.createBindGroupLayout(device)
	p.createBindGroup(device, p.bindGroupLayout)

	shaderModule, err := device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{Code: ShapesShaderCode},
	})
	if err != nil {
		panic(err)
	}
	defer shaderModule.Release()

	p.pipeline = p.createPipeline(device, shaderModule, scd, wgpu.PrimitiveTopology_TriangleList)

	return p
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
				Buffer:  p.Transform.Buffer,
				Offset:  0,
				Size:    uint64(unsafe.Sizeof(p.Transform.Matrix)),
			},
		},
		Label: "Polygon Bind Group",
	})
	if err != nil {
		panic(err)
	}
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
	if p.Transform.Buffer == nil {
		log.Println("Polygon RenderPass: vertexBuffer is nil")
		return
	}
	if p.pipeline == nil {
		log.Println("Polygon RenderPass: pipeline is nil or not set for the current render mode")
		return
	}

	encoder.SetPipeline(p.pipeline)
	encoder.SetBindGroup(0, p.bindGroup, nil)
	encoder.SetVertexBuffer(0, p.vertexBuffer, 0, wgpu.WholeSize)

	vertexCount := uint32(len(p.vertices))
	encoder.Draw(vertexCount, 1, 0, 0)
}

func (p *Polygon) Render() {
	if p.isDisposed {
		return
	}
	p.shouldeRender = true
	p.renderQueue.AddToRenderQueue(p)
}

// Dispose releases any resources used by the Triangle.
func (p *Polygon) Dispose() {
	if p.pipeline != nil {
		p.pipeline.Release()
		p.pipeline = nil
	}
	p.Transform.Destroy()
	if p.vertexBuffer != nil {
		p.vertexBuffer.Release()
		p.vertexBuffer = nil
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

func (p *Polygon) IsDisposed() bool {
	return p.isDisposed
}

func (p *Polygon) Hide() {
	p.shouldeRender = false
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

	p.createVertexBuffer()
}

func (p *Polygon) createVertexBuffer() {
	w, h := p.surface.GetSurfaceSize()
	p.vertexBuffer = common.CreateVertexBuffer(p.device, p.vertices, float32(w), float32(h))
}
