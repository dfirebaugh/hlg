package gui

// TextInputState holds the state for a text input widget.
// The caller owns this state and passes it to InputText calls.
type TextInputState struct {
	CursorPos      int
	SelectionStart int
	SelectionEnd   int
	focused        bool
}

// PanelState holds the state for a draggable panel.
// The caller owns this state and passes it to Panel calls.
type PanelState struct {
	X, Y      int
	Collapsed bool
	dragging  bool
	dragOffX  int
	dragOffY  int
}

// HasSelection returns true if there is a text selection.
func (s *TextInputState) HasSelection() bool {
	return s.SelectionStart != s.SelectionEnd
}

// ClearSelection clears any text selection.
func (s *TextInputState) ClearSelection() {
	s.SelectionStart = s.CursorPos
	s.SelectionEnd = s.CursorPos
}
