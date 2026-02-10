//go:build !js

package pipelines

import (
	"image/color"
	"log"
	"math"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/pkg/math/matrix"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

// PrimitiveShape implements graphics.Shape using the primitive buffer pipeline.
// It stores shapes as triangles with OpCodeSolid for simple fill rendering.
type PrimitiveShape struct {
	pb *PrimitiveBuffer

	vertexBuffer *wgpu.Buffer
	vertices     []graphics.PrimitiveVertex

	// Store screen-space positions for transforms
	screenPositions [][2]float32
	// Store per-vertex colors to preserve them during rebuild
	colors       [][4]float32
	screenWidth  int
	screenHeight int

	shouldRender bool
	isDisposed   bool
}

// NewPrimitiveShape creates a new shape that renders using the primitive buffer pipeline.
// vertices should already be in NDC with OpCodeSolid set.
// screenPositions contains the original screen-space coordinates for transform calculations.
func NewPrimitiveShape(pb *PrimitiveBuffer, vertices []graphics.PrimitiveVertex, screenPositions [][2]float32) *PrimitiveShape {
	if pb == nil {
		log.Fatal("PrimitiveBuffer is nil")
	}

	sw, sh := pb.GetSurfaceSize()

	p := &PrimitiveShape{
		pb:              pb,
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

	var err error
	p.vertexBuffer, err = p.pb.GetDevice().CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "PrimitiveShape Vertex Buffer",
		Contents: wgpu.ToBytes(p.vertices[:]),
		Usage:    wgpu.BufferUsage_Vertex | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		panic(err)
	}
}

func (p *PrimitiveShape) updateVertexBuffer() {
	if p.vertexBuffer == nil || len(p.vertices) == 0 {
		return
	}
	_ = p.pb.GetDevice().GetQueue().WriteBuffer(p.vertexBuffer, 0, wgpu.ToBytes(p.vertices[:]))
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

func (p *PrimitiveShape) handleScreenResize() {
	// Keep using original screen dimensions for NDC conversion
	// This ensures shapes stay proportionally placed when window resizes
	// (matching the behavior of gui.Draw which uses fixed dimensions)
}

func (p *PrimitiveShape) RenderPass(encoder *wgpu.RenderPassEncoder) {
	if encoder == nil || !p.shouldRender || p.isDisposed {
		return
	}
	if p.vertexBuffer == nil || len(p.vertices) == 0 {
		return
	}

	p.handleScreenResize()

	// Use the solid shape pipeline (vertex buffer based)
	encoder.SetPipeline(p.pb.solidShapePipeline)
	encoder.SetVertexBuffer(0, p.vertexBuffer, 0, wgpu.WholeSize)

	vertexCount := uint32(len(p.vertices))
	encoder.Draw(vertexCount, 1, 0, 0)
}

func (p *PrimitiveShape) Render() {
	if p.isDisposed {
		return
	}
	p.shouldRender = true
	p.pb.AddToRenderQueue(p)
}

func (p *PrimitiveShape) Dispose() {
	if p.vertexBuffer != nil {
		p.vertexBuffer.Release()
		p.vertexBuffer = nil
	}
	p.isDisposed = true
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

	// Update all stored colors and vertex colors
	for i := range p.vertices {
		p.vertices[i].Color = newColor
		if i < len(p.colors) {
			p.colors[i] = newColor
		}
	}
	p.updateVertexBuffer()
}

// Move translates the shape so its center is at (destX, destY) in screen coordinates.
func (p *PrimitiveShape) Move(destX, destY float32) {
	// Calculate current center
	center := p.calculateCenter()

	// Calculate offset
	dx := destX - center[0]
	dy := destY - center[1]

	// Move all screen positions
	for i := range p.screenPositions {
		p.screenPositions[i][0] += dx
		p.screenPositions[i][1] += dy
	}

	p.rebuildVertices()
}

// Rotate rotates the shape around its center by the given angle (in radians).
func (p *PrimitiveShape) Rotate(angle float32) matrix.Matrix {
	center := p.calculateCenter()
	cos := float32(math.Cos(float64(angle)))
	sin := float32(math.Sin(float64(angle)))

	for i := range p.screenPositions {
		// Translate to origin
		x := p.screenPositions[i][0] - center[0]
		y := p.screenPositions[i][1] - center[1]

		// Rotate
		newX := x*cos - y*sin
		newY := x*sin + y*cos

		// Translate back
		p.screenPositions[i][0] = newX + center[0]
		p.screenPositions[i][1] = newY + center[1]
	}

	p.rebuildVertices()
	return matrix.MatrixIdentity()
}

// Scale scales the shape from its center by the given factors.
func (p *PrimitiveShape) Scale(sx, sy float32) matrix.Matrix {
	center := p.calculateCenter()

	for i := range p.screenPositions {
		// Translate to origin
		x := p.screenPositions[i][0] - center[0]
		y := p.screenPositions[i][1] - center[1]

		// Scale
		x *= sx
		y *= sy

		// Translate back
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

// screenToNDC converts screen coordinates to normalized device coordinates.
func screenToNDC(x, y, screenWidth, screenHeight float32) [3]float32 {
	ndcX := (x/screenWidth)*2 - 1
	ndcY := 1 - (y/screenHeight)*2
	return [3]float32{ndcX, ndcY, 0}
}
