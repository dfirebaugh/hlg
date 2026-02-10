package input

type InputState struct {
	KeyState           map[Key]bool
	KeyJustPressed     map[Key]bool
	ButtonState        map[MouseButton]bool
	ButtonJustPressed  map[MouseButton]bool
	ButtonJustReleased map[MouseButton]bool
	CursorPosition     struct{ X, Y int }
	ScrollCallback     func(x, y float64)
	TypedRunes         []rune
}

func NewInputState() *InputState {
	return &InputState{
		KeyState:           make(map[Key]bool),
		KeyJustPressed:     make(map[Key]bool),
		ButtonState:        make(map[MouseButton]bool),
		ButtonJustPressed:  make(map[MouseButton]bool),
		ButtonJustReleased: make(map[MouseButton]bool),
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
	return exists && pressed
}

// PressKey simulates a key press
func (is *InputState) PressKey(keyCode Key) {
	is.KeyState[keyCode] = true
	if !is.KeyJustPressed[keyCode] {
		is.KeyJustPressed[keyCode] = true
	}
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
	justPressed, exists := is.ButtonJustPressed[buttonCode]
	return exists && justPressed
}

// IsButtonJustReleased returns true if the specified mouse button was just released
func (is *InputState) IsButtonJustReleased(buttonCode MouseButton) bool {
	justReleased, exists := is.ButtonJustReleased[buttonCode]
	return exists && justReleased
}

// PressButton simulates a mouse button press
func (is *InputState) PressButton(buttonCode MouseButton) {
	is.ButtonState[buttonCode] = true
	if !is.ButtonJustPressed[buttonCode] {
		is.ButtonJustPressed[buttonCode] = true
	}
}

// ReleaseButton simulates a mouse button release
func (is *InputState) ReleaseButton(buttonCode MouseButton) {
	// Track if button was pressed before release
	if is.ButtonState[buttonCode] {
		is.ButtonJustReleased[buttonCode] = true
	}
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

// ResetJustPressed resets the just pressed/released state at the end of a frame
func (is *InputState) ResetJustPressed() {
	for key := range is.KeyJustPressed {
		is.KeyJustPressed[key] = false
	}
	for button := range is.ButtonJustPressed {
		is.ButtonJustPressed[button] = false
	}
	for button := range is.ButtonJustReleased {
		is.ButtonJustReleased[button] = false
	}
	is.TypedRunes = is.TypedRunes[:0]
}

// AddTypedRune adds a rune to the typed runes list
func (is *InputState) AddTypedRune(r rune) {
	is.TypedRunes = append(is.TypedRunes, r)
}

// GetTypedRunes returns the runes typed this frame
func (is *InputState) GetTypedRunes() []rune {
	return is.TypedRunes
}
