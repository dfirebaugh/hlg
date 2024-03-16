package surface

import (
	"image"

	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/texture"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type textureHandle uintptr

type Texture struct {
	surface      *Surface
	handle       textureHandle
	gpuTexture   *texture.Texture
	renderQueue  *RenderQueue
	shouldRender bool
}

func NewTexture(surface *Surface, d *wgpu.Device, scd *wgpu.SwapChainDescriptor, img image.Image, renderQueue *RenderQueue) *Texture {
	gpuTexture, _ := texture.TextureFromImage(surface, d, scd, img, "label")
	t := &Texture{
		surface:     surface,
		gpuTexture:  gpuTexture,
		renderQueue: renderQueue,
	}

	return t
}

func (t *Texture) UpdateImage(img image.Image) error {
	return t.gpuTexture.UpdateImage(img)
}
func (t *Texture) SetHandle(h textureHandle) {
	t.handle = h
}
func (t *Texture) Handle() uintptr {
	return uintptr(t.handle)
}
func (t *Texture) SetShouldBeRendered(shouldRender bool) {
	t.shouldRender = shouldRender
}
func (t *Texture) Resize(width, height float32) {
	t.gpuTexture.Resize(width, height)
}
func (t *Texture) Move(x, y float32) {
	t.gpuTexture.Move(x, y)
}
func (t *Texture) Rotate(a, pivotX, pivotY float32) {
	t.gpuTexture.Rotate(a)
}
func (t *Texture) Scale(x, y float32) {
	t.gpuTexture.Scale(x, y)
}
func (t *Texture) FlipVertical() {
	t.gpuTexture.FlipVertical()
}
func (t *Texture) FlipHorizontal() {
	t.gpuTexture.FlipHorizontal()
}
func (t *Texture) Clip(minX, minY, maxX, maxY float32) {
	t.gpuTexture.SetClipRect(minX, minY, maxX, maxY)
}
func (t *Texture) RenderPass(pass *wgpu.RenderPassEncoder) {
	if !t.shouldRender {
		return
	}
	t.gpuTexture.RenderPass(pass)
}
func (t *Texture) Render() {
	t.SetShouldBeRendered(true)
	t.renderQueue.AddToRenderQueue(t)
}
func (t *Texture) Hide() {
	t.shouldRender = false
}
func (t *Texture) Dispose() {
	t.renderQueue.DisposeTexture(uintptr(t.handle))
}
func (t *Texture) IsDisposed() bool {
	return t.gpuTexture.IsDisposed()
}
