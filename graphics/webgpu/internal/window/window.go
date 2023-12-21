package window

import (
	"github.com/dfirebaugh/ggez/pkg/input"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Window struct {
	*glfw.Window
	aspectRatio float64
	eventChan   chan input.Event
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

	w.Window.SetSize(int(width), int(height))
	w.eventChan = make(chan input.Event, 100)

	w.Window.SetIconifyCallback(func(w *glfw.Window, iconified bool) {

	})

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
