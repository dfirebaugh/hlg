//go:build !js

package window

import (
	"github.com/dfirebaugh/hlg/pkg/input"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Window struct {
	*glfw.Window
	aspectRatio                 float64
	eventChan                   chan input.Event
	isDisposed                  bool
	targetWidth, targetHeight   int // logical size set by SetWindowSize
	currentWidth, currentHeight int // actual window size
}

func NewWindow(width, height int) (*Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	// Enable high-DPI/Retina framebuffer on macOS
	glfw.WindowHint(glfw.CocoaRetinaFramebuffer, glfw.True)
	glfw.WindowHint(glfw.ScaleToMonitor, glfw.True)

	w := &Window{
		eventChan:     make(chan input.Event, 100),
		isDisposed:    false,
		targetWidth:   width,
		targetHeight:  height,
		currentWidth:  width,
		currentHeight: height,
	}

	win, err := glfw.CreateWindow(width, height, "hlg-opengl", nil, nil)
	if err != nil {
		glfw.Terminate()
		return nil, err
	}
	w.Window = win

	win.MakeContextCurrent()

	return w, nil
}

func (w *Window) SetBorderlessWindowed(v bool) {
	if v {
		monitor := glfw.GetPrimaryMonitor()
		mode := monitor.GetVideoMode()
		w.Window.SetMonitor(monitor, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
	} else {
		w.Window.SetMonitor(nil, 100, 100, 800, 600, 0)
	}
}

func (w *Window) GetWindowPosition() (x int, y int) {
	return w.Window.GetPos()
}

func (w *Window) SetAspectRatio(numerator, denominator int) {
	w.aspectRatio = float64(numerator) / float64(denominator)
	w.Window.SetAspectRatio(numerator, denominator)
}

func (w *Window) DisableWindowResize() {
	if w.Window != nil {
		w.Window.SetAttrib(glfw.Resizable, glfw.False)
	}
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
		switch action {
		case glfw.Press:
			eventType = input.MousePress
		case glfw.Release:
			eventType = input.MouseRelease
		}
		w.eventChan <- input.Event{Type: eventType, MouseButton: input.MouseButton(button)}
		fn(w.eventChan)
	})

	w.Window.SetCursorPosCallback(func(window *glfw.Window, xpos, ypos float64) {
		// Translate from current window coordinates to logical coordinates
		x, y := w.windowToLogical(xpos, ypos)
		w.eventChan <- input.Event{Type: input.MouseMove, X: x, Y: y}
		fn(w.eventChan)
	})

	w.Window.SetKeyCallback(func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		var eventType input.EventType
		switch action {
		case glfw.Press:
			eventType = input.KeyPress
		case glfw.Release:
			eventType = input.KeyRelease
		}
		w.eventChan <- input.Event{Type: eventType, Key: input.Key(key)}
		fn(w.eventChan)
	})

	w.Window.SetCharCallback(func(window *glfw.Window, char rune) {
		w.eventChan <- input.Event{Type: input.CharInput, Rune: char}
		fn(w.eventChan)
	})
}

func (w *Window) SetResizedCallback(fn func(physicalWidth, physicalHeight uint32)) {
	// Use window size callback (not framebuffer) to match webgpu coordinate system
	w.Window.SetSizeCallback(func(window *glfw.Window, width, height int) {
		w.currentWidth = width
		w.currentHeight = height
		fn(uint32(width), uint32(height))
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

func (w *Window) GetSize() (int, int) {
	return w.GetFramebufferSize()
}

func (w *Window) GetFramebufferSize() (int, int) {
	if w == nil || w.isDisposed {
		return 0, 0
	}
	return w.Window.GetFramebufferSize()
}

func (w *Window) GetWindowSize() (int, int) {
	if w == nil || w.isDisposed {
		return 0, 0
	}
	return w.Window.GetSize()
}

func (w *Window) SetWindowSize(width int, height int) {
	w.targetWidth = width
	w.targetHeight = height
	w.currentWidth = width
	w.currentHeight = height
	w.SetAspectRatio(width, height)
	w.SetSize(width, height)
}

// windowToLogical translates window coordinates to logical coordinates
// This ensures mouse coordinates are consistent regardless of window resize
func (w *Window) windowToLogical(xpos, ypos float64) (int, int) {
	if w.currentWidth <= 0 || w.currentHeight <= 0 {
		return int(xpos), int(ypos)
	}
	if w.targetWidth <= 0 || w.targetHeight <= 0 {
		return int(xpos), int(ypos)
	}

	scaleX := float64(w.targetWidth) / float64(w.currentWidth)
	scaleY := float64(w.targetHeight) / float64(w.currentHeight)

	return int(xpos * scaleX), int(ypos * scaleY)
}

func (w *Window) DestroyWindow() {
	w.Destroy()
}

func (w *Window) Poll() bool {
	if w.isDisposed {
		return false
	}

	glfw.PollEvents()

	if w.Window != nil && w.Window.ShouldClose() {
		return false
	}

	return true
}

// SwapBuffers swaps the front and back buffers
func (w *Window) SwapBuffers() {
	if w.Window != nil {
		w.Window.SwapBuffers()
	}
}

func (w *Window) IsDisposed() bool {
	return w.isDisposed
}

func (w *Window) Destroy() {
	if w.isDisposed {
		return
	}
	w.isDisposed = true
	if w.Window != nil {
		w.Window.Destroy()
		w.Window = nil
	}
	glfw.Terminate()
}

func (w *Window) SetVSync(enabled bool) {
	if enabled {
		glfw.SwapInterval(1)
	} else {
		glfw.SwapInterval(0)
	}
}
