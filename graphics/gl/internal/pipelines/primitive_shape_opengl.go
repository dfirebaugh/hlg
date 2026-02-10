//go:build !js

package pipelines

import (
	"image/color"
	"math"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/glapi"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/shader"
	"github.com/dfirebaugh/hlg/pkg/math/matrix"
)

// PrimitiveShape implements graphics.Shape using the primitive buffer pipeline.
// It stores shapes as triangles with OpCodeSolid for simple fill rendering.
type PrimitiveShape struct {
	pb  *PrimitiveBuffer
	ctx *glapi.Context

	vao glapi.VertexArray
	vbo glapi.Buffer

	vertices []graphics.PrimitiveVertex

	// Store screen-space positions for transforms
	screenPositions [][2]float32
	// Store per-vertex colors to preserve them during rebuild
	colors       [][4]float32
	screenWidth  int
	screenHeight int

	shouldRender bool
	isDisposed   bool

	// Clip rect captured when Render() is called
	clipRect *[4]int
}

// NewPrimitiveShape creates a new shape that renders using the primitive buffer pipeline.
func NewPrimitiveShape(rq RenderQueue, vertices []graphics.PrimitiveVertex, screenPositions [][2]float32) *PrimitiveShape {
	if rq == nil {
		return &PrimitiveShape{isDisposed: true}
	}

	// Get primitive buffer from provider
	var pb *PrimitiveBuffer
	if provider, ok := rq.(PrimitiveBufferProvider); ok {
		pb = provider.GetPrimitiveBuffer()
	}
	if pb == nil {
		return &PrimitiveShape{isDisposed: true}
	}

	sw, sh := pb.GetSurfaceSize()

	p := &PrimitiveShape{
		pb:              pb,
		ctx:             pb.GetGL(),
		vertices:        vertices,
		screenPositions: screenPositions,
		screenWidth:     sw,
		screenHeight:    sh,
	}

	// Extract per-vertex colors to preserve them during rebuild
	p.colors = make([][4]float32, len(vertices))
	for i, v := range vertices {
		p.colors[i] = v.Color
	}

	p.createVertexBuffer()

	return p
}

func (p *PrimitiveShape) createVertexBuffer() {
	if len(p.vertices) == 0 {
		return
	}

	p.vao = p.ctx.CreateVertexArray()
	p.ctx.BindVertexArray(p.vao)

	p.vbo = p.ctx.CreateBuffer()
	p.ctx.BindBuffer(glapi.ARRAY_BUFFER, p.vbo)

	data := verticesToBytes(p.vertices)
	p.ctx.BufferData(glapi.ARRAY_BUFFER, data, glapi.DYNAMIC_DRAW)

	// PrimitiveVertex layout: Position[3], LocalPosition[2], OpCode, Radius, Color[4], TexCoords[2], HalfSize[2]
	stride := 60

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

func (p *PrimitiveShape) updateVertexBuffer() {
	if p.vbo == 0 || len(p.vertices) == 0 {
		return
	}
	p.ctx.BindBuffer(glapi.ARRAY_BUFFER, p.vbo)
	data := verticesToBytes(p.vertices)
	p.ctx.BufferSubData(glapi.ARRAY_BUFFER, 0, data)
}

func (p *PrimitiveShape) rebuildVertices() {
	sw := float32(p.screenWidth)
	sh := float32(p.screenHeight)

	for i := range p.vertices {
		if i < len(p.screenPositions) {
			p.vertices[i].Position = screenToNDC(p.screenPositions[i][0], p.screenPositions[i][1], sw, sh)
		}
		if i < len(p.colors) {
			p.vertices[i].Color = p.colors[i]
		}
	}
	p.updateVertexBuffer()
}

func (p *PrimitiveShape) GLRender() {
	if !p.shouldRender || p.isDisposed {
		return
	}
	if p.vao == 0 || len(p.vertices) == 0 {
		return
	}

	// Disable sRGB framebuffer for shapes - colors are already in sRGB space
	p.ctx.Disable(glapi.FRAMEBUFFER_SRGB)

	// Apply scissor if clip rect is set
	if p.clipRect != nil {
		// Get framebuffer size from viewport for HiDPI scaling
		viewport := p.ctx.GetViewport()
		fbWidth, fbHeight := viewport[2], viewport[3]

		screenW, screenH := p.pb.GetSurfaceSize()
		if screenW > 0 && screenH > 0 && fbWidth > 0 && fbHeight > 0 {
			// Scale from logical screen coordinates to framebuffer pixels
			scaleX := float64(fbWidth) / float64(screenW)
			scaleY := float64(fbHeight) / float64(screenH)

			x, y, w, h := p.clipRect[0], p.clipRect[1], p.clipRect[2], p.clipRect[3]

			scaledX := int(float64(x) * scaleX)
			scaledY := int(float64(y) * scaleY)
			scaledW := int(float64(w) * scaleX)
			scaledH := int(float64(h) * scaleY)

			// Flip Y: GL scissor origin is bottom-left, our API uses top-left
			glY := fbHeight - scaledY - scaledH

			p.ctx.Enable(glapi.SCISSOR_TEST)
			p.ctx.Scissor(scaledX, glY, scaledW, scaledH)
		}
	}

	program := p.pb.GetShaderManager().GetProgram(shader.PrimitiveBufferShader)
	p.ctx.UseProgram(program)

	// Set screen size uniform
	sw, sh := p.pb.GetSurfaceSize()
	screenSizeLoc := p.ctx.GetUniformLocation(program, "u_screen_size")
	p.ctx.Uniform2f(screenSizeLoc, float32(sw), float32(sh))

	// Set MSDF params (needed for shader even if not using MSDF)
	msdfParamsLoc := p.ctx.GetUniformLocation(program, "u_msdf_params")
	p.ctx.Uniform4fv(msdfParamsLoc, []float32{4.0, 1.0, 1.0, 0.0})

	p.ctx.BindVertexArray(p.vao)
	p.ctx.DrawArrays(glapi.TRIANGLES, 0, len(p.vertices))
	p.ctx.UnbindVertexArray()

	// Disable scissor after rendering
	if p.clipRect != nil {
		p.ctx.Disable(glapi.SCISSOR_TEST)
	}
}

func (p *PrimitiveShape) Render() {
	if p.isDisposed {
		return
	}
	p.shouldRender = true
	// Capture clip rect at time of Render() call
	if rq, ok := p.pb.surface.(RenderQueue); ok {
		p.clipRect = rq.GetCurrentClipRect()
		rq.AddToRenderQueue(p)
	}
}

func (p *PrimitiveShape) Dispose() {
	if p.isDisposed {
		return
	}
	p.isDisposed = true

	if p.vao != 0 {
		p.ctx.DeleteVertexArray(p.vao)
		p.vao = 0
	}
	if p.vbo != 0 {
		p.ctx.DeleteBuffer(p.vbo)
		p.vbo = 0
	}
}

func (p *PrimitiveShape) IsDisposed() bool {
	return p.isDisposed
}

func (p *PrimitiveShape) Hide() {
	p.shouldRender = false
}

func (p *PrimitiveShape) SetColor(c color.Color) {
	r, g, b, a := c.RGBA()
	newColor := [4]float32{
		float32(r) / 0xffff,
		float32(g) / 0xffff,
		float32(b) / 0xffff,
		float32(a) / 0xffff,
	}

	for i := range p.vertices {
		p.vertices[i].Color = newColor
		if i < len(p.colors) {
			p.colors[i] = newColor
		}
	}
	p.updateVertexBuffer()
}

func (p *PrimitiveShape) Move(destX, destY float32) {
	center := p.calculateCenter()

	dx := destX - center[0]
	dy := destY - center[1]

	for i := range p.screenPositions {
		p.screenPositions[i][0] += dx
		p.screenPositions[i][1] += dy
	}

	p.rebuildVertices()
}

func (p *PrimitiveShape) Rotate(angle float32) matrix.Matrix {
	center := p.calculateCenter()
	cos := float32(math.Cos(float64(angle)))
	sin := float32(math.Sin(float64(angle)))

	for i := range p.screenPositions {
		x := p.screenPositions[i][0] - center[0]
		y := p.screenPositions[i][1] - center[1]

		newX := x*cos - y*sin
		newY := x*sin + y*cos

		p.screenPositions[i][0] = newX + center[0]
		p.screenPositions[i][1] = newY + center[1]
	}

	p.rebuildVertices()
	return matrix.MatrixIdentity()
}

func (p *PrimitiveShape) Scale(sx, sy float32) matrix.Matrix {
	center := p.calculateCenter()

	for i := range p.screenPositions {
		x := p.screenPositions[i][0] - center[0]
		y := p.screenPositions[i][1] - center[1]

		x *= sx
		y *= sy

		p.screenPositions[i][0] = x + center[0]
		p.screenPositions[i][1] = y + center[1]
	}

	p.rebuildVertices()
	return matrix.MatrixIdentity()
}

func (p *PrimitiveShape) calculateCenter() [2]float32 {
	if len(p.screenPositions) == 0 {
		return [2]float32{0, 0}
	}

	var sumX, sumY float32
	for _, pos := range p.screenPositions {
		sumX += pos[0]
		sumY += pos[1]
	}
	count := float32(len(p.screenPositions))
	return [2]float32{sumX / count, sumY / count}
}

func screenToNDC(x, y, screenWidth, screenHeight float32) [3]float32 {
	ndcX := (x/screenWidth)*2 - 1
	ndcY := 1 - (y/screenHeight)*2
	return [3]float32{ndcX, ndcY, 0}
}

// Ensure PrimitiveShape implements graphics.Shape
var _ graphics.Shape = (*PrimitiveShape)(nil)
