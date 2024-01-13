package gpu

import (
	"fmt"
	"image/color"
	"log"

	"github.com/rajveermalviya/go-webgpu/wgpu"

	"github.com/dfirebaugh/ggez/graphics/webgpu/internal/window"
)

type RenderQueue interface {
	PrepareFrame()
	RenderFrame(pass *wgpu.RenderPassEncoder)
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
	*window.Window
	RenderQueue

	clearColor wgpu.Color
}

func NewRenderer(w *window.Window) (r *Renderer, err error) {
	wgpu.SetLogLevel(wgpu.LogLevel_Off)

	width, height := w.GetWindowSize()
	r = &Renderer{
		Window: w,
		windowSize: struct {
			Width  int
			Height int
		}{
			Width:  width,
			Height: height,
		},
	}
	err = r.setupDevice(w)

	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Recovered from panic: %v\n", r)
			if err != nil {
				r.Destroy()
			}
		}
	}()
	wgpu.SetLogLevel(wgpu.LogLevel_Error)

	w.SetCloseCallback(func() {
		r.Destroy()
	})

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
	if width <= 0 || height <= 0 {
		log.Println("Invalid dimensions for Resize")
		return
	}

	// Update window size
	r.windowSize.Width = width
	r.windowSize.Height = height

	// Update screen size used by the swap chain
	r.SetScreenSize(width, height)

	// Recreate swap chain with new dimensions
	var err error
	if r.SwapChain, err = r.createSwapChain(); err != nil {
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

func (r *Renderer) Clear(c color.Color) {
	red, green, blue, alpha := c.RGBA()
	r.clearColor = wgpu.Color{
		R: float64(red) / 0xffff,
		G: float64(green) / 0xffff,
		B: float64(blue) / 0xffff,
		A: float64(alpha) / 0xffff,
	}
}

func (r *Renderer) SurfaceIsOutdated() bool {
	if r.Window == nil {
		return true
	}
	currentWidth, currentHeight := r.Window.GetWindowSize()
	return currentWidth != int(r.SwapChainDescriptor.Width) || currentHeight != int(r.SwapChainDescriptor.Height)
}

func (r *Renderer) RecreateSwapChain() {
	if r.SwapChain != nil {
		r.SwapChain.Release()
		r.SwapChain = nil
	}

	width, height := r.Window.GetWindowSize()
	if width > 0 && height > 0 {
		r.windowSize.Width = width
		r.windowSize.Height = height
		r.SetScreenSize(width, height)

		var err error
		r.SwapChain, err = r.createSwapChain()
		if err != nil {
			fmt.Println("Failed to recreate swap chain:", err)
		}
	}
}

func (r *Renderer) Render() {
	if r.RenderQueue == nil {
		log.Println("RenderQueue is not set")
		return
	}

	width, height := r.Window.GetWindowSize()
	if width <= 0 || height <= 0 {
		return
	}

	if r.SwapChain == nil {
		log.Println("RenderQueue is not set")
		return
	}

	if r.SurfaceIsOutdated() {
		r.RecreateSwapChain()
	}

	r.RenderQueue.PrepareFrame()
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
	r.RenderQueue.RenderFrame(renderPass)
	renderPass.End()

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

func (r *Renderer) ScreenSize() (int, int) {
	return int(r.SwapChainDescriptor.Width), int(r.SwapChainDescriptor.Height)
}
