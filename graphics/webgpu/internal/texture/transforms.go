package texture

import "github.com/rajveermalviya/go-webgpu/wgpu"

// ScreenToNDC transforms screen space coordinates to NDC.
// screenWidth and screenHeight are the dimensions of the screen.
func ScreenToNDC(x, y, screenWidth, screenHeight float32) [3]float32 {
	// Normalize coordinates to [0, 1]
	normalizedX := x / screenWidth
	normalizedY := y / screenHeight

	// Map to NDC [-1, 1]
	ndcX := normalizedX*2 - 1
	ndcY := 1 - normalizedY*2 // Y is inverted in NDC

	return [3]float32{ndcX, ndcY, 0} // Assuming Z coordinate to be 0 for 2D
}

func (t *Texture) Resize(width, height float32) {
	t.ResizeInScreenSpace(width, height)
}

func (t *Texture) ResizeInScreenSpace(screenWidth, screenHeight float32) {
	scaleX := screenWidth / t.originalWidth
	scaleY := screenHeight / t.originalHeight
	uniformScale := min(scaleX, scaleY)

	t.transform = t.transform.Scale(uniformScale, uniformScale)

	t.updateTransformBuffer()
}

func (t *Texture) Move(dx, dy float32) {
	t.MoveInScreenSpace(dx, dy)
}

func (t *Texture) MoveInScreenSpace(screenX, screenY float32) {
	ndcCoords := ScreenToNDC(screenX, screenY, t.originalScreenWidth, t.originalScreenHeight)

	t.transform[12] = ndcCoords[0]
	t.transform[13] = ndcCoords[1]

	t.updateTransformBuffer()
}

func (t *Texture) Rotate(a float32) {
	t.transform = t.transform.Rotate(a)
	t.updateTransformBuffer()
}

func (t *Texture) Scale(x, y float32) {
	t.transform = t.transform.Scale(x, y)
	t.updateTransformBuffer()
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
