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
	Textures     map[textureHandle]*Texture
	renderQueue  []renderer.Renderable
	nextFrame    []renderer.Renderable
	currentFrame []renderer.Renderable
}

func NewRenderQueue(d *wgpu.Device, scd *wgpu.SwapChainDescriptor) *RenderQueue {
	return &RenderQueue{
		Device:              d,
		SwapChainDescriptor: scd,
		Textures:            make(map[textureHandle]*Texture),
		nextFrame:           []renderer.Renderable{},
		currentFrame:        []renderer.Renderable{},
	}
}

func (rq *RenderQueue) AddToRenderQueue(r renderer.Renderable) {
	rq.renderQueue = append(rq.renderQueue, r)
}
func (rq *RenderQueue) Pop() (renderer.Renderable, bool) {
	if len(rq.renderQueue) == 0 {
		return nil, false
	}

	renderable := rq.renderQueue[0]
	rq.renderQueue = rq.renderQueue[1:]

	return renderable, true
}

func (rq *RenderQueue) PrepareFrame() {
	if len(rq.nextFrame) > 64 {
		rq.nextFrame = rq.nextFrame[:64]
	}
	for {
		if renderable, ok := rq.Pop(); ok {
			rq.nextFrame = append(rq.nextFrame, renderable)
			continue
		}
		break
	}
	rq.currentFrame = rq.nextFrame
}

func (rq *RenderQueue) RenderFrame(pass *wgpu.RenderPassEncoder) {
	for _, renderable := range rq.currentFrame {
		renderable.RenderPass(pass)
	}
}

func (rq *RenderQueue) CreateTextureFromImage(img image.Image) (graphics.Texture, error) {
	tex := NewTexture(rq.Device, rq.SwapChainDescriptor, img, rq)
	handle := textureHandle(uintptr(unsafe.Pointer(tex)))
	tex.SetHandle(handle)
	rq.Textures[handle] = tex
	return tex, nil
}
func (rq *RenderQueue) UpdateTextureFromImage(texture graphics.Texture, img image.Image) {
	texture.UpdateImage(img)
}
func (rq *RenderQueue) DisposeTexture(h uintptr) {
	rq.Textures[textureHandle(h)].gpuTexture.Destroy()
	delete(rq.Textures, textureHandle(h))
}
