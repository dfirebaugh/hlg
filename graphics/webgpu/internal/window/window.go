package window

import (
	"github.com/dfirebaugh/ggez/pkg/input"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Window struct {
	*glfw.Window
	aspectRatio float64
	eventChan   chan input.Event
	isDisposed  bool
}

func NewWindow(width, height int) (*Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)

	w := &Window{
		eventChan:  make(chan input.Event, 100),
		isDisposed: false,
	}

	win, err := glfw.CreateWindow(640, 480, "go-webgpu with glfw", nil, nil)
	if err != nil {
		glfw.Terminate()
		w.DestroyWindow()
		return nil, err
	}
	w.Window = win
	win.SetSize(width, height)
	return w, nil
}

func (w *Window) SetAspectRatio(numerator, denominator int) {
	w.aspectRatio = float64(numerator) / float64(denominator)
	w.Window.SetAspectRatio(numerator, denominator)
}
func (w *Window) SetCloseCallback(fn func()) {
	w.Window.SetCloseCallback(func(w *glfw.Window) {
		defer w.Destroy()
		fn()
	})
}
func (w *Window) SetInputCallback(fn func(eventChan chan input.Event)) {
	w.Window.SetMouseButtonCallback(func(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		var eventType input.EventType
		if action == glfw.Press {
			eventType = input.MousePress
		} else if action == glfw.Release {
			eventType = input.MouseRelease
		}
		w.eventChan <- input.Event{Type: eventType, MouseButton: input.MouseButton(button)}
		fn(w.eventChan)
	})

	w.Window.SetCursorPosCallback(func(window *glfw.Window, xpos, ypos float64) {
		w.eventChan <- input.Event{Type: input.MouseMove, X: int(xpos), Y: int(ypos)}
		fn(w.eventChan)
	})

	w.Window.SetKeyCallback(func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		var eventType input.EventType
		if action == glfw.Press {
			eventType = input.KeyPress
		} else if action == glfw.Release {
			eventType = input.KeyRelease
		}
		w.eventChan <- input.Event{Type: eventType, Key: input.Key(key)}
		fn(w.eventChan)
	})
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
	if w.isDisposed {
		return 0, 0
	}
	return w.Window.GetSize()
}

func (w *Window) SetWindowSize(width int, height int) {
	w.SetAspectRatio(width, height)
	w.SetSize(width, height)
}

func (w *Window) SetScaleFactor(f int) {}

func (w *Window) DestroyWindow() {
	w.Destroy()
}

func (w *Window) Poll() bool {
	if w.isDisposed {
		return false
	}

	glfw.PollEvents()
	return true
}

func (w *Window) IsDisposed() bool {
	return w.isDisposed
}

func (w *Window) Destroy() {
	w.isDisposed = true
	if w.isDisposed {
		return
	}
	if w.Window != nil {
		w.Window.Destroy()
		w.Window = nil
	}
	glfw.Terminate()
}
