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
}

func NewTransform(surface Surface, device *wgpu.Device, scd *wgpu.SwapChainDescriptor, bufferLabel string) *Transform {
	t := &Transform{
		surface:             surface,
		Device:              device,
		SwapChainDescriptor: scd,
		Matrix:              matrix.MatrixIdentity(),
		label:               bufferLabel,
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

func (t *Transform) Resize(screenWidth, screenHeight float32) {
	sw, sh := t.surface.GetSurfaceSize()

	scaleFactor := float32(sw/sh) * 100
	scaleX := screenWidth / scaleFactor
	scaleY := screenHeight / scaleFactor

	t.Matrix = t.Scale(scaleX, scaleY)
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
