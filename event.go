package hlg

import "github.com/dfirebaugh/hlg/pkg/input"

func handleEvent(evt input.Event, state *input.InputState) {
	switch evt.Type {
	case input.KeyPress:
		state.PressKey(evt.Key)
	case input.KeyRelease:
		state.ReleaseKey(evt.Key)
	case input.MousePress:
		state.PressButton(evt.MouseButton)
	case input.MouseRelease:
		state.ReleaseButton(evt.MouseButton)
	case input.MouseMove:
		state.CursorPosition.X = evt.X
		state.CursorPosition.Y = evt.Y
	}
}
