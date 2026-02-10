//go:build !js

package renderer

import (
	"image"
	"image/draw"
	"math"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/glapi"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/shader"
	"github.com/dfirebaugh/hlg/pkg/math/matrix"
)

// Texture represents an OpenGL texture
type Texture struct {
	rq *RenderQueue

	textureID glapi.Texture
	vao       glapi.VertexArray
	vbo       glapi.Buffer
	ebo       glapi.Buffer

	originalWidth  float32
	originalHeight float32

	screenWidth, screenHeight int
	x, y                      float32
	scaleX, scaleY            float32

	transform matrix.Matrix
	flipInfo  [2]float32
	clipRect  [4]float32

	handle     textureHandle
	isDisposed bool

	shouldRender bool

	// Scissor clip rect captured when Render() is called
	scissorClipRect *[4]int
}

var textureIndices = []uint16{
	0, 1, 2,
	2, 1, 3,
}

// NewTexture creates a new texture from an image
func NewTexture(rq *RenderQueue, img image.Image) *Texture {
	r := img.Bounds()
	width := r.Dx()
	height := r.Dy()

	t := &Texture{
		rq:             rq,
		originalWidth:  float32(width),
		originalHeight: float32(height),
		scaleX:         1.0,
		scaleY:         1.0,
		transform:      matrix.MatrixIdentity(),
		clipRect:       [4]float32{0, 0, 1, 1},
	}

	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		rgbaImg = image.NewRGBA(r)
		draw.Draw(rgbaImg, r, img, image.Point{}, draw.Over)
	}

	ctx := rq.GetGL()

	t.textureID = ctx.CreateTexture()
	ctx.BindTexture(glapi.TEXTURE_2D, t.textureID)
	ctx.TexParameteri(glapi.TEXTURE_2D, glapi.TEXTURE_MIN_FILTER, glapi.NEAREST)
	ctx.TexParameteri(glapi.TEXTURE_2D, glapi.TEXTURE_MAG_FILTER, glapi.NEAREST)
	ctx.TexParameteri(glapi.TEXTURE_2D, glapi.TEXTURE_WRAP_S, glapi.CLAMP_TO_EDGE)
	ctx.TexParameteri(glapi.TEXTURE_2D, glapi.TEXTURE_WRAP_T, glapi.CLAMP_TO_EDGE)
	ctx.TexImage2D(glapi.TEXTURE_2D, 0, glapi.RGBA, width, height, 0, glapi.RGBA, glapi.UNSIGNED_BYTE, rgbaImg.Pix)

	t.createVertexBuffer()
	t.updateTransform()

	return t
}

func (t *Texture) createVertexBuffer() {
	clipWidth := (t.clipRect[2] - t.clipRect[0]) * t.originalWidth * t.scaleX
	clipHeight := (t.clipRect[3] - t.clipRect[1]) * t.originalHeight * t.scaleY

	x0 := t.x
	y0 := t.y
	x1 := t.x + clipWidth
	y1 := t.y + clipHeight

	// Vertices: position (3) + color (4) + tex_coords (2) = 9 floats per vertex
	vertices := []float32{
		// bottom-left
		x0, y1, 0, 1, 1, 1, 1, 0, 1,
		// bottom-right
		x1, y1, 0, 1, 1, 1, 1, 1, 1,
		// top-left
		x0, y0, 0, 1, 1, 1, 1, 0, 0,
		// top-right
		x1, y0, 0, 1, 1, 1, 1, 1, 0,
	}

	ctx := t.rq.GetGL()

	// Create VAO
	t.vao = ctx.CreateVertexArray()
	ctx.BindVertexArray(t.vao)

	// Create VBO
	t.vbo = ctx.CreateBuffer()
	ctx.BindBuffer(glapi.ARRAY_BUFFER, t.vbo)
	ctx.BufferData(glapi.ARRAY_BUFFER, float32SliceToBytes(vertices), glapi.DYNAMIC_DRAW)

	// Create EBO
	t.ebo = ctx.CreateBuffer()
	ctx.BindBuffer(glapi.ELEMENT_ARRAY_BUFFER, t.ebo)
	ctx.BufferData(glapi.ELEMENT_ARRAY_BUFFER, uint16SliceToBytes(textureIndices), glapi.STATIC_DRAW)

	// Vertex attributes
	stride := 9 * 4 // 9 floats * 4 bytes
	// Position
	ctx.VertexAttribPointer(0, 3, glapi.FLOAT, false, stride, 0)
	ctx.EnableVertexAttribArray(0)
	// Color
	ctx.VertexAttribPointer(1, 4, glapi.FLOAT, false, stride, 3*4)
	ctx.EnableVertexAttribArray(1)
	// TexCoords
	ctx.VertexAttribPointer(2, 2, glapi.FLOAT, false, stride, 7*4)
	ctx.EnableVertexAttribArray(2)

	ctx.UnbindVertexArray()
}

func (t *Texture) updateVertexBuffer() {
	clipWidth := (t.clipRect[2] - t.clipRect[0]) * t.originalWidth * t.scaleX
	clipHeight := (t.clipRect[3] - t.clipRect[1]) * t.originalHeight * t.scaleY

	x0 := t.x
	y0 := t.y
	x1 := t.x + clipWidth
	y1 := t.y + clipHeight

	vertices := []float32{
		x0, y1, 0, 1, 1, 1, 1, 0, 1,
		x1, y1, 0, 1, 1, 1, 1, 1, 1,
		x0, y0, 0, 1, 1, 1, 1, 0, 0,
		x1, y0, 0, 1, 1, 1, 1, 1, 0,
	}

	ctx := t.rq.GetGL()
	ctx.BindBuffer(glapi.ARRAY_BUFFER, t.vbo)
	ctx.BufferSubData(glapi.ARRAY_BUFFER, 0, float32SliceToBytes(vertices))
}

// Handle returns the texture handle
func (t *Texture) Handle() uintptr {
	return uintptr(t.handle)
}

// SetHandle sets the texture handle
func (t *Texture) SetHandle(h textureHandle) {
	t.handle = h
}

// UpdateImage updates the texture with a new image
func (t *Texture) UpdateImage(img image.Image) error {
	r := img.Bounds()
	width := r.Dx()
	height := r.Dy()

	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		rgbaImg = image.NewRGBA(r)
		draw.Draw(rgbaImg, r, img, image.Point{}, draw.Over)
	}

	ctx := t.rq.GetGL()

	if int(t.originalWidth) != width || int(t.originalHeight) != height {
		t.originalWidth = float32(width)
		t.originalHeight = float32(height)

		ctx.BindTexture(glapi.TEXTURE_2D, t.textureID)
		ctx.TexImage2D(glapi.TEXTURE_2D, 0, glapi.RGBA, width, height, 0, glapi.RGBA, glapi.UNSIGNED_BYTE, rgbaImg.Pix)
	} else {
		ctx.BindTexture(glapi.TEXTURE_2D, t.textureID)
		ctx.TexSubImage2D(glapi.TEXTURE_2D, 0, 0, 0, width, height, glapi.RGBA, glapi.UNSIGNED_BYTE, rgbaImg.Pix)
	}

	return nil
}

// SetShouldBeRendered sets whether the texture should be rendered
func (t *Texture) SetShouldBeRendered(shouldRender bool) {
	t.shouldRender = shouldRender
}

// Resize resizes the texture to fit specific dimensions
func (t *Texture) Resize(width, height float32) {
	clipWidth := (t.clipRect[2] - t.clipRect[0]) * t.originalWidth
	clipHeight := (t.clipRect[3] - t.clipRect[1]) * t.originalHeight

	t.scaleX = width / clipWidth
	t.scaleY = height / clipHeight
	t.updateVertexBuffer()
	t.updateTransform()
}

// Move moves the texture to a new position
func (t *Texture) Move(x, y float32) {
	t.x = x
	t.y = y
	t.updateVertexBuffer()
}

// Rotate rotates the texture (currently not implemented)
func (t *Texture) Rotate(a, pivotX, pivotY float32) {
	// Rotation transform
}

// Scale scales the texture
func (t *Texture) Scale(x, y float32) {
	t.scaleX *= x
	t.scaleY *= y
	t.updateVertexBuffer()
	t.updateTransform()
}

// FlipVertical flips the texture vertically
func (t *Texture) FlipVertical() {
	t.flipInfo[1] = 1.0 - t.flipInfo[1]
}

// FlipHorizontal flips the texture horizontally
func (t *Texture) FlipHorizontal() {
	t.flipInfo[0] = 1.0 - t.flipInfo[0]
}

// SetFlipHorizontal sets horizontal flip state
func (t *Texture) SetFlipHorizontal(shouldFlip bool) {
	if shouldFlip {
		t.flipInfo[0] = 1.0
	} else {
		t.flipInfo[0] = 0.0
	}
}

// SetFlipVertical sets vertical flip state
func (t *Texture) SetFlipVertical(shouldFlip bool) {
	if shouldFlip {
		t.flipInfo[1] = 1.0
	} else {
		t.flipInfo[1] = 0.0
	}
}

// Clip sets the clip rectangle in pixel coordinates
func (t *Texture) Clip(minX, minY, maxX, maxY float32) {
	t.clipRect = [4]float32{
		minX / t.originalWidth,
		minY / t.originalHeight,
		maxX / t.originalWidth,
		maxY / t.originalHeight,
	}
	t.updateVertexBuffer()
}

// Render adds the texture to the render queue
func (t *Texture) Render() {
	if t.isDisposed {
		return
	}
	t.shouldRender = true
	// Capture clip rect at time of Render() call
	t.scissorClipRect = t.rq.GetCurrentClipRect()
	t.rq.AddToRenderQueue(t)
}

// RenderToQueue adds the texture to a specific render queue
func (t *Texture) RenderToQueue(rq graphics.RenderQueue) {
	if t.isDisposed {
		return
	}
	t.shouldRender = true
	rq.AddToRenderQueue(t)
}

// GLRender performs the OpenGL rendering
func (t *Texture) GLRender() {
	if !t.shouldRender || t.isDisposed {
		return
	}

	t.handleScreenResize()

	ctx := t.rq.GetGL()

	// Apply scissor if clip rect is set
	if t.scissorClipRect != nil {
		viewport := ctx.GetViewport()
		fbWidth, fbHeight := viewport[2], viewport[3]

		screenW, screenH := t.rq.GetSurfaceSize()
		if screenW > 0 && screenH > 0 && fbWidth > 0 && fbHeight > 0 {
			scaleX := float64(fbWidth) / float64(screenW)
			scaleY := float64(fbHeight) / float64(screenH)

			x, y, w, h := t.scissorClipRect[0], t.scissorClipRect[1], t.scissorClipRect[2], t.scissorClipRect[3]

			scaledX := int(float64(x) * scaleX)
			scaledY := int(float64(y) * scaleY)
			scaledW := int(float64(w) * scaleX)
			scaledH := int(float64(h) * scaleY)

			glY := fbHeight - scaledY - scaledH

			ctx.Enable(glapi.SCISSOR_TEST)
			ctx.Scissor(scaledX, glY, scaledW, scaledH)
		}
	}

	program := t.rq.ShaderManager.GetProgram(shader.TextureShader)
	ctx.UseProgram(program)

	// Set uniforms
	transformLoc := ctx.GetUniformLocation(program, "u_transform")
	ctx.UniformMatrix4fv(transformLoc, true, t.transform[:])

	flipLoc := ctx.GetUniformLocation(program, "u_flip_info")
	ctx.Uniform2f(flipLoc, t.flipInfo[0], t.flipInfo[1])

	clipLoc := ctx.GetUniformLocation(program, "u_clip_rect")
	ctx.Uniform4f(clipLoc, t.clipRect[0], t.clipRect[1], t.clipRect[2], t.clipRect[3])

	texLoc := ctx.GetUniformLocation(program, "u_texture")
	ctx.Uniform1i(texLoc, 0)

	ctx.ActiveTexture(glapi.TEXTURE0)
	ctx.BindTexture(glapi.TEXTURE_2D, t.textureID)

	ctx.BindVertexArray(t.vao)
	ctx.DrawElements(glapi.TRIANGLES, 6, glapi.UNSIGNED_SHORT, 0)
	ctx.UnbindVertexArray()

	// Disable scissor after rendering
	if t.scissorClipRect != nil {
		ctx.Disable(glapi.SCISSOR_TEST)
	}
}

func (t *Texture) handleScreenResize() {
	w, h := t.rq.GetSurfaceSize()
	if w == t.screenWidth && h == t.screenHeight {
		return
	}
	t.screenWidth = w
	t.screenHeight = h
	t.updateVertexBuffer()
	t.updateTransform()
}

func (t *Texture) updateTransform() {
	sw, sh := t.rq.GetSurfaceSize()
	swf, shf := float32(sw), float32(sh)

	t.transform = matrixOrtho(0, swf, shf, 0, -1, 1)
}

func matrixOrtho(left, right, bottom, top, near, far float32) matrix.Matrix {
	rml := right - left
	tmb := top - bottom
	fmn := far - near

	return matrix.Matrix{
		2 / rml, 0, 0, -(right + left) / rml,
		0, 2 / tmb, 0, -(top + bottom) / tmb,
		0, 0, -2 / fmn, -(far + near) / fmn,
		0, 0, 0, 1,
	}
}

// Dispose cleans up texture resources
func (t *Texture) Dispose() {
	if t.isDisposed {
		return
	}
	t.isDisposed = true

	ctx := t.rq.GetGL()

	ctx.DeleteTexture(t.textureID)
	ctx.DeleteVertexArray(t.vao)
	ctx.DeleteBuffer(t.vbo)
	ctx.DeleteBuffer(t.ebo)
}

// IsDisposed returns whether the texture has been disposed
func (t *Texture) IsDisposed() bool {
	return t.isDisposed
}

// Helper functions for converting slices to bytes
func float32SliceToBytes(data []float32) []byte {
	bytes := make([]byte, len(data)*4)
	for i, f := range data {
		u := math.Float32bits(f)
		bytes[i*4+0] = byte(u)
		bytes[i*4+1] = byte(u >> 8)
		bytes[i*4+2] = byte(u >> 16)
		bytes[i*4+3] = byte(u >> 24)
	}
	return bytes
}

func uint16SliceToBytes(data []uint16) []byte {
	bytes := make([]byte, len(data)*2)
	for i, v := range data {
		bytes[i*2+0] = byte(v)
		bytes[i*2+1] = byte(v >> 8)
	}
	return bytes
}

// Ensure Texture implements graphics.Texture
var _ graphics.Texture = (*Texture)(nil)
