//go:build !js

package renderer

type Surface struct {
	*Renderer
	*RenderQueue

	Width  int
	Height int
}

func NewSurface(w, h int, renderTarget RenderTarget) *Surface {
	surface := &Surface{}
	s, err := NewRenderer(surface, w, h, renderTarget)
	if err != nil {
		panic(err)
	}
	rq := NewRenderQueue(surface, s.Device, s.SwapChainDescriptor)
	rq.SetPriority(100)
	rq.SetShouldClear(true)
	s.AddRenderQueue(rq)

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
}

func (s *Surface) SetSurfaceSize(w, h int) {
	s.Width = w
	s.Height = h
}

func (s *Surface) GetScreenSize() (int, int) {
	return s.Width, s.Height
}

func (s *Surface) SetVSync(enabled bool) {
	s.Renderer.SetVSync(enabled)
}

func (s *Surface) Destroy() {
	s.Renderer.Destroy()
	s.ReleaseShaders()
}

// GetCurrentClipRect disambiguates between embedded Renderer and RenderQueue.
func (s *Surface) GetCurrentClipRect() *[4]int {
	return s.Renderer.GetCurrentClipRect()
}
