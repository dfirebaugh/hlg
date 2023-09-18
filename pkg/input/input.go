package input

type InputDevice interface {
	PressKey(code uint32)
	ReleaseKey(code uint32)
	PressButton(code uint8)
	ReleaseButton(code uint8)
}
