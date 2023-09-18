package sdl

import (
	"fmt"
	"image"
	"image/color"
	"runtime"
	"unsafe"

	"github.com/dfirebaugh/ggez/pkg/input"
	"github.com/ebitengine/purego"
)

type graphicsBackend struct {
	WindowManager
	EventManager
	screenHeight        int
	screenWidth         int
	RendererHandle      uintptr
	WindowManagerHandle uintptr
}

const (
	SDL_INIT_VIDEO            = 0x00000020
	SDL_WINDOW_SHOWN          = 0x00000004
	SDL_RENDERER_ACCELERATED  = 0x00000002
	SDL_WINDOWPOS_CENTERED    = 0x2FFF0000
	SDL_WINDOW_RESIZABLE      = 0x00000020
	SDL_RENDERER_PRESENTVSYNC = 0x00000004
	SDL_BLENDMODE_NONE        = 0
	SDL_BLENDMODE_BLEND       = 1
	SDL_BLENDMODE_ADD         = 2
	SDL_BLENDMODE_MOD         = 3

	SDL_QUIT                = 0x100
	SDL_WINDOWEVENT         = 0x200
	SDL_WINDOWEVENT_RESIZED = 5
)

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
	libSDL uintptr
)

var (
	SDL_Init                 func(flags uint32)
	SDL_Quit                 func()
	SDL_CreateWindow         func(title string, x int, y int, w int, h int, flags uint32) uintptr
	SDL_DestroyWindow        func(window uintptr)
	SDL_CreateRenderer       func(window uintptr, index int, flags uint32) uintptr
	SDL_DestroyRenderer      func(renderer uintptr)
	SDL_SetRenderDrawColor   func(renderer uintptr, r uint8, g uint8, b uint8, a uint8)
	SDL_RenderDrawLine       func(renderer uintptr, x1, y1, x2, y2 int)
	SDL_RenderPresent        func(renderer uintptr)
	SDL_CreateTexture        func(renderer uintptr, format uint32, access int, w int, h int) uintptr
	SDL_UpdateTexture        func(texture uintptr, rect *SDL_Rect, pixels *uint8, pitch int)
	SDL_RenderCopy           func(renderer uintptr, src *SDL_Rect, dst *SDL_Rect)
	SDL_RenderCopyEx         func(renderer uintptr, texture uintptr, src uintptr, dst uintptr, uintptr, angle uintptr, flip int)
	SDL_DestroyTexture       func(texture uintptr)
	SDL_RenderClear          func(renderer uintptr)
	SDL_PollEvent            func(event *SDL_Event) int
	SDL_RenderSetLogicalSize func(renderer uintptr, w int, h int)
	SDL_SetTextureBlendMode  func(texture uintptr, blendMode uint8)
	SDL_SetWindowTitle       func(window uintptr, title string)
	SDL_GetRendererInfo      func(renderer uintptr, info *SDL_RendererInfo)
	SDL_GetNumRenderDrivers  func()
	SDL_GetPlatform          func()
	SDL_GetVersion           func(ver *SDL_Version)
	SDL_RenderFillRect       func(renderer uintptr, rect *SDL_Rect)
	SDL_RenderDrawPoint      func(renderer uintptr, x int, y int)
)

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
	purego.RegisterLibFunc(&SDL_CreateWindow, libSDL, "SDL_CreateWindow")
	purego.RegisterLibFunc(&SDL_DestroyWindow, libSDL, "SDL_DestroyWindow")
	purego.RegisterLibFunc(&SDL_CreateRenderer, libSDL, "SDL_CreateRenderer")
	purego.RegisterLibFunc(&SDL_DestroyRenderer, libSDL, "SDL_DestroyRenderer")
	purego.RegisterLibFunc(&SDL_SetRenderDrawColor, libSDL, "SDL_SetRenderDrawColor")
	purego.RegisterLibFunc(&SDL_RenderDrawLine, libSDL, "SDL_RenderDrawLine")
	purego.RegisterLibFunc(&SDL_RenderPresent, libSDL, "SDL_RenderPresent")
	purego.RegisterLibFunc(&SDL_CreateTexture, libSDL, "SDL_CreateTexture")
	purego.RegisterLibFunc(&SDL_UpdateTexture, libSDL, "SDL_UpdateTexture")
	purego.RegisterLibFunc(&SDL_RenderCopy, libSDL, "SDL_RenderCopy")
	purego.RegisterLibFunc(&SDL_RenderCopyEx, libSDL, "SDL_RenderCopyEx")
	purego.RegisterLibFunc(&SDL_DestroyTexture, libSDL, "SDL_DestroyTexture")
	purego.RegisterLibFunc(&SDL_RenderClear, libSDL, "SDL_RenderClear")
	purego.RegisterLibFunc(&SDL_PollEvent, libSDL, "SDL_PollEvent")
	purego.RegisterLibFunc(&SDL_RenderSetLogicalSize, libSDL, "SDL_RenderSetLogicalSize")
	purego.RegisterLibFunc(&SDL_SetTextureBlendMode, libSDL, "SDL_SetTextureBlendMode")
	purego.RegisterLibFunc(&SDL_SetWindowTitle, libSDL, "SDL_SetWindowTitle")
	purego.RegisterLibFunc(&SDL_GetRendererInfo, libSDL, "SDL_GetRendererInfo")
	purego.RegisterLibFunc(&SDL_GetNumRenderDrivers, libSDL, "SDL_GetNumRenderDrivers")
	purego.RegisterLibFunc(&SDL_GetPlatform, libSDL, "SDL_GetPlatform")
	purego.RegisterLibFunc(&SDL_GetVersion, libSDL, "SDL_GetVersion")
	purego.RegisterLibFunc(&SDL_RenderFillRect, libSDL, "SDL_RenderFillRect")
	purego.RegisterLibFunc(&SDL_RenderDrawPoint, libSDL, "SDL_RenderDrawPoint")
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

type SDL_RendererInfo struct {
	Name              *byte
	Flags             uint32
	NumTextureFormats uint32
	TextureFormats    [16]uint32
	MaxTextureWidth   int32
	MaxTextureHeight  int32
}

type SDL_Version struct {
	Major uint8
	Minor uint8
	Patch uint8
}

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

type SDL_Point struct {
	X, Y int32
}

type SDL_Rect struct {
	X, Y int32
	W, H int32
}

var (
	screenHeight = 600
	screenWidth  = 800
)

func New() (*graphicsBackend, error) {
	SDL_Init(SDL_INIT_VIDEO)
	g := &graphicsBackend{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}

	g.WindowManager.CreateWindow("", screenWidth, screenHeight)
	g.createRenderer()
	return g, nil
}

func (g *graphicsBackend) Close() {
	g.WindowManager.DestroyWindow()
	g.DestroyRenderer()
	SDL_Quit()
}

func (g *graphicsBackend) createRenderer() error {
	r := SDL_CreateRenderer(g.WindowManager.WindowManagerHandle, 0xFFFFFFFF, SDL_RENDERER_ACCELERATED)
	if r == 0 {
		return fmt.Errorf("SDL_CreateRenderer failed")
	}
	g.RendererHandle = r
	g.SetScreenSize(g.screenWidth, g.screenHeight)
	return nil
}

func (g *graphicsBackend) PrintPlatformAndVersion() {

}

type EventManager struct{}

// pollEvents polls for SDL events and returns false if a quit event is detected.
func (EventManager) PollEvents(inputDevice input.InputDevice) bool {
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

func rgbaFromColor(c color.Color) (r, g, b, a uint8) {
	R, G, B, A := c.RGBA()
	return uint8(R >> 8), uint8(G >> 8), uint8(B >> 8), uint8(A >> 8)
}

func (gb *graphicsBackend) setRenderDrawColor(c color.Color) {
	r, g, b, a := rgbaFromColor(c)
	SDL_SetRenderDrawColor(gb.RendererHandle, r, g, b, a)
}

func (g *graphicsBackend) Clear(c color.Color) {
	g.setRenderDrawColor(c)
	SDL_RenderClear(g.RendererHandle)
}

func (g *graphicsBackend) SetScreenSize(width, height int) {
	SDL_RenderSetLogicalSize(g.RendererHandle, width, height)
}

func (g *graphicsBackend) RenderPresent() {
	SDL_RenderPresent(g.RendererHandle)
}

func (g *graphicsBackend) DestroyRenderer() {
	SDL_DestroyRenderer(g.RendererHandle)
}

func (g *graphicsBackend) PrintRendererInfo() {
	var info SDL_RendererInfo

	SDL_GetRendererInfo(g.RendererHandle, &info)

	rendererName := goStringFromCString(info.Name)
	fmt.Printf("Backend: %s\n", rendererName)
}

// Convert a C-style string (null-terminated) to a Go string
func goStringFromCString(str *byte) string {
	if str == nil {
		return ""
	}
	var buffer []byte
	for i := str; *i != 0; i = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(i)) + 1)) {
		buffer = append(buffer, *i)
	}
	return string(buffer)
}

func (g graphicsBackend) DrawLine(x1, y1, x2, y2 int, c color.Color) {
	g.setRenderDrawColor(c)
	SDL_RenderDrawLine(g.RendererHandle, x1, y1, x2, y2)
}

func (g graphicsBackend) DrawPoint(x, y int, c color.Color) {
	g.setRenderDrawColor(c)
	SDL_RenderDrawPoint(g.RendererHandle, x, y)
}

// FillTriangle fills a triangle with a specific color.
func (g graphicsBackend) FillTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) {
	g.FillPolygon([]int{x1, x2, x3}, []int{y1, y2, y3}, c)
}

// DrawTriangle draws an unfilled triangle with a specific color.
func (g graphicsBackend) DrawTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) {
	g.DrawLine(x1, y1, x2, y2, c)
	g.DrawLine(x2, y2, x3, y3, c)
	g.DrawLine(x3, y3, x1, y1, c)
}

// FillPolygon fills a polygon with a specific color.
func (g graphicsBackend) FillPolygon(xPoints, yPoints []int, c color.Color) {
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
				SDL_RenderDrawPoint(g.RendererHandle, x, y)
			}
		}
	}
}

// DrawPolygon draws the outline of a polygon with a specific color.
func (g graphicsBackend) DrawPolygon(xPoints, yPoints []int, c color.Color) {
	numPoints := len(xPoints)
	if numPoints != len(yPoints) || numPoints < 3 {
		return
	}

	g.setRenderDrawColor(c)

	for i := 0; i < numPoints; i++ {
		nextIndex := (i + 1) % numPoints // Wrap around to close the polygon
		SDL_RenderDrawLine(g.RendererHandle, xPoints[i], yPoints[i], xPoints[nextIndex], yPoints[nextIndex])
	}
}

func (g graphicsBackend) pointInPolygon(x, y int, xPoints, yPoints []int) bool {
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

// FillRect fills a rectangle with a specific color.
func (g graphicsBackend) FillRect(x, y, width, height int, c color.Color) {
	g.setRenderDrawColor(c)
	rect := &SDL_Rect{
		X: int32(x),
		Y: int32(y),
		W: int32(width),
		H: int32(height),
	}
	SDL_RenderFillRect(g.RendererHandle, rect)
}

// DrawRect draws the outline of a rectangle with a specific color.
func (g graphicsBackend) DrawRect(x, y, width, height int, c color.Color) {
	g.setRenderDrawColor(c)

	// Top side
	SDL_RenderDrawLine(g.RendererHandle, x, y, x+width, y)
	// Bottom side
	SDL_RenderDrawLine(g.RendererHandle, x, y+height, x+width, y+height)
	// Left side
	SDL_RenderDrawLine(g.RendererHandle, x, y, x, y+height)
	// Right side
	SDL_RenderDrawLine(g.RendererHandle, x+width, y, x+width, y+height)
}

// FillCirc fills a circle with a specific color.
func (g graphicsBackend) FillCirc(x, y, radius int, c color.Color) {
	// Loop through each pixel in a circle and draw it
	for dx := -radius; dx <= radius; dx++ {
		for dy := -radius; dy <= radius; dy++ {
			if dx*dx+dy*dy <= radius*radius {
				g.DrawPoint(x+dx, y+dy, c)
			}
		}
	}
}

// DrawCirc draws the outline of a circle with a specific color.
func (g graphicsBackend) DrawCirc(xCenter, yCenter, radius int, c color.Color) {
	x := radius
	y := 0

	// Printing the initial point on the circle
	// at the end of a given radius
	g.DrawPoint(xCenter+x, yCenter-y, c)
	g.DrawPoint(xCenter-x, yCenter+y, c)

	// When the radius is zero, only a single point will be printed at the center
	if radius > 0 {
		g.DrawPoint(xCenter-x, yCenter-y, c)
		g.DrawPoint(xCenter+x, yCenter+y, c)
		g.DrawPoint(xCenter+y, yCenter-x, c)
		g.DrawPoint(xCenter-y, yCenter+x, c)
	}

	// Initial point on the circle at the end of radius
	p := 1 - radius
	for x > y {
		y++

		// Mid-point inside or on the perimeter of the circle
		if p <= 0 {
			p = p + 2*y + 1
		} else { // Mid-point outside the perimeter of the circle
			x--
			p = p + 2*y - 2*x + 1
		}

		// All perimeter points have only set of symmetric points
		if x < y {
			break
		}

		// Printing the generated point and its reflection in the other octants
		g.DrawPoint(xCenter+x, yCenter-y, c)
		g.DrawPoint(xCenter-x, yCenter+y, c)
		g.DrawPoint(xCenter-x, yCenter-y, c)
		g.DrawPoint(xCenter+x, yCenter+y, c)

		// If the generated point is on the line x = y, then the perimeter points have only one set of symmetric points
		if x != y {
			g.DrawPoint(xCenter+y, yCenter-x, c)
			g.DrawPoint(xCenter-y, yCenter+x, c)
			g.DrawPoint(xCenter-y, yCenter-x, c)
			g.DrawPoint(xCenter+y, yCenter+x, c)
		}
	}
}

type RendererFlip int

const (
	FLIP_NONE RendererFlip = iota
	FLIP_HORIZONTAL
	FLIP_VERTICAL
	FLIP_BOTH
)

// CreateTextureFromImage creates an SDL texture based on the provided image format.
func (g *graphicsBackend) CreateTextureFromImage(img image.Image) (uintptr, error) {
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
		pixelFormat = 0x16362004 // You can adjust this if necessary.
	default:
		return 0, fmt.Errorf("unsupported image type: %T", img)
	}

	texture := SDL_CreateTexture(g.RendererHandle, pixelFormat, 1, img.Bounds().Dx(), img.Bounds().Dy())
	if texture == 0 {
		return 0, fmt.Errorf("failed to create texture")
	}

	SDL_SetTextureBlendMode(texture, SDL_BLENDMODE_BLEND)

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
	SDL_UpdateTexture(texture, nil, pixels, pitch)

	return texture, nil
}

func (g *graphicsBackend) RenderTextureAt(texture uintptr, x, y, w, h int) {
	dstRect := &SDL_Rect{
		X: int32(x),
		Y: int32(y),
		W: int32(w),
		H: int32(h),
	}

	SDL_RenderCopy(texture, nil, dstRect)
}

// // RenderTexture renders a texture to the screen.
// func (g *graphicsBackend) RenderTexture(texture uintptr) {
// 	SDL_RenderCopy(g.RendererHandle, nil, nil)
// }

func (g *graphicsBackend) RenderTexture(texture uintptr, x, y, w, h int, angle float64, centerX, centerY int, flipType int) {
	dstRect := &SDL_Rect{
		X: int32(x),
		Y: int32(y),
		W: int32(w),
		H: int32(h),
	}

	SDL_RenderCopyEx(
		uintptr(g.RendererHandle),
		uintptr(texture),
		uintptr(unsafe.Pointer(nil)),
		uintptr(unsafe.Pointer(dstRect)),
		uintptr(*(*int)(unsafe.Pointer(&angle))),
		uintptr(unsafe.Pointer(&SDL_Point{X: int32(centerX), Y: int32(centerY)})),
		flipType,
	)
}

// DestroyTexture destroys an SDL texture.
func (g *graphicsBackend) DestroyTexture(texture uintptr) {
	SDL_DestroyTexture(texture)
}

type WindowManager struct {
	WindowManagerHandle uintptr
}

// CreateWindow creates an SDL window.
func (w *WindowManager) CreateWindow(title string, width, height int) (uintptr, error) {
	window := SDL_CreateWindow(
		title,
		SDL_WINDOWPOS_CENTERED, SDL_WINDOWPOS_CENTERED, width, height, SDL_WINDOW_SHOWN|SDL_WINDOW_RESIZABLE,
	)
	if window == 0 {
		return 0, fmt.Errorf("SDL_CreateWindow failed")
	}

	w.WindowManagerHandle = window
	w.SetWindowTitle(title)

	return window, nil
}

func (w *WindowManager) SetWindowTitle(title string) {
	if len(title) == 0 {
		return
	}
	if w.WindowManagerHandle == 0 {
		return
	}
	SDL_SetWindowTitle(w.WindowManagerHandle, title)
}

// DestroyWindow destroys an SDL window.
func (w *WindowManager) DestroyWindow() {
	SDL_DestroyWindow(w.WindowManagerHandle)
}
