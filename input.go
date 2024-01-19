package hlg

import "github.com/dfirebaugh/hlg/pkg/input"

func IsKeyPressed(keyCode input.Key) bool {
	return hlg.inputState.IsKeyPressed(keyCode)
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
func GetCursorPosition() (x, y int) {
	return hlg.inputState.GetCursorPosition()
}
func SetScrollCallback(cb func(x float64, y float64)) {
	hlg.inputState.SetScrollCallback(cb)
}
