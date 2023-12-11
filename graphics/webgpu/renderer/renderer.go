package renderer

import "github.com/rajveermalviya/go-webgpu/wgpu"

type Renderable interface {
	RenderPass(pass *wgpu.RenderPassEncoder)
	SetShouldBeRendered(souldRender bool)
}
