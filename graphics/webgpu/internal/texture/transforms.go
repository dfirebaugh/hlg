package texture

import (
	"log"
)

func (t *Texture) SetClipRect(minX, minY, maxX, maxY float32) {
	t.Transform.SetClipRect(minX, minY, maxX, maxY)
	if err := t.updateVertexBuffer(); err != nil {
		log.Printf("Failed to update vertex buffer: %v", err)
	}
}

func (t *Texture) Move(screenX float32, screenY float32) {
	t.MoveToScreenPosition(screenX, screenY)
}
