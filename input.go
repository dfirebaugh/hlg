package ggez

func IsKeyPressed(keyCode Key) bool {
	// return graphicsBackend.IsKeyPressed(uint32(keyCode))
	return false
}
func IsKeyJustPressed(keyCode Key) bool {
	// return graphicsBackend.IsKeyJustPressed(uint32(keyCode))
	return false
}
func PressKey(keyCode Key) {
	// graphicsBackend.PressKey(uint32(keyCode))
}
func ReleaseKey(keyCode Key) {
	// graphicsBackend.ReleaseKey(uint32(keyCode))
}
func IsButtonPressed(buttonCode MouseButton) bool {
	// return graphicsBackend.IsButtonPressed(uint8(buttonCode))
	return false
}
func IsButtonJustPressed(buttonCode MouseButton) bool {
	// return graphicsBackend.IsButtonJustPressed(uint8(buttonCode))
	return false
}
func PressButton(buttonCode MouseButton) {
	// graphicsBackend.PressButton(uint8(buttonCode))
}
func ReleaseButton(buttonCode MouseButton) {
	// graphicsBackend.ReleaseButton(uint8(buttonCode))
}
func GetCursorPosition() (x, y int) {
	// return graphicsBackend.GetCursorPosition()
	return 0, 0
}
func SetScrollCallback(cb func(x float64, y float64)) {
	// graphicsBackend.SetScrollCallback(cb)
}
