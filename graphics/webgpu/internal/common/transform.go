package common

import (
	"github.com/dfirebaugh/hlg/pkg/math/matrix"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Transform struct {
	surface Surface
	*wgpu.Device
	*wgpu.SwapChainDescriptor
	*wgpu.Buffer

	matrix.Matrix
	label string

	originalWidth  float32
	originalHeight float32
}

func NewTransform(surface Surface, device *wgpu.Device, scd *wgpu.SwapChainDescriptor, bufferLabel string, originalWidth, originalHeight float32) *Transform {
	t := &Transform{
		surface:             surface,
		Device:              device,
		SwapChainDescriptor: scd,
		Matrix:              matrix.MatrixIdentity(),
		label:               bufferLabel,
		originalWidth:       originalWidth,
		originalHeight:      originalHeight,
	}

	t.CreateBuffer()
	t.Update()

	return t
}

func (t *Transform) CreateBuffer() {
	var err error
	t.Buffer, err = t.Device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    t.label,
		Usage:    wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
		Contents: wgpu.ToBytes(t.Matrix[:]),
	})
	if err != nil {
		panic(err)
	}
}

func (t *Transform) Move(screenX, screenY float32) {
	sw, sh := t.surface.GetSurfaceSize()
	ndcCoords := ScreenToNDC(screenX, screenY, float32(sw), float32(sh))
	t.Matrix[12] = ndcCoords[0]
	t.Matrix[13] = ndcCoords[1]
	t.Update()
}

func (t *Transform) Scale(sx, sy float32) matrix.Matrix {
	t.Matrix = t.Matrix.Scale(sx, sy)
	t.Update()
	println("should scale and update")
	return t.Matrix
}

func (t *Transform) Resize(targetWidth, targetHeight float32) {
	scaleX := targetWidth / (t.originalWidth * 2)
	scaleY := targetHeight / (t.originalHeight * 2)

	// Maintain aspect ratio by choosing the smaller scale factor
	scaleFactor := scaleX
	if scaleY < scaleX {
		scaleFactor = scaleY
	}

	// Apply scale factor directly to the matrix elements
	t.Matrix[0] = scaleFactor
	t.Matrix[5] = scaleFactor
	t.Matrix[10] = 1 // Ensure no depth scaling
	t.Update()
}

func (t *Transform) Rotate(a float32) matrix.Matrix {
	t.Matrix = t.Matrix.Rotate(a)
	t.Update()
	return t.Matrix
}

func (t *Transform) Update() {
	t.Device.GetQueue().WriteBuffer(t.Buffer, 0, wgpu.ToBytes(t.Matrix[:]))
}

func (t *Transform) Destroy() {
	if t.Buffer == nil {
		return
	}
	t.Buffer.Release()
	t.Buffer = nil
}
