package ggez

import (
	"time"

	"github.com/dfirebaugh/ggez/pkg/renderer/sdl/libsdl"
)

type Pressable interface {
	IsPressed() bool
	IsJustPressed() bool
}

type Positionable interface {
	GetPosition() (int, int)
}

func GetMousePosition() (int, int) {
	return libsdl.MouseX, libsdl.MouseY
}

var DefaultInput *InputDevice

type Button struct {
	isPressed     bool
	lastPressedAt time.Time
}

func CursorPositionFloat() (float64, float64) {
	x, y := GetMousePosition()
	return float64(x), float64(y)
}

func IsLeftClickPressed() bool {
	return DefaultInput.IsButtonPressed(MouseButtonLeft)
}

func IsRightClickPressed() bool {
	return DefaultInput.IsButtonPressed(MouseButtonRight)
}

// NewKey creates a new instance of a Key.
func NewButton() *Button {
	return &Button{}
}

func (k *Button) Press() {
	k.isPressed = true
	k.lastPressedAt = time.Now()
}

func (k *Button) Release() {
	k.isPressed = false
}

// IsPressed checks if the key is currently pressed.
func (k *Button) IsPressed() bool {
	return k.isPressed
}

// IsJustPressed checks if the key was just pressed in the last 500ms.
func (k *Button) IsJustPressed() bool {
	return k.isPressed && time.Since(k.lastPressedAt) <= 500*time.Millisecond
}

type Key struct {
	isPressed     bool
	lastPressedAt time.Time
}

// NewKey creates a new instance of a Key.
func NewKey() *Key {
	return &Key{}
}

func (k *Key) Press() {
	k.isPressed = true
	k.lastPressedAt = time.Now()
}

func (k *Key) Release() {
	k.isPressed = false
}

// IsPressed checks if the key is currently pressed.
func (k *Key) IsPressed() bool {
	return k.isPressed
}

// IsJustPressed checks if the key was just pressed in the last 500ms.
func (k *Key) IsJustPressed() bool {
	return k.isPressed && time.Since(k.lastPressedAt) <= 500*time.Millisecond
}

type InputDevice struct {
	keys    map[uint32]*Key
	buttons map[uint8]*Button
}

func NewInputDevice() *InputDevice {
	ensureSetupCompletion()
	return &InputDevice{
		keys:    make(map[uint32]*Key),
		buttons: make(map[uint8]*Button),
	}
}

func (k *InputDevice) PressKey(code uint32) {
	if key, exists := k.keys[code]; exists {
		key.Press()
	}
}

func (k *InputDevice) ReleaseKey(code uint32) {
	if key, exists := k.keys[code]; exists {
		key.Release()
	}
}

func (k *InputDevice) IsKeyPressed(code uint32) bool {
	if key, exists := k.keys[code]; exists {
		return key.IsPressed()
	}
	return false
}

func (k *InputDevice) IsKeyJustPressed(code uint32) bool {
	if key, exists := k.keys[code]; exists {
		return key.IsJustPressed()
	}
	return false
}
func (k *InputDevice) IsButtonPressed(code uint8) bool {
	if key, exists := k.buttons[code]; exists {
		return key.IsPressed()
	}
	return false
}

func (k *InputDevice) IsButtonJustPressed(code uint8) bool {
	if key, exists := k.buttons[code]; exists {
		return key.IsJustPressed()
	}
	return false
}

func (k *InputDevice) PressButton(code uint8) {
	if key, exists := k.buttons[code]; exists {
		key.Press()
	}
}

func (k *InputDevice) ReleaseButton(code uint8) {
	if key, exists := k.buttons[code]; exists {
		key.Release()
	}
}
