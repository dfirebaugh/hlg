package texture

import (
	"fmt"
	"image"
	"image/draw"
	"unsafe"

	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/common"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Texture struct {
	surface common.Surface
	*wgpu.SwapChainDescriptor
	*wgpu.Device
	*wgpu.Texture
	*wgpu.TextureView
	*wgpu.Sampler
	*wgpu.BindGroup
	*wgpu.BindGroupLayout
	*wgpu.RenderPipeline
	vertexBuffer *wgpu.Buffer
	indexBuffer  *wgpu.Buffer

	numIndices uint32

	*common.Transform

	originalWidth  float32
	originalHeight float32

	flipHorizontal bool
	flipVertical   bool
	flipBuffer     *wgpu.Buffer
	flipMatrix     [2]float32

	clipRect   [4]float32 // minX, minY, maxX, maxY
	clipBuffer *wgpu.Buffer
	isDisposed bool
}

func TextureFromImage(surface common.Surface, d *wgpu.Device, scd *wgpu.SwapChainDescriptor, img image.Image, label string) (t *Texture, err error) {
	defer func() {
		if err != nil {
			t.Destroy()
			t = nil
		}
	}()
	r := img.Bounds()
	width := r.Dx()
	height := r.Dy()

	t = &Texture{
		Device:              d,
		surface:             surface,
		SwapChainDescriptor: scd,
		originalWidth:       float32(width),
		originalHeight:      float32(height),
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
	t.Texture, err = d.CreateTexture(&wgpu.TextureDescriptor{
		Label:         label,
		Size:          size,
		MipLevelCount: 1,
		SampleCount:   1,
		Dimension:     wgpu.TextureDimension_2D,
		Format:        wgpu.TextureFormat_RGBA8UnormSrgb,
		Usage:         wgpu.TextureUsage_TextureBinding | wgpu.TextureUsage_CopyDst,
	})
	if err != nil {
		return
	}

	d.GetQueue().WriteTexture(
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
		return
	}

	t.Sampler, err = d.CreateSampler(nil)
	if err != nil {
		return
	}

	err = t.createVertexBuffer()
	if err != nil {
		return nil, err
	}
	err = t.createIndexBuffer()
	if err != nil {
		return nil, err
	}

	t.Transform = common.NewTransform(surface, d, scd, "Texture Transform Buffer", float32(width), float32(height))

	err = t.createFlipBuffer()
	if err != nil {
		return nil, err
	}
	err = t.createClipBuffer()
	if err != nil {
		return nil, err
	}
	t.SetDefaultClip()

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
		t.Texture, err = t.Device.CreateTexture(&wgpu.TextureDescriptor{
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

	t.Device.GetQueue().WriteTexture(
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
	shader, err := t.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label: "texture.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{
			Code: TextureShaderCode,
		},
	})
	if err != nil {
		return err
	}
	defer shader.Release()

	renderPipelineLayout, err := t.CreatePipelineLayout(&wgpu.PipelineLayoutDescriptor{
		Label: "Render Pipeline Layout",
		BindGroupLayouts: []*wgpu.BindGroupLayout{
			t.BindGroupLayout,
		},
	})
	if err != nil {
		return err
	}
	defer renderPipelineLayout.Release()

	t.RenderPipeline, err = t.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Label:  "Render Pipeline",
		Layout: renderPipelineLayout,
		Vertex: wgpu.VertexState{
			Module:     shader,
			EntryPoint: "vs_main",
			Buffers:    []wgpu.VertexBufferLayout{VertexBufferLayout},
		},
		Fragment: &wgpu.FragmentState{
			Module:     shader,
			EntryPoint: "fs_main",
			Targets: []wgpu.ColorTargetState{{
				Format: t.SwapChainDescriptor.Format,
				Blend: &wgpu.BlendState{
					Color: wgpu.BlendComponent{
						SrcFactor: wgpu.BlendFactor_SrcAlpha,
						DstFactor: wgpu.BlendFactor_OneMinusSrcAlpha,
						Operation: wgpu.BlendOperation_Add,
					},
					Alpha: wgpu.BlendComponent{
						SrcFactor: wgpu.BlendFactor_One,
						DstFactor: wgpu.BlendFactor_Zero,
						Operation: wgpu.BlendOperation_Add,
					},
				}, WriteMask: wgpu.ColorWriteMask_All,
			}},
		},
		Primitive: wgpu.PrimitiveState{
			Topology:  wgpu.PrimitiveTopology_TriangleList,
			FrontFace: wgpu.FrontFace_CCW,
			CullMode:  wgpu.CullMode_Back,
		},
		Multisample: wgpu.MultisampleState{
			Count:                  1,
			Mask:                   0xFFFFFFFF,
			AlphaToCoverageEnabled: false,
		},
	})

	return err
}

func (t *Texture) createBindGroup() error {
	var err error
	t.BindGroupLayout, err = t.Device.CreateBindGroupLayout(&wgpu.BindGroupLayoutDescriptor{
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

	t.BindGroup, err = t.Device.CreateBindGroup(&wgpu.BindGroupDescriptor{
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
				Buffer:  t.flipBuffer,
				Offset:  0,
				Size:    uint64(unsafe.Sizeof(t.flipMatrix)),
			},
			{
				Binding: 4,
				Buffer:  t.clipBuffer,
				Offset:  0,
				Size:    uint64(unsafe.Sizeof(t.clipRect)),
			},
		},
		Label: "DiffuseBindGroup",
	})

	return err
}

func (t *Texture) createVertexBuffer() error {
	var err error
	sw, sh := t.surface.GetSurfaceSize()

	clipWidth := (t.clipRect[2] - t.clipRect[0]) * t.originalWidth
	clipHeight := (t.clipRect[3] - t.clipRect[1]) * t.originalHeight

	offsetX := (float32(sw) - clipWidth) / 2
	offsetY := (float32(sh) - clipHeight) / 2

	bottomLeft := common.ScreenToNDC(offsetX, offsetY+clipHeight, float32(sw), float32(sh))
	bottomRight := common.ScreenToNDC(offsetX+clipWidth, offsetY+clipHeight, float32(sw), float32(sh))
	topLeft := common.ScreenToNDC(offsetX, offsetY, float32(sw), float32(sh))
	topRight := common.ScreenToNDC(offsetX+clipWidth, offsetY, float32(sw), float32(sh))

	t.vertexBuffer, err = t.Device.CreateBufferInit(&wgpu.BufferInitDescriptor{
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
	t.indexBuffer, err = t.Device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Index Buffer",
		Contents: wgpu.ToBytes(INDICES[:]),
		Usage:    wgpu.BufferUsage_Index,
	})

	t.numIndices = uint32(len(INDICES))

	return err
}

func (t *Texture) createFlipBuffer() error {
	flipInfo := [2]float32{0.0, 0.0}

	var err error
	t.flipBuffer, err = t.Device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Flip Buffer",
		Usage:    wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
		Contents: wgpu.ToBytes(flipInfo[:]),
	})
	return err
}

func (t *Texture) createClipBuffer() error {
	clipInfo := [4]float32{0.0, 0.0, t.originalWidth, t.originalWidth}

	var err error
	t.clipBuffer, err = t.Device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Clip Buffer",
		Usage:    wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
		Contents: wgpu.ToBytes(clipInfo[:]),
	})
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
	if t.flipBuffer != nil {
		t.flipBuffer.Release()
		t.flipBuffer = nil
	}
	if t.clipBuffer != nil {
		t.clipBuffer.Release()
		t.clipBuffer = nil
	}

	t.Transform.Destroy()
	t.isDisposed = true
}

func (t *Texture) IsDisposed() bool {
	return t.isDisposed
}

func (t *Texture) updateFlipBuffer() {
	t.flipMatrix[0] = 0.0
	t.flipMatrix[1] = 0.0

	if t.flipHorizontal {
		t.flipMatrix[0] = 1.0
	}

	if t.flipVertical {
		t.flipMatrix[1] = 1.0
	}

	t.Device.GetQueue().WriteBuffer(t.flipBuffer, 0, wgpu.ToBytes(t.flipMatrix[:]))
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
