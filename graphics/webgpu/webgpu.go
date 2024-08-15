package webgpu

import (
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/surface"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/window"
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
	*surface.Surface
}

func NewGraphicsBackend(width, height int) (*GraphicsBackend, error) {
	w, err := window.NewWindow(width, height)
	if err != nil {
		panic(err)
	}

	gb := &GraphicsBackend{
		Window:  w,
		Surface: surface.New(width, height, w),
	}

	w.SetResizedCallback(func(physicalWidth, physicalHeight uint32) {
		gb.Resize(physicalWidth, physicalHeight)
	})

	w.SetCloseRequestedCallback(func() {
		w.Destroy()
		gb.Close()
	})

	return gb, nil
}

func (backend *GraphicsBackend) Close() {
	backend.Surface.Destroy()
	backend.Window.Destroy()
}

func (backend *GraphicsBackend) PollEvents() bool {
	return backend.Window.Poll()
}
