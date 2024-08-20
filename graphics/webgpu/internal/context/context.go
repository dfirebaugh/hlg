package context

import (
	"log"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/pipeline"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type pipelineManager interface {
	GetPipeline(key string, layout *wgpu.PipelineLayoutDescriptor, shaderModule *wgpu.ShaderModule, scd *wgpu.SwapChainDescriptor, topology wgpu.PrimitiveTopology) *wgpu.RenderPipeline
}

type RenderQueue interface {
	AddToRenderQueue(r graphics.Renderable)
	DisposeTexture(h uintptr)
}

type Surface interface {
	GetSurfaceSize() (int, int)
}

type RenderContext interface {
	Surface
	RenderQueue

	GetPipelineManager() *pipeline.PipelineManager
	GetDevice() *wgpu.Device
	GetSwapChainDescriptor() *wgpu.SwapChainDescriptor
}

type renderContext struct {
	Surface
	device              *wgpu.Device
	swapChainDescirptor *wgpu.SwapChainDescriptor
	*pipeline.PipelineManager
	RenderQueue
}

func NewRenderContext(surface Surface, device *wgpu.Device, scd *wgpu.SwapChainDescriptor, rq RenderQueue) RenderContext {
	if surface == nil || device == nil || scd == nil || rq == nil {
		log.Fatal("RenderContext dependencies are not initialized")
	}

	ctx := &renderContext{
		Surface:             surface,
		device:              device,
		swapChainDescirptor: scd,
		RenderQueue:         rq,
		PipelineManager:     pipeline.NewPipelineManager(device),
	}

	if ctx.RenderQueue == nil {
		log.Fatal("RenderQueue is not initialized")
	}

	return ctx
}

func (c *renderContext) GetDevice() *wgpu.Device {
	return c.device
}

func (c *renderContext) GetSwapChainDescriptor() *wgpu.SwapChainDescriptor {
	return c.swapChainDescirptor
}

func (c *renderContext) GetPipelineManager() *pipeline.PipelineManager {
	return c.PipelineManager
}
