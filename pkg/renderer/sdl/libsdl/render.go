package libsdl

import (
	"github.com/ebitengine/purego"
)

type SDL_RendererInfo struct {
	Name              *byte
	Flags             uint32
	NumTextureFormats uint32
	TextureFormats    [16]uint32
	MaxTextureWidth   int32
	MaxTextureHeight  int32
}

const (
	PIXELFORMAT_UNKNOWN     = 0x0
	PIXELFORMAT_INDEX1LSB   = 0x11100100
	PIXELFORMAT_INDEX1MSB   = 0x11200100
	PIXELFORMAT_INDEX4LSB   = 0x12100400
	PIXELFORMAT_INDEX4MSB   = 0x12200400
	PIXELFORMAT_INDEX8      = 0x13000801
	PIXELFORMAT_RGB332      = 0x14110801
	PIXELFORMAT_RGB444      = 0x15120c02
	PIXELFORMAT_RGB555      = 0x15130f02
	PIXELFORMAT_BGR555      = 0x15530f02
	PIXELFORMAT_ARGB4444    = 0x15321002
	PIXELFORMAT_RGBA4444    = 0x15421002
	PIXELFORMAT_ABGR4444    = 0x15721002
	PIXELFORMAT_BGRA4444    = 0x15821002
	PIXELFORMAT_ARGB1555    = 0x15331002
	PIXELFORMAT_RGBA5551    = 0x15441002
	PIXELFORMAT_ABGR1555    = 0x15731002
	PIXELFORMAT_BGRA5551    = 0x15841002
	PIXELFORMAT_RGB565      = 0x15151002
	PIXELFORMAT_BGR565      = 0x15551002
	PIXELFORMAT_RGB24       = 0x17101803
	PIXELFORMAT_BGR24       = 0x17401803
	PIXELFORMAT_RGB888      = 0x16161804
	PIXELFORMAT_RGBX8888    = 0x16261804
	PIXELFORMAT_BGR888      = 0x16561804
	PIXELFORMAT_BGRX8888    = 0x16661804
	PIXELFORMAT_ARGB8888    = 0x16362004
	PIXELFORMAT_RGBA8888    = 0x16462004
	PIXELFORMAT_ABGR8888    = 0x16762004
	PIXELFORMAT_BGRA8888    = 0x16862004
	PIXELFORMAT_ARGB2101010 = 0x16372004
	PIXELFORMAT_YV12        = 0x32315659
	PIXELFORMAT_IYUV        = 0x56555949
	PIXELFORMAT_YUY2        = 0x32595559
	PIXELFORMAT_UYVY        = 0x59565955
	PIXELFORMAT_YVYU        = 0x55595659
)

var (
	SDL_CreateRenderer       func(window uintptr, index int, flags uint32) uintptr
	SDL_DestroyRenderer      func(renderer uintptr)
	SDL_SetRenderDrawColor   func(renderer uintptr, r uint8, g uint8, b uint8, a uint8)
	SDL_RenderPresent        func(renderer uintptr)
	SDL_GetRendererInfo      func(renderer uintptr, info *SDL_RendererInfo)
	SDL_RenderClear          func(renderer uintptr)
	SDL_RenderCopy           func(renderer uintptr, src *SDL_Rect, dst *SDL_Rect)
	SDL_RenderCopyEx         func(renderer uintptr, texture uintptr, src uintptr, dst uintptr, uintptr, angle uintptr, flip int)
	SDL_RenderSetLogicalSize func(renderer uintptr, w int, h int)
)

func registerRenderFuncs() {
	purego.RegisterLibFunc(&SDL_CreateRenderer, libSDL, "SDL_CreateRenderer")
	purego.RegisterLibFunc(&SDL_DestroyRenderer, libSDL, "SDL_DestroyRenderer")
	purego.RegisterLibFunc(&SDL_SetRenderDrawColor, libSDL, "SDL_SetRenderDrawColor")
	purego.RegisterLibFunc(&SDL_RenderDrawLine, libSDL, "SDL_RenderDrawLine")
	purego.RegisterLibFunc(&SDL_RenderPresent, libSDL, "SDL_RenderPresent")
	purego.RegisterLibFunc(&SDL_RenderCopy, libSDL, "SDL_RenderCopy")
	purego.RegisterLibFunc(&SDL_RenderCopyEx, libSDL, "SDL_RenderCopyEx")
	purego.RegisterLibFunc(&SDL_RenderFillRect, libSDL, "SDL_RenderFillRect")
	purego.RegisterLibFunc(&SDL_RenderDrawPoint, libSDL, "SDL_RenderDrawPoint")
	purego.RegisterLibFunc(&SDL_RenderSetLogicalSize, libSDL, "SDL_RenderSetLogicalSize")
	purego.RegisterLibFunc(&SDL_RenderClear, libSDL, "SDL_RenderClear")
	purego.RegisterLibFunc(&SDL_GetRendererInfo, libSDL, "SDL_GetRendererInfo")
}
