package ggez

func IsKeyPressed(keyCode Key) bool {
	return graphicsBackend.IsKeyPressed(uint32(keyCode))
}
func IsKeyJustPressed(keyCode Key) bool {
	return graphicsBackend.IsKeyJustPressed(uint32(keyCode))
}
func PressKey(keyCode Key) {
	graphicsBackend.PressKey(uint32(keyCode))
}
func ReleaseKey(keyCode Key) {
	graphicsBackend.ReleaseKey(uint32(keyCode))
}
func IsButtonPressed(buttonCode MouseButton) bool {
	return graphicsBackend.IsButtonPressed(uint8(buttonCode))
}
func IsButtonJustPressed(buttonCode MouseButton) bool {
	return graphicsBackend.IsButtonJustPressed(uint8(buttonCode))
}
func PressButton(buttonCode MouseButton) {
	graphicsBackend.PressButton(uint8(buttonCode))
}
func ReleaseButton(buttonCode MouseButton) {
	graphicsBackend.ReleaseButton(uint8(buttonCode))
}
func GetCursorPosition() (x, y int) {
	return graphicsBackend.GetCursorPosition()
}
func SetScrollCallback(cb func(x float64, y float64)) {
	graphicsBackend.SetScrollCallback(cb)
}
