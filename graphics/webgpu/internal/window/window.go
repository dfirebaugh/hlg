package window

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Window struct {
	*glfw.Window
	size struct {
		Width  int
		Height int
	}
}

func NewWindow(width int, height int) (*Window, error) {
	w := &Window{}
	var err error

	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	w.Window, err = glfw.CreateWindow(640, 480, "go-webgpu with glfw", nil, nil)
	if err != nil {
		w.Destroy()
		return nil, err
	}

	w.size.Width = width
	w.size.Height = height
	w.Window.SetSize(int(width), int(height))

	return w, nil
}

func (w *Window) SetResizedCallback(fn func(physicalWidth, physicalHeight uint32, scaleFactor float64)) {
	w.Window.SetSizeCallback(func(window *glfw.Window, width, height int) {
		fn(uint32(width), uint32(height), 1)
	})
}
func (w *Window) SetCloseRequestedCallback(fn func()) {
	w.Window.SetCloseCallback(func(w *glfw.Window) {
		fn()
	})
}

func (w *Window) SetWindowTitle(title string) {
	w.SetTitle(title)
}

func (w *Window) GetWindowSize() (int, int) {
	return w.Window.GetSize()
}

func (w *Window) SetWindowSize(width int, height int) {
	w.SetSize(width, height)
}

func (w *Window) SetScaleFactor(f int) {}

func (w *Window) DestroyWindow() {
	w.Destroy()
}

func (w *Window) Poll() bool {
	glfw.PollEvents()
	return true
}

func (w *Window) Destroy() {
	if w.Window != nil {
		w.Window.Destroy()
		w.Window = nil
	}
	glfw.Terminate()
}
