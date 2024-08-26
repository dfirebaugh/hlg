package pipelines

import (
	"log"
	"unsafe"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/context"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/primitives"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/transforms"
	"github.com/google/uuid"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Renderable struct {
	Context context.RenderContext
	Shader  *wgpu.ShaderModule
	*transforms.Transform
	VertexBuffer    *wgpu.Buffer
	BindGroup       *wgpu.BindGroup
	BindGroupLayout *wgpu.BindGroupLayout
	Pipeline        *wgpu.RenderPipeline
	Uniforms        map[string]Uniform
	Vertices        []primitives.Vertex
	isDisposed      bool
	shouldRender    bool
}

type Uniform struct {
	Binding uint32
	Buffer  *wgpu.Buffer
	Size    uint64
}

func NewRenderable(ctx context.RenderContext, vertices []primitives.Vertex, shaderHandle int, uniforms map[string]Uniform) *Renderable {
	if ctx == nil {
		log.Fatal("RenderContext is nil")
	}

	r := &Renderable{
		Context:  ctx,
		Vertices: vertices,
		Uniforms: uniforms,
	}

	sw, sh := ctx.GetSurfaceSize()
	if sw == 0 || sh == 0 {
		log.Fatal("Surface size is invalid")
	}

	r.Transform = transforms.NewTransform(ctx, "Renderable Transform Buffer", float32(sw), float32(sh))
	if r.Transform == nil {
		log.Fatal("Failed to create Transform")
	}

	r.createVertexBuffer()
	if r.VertexBuffer == nil {
		log.Fatal("Vertex buffer is nil")
	}

	r.Shader = ctx.GetShader(graphics.ShaderHandle(shaderHandle))
	if r.Shader == nil {
		log.Fatal("Shader module is nil")
	}

	r.createBindGroupLayout()
	if r.BindGroupLayout == nil {
		log.Fatal("Bind group layout is nil")
	}

	r.createBindGroup()
	if r.BindGroup == nil {
		log.Fatal("Bind group is nil")
	}

	r.createPipeline()
	if r.Pipeline == nil {
		log.Fatal("Pipeline creation failed")
	}

	return r
}

func (r *Renderable) createVertexBuffer() {
	vertexBuffer, err := r.Context.GetDevice().CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Vertex Buffer",
		Contents: wgpu.ToBytes(r.Vertices[:]),
		Usage:    wgpu.BufferUsage_Vertex,
	})
	if err != nil {
		panic(err)
	}
	r.VertexBuffer = vertexBuffer
}

func (r *Renderable) createBindGroupLayout() {
	var entries []wgpu.BindGroupLayoutEntry
	for _, uniform := range r.Uniforms {
		entries = append(entries, wgpu.BindGroupLayoutEntry{
			Binding:    uniform.Binding,
			Visibility: wgpu.ShaderStage_Vertex | wgpu.ShaderStage_Fragment,
			Buffer: wgpu.BufferBindingLayout{
				Type: wgpu.BufferBindingType_Uniform,
			},
		})
	}

	var err error
	r.BindGroupLayout, err = r.Context.GetDevice().CreateBindGroupLayout(&wgpu.BindGroupLayoutDescriptor{
		Entries: entries,
		Label:   "Dynamic Bind Group Layout",
	})
	if err != nil {
		log.Fatal("Failed to create Bind Group Layout:", err)
	}
}

func (r *Renderable) createBindGroup() {
	var entries []wgpu.BindGroupEntry
	for _, uniform := range r.Uniforms {
		entries = append(entries, wgpu.BindGroupEntry{
			Binding: uniform.Binding,
			Buffer:  uniform.Buffer,
			Offset:  0,
			Size:    uniform.Size,
		})
	}

	var err error
	r.BindGroup, err = r.Context.GetDevice().CreateBindGroup(&wgpu.BindGroupDescriptor{
		Layout:  r.BindGroupLayout,
		Entries: entries,
		Label:   "Renderable Bind Group",
	})
	if err != nil {
		log.Fatal("Failed to create Bind Group:", err)
	}
}

func (r *Renderable) createPipeline() {
	pipelineName := uuid.New().String()

	r.Pipeline = r.Context.GetPipelineManager().GetPipeline(
		pipelineName,
		&wgpu.PipelineLayoutDescriptor{
			BindGroupLayouts: []*wgpu.BindGroupLayout{
				r.BindGroupLayout,
			},
		},
		r.Shader,
		r.Context.GetSwapChainDescriptor(),
		wgpu.PrimitiveTopology_TriangleList,
		[]wgpu.VertexBufferLayout{{
			ArrayStride: uint64(unsafe.Sizeof(primitives.Vertex{})),
			Attributes: []wgpu.VertexAttribute{
				{
					ShaderLocation: 0,
					Offset:         0,
					Format:         wgpu.VertexFormat_Float32x3,
				},
				{
					ShaderLocation: 1,
					Offset:         uint64(unsafe.Sizeof([3]float32{})),
					Format:         wgpu.VertexFormat_Float32x4,
				},
			},
		}},
	)
	if r.Pipeline == nil {
		log.Fatal("Failed to create Render Pipeline")
	}
}

func (r *Renderable) RenderPass(encoder *wgpu.RenderPassEncoder) {
	if !r.shouldRender || r.isDisposed {
		return
	}

	if r.Pipeline == nil {
		log.Fatal("RenderPass: Pipeline is nil")
	}
	if r.BindGroup == nil {
		log.Fatal("RenderPass: BindGroup is nil")
	}
	if r.VertexBuffer == nil {
		log.Fatal("RenderPass: VertexBuffer is nil")
	}

	encoder.SetPipeline(r.Pipeline)
	encoder.SetBindGroup(0, r.BindGroup, nil)
	encoder.SetVertexBuffer(0, r.VertexBuffer, 0, wgpu.WholeSize)

	vertexCount := uint32(len(r.Vertices))
	if vertexCount == 0 {
		log.Fatal("RenderPass: No vertices to draw")
	}
	encoder.Draw(vertexCount, 1, 0, 0)
}

func (r *Renderable) IsDisposed() bool {
	return r.isDisposed
}

func (r *Renderable) Hide() {
	r.shouldRender = false
}

func (r *Renderable) Render() {
	if r.isDisposed {
		return
	}
	r.shouldRender = true
	r.Context.AddToRenderQueue(r)
}

func (r *Renderable) Dispose() {
	r.isDisposed = true
	if r.VertexBuffer != nil {
		r.VertexBuffer.Release()
		r.VertexBuffer = nil
	}
	if r.BindGroup != nil {
		r.BindGroup.Release()
		r.BindGroup = nil
	}
	if r.BindGroupLayout != nil {
		r.BindGroupLayout.Release()
		r.BindGroupLayout = nil
	}
	if r.Pipeline != nil {
		r.Pipeline.Release()
		r.Pipeline = nil
	}
	if r.Shader != nil {
		r.Shader.Release()
		r.Shader = nil
	}
}

func (r *Renderable) UpdateUniforms(dataMap map[string][]byte) {
	for name, data := range dataMap {
		if uniform, exists := r.Uniforms[name]; exists {
			r.Context.GetDevice().GetQueue().WriteBuffer(uniform.Buffer, 0, data)
		} else {
			log.Printf("Uniform %s does not exist", name)
		}
	}
}

func (r *Renderable) UpdateUniform(name string, data []byte) {
	if uniform, exists := r.Uniforms[name]; exists {
		r.Context.GetDevice().GetQueue().WriteBuffer(uniform.Buffer, 0, data)
	} else {
		log.Printf("Uniform %s does not exist", name)
	}
}
