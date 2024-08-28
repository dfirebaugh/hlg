package hlg

import "github.com/dfirebaugh/hlg/pkg/input"

// IsKeyPressed checks if a key is currently pressed
func IsKeyPressed(keyCode input.Key) bool {
	return hlg.inputState.IsKeyPressed(keyCode)
}

func IsKeyJustPressed(keyCode input.Key) bool {
	return hlg.inputState.IsKeyJustPressed(keyCode)
}

func PressKey(keyCode input.Key) {
	hlg.inputState.PressKey(keyCode)
}

func ReleaseKey(keyCode input.Key) {
	hlg.inputState.ReleaseKey(keyCode)
}

func IsButtonPressed(buttonCode input.MouseButton) bool {
	return hlg.inputState.IsButtonPressed(buttonCode)
}

// IsButtonJustPressed checks if a mouse button was just pressed
func IsButtonJustPressed(buttonCode input.MouseButton) bool {
	return hlg.inputState.IsButtonJustPressed(buttonCode)
}

// PressButton simulates a mouse button press
func PressButton(buttonCode input.MouseButton) {
	hlg.inputState.PressButton(buttonCode)
}

func ReleaseButton(buttonCode input.MouseButton) {
	hlg.inputState.ReleaseButton(buttonCode)
}

func GetCursorPosition() (int, int) {
	return hlg.inputState.GetCursorPosition()
}

func SetScrollCallback(cb func(x float64, y float64)) {
	hlg.inputState.SetScrollCallback(cb)
}
