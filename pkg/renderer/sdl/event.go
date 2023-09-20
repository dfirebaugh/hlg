package sdl

import (
	"unsafe"

	"github.com/dfirebaugh/ggez/pkg/input"
	"github.com/ebitengine/purego"
)

type SDL_Event struct {
	Type uint32
	_    [52]byte // padding to match the size of SDL_Event in SDL2
}

type SDL_WindowEvent struct {
	Type      uint32
	Timestamp uint32
	WindowID  uint32
	Event     uint8
	_         [3]uint8 // padding
	Data1     int32
	Data2     int32
}

// Additional event structures if necessary, e.g.,
type SDL_KeyboardEvent struct {
	Type      uint32
	Timestamp uint32
	WindowID  uint32
	State     uint8
	Repeat    uint8
	_         [2]uint8
	Keysym    SDL_Keysym
}

type SDL_Keysym struct {
	Scancode uint32
	Sym      uint32
	Mod      uint16
	_        uint32
}

type SDL_MouseMotionEvent struct {
	Type      uint32
	Timestamp uint32
	WindowID  uint32
	Which     uint32
	State     uint32
	X         int32
	Y         int32
	XRel      int32
	YRel      int32
}

type SDL_MouseButtonEvent struct {
	Type      uint32
	Timestamp uint32
	WindowID  uint32
	Which     uint32
	Button    uint8
	State     uint8
	_         [2]uint8 // padding
	X         int32
	Y         int32
}

const (
	SDL_KEYDOWN         = 0x300
	SDL_KEYUP           = 0x301
	SDL_MOUSEMOTION     = 0x400
	SDL_MOUSEBUTTONDOWN = 0x401
	SDL_MOUSEBUTTONUP   = 0x402

	SDL_SCANCODE_W      = 0x1A
	SDL_SCANCODE_A      = 0x04
	SDL_SCANCODE_S      = 0x16
	SDL_SCANCODE_D      = 0x07
	SDL_SCANCODE_SPACE  = 0x2C
	SDL_SCANCODE_LCTRL  = 0xE0
	SDL_SCANCODE_LSHIFT = 0xE1
	SDL_SCANCODE_UP     = 0x52
	SDL_SCANCODE_DOWN   = 0x51
	SDL_SCANCODE_LEFT   = 0x50
	SDL_SCANCODE_RIGHT  = 0x4F
	SDL_SCANCODE_E      = 0x08
	SDL_SCANCODE_Q      = 0x14
	SDL_SCANCODE_1      = 0x1E
	SDL_SCANCODE_2      = 0x1F
	SDL_SCANCODE_3      = 0x20
	SDL_SCANCODE_4      = 0x21
	SDL_SCANCODE_5      = 0x22
	SDL_SCANCODE_6      = 0x23
	SDL_SCANCODE_7      = 0x24
	SDL_SCANCODE_8      = 0x25
	SDL_SCANCODE_9      = 0x26
	SDL_SCANCODE_0      = 0x27

	SDL_BUTTON_LEFT   = 1
	SDL_BUTTON_MIDDLE = 2
	SDL_BUTTON_RIGHT  = 3
)

var (
	MouseX, MouseY int
)

var (
	SDL_PollEvent func(event *SDL_Event) int
)

func registerEventFuncs() {
	purego.RegisterLibFunc(&SDL_PollEvent, libSDL, "SDL_PollEvent")
}

// pollEvents polls for SDL events and returns false if a quit event is detected.
func PollEvents(inputDevice input.InputDevice) bool {
	var event SDL_Event
	for {
		ret := SDL_PollEvent(&event)
		if ret == 0 {
			break
		}
		switch event.Type {
		case SDL_QUIT:
			return false
		case SDL_KEYDOWN:
			ke := (*SDL_KeyboardEvent)(unsafe.Pointer(&event))
			inputDevice.PressKey(ke.Keysym.Scancode)
		case SDL_KEYUP:
			ke := (*SDL_KeyboardEvent)(unsafe.Pointer(&event))
			inputDevice.ReleaseKey(ke.Keysym.Scancode)
		case SDL_MOUSEMOTION:
			me := (*SDL_MouseMotionEvent)(unsafe.Pointer(&event))
			MouseX = int(me.X)
			MouseY = int(me.Y)
		case SDL_MOUSEBUTTONDOWN:
			mbe := (*SDL_MouseButtonEvent)(unsafe.Pointer(&event))
			inputDevice.PressButton(mbe.Button)
		case SDL_MOUSEBUTTONUP:
			mbe := (*SDL_MouseButtonEvent)(unsafe.Pointer(&event))
			inputDevice.ReleaseButton(mbe.Button)
		}
	}
	return true
}
