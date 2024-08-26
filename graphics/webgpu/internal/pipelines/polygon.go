package pipelines

import (
	"image/color"
	"log"
	"unsafe"

	_ "embed"

	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/context"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/primitives"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/transforms"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

//go:embed shapes.wgsl
var ShapesShaderCode string

type Polygon struct {
	context.RenderContext

	bindGroup       *wgpu.BindGroup
	bindGroupLayout *wgpu.BindGroupLayout
	pipeline        *wgpu.RenderPipeline

	*transforms.Transform
	vertexBuffer *wgpu.Buffer
	vertices     []primitives.Vertex

	shouldeRender bool
	isDisposed    bool
}

func NewPolygon(ctx context.RenderContext, vertices []primitives.Vertex) *Polygon {
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

	p := &Polygon{
		RenderContext: ctx,
		vertices:      vertices,
	}

	sw, sh := p.GetSurfaceSize()

	if sw == 0 || sh == 0 {
		log.Fatal("Surface size is invalid")
	}

	p.Transform = transforms.NewTransform(ctx, "Polygon Transform Buffer", float32(sw), float32(sh))
	if p.Transform == nil {
		log.Fatal("Transform initialization failed")
	}

	p.createVertexBuffer()
	if p.vertexBuffer == nil {
		log.Fatal("Vertex buffer is nil")
	}

	p.createBindGroupLayout(p.GetDevice())
	if p.bindGroupLayout == nil {
		log.Fatal("Bind group layout is nil")
	}

	p.createBindGroup(p.GetDevice(), p.bindGroupLayout)
	if p.bindGroup == nil {
		log.Fatal("Bind group is nil")
	}

	shaderModule, err := p.GetDevice().CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{Code: ShapesShaderCode},
	})
	if err != nil {
		log.Fatal("Failed to create shader module:", err)
	}
	defer shaderModule.Release()

	p.pipeline = ctx.GetPipelineManager().GetPipeline("polygon",
		&wgpu.PipelineLayoutDescriptor{
			Label: "Render Pipeline Layout",
			BindGroupLayouts: []*wgpu.BindGroupLayout{
				p.bindGroupLayout,
			},
		},
		shaderModule,
		p.GetSwapChainDescriptor(), wgpu.PrimitiveTopology_TriangleList,
		[]wgpu.VertexBufferLayout{{
			ArrayStride: uint64(unsafe.Sizeof(primitives.Vertex{})),
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
	)

	if p.pipeline == nil {
		log.Fatal("Pipeline creation failed")
	}

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
	if p.isDisposed || p.IsOffScreen() {
		return
	}
	p.shouldeRender = true
	p.AddToRenderQueue(p)
}

func (p *Polygon) Dispose() {
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
	w, h := p.GetSurfaceSize()
	p.vertexBuffer = primitives.CreateVertexBuffer(p.GetDevice(), p.vertices, float32(w), float32(h))
}

func (p *Polygon) IsOffScreen() bool {
	for _, vertex := range p.vertices {
		x := vertex.Position[0]
		y := vertex.Position[1]

		sw, sh := p.GetSurfaceSize()
		if x >= 0 && x <= float32(sw) && y >= 0 && y <= float32(sh) {
			return false
		}
	}
	return true
}

func (p *Polygon) Move(destX, destY float32) {
	p.RecenterAndMove(p.vertices, destX, destY)
	p.vertexBuffer.Release()
	p.createVertexBuffer()
}
