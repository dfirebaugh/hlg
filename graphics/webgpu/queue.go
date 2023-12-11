package webgpu

import (
	"image"
	"unsafe"

	"github.com/dfirebaugh/ggez/graphics"
	"github.com/dfirebaugh/ggez/graphics/webgpu/renderer"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type textureHandle uintptr

type RenderQueue struct {
	*wgpu.Device
	*wgpu.SwapChainDescriptor
	Textures    map[textureHandle]*Texture
	renderQueue []renderer.Renderable
}

func NewRenderQueue(d *wgpu.Device, scd *wgpu.SwapChainDescriptor) *RenderQueue {
	return &RenderQueue{
		Device:              d,
		SwapChainDescriptor: scd,
		Textures:            make(map[textureHandle]*Texture),
	}
}

func (rq *RenderQueue) AddToRenderQueue(r renderer.Renderable) {
	rq.renderQueue = append(rq.renderQueue, r)
}

func (rq *RenderQueue) CreateTextureFromImage(img image.Image) (graphics.Texture, error) {
	tex := NewTexture(rq.Device, rq.SwapChainDescriptor, img, rq)
	handle := textureHandle(uintptr(unsafe.Pointer(tex)))
	tex.SetHandle(handle)
	rq.Textures[handle] = tex
	return tex, nil
}
func (rq *RenderQueue) UpdateTextureFromImage(texture graphics.Texture, img image.Image) {
}
func (rq *RenderQueue) DisposeTexture(h uintptr) {
	rq.Textures[textureHandle(h)].gpuTexture.Destroy()
	delete(rq.Textures, textureHandle(h))
}

func (rq *RenderQueue) ClearRenderQueue() {
	rq.renderQueue = []renderer.Renderable{}
}

func (rq RenderQueue) Renderables() []renderer.Renderable {
	return rq.renderQueue
}
