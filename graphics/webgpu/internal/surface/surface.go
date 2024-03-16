package surface

import (
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/gpu"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/window"
)

type Surface struct {
	*gpu.Renderer
	*RenderQueue

	Width  int
	Height int
}

func New(w, h int, win *window.Window) *Surface {
	surface := &Surface{}
	s, err := gpu.NewRenderer(surface, w, h, win)
	if err != nil {
		panic(err)
	}
	rq := NewRenderQueue(surface, s.Device, s.SwapChainDescriptor)
	s.SetRenderQueue(rq)

	surface.Renderer = s
	surface.RenderQueue = rq
	surface.Width = w
	surface.Height = h
	return surface
}

func (s *Surface) Resize(width, height uint32) {
	s.Renderer.Resize(int(width), int(height))
}

func (s *Surface) GetSurfaceSize() (int, int) {
	return s.Width, s.Height
}

func (s *Surface) SetScreenSize(w, h int) {
	s.SetSurfaceSize(w, h)
	// s.Renderer.SetScreenSize(w, h)
}

func (s *Surface) SetSurfaceSize(w, h int) {
	s.Width = w
	s.Height = h
}

func (s *Surface) GetScreenSize() (int, int) {
	return s.Width, s.Height
}

func (s *Surface) Destroy() {
	s.Renderer.Destroy()
}
