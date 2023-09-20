package libsdl

import (
	"fmt"
	"runtime"

	"github.com/ebitengine/purego"
)

const (
	SDL_RENDERER_ACCELERATED  = 0x00000002
	SDL_WINDOWPOS_CENTERED    = 0x2FFF0000
	SDL_RENDERER_PRESENTVSYNC = 0x00000004
	SDL_BLENDMODE_NONE        = 0
	SDL_BLENDMODE_BLEND       = 1
	SDL_BLENDMODE_ADD         = 2
	SDL_BLENDMODE_MOD         = 3

	SDL_QUIT        = 0x100
	SDL_WINDOWEVENT = 0x200

	SDL_INIT_TIMER          uint32 = 0x00000001
	SDL_INIT_AUDIO          uint32 = 0x00000010
	SDL_INIT_VIDEO          uint32 = 0x00000020
	SDL_INIT_JOYSTICK       uint32 = 0x00000200
	SDL_INIT_HAPTIC         uint32 = 0x00001000
	SDL_INIT_GAMECONTROLLER uint32 = 0x00002000
	SDL_INIT_EVENTS         uint32 = 0x00004000
	SDL_INIT_SENSOR         uint32 = 0x00008000
	SDL_INIT_NOPARACHUTE    uint32 = 0x00100000
	SDL_INIT_EVERYTHING     uint32 = SDL_INIT_TIMER | SDL_INIT_AUDIO | SDL_INIT_VIDEO | SDL_INIT_EVENTS |
		SDL_INIT_JOYSTICK | SDL_INIT_HAPTIC | SDL_INIT_GAMECONTROLLER | SDL_INIT_SENSOR
)

type SDL_GLprofile int

const (
	SDL_GL_CONTEXT_PROFILE_CORE SDL_GLprofile = iota + 1
	SDL_GL_CONTEXT_PROFILE_COMPATIBILITY
	SDL_GL_CONTEXT_PROFILE_ES
)

var (
	libSDL uintptr
)

var (
	SDL_Init        func(flags uint32) int
	SDL_Quit        func()
	SDL_GetPlatform func()
	SDL_GetVersion  func(ver *SDL_Version)
	SDL_GetError    func() string
)

type SDL_Version struct {
	Major uint8
	Minor uint8
	Patch uint8
}

func getSystemLibrary() string {
	switch runtime.GOOS {
	case "windows":
		return "SDL2.dll"
	case "darwin":
		return "libSDL2-2.0.dylib"
	case "linux":
		return "libSDL2-2.0.so.0"
	default:
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}
}

func init() {
	var err error
	libSDL, err = openLibrary(getSystemLibrary())
	if err != nil {
		panic(err)
	}
	purego.RegisterLibFunc(&SDL_Init, libSDL, "SDL_Init")
	purego.RegisterLibFunc(&SDL_Quit, libSDL, "SDL_Quit")
	purego.RegisterLibFunc(&SDL_GetPlatform, libSDL, "SDL_GetPlatform")
	purego.RegisterLibFunc(&SDL_GetVersion, libSDL, "SDL_GetVersion")
	purego.RegisterLibFunc(&SDL_GetError, libSDL, "SDL_GetError")
	registerWindowFuncs()
	registerEventFuncs()
	registerTextureFuncs()
	registerRenderFuncs()
	registerGLFuncs()
	registerShapeFuncs()
}
