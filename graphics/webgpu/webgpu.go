package webgpu

import (
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
