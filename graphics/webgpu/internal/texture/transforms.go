package texture

import "github.com/rajveermalviya/go-webgpu/wgpu"

func (t *Texture) Resize(width, height float32) {
	t.ResizeInScreenSpace(width, height)
}

func (t *Texture) ResizeInScreenSpace(screenWidth, screenHeight float32) {
	ndcWidth := (2.0 * screenWidth) / float32(t.SwapChainDescriptor.Width)
	ndcHeight := (2.0 * screenHeight) / float32(t.SwapChainDescriptor.Height)

	scaleX := ndcWidth
	scaleY := ndcHeight

	t.transform = t.transform.Scale(scaleX, scaleY)
	t.UpdateTransformBuffer()
}

func (t *Texture) MoveInScreenSpace(screenX, screenY float32) {
	clipWidth, clipHeight := t.GetCurrentClipSize()

	ndcX := (2.0 * (screenX + clipWidth/2) / float32(t.SwapChainDescriptor.Width)) - 1.0
	ndcY := 1.0 - (2.0 * (screenY + clipHeight/2) / float32(t.SwapChainDescriptor.Height))

	t.transform[12] = ndcX
	t.transform[13] = ndcY

	t.UpdateTransformBuffer()
}

func (t *Texture) Rotate(a float32) {
	t.transform = t.transform.Rotate(a)
	t.UpdateTransformBuffer()
}

func (t *Texture) Scale(x, y float32) {
	t.transform = t.transform.Scale(x, y)
	t.UpdateTransformBuffer()
}

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
	scaleX := t.transform[0]
	scaleY := t.transform[5]

	currentWidth := t.originalWidth * scaleX
	currentHeight := t.originalHeight * scaleY

	return currentWidth, currentHeight
}

func (t *Texture) GetCurrentClipSize() (float32, float32) {
	scaleX := t.transform[0]
	scaleY := t.transform[5]

	clipWidth := ((t.clipRect[2] * t.originalWidth) - (t.clipRect[0] * t.originalWidth)) * scaleX
	clipHeight := ((t.clipRect[3] * t.originalHeight) - (t.clipRect[1] * t.originalHeight)) * scaleY

	return clipWidth, clipHeight
}
