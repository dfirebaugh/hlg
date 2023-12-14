package ggez

import "github.com/dfirebaugh/ggez/pkg/input"

func IsKeyPressed(keyCode input.Key) bool {
	return ggez.inputState.IsKeyPressed(keyCode)
}
func PressKey(keyCode input.Key) {
	ggez.inputState.PressKey(keyCode)
}
func ReleaseKey(keyCode input.Key) {
	ggez.inputState.ReleaseKey(keyCode)
}
func IsButtonPressed(buttonCode input.MouseButton) bool {
	return ggez.inputState.IsButtonPressed(buttonCode)
}
func PressButton(buttonCode input.MouseButton) {
	ggez.inputState.PressButton(buttonCode)
}
func ReleaseButton(buttonCode input.MouseButton) {
	ggez.inputState.ReleaseButton(buttonCode)
}
func GetCursorPosition() (x, y int) {
	return ggez.inputState.GetCursorPosition()
}
func SetScrollCallback(cb func(x float64, y float64)) {
	ggez.inputState.SetScrollCallback(cb)
}
