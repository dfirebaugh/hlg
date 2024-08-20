package transforms

import (
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/context"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/primitives"
	"github.com/dfirebaugh/hlg/pkg/math/matrix"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Transform struct {
	context.RenderContext

	*wgpu.Buffer

	matrix.Matrix
	label string

	originalWidth  float32
	originalHeight float32

	flipHorizontal bool
	flipVertical   bool
	FlipBuffer     *wgpu.Buffer
	FlipMatrix     [2]float32

	ClipRect   [4]float32 // minX, minY, maxX, maxY
	ClipBuffer *wgpu.Buffer
}

func NewTransform(ctx context.RenderContext, bufferLabel string, originalWidth, originalHeight float32) *Transform {
	t := &Transform{
		RenderContext:  ctx,
		Matrix:         matrix.MatrixIdentity(),
		label:          bufferLabel,
		originalWidth:  originalWidth,
		originalHeight: originalHeight,
	}

	t.CreateBuffer()
	t.createFlipBuffer()
	t.createClipBuffer()
	t.SetDefaultClip()
	t.Update()

	return t
}

func (t *Transform) CreateBuffer() {
	var err error
	t.Buffer, err = t.GetDevice().CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    t.label,
		Usage:    wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
		Contents: wgpu.ToBytes(t.Matrix[:]),
	})
	if err != nil {
		panic(err)
	}
}

func (t *Transform) createClipBuffer() error {
	clipInfo := [4]float32{0.0, 0.0, t.originalWidth, t.originalWidth}

	var err error
	t.ClipBuffer, err = t.GetDevice().CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Clip Buffer",
		Usage:    wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
		Contents: wgpu.ToBytes(clipInfo[:]),
	})
	return err
}

func (t *Transform) createFlipBuffer() error {
	flipInfo := [2]float32{0.0, 0.0}

	var err error
	t.FlipBuffer, err = t.GetDevice().CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Flip Buffer",
		Usage:    wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
		Contents: wgpu.ToBytes(flipInfo[:]),
	})
	return err
}

func (t *Transform) RecenterAndMove(vertices []primitives.Vertex, destX, destY float32) {
	center := primitives.CalculateCenter(vertices)
	for i := range vertices {
		vertices[i].Position[0] += destX - center[0]
		vertices[i].Position[1] += destY - center[1]
	}
}

func (t *Transform) MoveToScreenPosition(screenX, screenY float32) {
	sw, sh := t.GetSurfaceSize()
	ndcCoords := primitives.ScreenToNDC(screenX, screenY, float32(sw), float32(sh))
	t.Matrix[12] = ndcCoords[0]
	t.Matrix[13] = ndcCoords[1]
	t.Update()
}

func (t *Transform) Scale(sx, sy float32) matrix.Matrix {
	t.Matrix = t.Matrix.Scale(sx, sy)
	t.Update()
	return t.Matrix
}

func (t *Transform) Rotate(a float32) matrix.Matrix {
	t.Matrix = t.Matrix.Rotate(a)
	t.Update()
	return t.Matrix
}

func (t *Transform) FlipVertical() {
	t.flipVertical = !t.flipVertical
	t.updateFlipBuffer()
}

func (t *Transform) FlipHorizontal() {
	t.flipHorizontal = !t.flipHorizontal
	t.updateFlipBuffer()
}

func (t *Transform) SetFlipHorizontal(shouldFlip bool) {
	t.flipHorizontal = shouldFlip
	t.updateFlipBuffer()
}

func (t *Transform) SetFlipVertical(shouldFlip bool) {
	t.flipVertical = shouldFlip
	t.updateFlipBuffer()
}

func (t *Transform) updateFlipBuffer() {
	t.FlipMatrix[0] = 0.0
	t.FlipMatrix[1] = 0.0

	if t.flipHorizontal {
		t.FlipMatrix[0] = 1.0
	}

	if t.flipVertical {
		t.FlipMatrix[1] = 1.0
	}

	t.GetDevice().GetQueue().WriteBuffer(t.FlipBuffer, 0, wgpu.ToBytes(t.FlipMatrix[:]))
}

func (t *Transform) SetDefaultClip() {
	t.SetClipRect(0, 0, t.originalWidth, t.originalHeight)
}

func (t *Transform) SetClipRect(minX, minY, maxX, maxY float32) {
	t.ClipRect = [4]float32{
		minX / t.originalWidth,
		minY / t.originalHeight,
		maxX / t.originalWidth,
		maxY / t.originalHeight,
	}

	t.GetDevice().GetQueue().WriteBuffer(t.ClipBuffer, 0, wgpu.ToBytes(t.ClipRect[:]))
}

func (t *Transform) GetCurrentSize() (float32, float32) {
	scaleX := t.Matrix[0]
	scaleY := t.Matrix[5]

	currentWidth := t.originalWidth * scaleX
	currentHeight := t.originalHeight * scaleY

	return currentWidth, currentHeight
}

func (t *Transform) GetCurrentClipSize() (float32, float32) {
	scaleX := t.Matrix[0]
	scaleY := t.Matrix[5]

	clipWidth := ((t.ClipRect[2] * t.originalWidth) - (t.ClipRect[0] * t.originalWidth)) * scaleX
	clipHeight := ((t.ClipRect[3] * t.originalHeight) - (t.ClipRect[1] * t.originalHeight)) * scaleY

	return clipWidth, clipHeight
}

func (t *Transform) Resize(targetWidth, targetHeight float32) {
	clipWidth := (t.ClipRect[2] - t.ClipRect[0]) * t.originalWidth
	clipHeight := (t.ClipRect[3] - t.ClipRect[1]) * t.originalHeight

	scaleX := targetWidth / clipWidth
	scaleY := targetHeight / clipHeight

	t.Matrix[0] = scaleX
	t.Matrix[5] = scaleY
	t.Matrix[10] = 1
	t.Update()
}

func (t *Transform) Update() {
	t.GetDevice().GetQueue().WriteBuffer(t.Buffer, 0, wgpu.ToBytes(t.Matrix[:]))
}

func (t *Transform) Destroy() {
	if t.Buffer == nil {
		return
	}
	t.Buffer.Release()
	t.Buffer = nil
	if t.FlipBuffer != nil {
		t.FlipBuffer.Release()
		t.FlipBuffer = nil
	}
	if t.ClipBuffer != nil {
		t.ClipBuffer.Release()
		t.ClipBuffer = nil
	}
}
