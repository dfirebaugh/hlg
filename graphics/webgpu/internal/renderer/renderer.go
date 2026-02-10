//go:build !js

package renderer

import (
	"fmt"
	"image/color"
	"log"

	"github.com/rajveermalviya/go-webgpu/wgpu"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/context"
)

type RenderTarget interface {
	GetSize() (int, int)
	GetSurfaceDescriptor() *wgpu.SurfaceDescriptor
}

type Renderer struct {
	*wgpu.Surface
	*wgpu.Device
	*wgpu.SwapChain
	*wgpu.SwapChainDescriptor
	renderTarget RenderTarget
	surface      context.Surface
	clearColor   wgpu.Color

	RenderQueues []*RenderQueue

	clipRectStack [][4]int
}

func NewRenderer(s context.Surface, width, height int, renderTarget RenderTarget) (r *Renderer, err error) {
	wgpu.SetLogLevel(wgpu.LogLevel_Error)
	r = &Renderer{
		surface:      s,
		renderTarget: renderTarget,
	}

	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Recovered from panic: %v\n", rec)
			r.Destroy()
			err = fmt.Errorf("panic during renderer setup: %v", rec)
		}
	}()

	_ = r.setupDevice()

	return r, err
}

func (r *Renderer) AddRenderQueue(rq *RenderQueue) {
	r.RenderQueues = append(r.RenderQueues, rq)
}

func (r *Renderer) CreateRenderQueue() graphics.RenderQueue {
	rq := NewRenderQueue(r.surface, r.Device, r.SwapChainDescriptor)
	r.RenderQueues = append(r.RenderQueues, rq)
	return rq
}

func (r *Renderer) setupDevice() error {
	var err error
	instance := wgpu.CreateInstance(nil)
	defer instance.Release()
	r.Surface = instance.CreateSurface(r.renderTarget.GetSurfaceDescriptor())

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

	ww, wh := r.renderTarget.GetSize()

	r.SwapChainDescriptor = &wgpu.SwapChainDescriptor{
		Usage:       wgpu.TextureUsage_RenderAttachment,
		Format:      r.Surface.GetPreferredFormat(adapter),
		Width:       uint32(ww),
		Height:      uint32(wh),
		PresentMode: wgpu.PresentMode_Immediate,
	}
	r.SwapChain, err = r.Device.CreateSwapChain(r.Surface, r.SwapChainDescriptor)

	return err
}

func (r *Renderer) Resize(width int, height int) {
	if width <= 0 || height <= 0 {
		log.Println("Invalid dimensions for Resize")
		return
	}

	r.surface.SetSurfaceSize(width, height)

	r.SetScreenSize(width, height)
	if r.SwapChain != nil {
		r.SwapChain.Release()
		r.SwapChain = nil
	}

	var err error
	r.SwapChain, err = r.createSwapChain()
	if err != nil {
		panic(err)
	}
}

func (r *Renderer) createSwapChain() (*wgpu.SwapChain, error) {
	if r.Device == nil {
		return nil, fmt.Errorf("device is nil")
	}

	return r.Device.CreateSwapChain(r.Surface, r.SwapChainDescriptor)
}

func (r *Renderer) SetScreenSize(width int, height int) {
	r.SwapChainDescriptor.Width = uint32(width)
	r.SwapChainDescriptor.Height = uint32(height)
}

func (r *Renderer) SetVSync(enabled bool) {
	var mode wgpu.PresentMode
	if enabled {
		mode = wgpu.PresentMode_Fifo
	} else {
		mode = wgpu.PresentMode_Immediate
	}

	if r.SwapChainDescriptor.PresentMode == mode {
		return
	}

	r.SwapChainDescriptor.PresentMode = mode
	r.RecreateSwapChain()
}

func (r *Renderer) Clear(c color.Color) {
	for _, rq := range r.RenderQueues {
		if !rq.shouldClear {
			continue
		}
		rq.RenderClear()
	}
	red, green, blue, alpha := c.RGBA()
	r.clearColor = wgpu.Color{
		R: float64(red) / 0xffff,
		G: float64(green) / 0xffff,
		B: float64(blue) / 0xffff,
		A: float64(alpha) / 0xffff,
	}
}

func (r *Renderer) SurfaceIsOutdated() bool {
	if r.renderTarget == nil {
		return true
	}
	currentWidth, currentHeight := r.renderTarget.GetSize()
	return currentWidth != int(r.SwapChainDescriptor.Width) || currentHeight != int(r.SwapChainDescriptor.Height)
}

func (r *Renderer) RecreateSwapChain() {
	if r.SwapChain != nil {
		r.SwapChain.Release()
		r.SwapChain = nil
	}

	width, height := r.renderTarget.GetSize()
	if width > 0 && height > 0 {
		r.SetScreenSize(width, height)

		var err error
		r.SwapChain, err = r.createSwapChain()
		if err != nil {
			fmt.Println("Failed to recreate swap chain:", err)
		}
	}
}

func (r *Renderer) Render() {
	width, height := r.renderTarget.GetSize()
	if width <= 0 || height <= 0 {
		return
	}

	if r.SwapChain == nil {
		log.Println("swapChain is not set")
		return
	}

	if r.SurfaceIsOutdated() {
		r.RecreateSwapChain()
	}

	for _, rq := range r.RenderQueues {
		rq.PrepareFrame()
	}
	view, err := r.SwapChain.GetCurrentTextureView()
	if err != nil {
		fmt.Println("Error getting texture view:", err)

		if r.SurfaceIsOutdated() {
			r.RecreateSwapChain()
		}
		return
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
	for _, rq := range r.RenderQueues {
		rq.RenderFrame(renderPass)
	}
	_ = renderPass.End()

	cmdBuffer, err := encoder.Finish(nil)
	if err != nil {
		panic(err.Error())
	}
	defer cmdBuffer.Release()

	r.Device.GetQueue().Submit(cmdBuffer)
	r.SwapChain.Present()
}

func (r *Renderer) Destroy() {
	if r.SwapChain != nil {
		r.SwapChain.Release()
		r.SwapChain = nil
	}
	if r.SwapChainDescriptor != nil {
		r.SwapChainDescriptor = nil
	}
	if r.Device != nil && r.Device.GetQueue() != nil {
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

func (r *Renderer) PushClipRect(x, y, width, height int) {
	r.clipRectStack = append(r.clipRectStack, [4]int{x, y, width, height})
}

func (r *Renderer) PopClipRect() {
	if len(r.clipRectStack) == 0 {
		return
	}
	r.clipRectStack = r.clipRectStack[:len(r.clipRectStack)-1]
}

func (r *Renderer) GetCurrentClipRect() *[4]int {
	if len(r.clipRectStack) == 0 {
		return nil
	}
	rect := r.clipRectStack[len(r.clipRectStack)-1]
	return &rect
}
