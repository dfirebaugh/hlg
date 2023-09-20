package libsdl

import (
	"github.com/ebitengine/purego"
)

type SDL_Point struct {
	X, Y int32
}

type SDL_Rect struct {
	X, Y int32
	W, H int32
}

var (
	SDL_RenderFillRect  func(renderer uintptr, rect *SDL_Rect)
	SDL_RenderDrawRect  func(renderer uintptr, rect *SDL_Rect)
	SDL_RenderDrawPoint func(renderer uintptr, x int, y int)
	SDL_RenderDrawLine  func(renderer uintptr, x1, y1, x2, y2 int)
)

func registerShapeFuncs() {
	purego.RegisterLibFunc(&SDL_RenderFillRect, libSDL, "SDL_RenderFillRect")
	purego.RegisterLibFunc(&SDL_RenderDrawRect, libSDL, "SDL_RenderDrawRect")
	purego.RegisterLibFunc(&SDL_RenderDrawPoint, libSDL, "SDL_RenderDrawPoint")
	purego.RegisterLibFunc(&SDL_RenderDrawLine, libSDL, "SDL_RenderDrawLine")
}
