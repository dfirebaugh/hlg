package gl

import (
	"time"

	"github.com/dfirebaugh/ggez/pkg/graphics"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func (g *GLRenderer) keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		g.InputManager.PressKey(uint32(key))
	} else if action == glfw.Release {
		g.InputManager.ReleaseKey(uint32(key))
	}

	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

// Ensure that our InputDeviceGlfw struct satisfies the InputManager interface.
var _ graphics.InputManager = (*InputDeviceGlfw)(nil)

type InputDeviceGlfw struct {
	window *glfw.Window
	keys   map[glfw.Key]*Key
}

type Key struct {
	isPressed       bool
	wasJustPressed  bool
	wasJustReleased bool
	lastPressedAt   time.Time
	lastReleasedAt  time.Time
}

func NewInputDeviceGlfw(window *glfw.Window) *InputDeviceGlfw {
	return &InputDeviceGlfw{
		window: window,
		keys:   make(map[glfw.Key]*Key),
	}
}

func (id *InputDeviceGlfw) PressKey(keyCode uint32) {
	key := glfw.Key(keyCode)
	if k, exists := id.keys[key]; exists {
		// Avoiding updating if key is already pressed
		if !k.isPressed {
			k.isPressed = true
			k.wasJustPressed = true
			k.lastPressedAt = time.Now()
		}
	} else {
		id.keys[key] = &Key{
			isPressed:      true,
			wasJustPressed: true,
			lastPressedAt:  time.Now(),
		}
	}
}

func (id *InputDeviceGlfw) ReleaseKey(keyCode uint32) {
	key := glfw.Key(keyCode)
	if k, exists := id.keys[key]; exists {
		k.isPressed = false
		k.wasJustReleased = true
		k.lastReleasedAt = time.Now()
	}
}

func (id *InputDeviceGlfw) IsKeyJustPressed(keyCode uint32) bool {
	key := glfw.Key(keyCode)

	if k, exists := id.keys[key]; exists {
		return k.wasJustPressed
	}
	return false
}

func (id *InputDeviceGlfw) IsKeyJustReleased(keyCode uint32) bool {
	key := glfw.Key(keyCode)

	if k, exists := id.keys[key]; exists {
		return k.wasJustReleased
	}
	return false
}

func (id *InputDeviceGlfw) Update() {
	for _, k := range id.keys {
		k.wasJustPressed = false
		k.wasJustReleased = false
	}
}

func (id *InputDeviceGlfw) IsKeyPressed(keyCode uint32) bool {
	key := glfw.Key(keyCode)
	state := id.window.GetKey(key)
	return state == glfw.Press
}
func (id *InputDeviceGlfw) IsButtonPressed(buttonCode uint8) bool {
	state := id.window.GetMouseButton(glfw.MouseButton(buttonCode))
	return state == glfw.Press
}

func (id *InputDeviceGlfw) IsButtonJustPressed(buttonCode uint8) bool {
	return false
}

func (id *InputDeviceGlfw) PressButton(buttonCode uint8) {
}

func (id *InputDeviceGlfw) ReleaseButton(buttonCode uint8) {
}

func (id *InputDeviceGlfw) GetCursorPosition() (x, y int) {
	xf, yf := id.window.GetCursorPos()
	return int(xf), int(yf)
}
