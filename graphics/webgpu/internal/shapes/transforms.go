package shapes

import (
	"github.com/dfirebaugh/hlg/pkg/math/matrix"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func (p *Polygon) Move(destX, destY float32) {
	for i := range p.vertices {
		p.vertices[i].Position[0] += destX - p.center[0]
		p.vertices[i].Position[1] += destY - p.center[1]
	}

	p.vertexBuffer.Release()
	p.vertexBuffer = createVertexBuffer(p.device, p.vertices, float32(p.SwapChainDescriptor.Width), float32(p.SwapChainDescriptor.Height))
	p.updateTransformBuffer()
	p.center = [2]float32{destX, destY}
}

func (p *Polygon) Rotate(angle float32) {
	p.transform = p.transform.Rotate(angle)
	p.updateTransformBuffer()
}

func (p *Polygon) Scale(sx, sy float32) {
	p.transform = p.transform.Scale(sx, sy)
	p.updateTransformBuffer()
}

func (p *Polygon) Resize(newWidth, newHeight float32) {
	minX, maxX, minY, maxY := p.vertices[0].Position[0], p.vertices[0].Position[0], p.vertices[0].Position[1], p.vertices[0].Position[1]
	for _, vertex := range p.vertices[1:] {
		if vertex.Position[0] < minX {
			minX = vertex.Position[0]
		}
		if vertex.Position[0] > maxX {
			maxX = vertex.Position[0]
		}
		if vertex.Position[1] < minY {
			minY = vertex.Position[1]
		}
		if vertex.Position[1] > maxY {
			maxY = vertex.Position[1]
		}
	}

	currentWidth := maxX - minX
	currentHeight := maxY - minY

	scaleX := newWidth / currentWidth
	scaleY := newHeight / currentHeight

	centerX := (minX + maxX) / 2
	centerY := (minY + maxY) / 2

	for i := range p.vertices {
		dx := (p.vertices[i].Position[0] - centerX) * scaleX
		dy := (p.vertices[i].Position[1] - centerY) * scaleY
		p.vertices[i].Position[0] = centerX + dx
		p.vertices[i].Position[1] = centerY + dy
	}

	p.vertexBuffer.Release()
	p.vertexBuffer = createVertexBuffer(p.device, p.vertices, float32(p.SwapChainDescriptor.Width), float32(p.SwapChainDescriptor.Height))

	p.updateTransformBuffer()
}

func (p *Polygon) updateTransformBuffer() {
	if p.transform.IsZero() {
		return
	}
	p.device.GetQueue().WriteBuffer(p.transformBuffer, 0, wgpu.ToBytes(p.transform[:]))
}

func (p *Polygon) createTransformBuffer() {
	var err error
	p.transform = matrix.MatrixIdentity()
	p.transformBuffer, err = p.device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Polygon Transform Buffer",
		Usage:    wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
		Contents: wgpu.ToBytes(p.transform[:]),
	})
	if err != nil {
		panic(err)
	}
	p.updateTransformBuffer()
}
