package gui

import (
	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/pkg/input"
)

// DefaultInputContext wraps hlg's input functions to implement InputContext.
type DefaultInputContext struct {
	charInput []rune
}

// NewDefaultInputContext creates an InputContext that uses hlg's input functions.
func NewDefaultInputContext() *DefaultInputContext {
	return &DefaultInputContext{}
}

// MousePosition returns the current mouse position.
func (d *DefaultInputContext) MousePosition() (int, int) {
	return hlg.GetCursorPosition()
}

// IsButtonPressed returns true if the button is currently pressed.
func (d *DefaultInputContext) IsButtonPressed(button input.MouseButton) bool {
	return hlg.IsButtonPressed(button)
}

// IsButtonJustPressed returns true if the button was just pressed this frame.
func (d *DefaultInputContext) IsButtonJustPressed(button input.MouseButton) bool {
	return hlg.IsButtonJustPressed(button)
}

// IsButtonJustReleased returns true if the button was just released this frame.
func (d *DefaultInputContext) IsButtonJustReleased(button input.MouseButton) bool {
	return hlg.IsButtonJustReleased(button)
}

// IsKeyPressed returns true if the key is currently pressed.
func (d *DefaultInputContext) IsKeyPressed(key input.Key) bool {
	return hlg.IsKeyPressed(key)
}

// IsKeyJustPressed returns true if the key was just pressed this frame.
func (d *DefaultInputContext) IsKeyJustPressed(key input.Key) bool {
	return hlg.IsKeyJustPressed(key)
}

// GetCharInput returns runes typed this frame.
func (d *DefaultInputContext) GetCharInput() []rune {
	return d.charInput
}

// Update must be called once per frame to update input state.
func (d *DefaultInputContext) Update() {
	d.charInput = hlg.GetTypedRunes()
}
