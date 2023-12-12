package gpu

import (
	"image/color"

	"github.com/rajveermalviya/go-webgpu/wgpu"

	"github.com/dfirebaugh/ggez/graphics/webgpu/internal/window"
	"github.com/dfirebaugh/ggez/graphics/webgpu/renderer"
)

type RenderQueue interface {
	Renderables() []renderer.Renderable
	ClearRenderQueue()
}

type Renderer struct {
	windowSize struct {
		Width  int
		Height int
	}
	*wgpu.Surface
	*wgpu.Device
	*wgpu.SwapChain
	*wgpu.SwapChainDescriptor
	RenderQueue

	clearColor wgpu.Color
}

func NewRenderer(w *window.Window) (r *Renderer, err error) {
	wgpu.SetLogLevel(wgpu.LogLevel_Off)

	defer func() {
		if err != nil {
			r.Destroy()
			r = nil
		}
	}()
	width, height := w.GetWindowSize()
	r = &Renderer{
		windowSize: struct {
			Width  int
			Height int
		}{
			Width:  width,
			Height: height,
		},
	}
	err = r.setupDevice(w)

	return r, err
}

func (r *Renderer) SetRenderQueue(rq RenderQueue) {
	r.RenderQueue = rq
}

func (r *Renderer) setupDevice(w *window.Window) error {
	var err error
	instance := wgpu.CreateInstance(nil)
	defer instance.Release()
	r.Surface = instance.CreateSurface(window.GetSurfaceDescriptor(w.Window))

	adapter, err := instance.RequestAdapter(&wgpu.RequestAdapterOptions{
		CompatibleSurface: r.Surface,
	})
	if err != nil {
		return err
	}
	defer adapter.Release()

	r.Device, err = adapter.RequestDevice(nil)
	if err != nil {
		return err
	}

	r.SwapChainDescriptor = &wgpu.SwapChainDescriptor{
		Usage:       wgpu.TextureUsage_RenderAttachment,
		Format:      r.Surface.GetPreferredFormat(adapter),
		Width:       uint32(r.windowSize.Width),
		Height:      uint32(r.windowSize.Height),
		PresentMode: wgpu.PresentMode_Fifo,
	}
	r.SwapChain, err = r.Device.CreateSwapChain(r.Surface, r.SwapChainDescriptor)

	return err
}

func (r *Renderer) Resize(width int, height int) {
	if width > 0 && height > 0 {
		r.windowSize.Width = width
		r.windowSize.Height = height
		r.SetScreenSize(width, height)

		if r.SwapChain != nil {
			r.SwapChain.Release()
		}
		var err error
		r.SwapChain, err = r.Device.CreateSwapChain(r.Surface, r.SwapChainDescriptor)
		if err != nil {
			panic(err)
		}
	}
}

func (r *Renderer) SetScreenSize(width int, height int) {
	r.SwapChainDescriptor.Width = uint32(width)
	r.SwapChainDescriptor.Height = uint32(height)
}

func (r *Renderer) Clear(c color.Color) {
	red, green, blue, alpha := c.RGBA()
	r.clearColor = wgpu.Color{
		R: float64(red) / 0xffff,
		G: float64(green) / 0xffff,
		B: float64(blue) / 0xffff,
		A: float64(alpha) / 0xffff,
	}
}
func (r *Renderer) Render() {
	view, err := r.SwapChain.GetCurrentTextureView()
	if err != nil {
		panic(err.Error())
	}
	defer view.Release()

	if view == nil {
		println("view is nil")
		return
	}

	encoder, err := r.Device.CreateCommandEncoder(nil)
	if err != nil {
		panic(err.Error())
	}
	defer encoder.Release()

	renderPass := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
		ColorAttachments: []wgpu.RenderPassColorAttachment{{
			View:       view,
			LoadOp:     wgpu.LoadOp_Clear,
			ClearValue: r.clearColor,
			StoreOp:    wgpu.StoreOp_Store,
		}},
	})
	defer renderPass.Release()

	for _, renderable := range r.RenderQueue.Renderables() {
		renderable.RenderPass(renderPass)
	}
	renderPass.End()

	cmdBuffer, err := encoder.Finish(nil)
	if err != nil {
		panic(err.Error())
	}
	defer cmdBuffer.Release()

	r.Device.GetQueue().Submit(cmdBuffer)
	r.SwapChain.Present()

	r.ClearRenderQueue()
}

func (r *Renderer) Destroy() {
	if r.SwapChain != nil {
		r.SwapChain.Release()
		r.SwapChain = nil
	}
	if r.SwapChainDescriptor != nil {
		r.SwapChainDescriptor = nil
	}
	if r.Device.GetQueue() != nil {
		r.Device.GetQueue().Release()
	}
	if r.Device != nil {
		r.Device.Release()
		r.Device = nil
	}
	if r.Surface != nil {
		r.Surface.Release()
		r.Surface = nil
	}
}

func (r *Renderer) ScreenHeight() int {
	return int(r.SwapChainDescriptor.Height)
}

func (r *Renderer) ScreenWidth() int {
	return int(r.SwapChainDescriptor.Width)
}
