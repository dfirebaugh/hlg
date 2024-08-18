package shapes


func (p *Polygon) Move(destX, destY float32) {
	p.RecenterAndMove(p.vertices, destX, destY)
	p.vertexBuffer.Release()
	p.createVertexBuffer()
}
