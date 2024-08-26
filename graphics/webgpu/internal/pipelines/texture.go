package pipelines

import (
	_ "embed"
	"fmt"
	"image"
	"image/draw"
	"log"
	"unsafe"

	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/context"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/primitives"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/transforms"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

//go:embed texture.wgsl
var TextureShaderCode string

type Vertex struct {
	position  [3]float32
	texCoords [2]float32
}

var VertexBufferLayout = wgpu.VertexBufferLayout{
	ArrayStride: uint64(unsafe.Sizeof(Vertex{})),
	StepMode:    wgpu.VertexStepMode_Vertex,
	Attributes: []wgpu.VertexAttribute{
		{
			Offset:         0,
			ShaderLocation: 0,
			Format:         wgpu.VertexFormat_Float32x3,
		},
		{
			Offset:         uint64(unsafe.Sizeof([3]float32{})),
			ShaderLocation: 1,
			Format:         wgpu.VertexFormat_Float32x2,
		},
	},
}

var INDICES = [...]uint16{
	0, 1, 2, // first triangle
	2, 1, 3, // second triangle
}

type Texture struct {
	context.RenderContext

	*wgpu.Texture
	*wgpu.TextureView
	*wgpu.Sampler
	*wgpu.BindGroup
	*wgpu.BindGroupLayout
	*wgpu.RenderPipeline
	vertexBuffer *wgpu.Buffer
	indexBuffer  *wgpu.Buffer

	numIndices uint32

	*transforms.Transform

	originalWidth  float32
	originalHeight float32

	isDisposed bool
}

func TextureFromImage(ctx context.RenderContext, img image.Image, label string) (*Texture, error) {
	r := img.Bounds()
	width := r.Dx()
	height := r.Dy()

	t := &Texture{
		RenderContext:  ctx,
		originalWidth:  float32(width),
		originalHeight: float32(height),
	}

	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		rgbaImg = image.NewRGBA(r)
		draw.Draw(rgbaImg, r, img, image.Point{}, draw.Over)
	}

	size := wgpu.Extent3D{
		Width:              uint32(width),
		Height:             uint32(height),
		DepthOrArrayLayers: 1,
	}
	var err error
	t.Texture, err = ctx.GetDevice().CreateTexture(&wgpu.TextureDescriptor{
		Label:         label,
		Size:          size,
		MipLevelCount: 1,
		SampleCount:   1,
		Dimension:     wgpu.TextureDimension_2D,
		Format:        wgpu.TextureFormat_RGBA8UnormSrgb,
		Usage:         wgpu.TextureUsage_TextureBinding | wgpu.TextureUsage_CopyDst,
	})
	if err != nil {
		return nil, err
	}

	ctx.GetDevice().GetQueue().WriteTexture(
		&wgpu.ImageCopyTexture{
			Aspect:   wgpu.TextureAspect_All,
			Texture:  t.Texture,
			MipLevel: 0,
			Origin:   wgpu.Origin3D{X: 0, Y: 0, Z: 0},
		},
		rgbaImg.Pix,
		&wgpu.TextureDataLayout{
			Offset:       0,
			BytesPerRow:  4 * uint32(width),
			RowsPerImage: uint32(height),
		},
		&size,
	)

	t.TextureView, err = t.Texture.CreateView(nil)
	if err != nil {
		return nil, err
	}

	t.Sampler, err = ctx.GetDevice().CreateSampler(nil)
	if err != nil {
		return nil, err
	}

	t.Transform = transforms.NewTransform(ctx, "Texture Transform Buffer", float32(width), float32(height))
	err = t.createVertexBuffer()
	if err != nil {
		return nil, err
	}
	err = t.createIndexBuffer()
	if err != nil {
		return nil, err
	}

	err = t.createBindGroup()
	if err != nil {
		return nil, err
	}
	err = t.createPipeline()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Texture) UpdateImage(img image.Image) error {
	r := img.Bounds()
	width := r.Dx()
	height := r.Dy()

	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		rgbaImg = image.NewRGBA(r)
		draw.Draw(rgbaImg, r, img, image.Point{}, draw.Over)
	}

	if int(t.originalWidth) != width || int(t.originalHeight) != height {
		size := wgpu.Extent3D{
			Width:              uint32(width),
			Height:             uint32(height),
			DepthOrArrayLayers: 1,
		}

		if t.Texture != nil {
			t.Texture.Release()
		}

		var err error
		t.Texture, err = t.GetDevice().CreateTexture(&wgpu.TextureDescriptor{
			Label:         "UpdatedTexture",
			Size:          size,
			MipLevelCount: 1,
			SampleCount:   1,
			Dimension:     wgpu.TextureDimension_2D,
			Format:        wgpu.TextureFormat_RGBA8UnormSrgb,
			Usage:         wgpu.TextureUsage_TextureBinding | wgpu.TextureUsage_CopyDst,
		})
		if err != nil {
			return err
		}

		t.originalWidth = float32(width)
		t.originalHeight = float32(height)
	}

	size := wgpu.Extent3D{
		Width:              uint32(width),
		Height:             uint32(height),
		DepthOrArrayLayers: 1,
	}

	t.GetDevice().GetQueue().WriteTexture(
		&wgpu.ImageCopyTexture{
			Aspect:   wgpu.TextureAspect_All,
			Texture:  t.Texture,
			MipLevel: 0,
			Origin:   wgpu.Origin3D{X: 0, Y: 0, Z: 0},
		},
		rgbaImg.Pix,
		&wgpu.TextureDataLayout{
			Offset:       0,
			BytesPerRow:  4 * uint32(width),
			RowsPerImage: uint32(height),
		},
		&size,
	)

	return nil
}

func (t *Texture) createPipeline() error {
	shader, err := t.GetDevice().CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label: "texture.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{
			Code: TextureShaderCode,
		},
	})
	if err != nil {
		return err
	}
	defer shader.Release()

	t.RenderPipeline = t.GetPipelineManager().GetPipeline(
		"texture-pipeline",
		&wgpu.PipelineLayoutDescriptor{
			Label: "Render Pipeline Layout",
			BindGroupLayouts: []*wgpu.BindGroupLayout{
				t.BindGroupLayout,
			},
		},
		shader,
		t.GetSwapChainDescriptor(),
		wgpu.PrimitiveTopology_TriangleList,
		[]wgpu.VertexBufferLayout{VertexBufferLayout},
	)

	return nil
}

func (t *Texture) createBindGroup() error {
	var err error
	t.BindGroupLayout, err = t.GetDevice().CreateBindGroupLayout(&wgpu.BindGroupLayoutDescriptor{
		Entries: []wgpu.BindGroupLayoutEntry{
			{
				Binding:    0,
				Visibility: wgpu.ShaderStage_Fragment,
				Texture: wgpu.TextureBindingLayout{
					Multisampled:  false,
					ViewDimension: wgpu.TextureViewDimension_2D,
					SampleType:    wgpu.TextureSampleType_Float,
				},
			},
			{
				Binding:    1,
				Visibility: wgpu.ShaderStage_Fragment,
				Sampler: wgpu.SamplerBindingLayout{
					Type: wgpu.SamplerBindingType_Filtering,
				},
			},
			{
				Binding:    2,
				Visibility: wgpu.ShaderStage_Vertex,
				Buffer: wgpu.BufferBindingLayout{
					Type: wgpu.BufferBindingType_Uniform,
				},
			},
			{
				Binding:    3,
				Visibility: wgpu.ShaderStage_Vertex,
				Buffer: wgpu.BufferBindingLayout{
					Type: wgpu.BufferBindingType_Uniform,
				},
			},
			{
				Binding:    4,
				Visibility: wgpu.ShaderStage_Vertex,
				Buffer: wgpu.BufferBindingLayout{
					Type: wgpu.BufferBindingType_Uniform,
				},
			},
		},
		Label: "TextureBindGroupLayout",
	})
	if err != nil {
		return err
	}

	t.BindGroup, err = t.GetDevice().CreateBindGroup(&wgpu.BindGroupDescriptor{
		Layout: t.BindGroupLayout,
		Entries: []wgpu.BindGroupEntry{
			{
				Binding:     0,
				TextureView: t.TextureView,
			},
			{
				Binding: 1,
				Sampler: t.Sampler,
			},
			{
				Binding: 2,
				Buffer:  t.Transform.Buffer,
				Offset:  0,
				Size:    uint64(unsafe.Sizeof(t.Transform.Matrix)),
			},
			{
				Binding: 3,
				Buffer:  t.Transform.FlipBuffer,
				Offset:  0,
				Size:    uint64(unsafe.Sizeof(t.Transform.FlipMatrix)),
			},
			{
				Binding: 4,
				Buffer:  t.Transform.ClipBuffer,
				Offset:  0,
				Size:    uint64(unsafe.Sizeof(t.Transform.ClipRect)),
			},
		},
		Label: "DiffuseBindGroup",
	})

	return err
}

func (t *Texture) createVertexBuffer() error {
	var err error
	sw, sh := t.GetSurfaceSize()

	clipWidth := (t.ClipRect[2] - t.ClipRect[0]) * t.originalWidth
	clipHeight := (t.ClipRect[3] - t.ClipRect[1]) * t.originalHeight

	offsetX := (float32(sw) - clipWidth) / 2
	offsetY := (float32(sh) - clipHeight) / 2

	bottomLeft := primitives.ScreenToNDC(offsetX, offsetY+clipHeight, float32(sw), float32(sh))
	bottomRight := primitives.ScreenToNDC(offsetX+clipWidth, offsetY+clipHeight, float32(sw), float32(sh))
	topLeft := primitives.ScreenToNDC(offsetX, offsetY, float32(sw), float32(sh))
	topRight := primitives.ScreenToNDC(offsetX+clipWidth, offsetY, float32(sw), float32(sh))

	t.vertexBuffer, err = t.GetDevice().CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label: "Vertex Buffer",
		Contents: wgpu.ToBytes(
			[]Vertex{
				{
					position:  bottomLeft,
					texCoords: [2]float32{0.0, 1.0},
				},
				{
					position:  bottomRight,
					texCoords: [2]float32{1.0, 1.0},
				},
				{
					position:  topLeft,
					texCoords: [2]float32{0.0, 0.0},
				},
				{
					position:  topRight,
					texCoords: [2]float32{1.0, 0.0},
				},
			}),
		Usage: wgpu.BufferUsage_Vertex,
	})

	return err
}

func (t *Texture) updateVertexBuffer() error {
	if t.vertexBuffer != nil {
		t.vertexBuffer.Release()
	}

	err := t.createVertexBuffer()
	if err != nil {
		return fmt.Errorf("failed to update vertex buffer: %w", err)
	}

	return nil
}

func (t *Texture) createIndexBuffer() error {
	var err error
	t.indexBuffer, err = t.GetDevice().CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Index Buffer",
		Contents: wgpu.ToBytes(INDICES[:]),
		Usage:    wgpu.BufferUsage_Index,
	})

	t.numIndices = uint32(len(INDICES))

	return err
}

func (t *Texture) Destroy() {
	if t.Sampler != nil {
		t.Sampler.Release()
		t.Sampler = nil
	}
	if t.TextureView != nil {
		t.TextureView.Release()
		t.TextureView = nil
	}
	if t.Texture != nil {
		t.Texture.Release()
		t.Texture = nil
	}
	if t.indexBuffer != nil {
		t.indexBuffer.Release()
		t.indexBuffer = nil
	}
	if t.vertexBuffer != nil {
		t.vertexBuffer.Release()
		t.vertexBuffer = nil
	}
	if t.BindGroup != nil {
		t.BindGroup.Release()
		t.BindGroup = nil
	}
	if t.BindGroupLayout != nil {
		t.BindGroupLayout.Release()
		t.BindGroupLayout = nil
	}
	if t.Texture != nil {
		t.Texture.Destroy()
		t.Texture = nil
	}

	t.Transform.Destroy()
	t.isDisposed = true
}

func (t *Texture) IsDisposed() bool {
	return t.isDisposed
}

func (t *Texture) RenderPass(pass *wgpu.RenderPassEncoder) {
	if t.isDisposed {
		return
	}
	pass.SetPipeline(t.RenderPipeline)
	pass.SetBindGroup(0, t.BindGroup, nil)
	pass.SetVertexBuffer(0, t.vertexBuffer, 0, wgpu.WholeSize)
	pass.SetIndexBuffer(t.indexBuffer, wgpu.IndexFormat_Uint16, 0, wgpu.WholeSize)
	pass.DrawIndexed(t.numIndices, 1, 0, 0, 0)
}

func (t *Texture) SetClipRect(minX, minY, maxX, maxY float32) {
	t.Transform.SetClipRect(minX, minY, maxX, maxY)
	if err := t.updateVertexBuffer(); err != nil {
		log.Printf("Failed to update vertex buffer: %v", err)
	}
}

func (t *Texture) Move(screenX float32, screenY float32) {
	t.MoveToScreenPosition(screenX, screenY)
}
