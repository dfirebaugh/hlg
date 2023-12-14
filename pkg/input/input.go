package input

type InputState struct {
	KeyState          map[Key]bool
	KeyJustPressed    map[Key]bool
	ButtonState       map[MouseButton]bool
	ButtonJustPressed map[MouseButton]bool
	CursorPosition    struct{ X, Y int }
	ScrollCallback    func(x, y float64)
}

func NewInputState() *InputState {
	return &InputState{
		KeyState:          make(map[Key]bool),
		KeyJustPressed:    make(map[Key]bool),
		ButtonState:       make(map[MouseButton]bool),
		ButtonJustPressed: make(map[MouseButton]bool),
	}
}

// IsKeyPressed returns true if the specified key is currently pressed
func (is *InputState) IsKeyPressed(keyCode Key) bool {
	pressed, exists := is.KeyState[keyCode]
	return exists && pressed
}

// IsKeyJustPressed returns true if the specified key was just pressed
func (is *InputState) IsKeyJustPressed(keyCode Key) bool {
	pressed, exists := is.KeyJustPressed[keyCode]
	if exists {
		is.KeyJustPressed[keyCode] = false // Reset the state
	}
	return exists && pressed
}

// PressKey simulates a key press
func (is *InputState) PressKey(keyCode Key) {
	is.KeyState[keyCode] = true
	is.KeyJustPressed[keyCode] = true
}

// ReleaseKey simulates a key release
func (is *InputState) ReleaseKey(keyCode Key) {
	is.KeyState[keyCode] = false
}

// IsButtonPressed returns true if the specified mouse button is currently pressed
func (is *InputState) IsButtonPressed(buttonCode MouseButton) bool {
	pressed, exists := is.ButtonState[buttonCode]
	return exists && pressed
}

// IsButtonJustPressed returns true if the specified mouse button was just pressed
func (is *InputState) IsButtonJustPressed(buttonCode MouseButton) bool {
	pressed, exists := is.ButtonJustPressed[buttonCode]
	if exists {
		is.ButtonJustPressed[buttonCode] = false // Reset the state
	}
	return exists && pressed
}

// PressButton simulates a mouse button press
func (is *InputState) PressButton(buttonCode MouseButton) {
	is.ButtonState[buttonCode] = true
	is.ButtonJustPressed[buttonCode] = true
}

// ReleaseButton simulates a mouse button release
func (is *InputState) ReleaseButton(buttonCode MouseButton) {
	is.ButtonState[buttonCode] = false
}

// GetCursorPosition returns the current cursor position
func (is *InputState) GetCursorPosition() (x, y int) {
	return is.CursorPosition.X, is.CursorPosition.Y
}

// SetScrollCallback sets a callback function for scroll events
func (is *InputState) SetScrollCallback(cb func(x, y float64)) {
	is.ScrollCallback = cb
}

func (is *InputState) ResetJustPressed() {
	for key := range is.KeyJustPressed {
		is.KeyJustPressed[key] = false
	}
	for button := range is.ButtonJustPressed {
		is.ButtonJustPressed[button] = false
	}
}
