package webgpu

import (
	"github.com/dfirebaugh/ggez/graphics/webgpu/internal/gpu"
	"github.com/dfirebaugh/ggez/graphics/webgpu/internal/window"
)

// GraphicsBackend
// Renderer
// WindowManager
// TextureManager
// EventManager
// ShapeRenderer
// Model
// ModelRenderer
// Texture
// Shader
// ShaderProgram
// Camera

type GraphicsBackend struct {
	*window.Window
	*gpu.Renderer
	*RenderQueue
}

func NewGraphicsBackend() (*GraphicsBackend, error) {
	w, err := window.NewWindow(600, 412)
	if err != nil {
		panic(err)
	}

	s, err := gpu.NewRenderer(w)
	if err != nil {
		panic(err)
	}
	rq := NewRenderQueue(s.Device, s.SwapChainDescriptor)
	s.SetRenderQueue(rq)

	w.SetResizedCallback(func(physicalWidth, physicalHeight uint32, scaleFactor float64) {
		s.Resize(int(physicalWidth), int(physicalHeight))
	})

	w.SetCloseRequestedCallback(func() {
		w.Destroy()
	})

	return &GraphicsBackend{
		Window:      w,
		Renderer:    s,
		RenderQueue: rq,
	}, nil
}

func (backend *GraphicsBackend) Close() {
	backend.Renderer.Destroy()
	backend.Window.Destroy()
}

func (backend *GraphicsBackend) PollEvents() bool {
	return backend.Window.Poll()
}

func (backend *GraphicsBackend) PrintPlatformAndVersion() {}
