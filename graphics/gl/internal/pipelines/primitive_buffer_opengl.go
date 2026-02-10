//go:build !js

package pipelines

import (
	"image"
	"image/draw"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/glapi"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/shader"
)

// PrimitiveBuffer implements the primitive buffer rendering using OpenGL
type PrimitiveBuffer struct {
	ctx           *glapi.Context
	surface       Surface
	shaderManager *shader.ShaderManager

	vao          glapi.VertexArray
	vbo          glapi.Buffer
	vertices     []graphics.PrimitiveVertex
	clipRectRuns []ClipRectRun
	verticesCap  int
	screenWidth  int
	screenHeight int

	// MSDF atlas resources
	msdfTextureID glapi.Texture
	msdfParams    [4]float32

	// Cached viewport for scissor calculations
	cachedFBWidth  int
	cachedFBHeight int

	isDisposed bool
}

func NewPrimitiveBuffer(ctx *glapi.Context, surface Surface, sm *shader.ShaderManager) *PrimitiveBuffer {
	sw, sh := surface.GetSurfaceSize()

	p := &PrimitiveBuffer{
		ctx:           ctx,
		surface:       surface,
		shaderManager: sm,
		verticesCap:   1024 * 6,
		msdfParams:    [4]float32{4.0, 1.0, 1.0, 0.0},
		screenWidth:   sw,
		screenHeight:  sh,
	}

	p.createVertexBuffer()
	p.createMSDFTexture()

	return p
}

func (p *PrimitiveBuffer) createVertexBuffer() {
	p.vao = p.ctx.CreateVertexArray()
	p.ctx.BindVertexArray(p.vao)

	p.vbo = p.ctx.CreateBuffer()
	p.ctx.BindBuffer(glapi.ARRAY_BUFFER, p.vbo)

	// PrimitiveVertex: Position[3], LocalPosition[2], OpCode, Radius, Color[4], TexCoords[2], HalfSize[2]
	// Total: 15 floats = 60 bytes
	stride := 60

	// Allocate buffer with initial capacity
	p.ctx.BufferDataSize(glapi.ARRAY_BUFFER, p.verticesCap*60, glapi.DYNAMIC_DRAW)

	// Position (location 0)
	p.ctx.VertexAttribPointer(0, 3, glapi.FLOAT, false, stride, 0)
	p.ctx.EnableVertexAttribArray(0)
	// LocalPosition (location 1)
	p.ctx.VertexAttribPointer(1, 2, glapi.FLOAT, false, stride, 12)
	p.ctx.EnableVertexAttribArray(1)
	// OpCode (location 2)
	p.ctx.VertexAttribPointer(2, 1, glapi.FLOAT, false, stride, 20)
	p.ctx.EnableVertexAttribArray(2)
	// Radius (location 3)
	p.ctx.VertexAttribPointer(3, 1, glapi.FLOAT, false, stride, 24)
	p.ctx.EnableVertexAttribArray(3)
	// Color (location 4)
	p.ctx.VertexAttribPointer(4, 4, glapi.FLOAT, false, stride, 28)
	p.ctx.EnableVertexAttribArray(4)
	// TexCoords (location 5)
	p.ctx.VertexAttribPointer(5, 2, glapi.FLOAT, false, stride, 44)
	p.ctx.EnableVertexAttribArray(5)
	// HalfSize (location 6)
	p.ctx.VertexAttribPointer(6, 2, glapi.FLOAT, false, stride, 52)
	p.ctx.EnableVertexAttribArray(6)

	p.ctx.UnbindVertexArray()
}

func (p *PrimitiveBuffer) createMSDFTexture() {
	// Create a 1x1 placeholder texture
	p.msdfTextureID = p.ctx.CreateTexture()
	p.ctx.BindTexture(glapi.TEXTURE_2D, p.msdfTextureID)
	p.ctx.TexParameteri(glapi.TEXTURE_2D, glapi.TEXTURE_MIN_FILTER, glapi.LINEAR)
	p.ctx.TexParameteri(glapi.TEXTURE_2D, glapi.TEXTURE_MAG_FILTER, glapi.LINEAR)
	p.ctx.TexParameteri(glapi.TEXTURE_2D, glapi.TEXTURE_WRAP_S, glapi.CLAMP_TO_EDGE)
	p.ctx.TexParameteri(glapi.TEXTURE_2D, glapi.TEXTURE_WRAP_T, glapi.CLAMP_TO_EDGE)

	placeholder := []byte{0, 0, 0, 0}
	p.ctx.TexImage2D(glapi.TEXTURE_2D, 0, glapi.RGBA, 1, 1, 0, glapi.RGBA, glapi.UNSIGNED_BYTE, placeholder)
}

func (p *PrimitiveBuffer) SetMSDFAtlas(atlasImg image.Image, pxRange float64) {
	r := atlasImg.Bounds()
	width := r.Dx()
	height := r.Dy()

	rgbaImg, ok := atlasImg.(*image.RGBA)
	if !ok {
		rgbaImg = image.NewRGBA(r)
		draw.Draw(rgbaImg, r, atlasImg, image.Point{}, draw.Over)
	}

	p.ctx.BindTexture(glapi.TEXTURE_2D, p.msdfTextureID)
	p.ctx.TexImage2D(glapi.TEXTURE_2D, 0, glapi.RGBA, width, height, 0, glapi.RGBA, glapi.UNSIGNED_BYTE, rgbaImg.Pix)

	p.msdfParams[0] = float32(pxRange)
	p.msdfParams[1] = float32(width)
	p.msdfParams[2] = float32(height)
}

func (p *PrimitiveBuffer) SetMSDFMode(mode int) {
	p.msdfParams[3] = float32(mode)
}

func (p *PrimitiveBuffer) EnableSnapMSDFToPixels(_ bool) {
	// Reserved for future use
}

func (p *PrimitiveBuffer) UpdateVertexBuffer(vertices []graphics.PrimitiveVertex) {
	p.clipRectRuns = p.clipRectRuns[:0]
	p.updateVertexBufferInternal(vertices)
}

// clipRectsEqual compares two clip rects by value (handles nil cases)
func clipRectsEqual(a, b *[4]int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func (p *PrimitiveBuffer) UpdateVertexBufferWithClipRects(vertices []graphics.PrimitiveVertex, clipRects []*[4]int) {
	if len(vertices) == 0 || len(clipRects) == 0 {
		p.clipRectRuns = p.clipRectRuns[:0]
		p.updateVertexBufferInternal(vertices)
		return
	}

	const vertsPerPrim = 6
	numPrims := len(vertices) / vertsPerPrim

	p.clipRectRuns = p.clipRectRuns[:0]

	// Build runs of CONSECUTIVE primitives with the same clip rect.
	// This preserves draw order while still batching consecutive primitives.
	var currentRun *ClipRectRun
	for i := range numPrims {
		vertIdx := i * vertsPerPrim
		var clipRect *[4]int
		if vertIdx < len(clipRects) {
			clipRect = clipRects[vertIdx]
		}

		// Check if this primitive continues the current run
		if currentRun != nil && clipRectsEqual(currentRun.ClipRect, clipRect) {
			currentRun.Count += vertsPerPrim
		} else {
			// Start a new run
			p.clipRectRuns = append(p.clipRectRuns, ClipRectRun{
				ClipRect: clipRect,
				StartIdx: vertIdx,
				Count:    vertsPerPrim,
			})
			currentRun = &p.clipRectRuns[len(p.clipRectRuns)-1]
		}
	}

	// No need to reorder vertices - they're already in the correct order
	p.updateVertexBufferInternal(vertices)
}

func (p *PrimitiveBuffer) updateVertexBufferInternal(vertices []graphics.PrimitiveVertex) {
	if len(vertices) == 0 {
		p.vertices = p.vertices[:0]
		return
	}

	// Copy vertices
	if cap(p.vertices) < len(vertices) {
		p.vertices = make([]graphics.PrimitiveVertex, len(vertices))
	} else {
		p.vertices = p.vertices[:len(vertices)]
	}
	copy(p.vertices, vertices)

	p.ctx.BindBuffer(glapi.ARRAY_BUFFER, p.vbo)

	// Grow capacity if needed
	if len(p.vertices) > p.verticesCap {
		p.verticesCap = len(p.vertices) * 2
	}

	// Orphan: signal driver we don't need old data
	p.ctx.BufferDataSize(glapi.ARRAY_BUFFER, p.verticesCap*60, glapi.DYNAMIC_DRAW)

	// Upload vertex data
	data := verticesToBytes(p.vertices)
	p.ctx.BufferSubData(glapi.ARRAY_BUFFER, 0, data)
}

func (p *PrimitiveBuffer) UpdatePrimitives(primitives []graphics.Primitive) {
	if len(primitives) == 0 {
		p.vertices = p.vertices[:0]
		return
	}

	sw, sh := p.surface.GetSurfaceSize()
	vertices := convertPrimitivesToVertices(primitives, float32(sw), float32(sh))
	p.UpdateVertexBuffer(vertices)
}

func (p *PrimitiveBuffer) UpdateScreenSize(width, height int) {
	p.screenWidth = width
	p.screenHeight = height
}

func (p *PrimitiveBuffer) GLRender() {
	if p.isDisposed || len(p.vertices) == 0 {
		return
	}

	// Disable sRGB framebuffer for primitives
	p.ctx.Disable(glapi.FRAMEBUFFER_SRGB)

	// Update screen size
	sw, sh := p.surface.GetSurfaceSize()
	if sw > 0 && sh > 0 {
		p.screenWidth = sw
		p.screenHeight = sh
	}

	program := p.shaderManager.GetProgram(shader.PrimitiveBufferShader)
	p.ctx.UseProgram(program)

	// Set uniforms
	screenSizeLoc := p.ctx.GetUniformLocation(program, "u_screen_size")
	p.ctx.Uniform2f(screenSizeLoc, float32(p.screenWidth), float32(p.screenHeight))

	msdfParamsLoc := p.ctx.GetUniformLocation(program, "u_msdf_params")
	p.ctx.Uniform4fv(msdfParamsLoc, p.msdfParams[:])

	msdfAtlasLoc := p.ctx.GetUniformLocation(program, "u_msdf_atlas")
	p.ctx.Uniform1i(msdfAtlasLoc, 0)

	p.ctx.ActiveTexture(glapi.TEXTURE0)
	p.ctx.BindTexture(glapi.TEXTURE_2D, p.msdfTextureID)

	p.ctx.BindVertexArray(p.vao)

	// Cache viewport size for scissor calculations
	viewport := p.ctx.GetViewport()
	p.cachedFBWidth, p.cachedFBHeight = viewport[2], viewport[3]

	// Render vertices grouped by clip rect
	p.renderWithClipRects()

	p.ctx.UnbindVertexArray()
	p.ctx.Disable(glapi.SCISSOR_TEST)
}

func (p *PrimitiveBuffer) renderWithClipRects() {
	if len(p.vertices) == 0 {
		return
	}

	p.ctx.Disable(glapi.SCISSOR_TEST)

	if len(p.clipRectRuns) == 0 {
		p.ctx.DrawArrays(glapi.TRIANGLES, 0, len(p.vertices))
		return
	}

	for _, run := range p.clipRectRuns {
		p.applyScissor(run.ClipRect)
		p.ctx.DrawArrays(glapi.TRIANGLES, run.StartIdx, run.Count)
	}
}

func (p *PrimitiveBuffer) applyScissor(clipRect *[4]int) {
	if clipRect == nil {
		p.ctx.Disable(glapi.SCISSOR_TEST)
		return
	}

	fbWidth, fbHeight := p.cachedFBWidth, p.cachedFBHeight
	screenW, screenH := p.surface.GetSurfaceSize()
	if screenW <= 0 || screenH <= 0 {
		p.ctx.Disable(glapi.SCISSOR_TEST)
		return
	}

	scaleX := float64(fbWidth) / float64(screenW)
	scaleY := float64(fbHeight) / float64(screenH)

	x, y, w, h := clipRect[0], clipRect[1], clipRect[2], clipRect[3]

	scaledX := int(float64(x) * scaleX)
	scaledY := int(float64(y) * scaleY)
	scaledW := int(float64(w) * scaleX)
	scaledH := int(float64(h) * scaleY)

	glY := fbHeight - scaledY - scaledH

	p.ctx.Enable(glapi.SCISSOR_TEST)
	p.ctx.Scissor(scaledX, glY, scaledW, scaledH)
}

func (p *PrimitiveBuffer) Render() {}

func (p *PrimitiveBuffer) Dispose() {
	if p.isDisposed {
		return
	}
	p.isDisposed = true

	p.ctx.DeleteVertexArray(p.vao)
	p.ctx.DeleteBuffer(p.vbo)
	p.ctx.DeleteTexture(p.msdfTextureID)
}

func (p *PrimitiveBuffer) IsDisposed() bool {
	return p.isDisposed
}

func (p *PrimitiveBuffer) GetSurfaceSize() (int, int) {
	return p.surface.GetSurfaceSize()
}

func (p *PrimitiveBuffer) GetVertices() []graphics.PrimitiveVertex {
	return p.vertices
}

func (p *PrimitiveBuffer) FlushImmediate() {
	if p.isDisposed || len(p.vertices) == 0 {
		return
	}

	p.ctx.Enable(glapi.BLEND)
	p.ctx.BlendFunc(glapi.SRC_ALPHA, glapi.ONE_MINUS_SRC_ALPHA)

	p.GLRender()

	p.vertices = p.vertices[:0]
}

func (p *PrimitiveBuffer) GetShaderManager() *shader.ShaderManager {
	return p.shaderManager
}

func (p *PrimitiveBuffer) GetGL() *glapi.Context {
	return p.ctx
}
