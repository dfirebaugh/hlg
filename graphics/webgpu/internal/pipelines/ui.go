//go:build !js

package pipelines

import (
	"image"
	"image/draw"
	"log"
	"unsafe"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/context"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/shader"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

// PrimitiveBuffer implements the "one draw call UI" pattern from:
// https://ruby0x1.github.io/machinery_blog_archive/post/ui-rendering-using-primitive-buffers/
//
// Key insight: Store one compact primitive per shape in a storage buffer.
// The vertex shader constructs 6 vertices per primitive using vertex_index.
// Memory: 48 bytes/primitive vs 312 bytes (6 vertices x 52 bytes) = ~6.5x reduction
type PrimitiveBuffer struct {
	context.RenderContext

	pipeline        *wgpu.RenderPipeline
	bindGroupLayout *wgpu.BindGroupLayout
	bindGroup       *wgpu.BindGroup

	// Separate pipeline for vertex-buffer based shapes (PrimitiveShape)
	solidShapePipeline *wgpu.RenderPipeline

	storageBuffer *wgpu.Buffer
	primitives    []graphics.Primitive
	primitivesCap int // capacity to avoid frequent reallocations

	// Uniform buffer for screen size
	screenSizeBuffer *wgpu.Buffer
	screenSize       [2]float32

	// MSDF atlas resources
	msdfTexture      *wgpu.Texture
	msdfTextureView  *wgpu.TextureView
	msdfSampler      *wgpu.Sampler
	msdfParamsBuffer *wgpu.Buffer
	msdfParams       [4]float32 // px_range, tex_width, tex_height, unused

	isDisposed bool
}

func NewPrimitiveBuffer(ctx context.RenderContext, primitives []graphics.Primitive) *PrimitiveBuffer {
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

	sw, sh := ctx.GetSurfaceSize()

	p := &PrimitiveBuffer{
		RenderContext: ctx,
		primitives:    primitives,
		primitivesCap: max(len(primitives), 1024), // initial capacity
		msdfParams:    [4]float32{4.0, 1.0, 1.0, 0.0},
		screenSize:    [2]float32{float32(sw), float32(sh)},
	}

	p.createScreenSizeBuffer()
	p.createMSDFResources()
	p.createBindGroupLayout()

	// Create storage buffer and bind group
	p.createStorageBuffer()

	// Create pipeline (no vertex buffer layout - vertices constructed in shader)
	p.pipeline = ctx.GetPipelineManager().GetPipeline("primitive buffer",
		&wgpu.PipelineLayoutDescriptor{
			Label: "Primitive Buffer Pipeline Layout",
			BindGroupLayouts: []*wgpu.BindGroupLayout{
				p.bindGroupLayout,
			},
		},
		p.GetShader(shader.PrimitiveBufferShader),
		p.GetSwapChainDescriptor(),
		wgpu.PrimitiveTopology_TriangleList,
		[]wgpu.VertexBufferLayout{}, // No vertex buffers - using storage buffer
	)

	if p.pipeline == nil {
		log.Fatal("Pipeline creation failed")
	}

	// Create solid shape pipeline (vertex buffer based, for PrimitiveShape)
	// PrimitiveVertex layout: Position[3], LocalPosition[2], OpCode, Radius, Color[4], TexCoords[2]
	// Total: 3+2+1+1+4+2 = 13 floats = 52 bytes
	solidShapeLayout := []wgpu.VertexBufferLayout{
		{
			ArrayStride: 52, // sizeof(PrimitiveVertex)
			StepMode:    wgpu.VertexStepMode_Vertex,
			Attributes: []wgpu.VertexAttribute{
				{Format: wgpu.VertexFormat_Float32x3, Offset: 0, ShaderLocation: 0},  // position
				{Format: wgpu.VertexFormat_Float32x2, Offset: 12, ShaderLocation: 1}, // local_pos
				{Format: wgpu.VertexFormat_Float32, Offset: 20, ShaderLocation: 2},   // op_code
				{Format: wgpu.VertexFormat_Float32, Offset: 24, ShaderLocation: 3},   // radius
				{Format: wgpu.VertexFormat_Float32x4, Offset: 28, ShaderLocation: 4}, // color
				{Format: wgpu.VertexFormat_Float32x2, Offset: 44, ShaderLocation: 5}, // tex_coords
			},
		},
	}

	p.solidShapePipeline = ctx.GetPipelineManager().GetPipeline("solid shape",
		&wgpu.PipelineLayoutDescriptor{
			Label:            "Solid Shape Pipeline Layout",
			BindGroupLayouts: []*wgpu.BindGroupLayout{},
		},
		p.GetShader(shader.SolidShapeShader),
		p.GetSwapChainDescriptor(),
		wgpu.PrimitiveTopology_TriangleList,
		solidShapeLayout,
	)

	if p.solidShapePipeline == nil {
		log.Fatal("Solid shape pipeline creation failed")
	}

	return p
}

// NewPrimitiveBufferCompat creates a PrimitiveBuffer from the old PrimitiveVertex format
// for backwards compatibility during migration
func NewPrimitiveBufferCompat(ctx context.RenderContext, vertices []graphics.PrimitiveVertex, layout graphics.VertexBufferLayout) *PrimitiveBuffer {
	sw, sh := ctx.GetSurfaceSize()
	primitives := convertVerticesToPrimitives(vertices, float32(sw), float32(sh))
	return NewPrimitiveBuffer(ctx, primitives)
}

func (p *PrimitiveBuffer) createStorageBuffer() {
	var err error
	// Create buffer with initial capacity
	primitiveSize := uint64(unsafe.Sizeof(graphics.Primitive{})) // 64 bytes
	bufferSize := uint64(p.primitivesCap) * primitiveSize
	if bufferSize == 0 {
		bufferSize = primitiveSize // minimum size for one primitive
	}

	p.storageBuffer, err = p.GetDevice().CreateBuffer(&wgpu.BufferDescriptor{
		Label: "Primitive Storage Buffer",
		Size:  bufferSize,
		Usage: wgpu.BufferUsage_Storage | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		panic(err)
	}

	// Create bind group
	p.createBindGroup()

	// Upload initial data if any
	if len(p.primitives) > 0 {
		_ = p.GetDevice().GetQueue().WriteBuffer(p.storageBuffer, 0, wgpu.ToBytes(p.primitives))
	}
}

func (p *PrimitiveBuffer) createScreenSizeBuffer() {
	var err error
	p.screenSizeBuffer, err = p.GetDevice().CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Screen Size Buffer",
		Contents: wgpu.ToBytes(p.screenSize[:]),
		Usage:    wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		panic(err)
	}
}

func (p *PrimitiveBuffer) createMSDFResources() {
	var err error

	// Create a 1x1 placeholder texture
	p.msdfTexture, err = p.GetDevice().CreateTexture(&wgpu.TextureDescriptor{
		Label: "MSDF Placeholder Texture",
		Size: wgpu.Extent3D{
			Width:              1,
			Height:             1,
			DepthOrArrayLayers: 1,
		},
		MipLevelCount: 1,
		SampleCount:   1,
		Dimension:     wgpu.TextureDimension_2D,
		Format:        wgpu.TextureFormat_RGBA8Unorm,
		Usage:         wgpu.TextureUsage_TextureBinding | wgpu.TextureUsage_CopyDst,
	})
	if err != nil {
		log.Fatalf("Failed to create placeholder MSDF texture: %v", err)
	}

	_ = p.GetDevice().GetQueue().WriteTexture(
		&wgpu.ImageCopyTexture{
			Aspect:   wgpu.TextureAspect_All,
			Texture:  p.msdfTexture,
			MipLevel: 0,
			Origin:   wgpu.Origin3D{X: 0, Y: 0, Z: 0},
		},
		[]byte{0, 0, 0, 0},
		&wgpu.TextureDataLayout{
			Offset:       0,
			BytesPerRow:  4,
			RowsPerImage: 1,
		},
		&wgpu.Extent3D{Width: 1, Height: 1, DepthOrArrayLayers: 1},
	)

	p.msdfTextureView, err = p.msdfTexture.CreateView(nil)
	if err != nil {
		log.Fatalf("Failed to create MSDF texture view: %v", err)
	}

	p.msdfSampler, err = p.GetDevice().CreateSampler(&wgpu.SamplerDescriptor{
		AddressModeU:   wgpu.AddressMode_ClampToEdge,
		AddressModeV:   wgpu.AddressMode_ClampToEdge,
		AddressModeW:   wgpu.AddressMode_ClampToEdge,
		MagFilter:      wgpu.FilterMode_Linear,
		MinFilter:      wgpu.FilterMode_Linear,
		MipmapFilter:   wgpu.MipmapFilterMode_Linear,
		MaxAnisotrophy: 1,
	})
	if err != nil {
		log.Fatalf("Failed to create MSDF sampler: %v", err)
	}

	p.msdfParamsBuffer, err = p.GetDevice().CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "MSDF Params Buffer",
		Usage:    wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
		Contents: wgpu.ToBytes(p.msdfParams[:]),
	})
	if err != nil {
		log.Fatalf("Failed to create MSDF params buffer: %v", err)
	}
}

func (p *PrimitiveBuffer) createBindGroupLayout() {
	var err error
	p.bindGroupLayout, err = p.GetDevice().CreateBindGroupLayout(&wgpu.BindGroupLayoutDescriptor{
		Label: "Primitive Buffer Bind Group Layout",
		Entries: []wgpu.BindGroupLayoutEntry{
			{
				// Storage buffer for primitives
				Binding:    0,
				Visibility: wgpu.ShaderStage_Vertex,
				Buffer: wgpu.BufferBindingLayout{
					Type: wgpu.BufferBindingType_ReadOnlyStorage,
				},
			},
			{
				// Uniform buffer for screen size
				Binding:    1,
				Visibility: wgpu.ShaderStage_Vertex | wgpu.ShaderStage_Fragment,
				Buffer: wgpu.BufferBindingLayout{
					Type: wgpu.BufferBindingType_Uniform,
				},
			},
			{
				// MSDF texture
				Binding:    2,
				Visibility: wgpu.ShaderStage_Fragment,
				Texture: wgpu.TextureBindingLayout{
					Multisampled:  false,
					ViewDimension: wgpu.TextureViewDimension_2D,
					SampleType:    wgpu.TextureSampleType_Float,
				},
			},
			{
				// MSDF sampler
				Binding:    3,
				Visibility: wgpu.ShaderStage_Fragment,
				Sampler: wgpu.SamplerBindingLayout{
					Type: wgpu.SamplerBindingType_Filtering,
				},
			},
			{
				// MSDF params uniform
				Binding:    4,
				Visibility: wgpu.ShaderStage_Vertex | wgpu.ShaderStage_Fragment,
				Buffer: wgpu.BufferBindingLayout{
					Type: wgpu.BufferBindingType_Uniform,
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create bind group layout: %v", err)
	}
}

func (p *PrimitiveBuffer) createBindGroup() {
	var err error

	primitiveSize := uint64(unsafe.Sizeof(graphics.Primitive{})) // 64 bytes
	storageBufferSize := uint64(p.primitivesCap) * primitiveSize
	if storageBufferSize == 0 {
		storageBufferSize = primitiveSize
	}

	p.bindGroup, err = p.GetDevice().CreateBindGroup(&wgpu.BindGroupDescriptor{
		Label:  "Primitive Buffer Bind Group",
		Layout: p.bindGroupLayout,
		Entries: []wgpu.BindGroupEntry{
			{
				Binding: 0,
				Buffer:  p.storageBuffer,
				Offset:  0,
				Size:    storageBufferSize,
			},
			{
				Binding: 1,
				Buffer:  p.screenSizeBuffer,
				Offset:  0,
				Size:    uint64(unsafe.Sizeof(p.screenSize)),
			},
			{
				Binding:     2,
				TextureView: p.msdfTextureView,
			},
			{
				Binding: 3,
				Sampler: p.msdfSampler,
			},
			{
				Binding: 4,
				Buffer:  p.msdfParamsBuffer,
				Offset:  0,
				Size:    uint64(unsafe.Sizeof(p.msdfParams)),
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create bind group: %v", err)
	}
}

// SetMSDFAtlas sets the MSDF atlas texture for text rendering
func (p *PrimitiveBuffer) SetMSDFAtlas(atlasImg image.Image, pxRange float64) {
	r := atlasImg.Bounds()
	width := r.Dx()
	height := r.Dy()

	rgbaImg, ok := atlasImg.(*image.RGBA)
	if !ok {
		rgbaImg = image.NewRGBA(r)
		draw.Draw(rgbaImg, r, atlasImg, image.Point{}, draw.Over)
	}

	if p.msdfTextureView != nil {
		p.msdfTextureView.Release()
	}
	if p.msdfTexture != nil {
		p.msdfTexture.Release()
	}

	size := wgpu.Extent3D{
		Width:              uint32(width),
		Height:             uint32(height),
		DepthOrArrayLayers: 1,
	}

	var err error
	p.msdfTexture, err = p.GetDevice().CreateTexture(&wgpu.TextureDescriptor{
		Label:         "MSDF Atlas Texture",
		Size:          size,
		MipLevelCount: 1,
		SampleCount:   1,
		Dimension:     wgpu.TextureDimension_2D,
		Format:        wgpu.TextureFormat_RGBA8Unorm,
		Usage:         wgpu.TextureUsage_TextureBinding | wgpu.TextureUsage_CopyDst,
	})
	if err != nil {
		log.Printf("Failed to create MSDF atlas texture: %v", err)
		return
	}

	if err = p.GetDevice().GetQueue().WriteTexture(
		&wgpu.ImageCopyTexture{
			Aspect:   wgpu.TextureAspect_All,
			Texture:  p.msdfTexture,
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
	); err != nil {
		log.Printf("Failed to write MSDF atlas texture: %v", err)
		return
	}

	p.msdfTextureView, err = p.msdfTexture.CreateView(nil)
	if err != nil {
		log.Printf("Failed to create MSDF texture view: %v", err)
		return
	}

	p.msdfParams = [4]float32{float32(pxRange), float32(width), float32(height), p.msdfParams[3]}
	_ = p.GetDevice().GetQueue().WriteBuffer(p.msdfParamsBuffer, 0, wgpu.ToBytes(p.msdfParams[:]))

	// Recreate bind group with new texture
	if p.bindGroup != nil {
		p.bindGroup.Release()
	}
	p.createBindGroup()
}

// SetMSDFMode sets the MSDF rendering mode.
// Mode 0: median(RGB) - MSDF reconstruction for sharp corners (default)
// Mode 1: alpha channel only (true SDF fallback)
// Mode 2: visualize RGB channels directly (for debugging atlas)
func (p *PrimitiveBuffer) SetMSDFMode(mode int) {
	p.msdfParams[3] = float32(mode)
	if p.msdfParamsBuffer != nil {
		_ = p.GetDevice().GetQueue().WriteBuffer(p.msdfParamsBuffer, 0, wgpu.ToBytes(p.msdfParams[:]))
	}
}

// UpdatePrimitives updates the storage buffer with new primitives
func (p *PrimitiveBuffer) UpdatePrimitives(primitives []graphics.Primitive) {
	if len(primitives) == 0 {
		p.primitives = p.primitives[:0]
		return
	}

	// Make a copy of the primitives to ensure data stability
	if cap(p.primitives) < len(primitives) {
		p.primitives = make([]graphics.Primitive, len(primitives))
	} else {
		p.primitives = p.primitives[:len(primitives)]
	}
	copy(p.primitives, primitives)

	// Check if we need to grow the buffer
	if len(p.primitives) > p.primitivesCap {
		p.primitivesCap = len(p.primitives) * 2 // grow by 2x

		// Release old resources
		if p.bindGroup != nil {
			p.bindGroup.Release()
			p.bindGroup = nil
		}
		if p.storageBuffer != nil {
			p.storageBuffer.Release()
			p.storageBuffer = nil
		}

		// Create new buffer and bind group
		p.createStorageBuffer()
	}

	// Write to the storage buffer
	_ = p.GetDevice().GetQueue().WriteBuffer(p.storageBuffer, 0, wgpu.ToBytes(p.primitives))
}

// UpdateScreenSize updates the screen size uniform
func (p *PrimitiveBuffer) UpdateScreenSize(width, height int) {
	p.screenSize = [2]float32{float32(width), float32(height)}
	_ = p.GetDevice().GetQueue().WriteBuffer(p.screenSizeBuffer, 0, wgpu.ToBytes(p.screenSize[:]))
}

func (p *PrimitiveBuffer) RenderPass(encoder *wgpu.RenderPassEncoder) {
	if encoder == nil || p.isDisposed {
		return
	}
	if p.pipeline == nil {
		return
	}
	if len(p.primitives) == 0 {
		return
	}

	// Update screen size every frame (initial size is 0x0 before window is ready)
	sw, sh := p.GetSurfaceSize()
	if sw > 0 && sh > 0 {
		p.UpdateScreenSize(sw, sh)
	}

	encoder.SetPipeline(p.pipeline)
	encoder.SetBindGroup(0, p.bindGroup, nil)

	// Draw 6 vertices per primitive (2 triangles)
	// No vertex buffer - vertices constructed in shader from storage buffer
	vertexCount := uint32(len(p.primitives) * 6)
	encoder.Draw(vertexCount, 1, 0, 0)
}

func (p *PrimitiveBuffer) Render() {
	// Rendering handled by RenderPass
}

func (p *PrimitiveBuffer) Dispose() {
	if p.storageBuffer != nil {
		p.storageBuffer.Release()
		p.storageBuffer = nil
	}
	if p.bindGroup != nil {
		p.bindGroup.Release()
		p.bindGroup = nil
	}
	if p.screenSizeBuffer != nil {
		p.screenSizeBuffer.Release()
		p.screenSizeBuffer = nil
	}
	if p.msdfParamsBuffer != nil {
		p.msdfParamsBuffer.Release()
		p.msdfParamsBuffer = nil
	}
	if p.msdfSampler != nil {
		p.msdfSampler.Release()
		p.msdfSampler = nil
	}
	if p.msdfTextureView != nil {
		p.msdfTextureView.Release()
		p.msdfTextureView = nil
	}
	if p.msdfTexture != nil {
		p.msdfTexture.Release()
		p.msdfTexture = nil
	}
	if p.bindGroupLayout != nil {
		p.bindGroupLayout.Release()
		p.bindGroupLayout = nil
	}
	p.isDisposed = true
}

func (p *PrimitiveBuffer) IsDisposed() bool {
	return p.isDisposed
}

// Legacy compatibility - convert PrimitiveVertex to storage buffer approach
func (p *PrimitiveBuffer) UpdateVertexBuffer(vertices []graphics.PrimitiveVertex) {
	// Get current screen size (don't rely on cached value which might be stale)
	sw, sh := p.GetSurfaceSize()
	if sw == 0 || sh == 0 {
		// Fallback to cached value if surface not yet initialized
		sw, sh = int(p.screenSize[0]), int(p.screenSize[1])
	} else {
		// Update cached screen size
		p.screenSize[0] = float32(sw)
		p.screenSize[1] = float32(sh)
	}

	// Convert to new Primitive format
	primitives := convertVerticesToPrimitives(vertices, float32(sw), float32(sh))
	p.UpdatePrimitives(primitives)
}

var snapMSDFToPixels bool

// EnableSnapMSDFToPixels controls whether MSDF primitives are snapped to integer pixels.
// This is a debugging aid to determine if subpixel placement is a primary cause of fuzzy text.
func (p *PrimitiveBuffer) EnableSnapMSDFToPixels(enable bool) {
	snapMSDFToPixels = enable
}

func convertVerticesToPrimitives(vertices []graphics.PrimitiveVertex, screenW, screenH float32) []graphics.Primitive {
	primitives := make([]graphics.Primitive, 0, len(vertices)/6)

	for i := 0; i+6 <= len(vertices); i += 6 {
		v := vertices[i : i+6]

		// Find min/max screen positions from all 6 vertices
		// This handles different vertex orderings across op codes
		minX, minY := float32(1e9), float32(1e9)
		maxX, maxY := float32(-1e9), float32(-1e9)
		for j := 0; j < 6; j++ {
			// NDC to screen
			sx := (v[j].Position[0] + 1) / 2 * screenW
			sy := (1 - v[j].Position[1]) / 2 * screenH
			if sx < minX {
				minX = sx
			}
			if sx > maxX {
				maxX = sx
			}
			if sy < minY {
				minY = sy
			}
			if sy > maxY {
				maxY = sy
			}
		}

		w := maxX - minX
		h := maxY - minY
		opCode := v[0].OpCode

		// Skip primitives with invalid sizes (e.g., space characters in text, or uninitialized data)
		if w <= 0 || h <= 0 || w > 10000 || h > 10000 ||
			minX < -10000 || maxX > 10000 || minY < -10000 || maxY > 10000 {
			continue
		}

		var extra [4]float32

		if opCode == graphics.OpCodeMSDF {
			// For MSDF: find UV bounds from vertices
			// UV coordinates vary per vertex, find min/max
			minU, minV := float32(1e9), float32(1e9)
			maxU, maxV := float32(-1e9), float32(-1e9)
			for j := 0; j < 6; j++ {
				u, uv := v[j].TexCoords[0], v[j].TexCoords[1]
				if u < minU {
					minU = u
				}
				if u > maxU {
					maxU = u
				}
				if uv < minV {
					minV = uv
				}
				if uv > maxV {
					maxV = uv
				}
			}
			// Store UV base (top-left) and UV size
			extra = [4]float32{minU, minV, maxU - minU, maxV - minV}
		} else {
			// For shapes: store half_size from TexCoords (all vertices have same value)
			extra = [4]float32{v[0].TexCoords[0], v[0].TexCoords[1], 0, 0}
		}

		if snapMSDFToPixels && opCode == graphics.OpCodeMSDF {
			minX = float32(int(minX + 0.5))
			minY = float32(int(minY + 0.5))
			maxX = float32(int(maxX + 0.5))
			maxY = float32(int(maxY + 0.5))
			w = maxX - minX
			h = maxY - minY
			if w <= 0 || h <= 0 {
				continue
			}
		}

		primitives = append(primitives, graphics.Primitive{
			X:      minX,
			Y:      minY,
			W:      w,
			H:      h,
			Color:  v[0].Color,
			Radius: v[0].Radius,
			OpCode: opCode,
			Extra:  extra,
		})
	}

	return primitives
}
