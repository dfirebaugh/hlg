package libsdl

import (
	"github.com/ebitengine/purego"
)

var (
	SDL_CreateTexture       func(renderer uintptr, format uint32, access int, w int, h int) uintptr
	SDL_UpdateTexture       func(texture uintptr, rect *SDL_Rect, pixels *uint8, pitch int)
	SDL_DestroyTexture      func(texture uintptr)
	SDL_SetTextureBlendMode func(texture uintptr, blendMode uint8)
)

func registerTextureFuncs() {
	purego.RegisterLibFunc(&SDL_CreateTexture, libSDL, "SDL_CreateTexture")
	purego.RegisterLibFunc(&SDL_UpdateTexture, libSDL, "SDL_UpdateTexture")
	purego.RegisterLibFunc(&SDL_DestroyTexture, libSDL, "SDL_DestroyTexture")
	purego.RegisterLibFunc(&SDL_SetTextureBlendMode, libSDL, "SDL_SetTextureBlendMode")
}

type RendererFlip int

const (
	FLIP_NONE RendererFlip = iota
	FLIP_HORIZONTAL
	FLIP_VERTICAL
	FLIP_BOTH
)
