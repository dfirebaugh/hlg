package libsdl

import (
	"github.com/ebitengine/purego"
)

type SDL_GLContext struct{}
type SDL_Window struct{}

type SDL_GLattr int

const (
	SDL_GL_RED_SIZE SDL_GLattr = iota
	SDL_GL_GREEN_SIZE
	SDL_GL_BLUE_SIZE
	SDL_GL_ALPHA_SIZE
	SDL_GL_BUFFER_SIZE
	SDL_GL_DOUBLEBUFFER
	SDL_GL_DEPTH_SIZE
	SDL_GL_STENCIL_SIZE
	SDL_GL_ACCUM_RED_SIZE
	SDL_GL_ACCUM_GREEN_SIZE
	SDL_GL_ACCUM_BLUE_SIZE
	SDL_GL_ACCUM_ALPHA_SIZE
	SDL_GL_STEREO
	SDL_GL_MULTISAMPLEBUFFERS
	SDL_GL_MULTISAMPLESAMPLES
	SDL_GL_ACCELERATED_VISUAL
	SDL_GL_RETAINED_BACKING
	SDL_GL_CONTEXT_MAJOR_VERSION
	SDL_GL_CONTEXT_MINOR_VERSION
	SDL_GL_CONTEXT_FLAGS
	SDL_GL_CONTEXT_PROFILE_MASK
	SDL_GL_SHARE_WITH_CURRENT_CONTEXT
	SDL_GL_FRAMEBUFFER_SRGB_CAPABLE
	SDL_GL_CONTEXT_RELEASE_BEHAVIOR
	SDL_GL_CONTEXT_EGL
)

var (
	SDL_GL_BindTexture        func(texture uintptr, w float64, h float64) int
	SDL_GL_CreateContext      func(window uintptr) uintptr
	SDL_GL_DeleteContext      func(ctx uintptr)
	SDL_GL_ExtensionSupported func(extension string) bool
	SDL_GL_GetAttribute       func(attr SDL_GLattr, i uintptr) int
	SDL_GL_GetCurrentContext  func() uintptr
	SDL_GL_GetCurrentWindow   func() uintptr
	SDL_GL_GetDrawableSize    func(window uintptr, w int, h int)
	SDL_GL_GetProcAddress     func(proc string) uintptr
	SDL_GL_GetSwapInterval    func() int
	SDL_GL_LoadLibrary        func(path string) int
	SDL_GL_MakeCurrent        func(window uintptr, ctx uintptr) int
	SDL_GL_ResetAttributes    func()
	SDL_GL_SetAttribute       func(attr SDL_GLattr, value int) int
	SDL_GL_SetSwapInterval    func(interval int)
	SDL_GL_SwapWindow         func(window uintptr)
	SDL_GL_UnbindTexture      func(texture uintptr) int
	SDL_GL_UnloadLibrary      func()
)

func registerGLFuncs() {
	purego.RegisterLibFunc(&SDL_GL_BindTexture, libSDL, "SDL_GL_BindTexture")
	purego.RegisterLibFunc(&SDL_GL_CreateContext, libSDL, "SDL_GL_CreateContext")
	purego.RegisterLibFunc(&SDL_GL_DeleteContext, libSDL, "SDL_GL_DeleteContext")
	purego.RegisterLibFunc(&SDL_GL_ExtensionSupported, libSDL, "SDL_GL_ExtensionSupported")
	purego.RegisterLibFunc(&SDL_GL_GetAttribute, libSDL, "SDL_GL_GetAttribute")
	purego.RegisterLibFunc(&SDL_GL_GetCurrentContext, libSDL, "SDL_GL_GetCurrentContext")
	purego.RegisterLibFunc(&SDL_GL_GetCurrentWindow, libSDL, "SDL_GL_GetCurrentWindow")
	purego.RegisterLibFunc(&SDL_GL_GetDrawableSize, libSDL, "SDL_GL_GetDrawableSize")
	purego.RegisterLibFunc(&SDL_GL_GetProcAddress, libSDL, "SDL_GL_GetProcAddress")
	purego.RegisterLibFunc(&SDL_GL_GetSwapInterval, libSDL, "SDL_GL_GetSwapInterval")
	purego.RegisterLibFunc(&SDL_GL_LoadLibrary, libSDL, "SDL_GL_LoadLibrary")
	purego.RegisterLibFunc(&SDL_GL_MakeCurrent, libSDL, "SDL_GL_MakeCurrent")
	purego.RegisterLibFunc(&SDL_GL_ResetAttributes, libSDL, "SDL_GL_ResetAttributes")
	purego.RegisterLibFunc(&SDL_GL_SetAttribute, libSDL, "SDL_GL_SetAttribute")
	purego.RegisterLibFunc(&SDL_GL_SetSwapInterval, libSDL, "SDL_GL_SetSwapInterval")
	purego.RegisterLibFunc(&SDL_GL_SwapWindow, libSDL, "SDL_GL_SwapWindow")
	purego.RegisterLibFunc(&SDL_GL_UnbindTexture, libSDL, "SDL_GL_UnbindTexture")
	purego.RegisterLibFunc(&SDL_GL_UnloadLibrary, libSDL, "SDL_GL_UnloadLibrary")
}
