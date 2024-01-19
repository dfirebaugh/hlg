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
