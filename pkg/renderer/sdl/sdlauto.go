package sdl

import (
	"fmt"
	"image"
	"image/color"
	"unsafe"

	"github.com/dfirebaugh/ggez/pkg/input"
	"github.com/dfirebaugh/ggez/pkg/renderer"
	"github.com/dfirebaugh/ggez/pkg/renderer/sdl/libsdl"
)

// SDLAutoRenderer automatically determines what backend graphics api to use
type SDLAutoRenderer struct {
	renderer.WindowManager
	renderer.EventManager
	screenHeight        int
	screenWidth         int
	RendererHandle      uintptr
	WindowManagerHandle uintptr
}

func New() (*SDLAutoRenderer, error) {
	libsdl.SDL_Init(libsdl.SDL_INIT_VIDEO)
	g := &SDLAutoRenderer{}

	g.CreateWindow("ggez", 800, 600)
	g.createRenderer()
	return g, nil
}

func (g *SDLAutoRenderer) Close() {
	g.DestroyWindow()
	g.DestroyRenderer()
	libsdl.SDL_Quit()
}

func (g *SDLAutoRenderer) PrintPlatformAndVersion() {

}

func (g *SDLAutoRenderer) createRenderer() error {
	r := libsdl.SDL_CreateRenderer(g.WindowManagerHandle, 0xFFFFFFFF, libsdl.SDL_RENDERER_ACCELERATED)
	if r == 0 {
		return fmt.Errorf("libsdl.SDL_CreateRenderer failed")
	}
	g.RendererHandle = r
	g.SetScreenSize(g.screenWidth, g.screenHeight)
	return nil
}

func (g *SDLAutoRenderer) Render() {
	libsdl.SDL_RenderPresent(g.RendererHandle)
}

func (g *SDLAutoRenderer) DestroyRenderer() {
	libsdl.SDL_DestroyRenderer(g.RendererHandle)
}

func (g *SDLAutoRenderer) PrintRendererInfo() {
	var info libsdl.SDL_RendererInfo

	libsdl.SDL_GetRendererInfo(g.RendererHandle, &info)

	rendererName := goStringFromCString(info.Name)
	fmt.Printf("Backend: %s\n", rendererName)
}

func (g *SDLAutoRenderer) Clear(c color.Color) {
	g.setRenderDrawColor(c)
	libsdl.SDL_RenderClear(g.RendererHandle)
}

func (w *SDLAutoRenderer) CreateWindow(title string, width, height int) (uintptr, error) {
	window := libsdl.SDL_CreateWindow(
		title,
		libsdl.SDL_WINDOWPOS_CENTERED, libsdl.SDL_WINDOWPOS_CENTERED, width, height, libsdl.SDL_WINDOW_SHOWN|libsdl.SDL_WINDOW_RESIZABLE,
	)
	if window == 0 {
		return 0, fmt.Errorf("libsdl.SDL_CreateWindow failed")
	}

	w.WindowManagerHandle = window
	w.SetWindowTitle(title)

	return window, nil
}

func (w *SDLAutoRenderer) SetWindowTitle(title string) {
	if len(title) == 0 {
		return
	}
	if w.WindowManagerHandle == 0 {
		return
	}
	libsdl.SDL_SetWindowTitle(w.WindowManagerHandle, title)
}

func (w *SDLAutoRenderer) SetScaleFactor(f int) {

}

func (w *SDLAutoRenderer) DestroyWindow() {
	libsdl.SDL_DestroyWindow(w.WindowManagerHandle)
}

func (gb *SDLAutoRenderer) setRenderDrawColor(c color.Color) {
	r, g, b, a := rgbaFromColor(c)
	libsdl.SDL_SetRenderDrawColor(gb.RendererHandle, r, g, b, a)
}

func (g *SDLAutoRenderer) SetScreenSize(width, height int) {
	libsdl.SDL_RenderSetLogicalSize(g.RendererHandle, width, height)
}

func (g SDLAutoRenderer) DrawLine(x1, y1, x2, y2 int, c color.Color) {
	g.setRenderDrawColor(c)
	libsdl.SDL_RenderDrawLine(g.RendererHandle, x1, y1, x2, y2)
}

func (g SDLAutoRenderer) DrawPoint(x, y int, c color.Color) {
	g.setRenderDrawColor(c)
	libsdl.SDL_RenderDrawPoint(g.RendererHandle, x, y)
}

func (g SDLAutoRenderer) FillTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) {
	g.FillPolygon([]int{x1, x2, x3}, []int{y1, y2, y3}, c)
}

func (g SDLAutoRenderer) DrawTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) {
	g.DrawLine(x1, y1, x2, y2, c)
	g.DrawLine(x2, y2, x3, y3, c)
	g.DrawLine(x3, y3, x1, y1, c)
}

func (g SDLAutoRenderer) FillPolygon(xPoints, yPoints []int, c color.Color) {
	numPoints := len(xPoints)
	if numPoints != len(yPoints) || numPoints < 3 {
		return
	}

	g.setRenderDrawColor(c)

	xMin := xPoints[0]
	xMax := xPoints[0]
	yMin := yPoints[0]
	yMax := yPoints[0]

	for i := 1; i < numPoints; i++ {
		if xPoints[i] < xMin {
			xMin = xPoints[i]
		} else if xPoints[i] > xMax {
			xMax = xPoints[i]
		}

		if yPoints[i] < yMin {
			yMin = yPoints[i]
		} else if yPoints[i] > yMax {
			yMax = yPoints[i]
		}
	}

	for x := xMin; x <= xMax; x++ {
		for y := yMin; y <= yMax; y++ {
			if g.pointInPolygon(x, y, xPoints, yPoints) {
				libsdl.SDL_RenderDrawPoint(g.RendererHandle, x, y)
			}
		}
	}
}

func (g SDLAutoRenderer) DrawPolygon(xPoints, yPoints []int, c color.Color) {
	numPoints := len(xPoints)
	if numPoints != len(yPoints) || numPoints < 3 {
		return
	}

	g.setRenderDrawColor(c)

	for i := 0; i < numPoints; i++ {
		nextIndex := (i + 1) % numPoints
		libsdl.SDL_RenderDrawLine(g.RendererHandle, xPoints[i], yPoints[i], xPoints[nextIndex], yPoints[nextIndex])
	}
}

func (g SDLAutoRenderer) pointInPolygon(x, y int, xPoints, yPoints []int) bool {
	numPoints := len(xPoints)
	oddNodes := false
	j := numPoints - 1

	for i := 0; i < numPoints; i++ {
		if yPoints[i] < y && yPoints[j] >= y || yPoints[j] < y && yPoints[i] >= y {
			if xPoints[i]+(y-yPoints[i])/(yPoints[j]-yPoints[i])*(xPoints[j]-xPoints[i]) < x {
				oddNodes = !oddNodes
			}
		}
		j = i
	}

	return oddNodes
}

func (g SDLAutoRenderer) FillRect(x, y, width, height int, c color.Color) {
	g.setRenderDrawColor(c)
	rect := &libsdl.SDL_Rect{
		X: int32(x),
		Y: int32(y),
		W: int32(width),
		H: int32(height),
	}
	libsdl.SDL_RenderFillRect(g.RendererHandle, rect)
}

func (g SDLAutoRenderer) DrawRect(x, y, width, height int, c color.Color) {
	g.setRenderDrawColor(c)

	// Top side
	libsdl.SDL_RenderDrawLine(g.RendererHandle, x, y, x+width, y)
	// Bottom side
	libsdl.SDL_RenderDrawLine(g.RendererHandle, x, y+height, x+width, y+height)
	// Left side
	libsdl.SDL_RenderDrawLine(g.RendererHandle, x, y, x, y+height)
	// Right side
	libsdl.SDL_RenderDrawLine(g.RendererHandle, x+width, y, x+width, y+height)
}

func (g SDLAutoRenderer) FillCirc(x, y, radius int, c color.Color) {
	for dx := -radius; dx <= radius; dx++ {
		for dy := -radius; dy <= radius; dy++ {
			if dx*dx+dy*dy <= radius*radius {
				g.DrawPoint(x+dx, y+dy, c)
			}
		}
	}
}

func (g SDLAutoRenderer) DrawCirc(xCenter, yCenter, radius int, c color.Color) {
	x := radius
	y := 0

	g.DrawPoint(xCenter+x, yCenter-y, c)
	g.DrawPoint(xCenter-x, yCenter+y, c)

	if radius > 0 {
		g.DrawPoint(xCenter-x, yCenter-y, c)
		g.DrawPoint(xCenter+x, yCenter+y, c)
		g.DrawPoint(xCenter+y, yCenter-x, c)
		g.DrawPoint(xCenter-y, yCenter+x, c)
	}

	p := 1 - radius
	for x > y {
		y++

		if p <= 0 {
			p = p + 2*y + 1
		} else {
			x--
			p = p + 2*y - 2*x + 1
		}

		if x < y {
			break
		}

		g.DrawPoint(xCenter+x, yCenter-y, c)
		g.DrawPoint(xCenter-x, yCenter+y, c)
		g.DrawPoint(xCenter-x, yCenter-y, c)
		g.DrawPoint(xCenter+x, yCenter+y, c)

		if x != y {
			g.DrawPoint(xCenter+y, yCenter-x, c)
			g.DrawPoint(xCenter-y, yCenter+x, c)
			g.DrawPoint(xCenter-y, yCenter-x, c)
			g.DrawPoint(xCenter+y, yCenter+x, c)
		}
	}
}

func (g *SDLAutoRenderer) CreateTextureFromImage(img image.Image) (uintptr, error) {
	if g.RendererHandle == 0 {
		return 0, fmt.Errorf("renderer doesn't exist")
	}
	var pixelFormat uint32
	if img == nil {
		return 0, fmt.Errorf("error creating texture because image is nil")
	}

	switch img := img.(type) {
	case *image.RGBA:
		pixelFormat = 0x16762004
	case *image.NRGBA:
		pixelFormat = 0x16362004
	default:
		return 0, fmt.Errorf("unsupported image type: %T", img)
	}

	texture := libsdl.SDL_CreateTexture(g.RendererHandle, pixelFormat, 1, img.Bounds().Dx(), img.Bounds().Dy())
	if texture == 0 {
		return 0, fmt.Errorf("failed to create texture")
	}

	libsdl.SDL_SetTextureBlendMode(texture, libsdl.SDL_BLENDMODE_BLEND)

	var pixels *uint8
	switch pImg := img.(type) {
	case *image.RGBA:
		pixels = &pImg.Pix[0]
	case *image.NRGBA:
		pixels = &pImg.Pix[0]
	default:
		// pixels = unsafe.Pointer(&pImg.Pix[0])
	}

	pitch := img.Bounds().Dx() * 4 // 4 bytes per pixel for RGBA
	libsdl.SDL_UpdateTexture(texture, nil, pixels, pitch)

	return texture, nil
}

func (g *SDLAutoRenderer) RenderTextureAt(texture uintptr, x, y, w, h int) {
	dstRect := &libsdl.SDL_Rect{
		X: int32(x),
		Y: int32(y),
		W: int32(w),
		H: int32(h),
	}

	libsdl.SDL_RenderCopy(texture, nil, dstRect)
}

func (g *SDLAutoRenderer) RenderTexture(texture uintptr, x, y, w, h int, angle float64, centerX, centerY int, flipType int) {
	dstRect := &libsdl.SDL_Rect{
		X: int32(x),
		Y: int32(y),
		W: int32(w),
		H: int32(h),
	}

	libsdl.SDL_RenderCopyEx(
		uintptr(g.RendererHandle),
		uintptr(texture),
		uintptr(unsafe.Pointer(nil)),
		uintptr(unsafe.Pointer(dstRect)),
		uintptr(*(*int)(unsafe.Pointer(&angle))),
		uintptr(unsafe.Pointer(&libsdl.SDL_Point{X: int32(centerX), Y: int32(centerY)})),
		flipType,
	)
}

func (g *SDLAutoRenderer) DestroyTexture(texture uintptr) {
	libsdl.SDL_DestroyTexture(texture)
}

func (g *SDLAutoRenderer) PollEvents(inputDevice input.InputDevice) bool {
	return libsdl.PollEvents(inputDevice)
}
