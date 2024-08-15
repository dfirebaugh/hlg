package input

type EventType int

const (
	KeyPress EventType = iota
	KeyRelease
	MousePress
	MouseRelease
	MouseMove
)

type Event struct {
	Type        EventType
	Key         Key
	MouseButton MouseButton
	X, Y        int
}
