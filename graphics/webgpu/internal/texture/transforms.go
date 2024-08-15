package texture

import "github.com/rajveermalviya/go-webgpu/wgpu"

func (t *Texture) FlipVertical() {
	t.flipVertical = !t.flipVertical
	t.updateFlipBuffer()
}

func (t *Texture) FlipHorizontal() {
	t.flipHorizontal = !t.flipHorizontal
	t.updateFlipBuffer()
}

func (t *Texture) SetDefaultClip() {
	t.SetClipRect(0, 0, t.originalWidth, t.originalHeight)
}

func (t *Texture) SetClipRect(minX, minY, maxX, maxY float32) {
	t.clipRect = [4]float32{
		minX / t.originalWidth,
		minY / t.originalHeight,
		maxX / t.originalWidth,
		maxY / t.originalHeight,
	}

	t.Device.GetQueue().WriteBuffer(t.clipBuffer, 0, wgpu.ToBytes(t.clipRect[:]))
}

func (t *Texture) GetCurrentSize() (float32, float32) {
	scaleX := t.Transform.Matrix[0]
	scaleY := t.Transform.Matrix[5]

	currentWidth := t.originalWidth * scaleX
	currentHeight := t.originalHeight * scaleY

	return currentWidth, currentHeight
}

func (t *Texture) GetCurrentClipSize() (float32, float32) {
	scaleX := t.Transform.Matrix[0]
	scaleY := t.Transform.Matrix[5]

	clipWidth := ((t.clipRect[2] * t.originalWidth) - (t.clipRect[0] * t.originalWidth)) * scaleX
	clipHeight := ((t.clipRect[3] * t.originalHeight) - (t.clipRect[1] * t.originalHeight)) * scaleY

	return clipWidth, clipHeight
}
