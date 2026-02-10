//go:build !js

package pipelines

import (
	"log"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/context"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

// vertexFormatMap pre-computed map for O(1) vertex format lookup
var vertexFormatMap = map[string]wgpu.VertexFormat{
	"uint8x2":   wgpu.VertexFormat_Uint8x2,
	"uint8x4":   wgpu.VertexFormat_Uint8x4,
	"sint8x2":   wgpu.VertexFormat_Sint8x2,
	"sint8x4":   wgpu.VertexFormat_Sint8x4,
	"unorm8x2":  wgpu.VertexFormat_Unorm8x2,
	"unorm8x4":  wgpu.VertexFormat_Unorm8x4,
	"snorm8x2":  wgpu.VertexFormat_Snorm8x2,
	"snorm8x4":  wgpu.VertexFormat_Snorm8x4,
	"uint16x2":  wgpu.VertexFormat_Uint16x2,
	"uint16x4":  wgpu.VertexFormat_Uint16x4,
	"sint16x2":  wgpu.VertexFormat_Sint16x2,
	"sint16x4":  wgpu.VertexFormat_Sint16x4,
	"unorm16x2": wgpu.VertexFormat_Unorm16x2,
	"unorm16x4": wgpu.VertexFormat_Unorm16x4,
	"snorm16x2": wgpu.VertexFormat_Snorm16x2,
	"snorm16x4": wgpu.VertexFormat_Snorm16x4,
	"float16x2": wgpu.VertexFormat_Float16x2,
	"float16x4": wgpu.VertexFormat_Float16x4,
	"float32":   wgpu.VertexFormat_Float32,
	"float32x2": wgpu.VertexFormat_Float32x2,
	"float32x3": wgpu.VertexFormat_Float32x3,
	"float32x4": wgpu.VertexFormat_Float32x4,
	"uint32":    wgpu.VertexFormat_Uint32,
	"uint32x2":  wgpu.VertexFormat_Uint32x2,
	"uint32x3":  wgpu.VertexFormat_Uint32x3,
	"uint32x4":  wgpu.VertexFormat_Uint32x4,
	"sint32":    wgpu.VertexFormat_Sint32,
	"sint32x2":  wgpu.VertexFormat_Sint32x2,
	"sint32x3":  wgpu.VertexFormat_Sint32x3,
	"sint32x4":  wgpu.VertexFormat_Sint32x4,
}

// Renderable structure, now with dynamic vertex handling
type Renderable struct {
	context.RenderContext
	Shader             *wgpu.ShaderModule
	VertexBuffer       *wgpu.Buffer
	vertexBufferLayout *wgpu.VertexBufferLayout
	BindGroup          *wgpu.BindGroup
	BindGroupLayout    *wgpu.BindGroupLayout
	Pipeline           *wgpu.RenderPipeline
	Uniforms           map[string]Uniform
	isDisposed         bool
	shouldRender       bool
	vertexData         []byte
}

type Uniform struct {
	Binding uint32
	Buffer  *wgpu.Buffer
	Size    uint64
}

func NewRenderable(ctx context.RenderContext, vertexData []byte, layout graphics.VertexBufferLayout, shaderHandle int, uniforms map[string]Uniform) *Renderable {
	if ctx == nil {
		log.Fatal("RenderContext is nil")
	}

	r := &Renderable{
		RenderContext: ctx,
		vertexData:    vertexData,
		Uniforms:      uniforms,
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

	r.createPipeline(layout)
	if r.Pipeline == nil {
		log.Fatal("Pipeline creation failed")
	}

	return r
}

func translateVertexBufferLayout(layout graphics.VertexBufferLayout) wgpu.VertexBufferLayout {
	var translatedAttributes []wgpu.VertexAttribute
	for _, attr := range layout.Attributes {
		format := translateVertexFormat(attr.Format)
		translatedAttributes = append(translatedAttributes, wgpu.VertexAttribute{
			ShaderLocation: attr.ShaderLocation,
			Offset:         attr.Offset,
			Format:         format,
		})
	}

	return wgpu.VertexBufferLayout{
		ArrayStride: layout.ArrayStride,
		Attributes:  translatedAttributes,
	}
}

func translateVertexFormat(customFormat string) wgpu.VertexFormat {
	if format, ok := vertexFormatMap[customFormat]; ok {
		return format
	}
	log.Fatalf("Unknown vertex format: %s", customFormat)
	return wgpu.VertexFormat_Float32x4
}

func (r *Renderable) createVertexBuffer() {
	var err error
	r.VertexBuffer, err = r.GetDevice().CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Dynamic Vertex Buffer",
		Contents: r.vertexData,
		Usage:    wgpu.BufferUsage_Vertex,
	})
	if err != nil {
		log.Fatal("Failed to create Vertex Buffer:", err)
	}
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
	r.BindGroupLayout, err = r.GetDevice().CreateBindGroupLayout(&wgpu.BindGroupLayoutDescriptor{
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
	r.BindGroup, err = r.GetDevice().CreateBindGroup(&wgpu.BindGroupDescriptor{
		Layout:  r.BindGroupLayout,
		Entries: entries,
		Label:   "Renderable Bind Group",
	})
	if err != nil {
		log.Fatal("Failed to create Bind Group:", err)
	}
}

func (r *Renderable) createPipeline(layout graphics.VertexBufferLayout) {
	pipelineName := "user_defined_pipeline"

	translatedLayout := translateVertexBufferLayout(layout)
	r.vertexBufferLayout = &translatedLayout

	r.Pipeline = r.GetPipelineManager().GetPipeline(
		pipelineName,
		&wgpu.PipelineLayoutDescriptor{
			BindGroupLayouts: []*wgpu.BindGroupLayout{
				r.BindGroupLayout,
			},
		},
		r.Shader,
		r.GetSwapChainDescriptor(),
		wgpu.PrimitiveTopology_TriangleList,
		[]wgpu.VertexBufferLayout{translatedLayout},
	)
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

	vertexCount := uint32(len(r.vertexData) / int(r.vertexBufferLayout.ArrayStride))
	if vertexCount == 0 {
		log.Fatal("RenderPass: No vertices to draw")
	}
	encoder.Draw(vertexCount, 1, 0, 0)
}

func (r *Renderable) IsDisposed() bool {
	return r.isDisposed
}

func (r *Renderable) Render() {
	if r.isDisposed {
		return
	}
	r.shouldRender = true
	r.AddToRenderQueue(r)
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
			_ = r.GetDevice().GetQueue().WriteBuffer(uniform.Buffer, 0, data)
		} else {
			log.Printf("Uniform %s does not exist", name)
		}
	}
}

func (r *Renderable) UpdateUniform(name string, data []byte) {
	if uniform, exists := r.Uniforms[name]; exists {
		_ = r.GetDevice().GetQueue().WriteBuffer(uniform.Buffer, 0, data)
	} else {
		log.Printf("Uniform %s does not exist", name)
	}
}
