//go:build js && wasm

package canvas

import (
	"syscall/js"

	"github.com/dfirebaugh/hlg/pkg/input"
)

// setupEventListeners sets up DOM event listeners for mouse and keyboard
func (c *Canvas) setupEventListeners() {
	// Mouse button events
	c.element.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		button := event.Get("button").Int()
		c.eventChan <- input.Event{
			Type:        input.MousePress,
			MouseButton: input.MouseButton(button),
		}
		if c.inputCallback != nil {
			c.inputCallback(c.eventChan)
		}
		return nil
	}))

	c.element.Call("addEventListener", "mouseup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		button := event.Get("button").Int()
		c.eventChan <- input.Event{
			Type:        input.MouseRelease,
			MouseButton: input.MouseButton(button),
		}
		if c.inputCallback != nil {
			c.inputCallback(c.eventChan)
		}
		return nil
	}))

	// Mouse move event
	c.element.Call("addEventListener", "mousemove", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		// Use getBoundingClientRect for accurate coordinates with CSS transforms
		rect := c.element.Call("getBoundingClientRect")
		clientX := event.Get("clientX").Float()
		clientY := event.Get("clientY").Float()
		rectLeft := rect.Get("left").Float()
		rectTop := rect.Get("top").Float()
		rectWidth := rect.Get("width").Float()
		rectHeight := rect.Get("height").Float()

		// Calculate position relative to canvas element
		relX := clientX - rectLeft
		relY := clientY - rectTop

		// Scale to logical coordinates using the actual rendered size from getBoundingClientRect
		// This ensures correct translation even immediately after window resize
		x, y := c.cssToLogicalWithRect(relX, relY, rectWidth, rectHeight)

		c.eventChan <- input.Event{
			Type: input.MouseMove,
			X:    x,
			Y:    y,
		}
		if c.inputCallback != nil {
			c.inputCallback(c.eventChan)
		}
		return nil
	}))

	// Keyboard events (on document for global capture)
	doc := js.Global().Get("document")

	doc.Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		keyCode := event.Get("keyCode").Int()
		c.eventChan <- input.Event{
			Type: input.KeyPress,
			Key:  convertKeyCode(keyCode),
		}
		if c.inputCallback != nil {
			c.inputCallback(c.eventChan)
		}
		return nil
	}))

	doc.Call("addEventListener", "keyup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		keyCode := event.Get("keyCode").Int()
		c.eventChan <- input.Event{
			Type: input.KeyRelease,
			Key:  convertKeyCode(keyCode),
		}
		if c.inputCallback != nil {
			c.inputCallback(c.eventChan)
		}
		return nil
	}))

	// Character input via keypress
	doc.Call("addEventListener", "keypress", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		key := event.Get("key").String()
		if len(key) == 1 {
			c.eventChan <- input.Event{
				Type: input.CharInput,
				Rune: rune(key[0]),
			}
			if c.inputCallback != nil {
				c.inputCallback(c.eventChan)
			}
		}
		return nil
	}))
}

// setupResizeListener sets up window resize listener
func (c *Canvas) setupResizeListener() {
	// Listen for window resize events to re-fit the canvas
	js.Global().Get("window").Call("addEventListener", "resize", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		c.handleResize()
		return nil
	}))
}

func (c *Canvas) handleResize() {
	// Re-fit the canvas to the window while preserving aspect ratio
	c.fitToWindow()
}

// convertKeyCode converts JavaScript keyCode to input.Key
// This is a basic mapping - extend as needed
func convertKeyCode(keyCode int) input.Key {
	switch keyCode {
	// Letter keys
	case 65:
		return input.KeyA
	case 66:
		return input.KeyB
	case 67:
		return input.KeyC
	case 68:
		return input.KeyD
	case 69:
		return input.KeyE
	case 70:
		return input.KeyF
	case 71:
		return input.KeyG
	case 72:
		return input.KeyH
	case 73:
		return input.KeyI
	case 74:
		return input.KeyJ
	case 75:
		return input.KeyK
	case 76:
		return input.KeyL
	case 77:
		return input.KeyM
	case 78:
		return input.KeyN
	case 79:
		return input.KeyO
	case 80:
		return input.KeyP
	case 81:
		return input.KeyQ
	case 82:
		return input.KeyR
	case 83:
		return input.KeyS
	case 84:
		return input.KeyT
	case 85:
		return input.KeyU
	case 86:
		return input.KeyV
	case 87:
		return input.KeyW
	case 88:
		return input.KeyX
	case 89:
		return input.KeyY
	case 90:
		return input.KeyZ
	// Number keys
	case 48:
		return input.Key0
	case 49:
		return input.Key1
	case 50:
		return input.Key2
	case 51:
		return input.Key3
	case 52:
		return input.Key4
	case 53:
		return input.Key5
	case 54:
		return input.Key6
	case 55:
		return input.Key7
	case 56:
		return input.Key8
	case 57:
		return input.Key9
	// Special keys
	case 13:
		return input.KeyEnter
	case 27:
		return input.KeyEscape
	case 8:
		return input.KeyBackspace
	case 9:
		return input.KeyTab
	case 32:
		return input.KeySpace
	case 37:
		return input.KeyLeft
	case 38:
		return input.KeyUp
	case 39:
		return input.KeyRight
	case 40:
		return input.KeyDown
	case 16:
		return input.KeyLeftShift
	case 17:
		return input.KeyLeftControl
	case 18:
		return input.KeyLeftAlt
	default:
		return input.Key(keyCode)
	}
}
