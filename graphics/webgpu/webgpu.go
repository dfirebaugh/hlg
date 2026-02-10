//go:build !js

package webgpu

import (
	"image"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/renderer"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/window"
)

type GraphicsBackend struct {
	*window.Window
	*renderer.Surface
}

func NewGraphicsBackend(width, height int) (*GraphicsBackend, error) {
	w, err := window.NewWindow(width, height)
	if err != nil {
		panic(err)
	}

	gb := &GraphicsBackend{
		Window:  w,
		Surface: renderer.NewSurface(width, height, w),
	}

	w.SetResizedCallback(func(physicalWidth, physicalHeight uint32) {
		gb.Resize(physicalWidth, physicalHeight)
	})

	// Note: Don't destroy resources in the close callback - it's called from within
	// glfwPollEvents(). Cleanup happens after the main loop exits via Close().

	return gb, nil
}

func (backend *GraphicsBackend) Close() {
	if backend.Surface != nil {
		backend.Surface.Destroy()
	}
	if backend.Window != nil {
		backend.Window.Destroy()
	}
}

func (backend *GraphicsBackend) PollEvents() bool {
	return backend.Window.Poll()
}

// CreateMSDFAtlas creates a new MSDF atlas from an image.
// Implements graphics.FontManager interface.
func (backend *GraphicsBackend) CreateMSDFAtlas(atlasImg image.Image, distanceRange float64) (graphics.MSDFAtlas, error) {
	return backend.RenderQueue.CreateMSDFAtlas(atlasImg, distanceRange)
}
