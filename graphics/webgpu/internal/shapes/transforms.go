package shapes

import "github.com/dfirebaugh/hlg/graphics/webgpu/internal/common"

func (p *Polygon) Move(destX, destY float32) {
	center := common.CalculateCenter(p.vertices)
	for i := range p.vertices {
		p.vertices[i].Position[0] += destX - center[0]
		p.vertices[i].Position[1] += destY - center[1]
	}

	p.vertexBuffer.Release()
	p.createVertexBuffer()
}
