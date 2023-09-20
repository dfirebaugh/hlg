package ggez

import "github.com/dfirebaugh/ggez/pkg/renderer/sdl"

// Keyboard Keys
const (
	KeyW      = sdl.SDL_SCANCODE_W
	KeyA      = sdl.SDL_SCANCODE_A
	KeyS      = sdl.SDL_SCANCODE_S
	KeyD      = sdl.SDL_SCANCODE_D
	KeySPACE  = sdl.SDL_SCANCODE_SPACE
	KeyLCTRL  = sdl.SDL_SCANCODE_LCTRL
	KeyLSHIFT = sdl.SDL_SCANCODE_LSHIFT
	KeyUP     = sdl.SDL_SCANCODE_UP
	KeyDOWN   = sdl.SDL_SCANCODE_DOWN
	KeyLEFT   = sdl.SDL_SCANCODE_LEFT
	KeyRIGHT  = sdl.SDL_SCANCODE_RIGHT
	KeyE      = sdl.SDL_SCANCODE_E
	KeyQ      = sdl.SDL_SCANCODE_Q
	Key1      = sdl.SDL_SCANCODE_1
	Key2      = sdl.SDL_SCANCODE_2
	Key3      = sdl.SDL_SCANCODE_3
	Key4      = sdl.SDL_SCANCODE_4
	Key5      = sdl.SDL_SCANCODE_5
	Key6      = sdl.SDL_SCANCODE_6
	Key7      = sdl.SDL_SCANCODE_7
	Key8      = sdl.SDL_SCANCODE_8
	Key9      = sdl.SDL_SCANCODE_9
	Key0      = sdl.SDL_SCANCODE_0
)

// Mouse Buttons
const (
	MouseButtonLeft   uint8 = sdl.SDL_BUTTON_LEFT
	MouseButtonRight  uint8 = sdl.SDL_BUTTON_RIGHT
	MouseButtonMiddle uint8 = sdl.SDL_BUTTON_MIDDLE
)

func setupDefaultInput() {
	ensureSetupCompletion()

	DefaultInput = &InputDevice{
		keys:    make(map[uint32]*Key),
		buttons: make(map[uint8]*Button),
	}

	DefaultInput.keys[KeyW] = NewKey()
	DefaultInput.keys[KeyA] = NewKey()
	DefaultInput.keys[KeyS] = NewKey()
	DefaultInput.keys[KeyD] = NewKey()
	DefaultInput.keys[KeySPACE] = NewKey()
	DefaultInput.keys[KeyLCTRL] = NewKey()
	DefaultInput.keys[KeyLSHIFT] = NewKey()
	DefaultInput.keys[KeyUP] = NewKey()
	DefaultInput.keys[KeyDOWN] = NewKey()
	DefaultInput.keys[KeyLEFT] = NewKey()
	DefaultInput.keys[KeyRIGHT] = NewKey()
	DefaultInput.keys[KeyE] = NewKey()
	DefaultInput.keys[KeyQ] = NewKey()
	DefaultInput.keys[Key1] = NewKey()
	DefaultInput.keys[Key2] = NewKey()
	DefaultInput.keys[Key3] = NewKey()
	DefaultInput.keys[Key4] = NewKey()
	DefaultInput.keys[Key5] = NewKey()
	DefaultInput.keys[Key6] = NewKey()
	DefaultInput.keys[Key7] = NewKey()
	DefaultInput.keys[Key8] = NewKey()
	DefaultInput.keys[Key9] = NewKey()
	DefaultInput.keys[Key0] = NewKey()

	DefaultInput.buttons[MouseButtonLeft] = NewButton()
	DefaultInput.buttons[MouseButtonRight] = NewButton()
	DefaultInput.buttons[MouseButtonMiddle] = NewButton()
}
