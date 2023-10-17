package graphics

// InputManager is an interface defining methods for handling user input.
type InputManager interface {
	IsKeyPressed(keyCode uint32) bool
	IsKeyJustPressed(keyCode uint32) bool
	PressKey(keyCode uint32)
	ReleaseKey(keyCode uint32)

	IsButtonPressed(buttonCode uint8) bool
	IsButtonJustPressed(buttonCode uint8) bool
	PressButton(buttonCode uint8)
	ReleaseButton(buttonCode uint8)

	GetCursorPosition() (x, y int)

	SetScrollCallback(cb func(x float64, y float64))

	Update()
}
