package hlg

import "github.com/dfirebaugh/hlg/pkg/input"

var previousKeyboardState = map[input.Key]bool{}

func IsKeyPressed(keyCode input.Key) bool {
	return hlg.inputState.IsKeyPressed(keyCode)
}

func IsKeyJustPressed(keyCode input.Key) bool {
	isPressed := hlg.inputState.IsKeyPressed(keyCode)
	wasPressed := previousKeyboardState[keyCode]

	previousKeyboardState[keyCode] = isPressed

	return isPressed && !wasPressed
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

func PressButton(buttonCode input.MouseButton) {
	hlg.inputState.PressButton(buttonCode)
}

func ReleaseButton(buttonCode input.MouseButton) {
	hlg.inputState.ReleaseButton(buttonCode)
}

func GetCursorPosition() (int, int) {
	sw, sh := hlg.graphicsBackend.GetScreenSize()
	w, h := hlg.graphicsBackend.GetWindowSize()

	x, y := hlg.inputState.GetCursorPosition()
	scaleX := float64(sw) / float64(w)
	scaleY := float64(sh) / float64(h)

	virtualX := int(float64(x) * scaleX)
	virtualY := int(float64(y) * scaleY)

	return virtualX, virtualY
}

func SetScrollCallback(cb func(x float64, y float64)) {
	hlg.inputState.SetScrollCallback(cb)
}
