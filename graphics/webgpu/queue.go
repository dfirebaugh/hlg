package webgpu

import (
	"image"
	"image/color"
	"unsafe"

	"github.com/dfirebaugh/ggez/graphics"
	"github.com/dfirebaugh/ggez/graphics/webgpu/internal/shapes"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type textureHandle uintptr

type RenderQueue struct {
	*wgpu.Device
	*wgpu.SwapChainDescriptor
	Textures     map[textureHandle]*Texture
	renderQueue  []graphics.Renderable
	nextFrame    []graphics.Renderable
	currentFrame []graphics.Renderable
}

func NewRenderQueue(d *wgpu.Device, scd *wgpu.SwapChainDescriptor) *RenderQueue {
	return &RenderQueue{
		Device:              d,
		SwapChainDescriptor: scd,
		Textures:            make(map[textureHandle]*Texture),
		nextFrame:           []graphics.Renderable{},
		currentFrame:        []graphics.Renderable{},
	}
}

func (rq *RenderQueue) AddToRenderQueue(r graphics.Renderable) {
	rq.renderQueue = append(rq.renderQueue, r)
}
func (rq *RenderQueue) Pop() (graphics.Renderable, bool) {
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

// AddTriangle creates a new Triangle renderable and adds it to the RenderQueue.
// It returns a reference to the created Triangle.
func (rq *RenderQueue) AddTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) graphics.Renderable {
	r, g, b, a := c.RGBA()
	triangle := shapes.NewTriangle(rq.Device, rq.SwapChainDescriptor, rq, []shapes.Vertex{
		{
			Position: [3]float32{float32(x1), float32(y1), 0},
			Color:    [4]float32{float32(r) / 0xffff, float32(g) / 0xffff, float32(b) / 0xffff, float32(a) / 0xffff},
		},
		{
			Position: [3]float32{float32(x2), float32(y2), 0},
			Color:    [4]float32{float32(r) / 0xffff, float32(g) / 0xffff, float32(b) / 0xffff, float32(a) / 0xffff},
		},
		{
			Position: [3]float32{float32(x3), float32(y3), 0},
			Color:    [4]float32{float32(r) / 0xffff, float32(g) / 0xffff, float32(b) / 0xffff, float32(a) / 0xffff},
		},
	}, shapes.RenderFilled)

	rq.AddToRenderQueue(triangle)

	return triangle
}
