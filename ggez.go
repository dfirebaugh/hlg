package ggez

import (
	"fmt"
	"image/color"
	"runtime"

	"github.com/dfirebaugh/ggez/pkg/fb"
	"github.com/dfirebaugh/ggez/pkg/graphics"
	"github.com/dfirebaugh/ggez/pkg/graphics/gl"
)

type Runner struct {
}

var (
	screenWidth  = 240
	screenHeight = 160
)

var (
	graphicsBackend graphics.GraphicsBackend

	windowTitle string
	uifb        = fb.New(screenWidth, screenHeight)
	uiTexture   uintptr

	ConfiguredRenderer RendererType

	hasSetupCompleted = false
	fpsCounter        *FPSCounter
)

type RendererType uint

const (
	// SDLAutoRenderer let's SDL pick the appropriate renderer
	SDLAutoRenderer = iota
	// Uses OpenGL Renderer
	GLRenderer
)

func SetRenderer(t RendererType) {
	ConfiguredRenderer = t
}

func Setup(t RendererType) {
	SetRenderer(t)
	runtime.LockOSThread()

	graphicsBackend, _ = gl.New()
	hasSetupCompleted = true

	SetTitle("ggez")
	SetScreenSize(screenWidth, screenHeight)
	SetScaleFactor(3)

	uiTexture, _ = graphicsBackend.CreateTextureFromImage(uifb.ToImage())
}

func ensureSetupCompletion() {
	if !hasSetupCompleted {
		Setup(ConfiguredRenderer)
	}
}

func Update(updateFn func()) {
	ensureSetupCompletion()
	defer close()
	fpsCounter = NewFPSCounter()
	uiTexture, _ = graphicsBackend.CreateTextureFromImage(uifb.ToImage())
	for {
		if !graphicsBackend.PollEvents() {
			break
		}
		updateFn()

		// graphicsBackend.RenderTexture(uiTexture, 0, 0, screenWidth, screenHeight, 0, 0, 0, 0)
		graphicsBackend.Render()

		calculateFPS()
	}
	graphicsBackend.DestroyTexture(uiTexture)
}

func calculateFPS() {
	fpsCounter.Frame()
	fps := fpsCounter.GetFPS()
	title := windowTitle
	if fps != 0 && fpsEnabled {
		title = fmt.Sprintf("%s -- FPS: %d\n", title, int(fps))
	}
	graphicsBackend.SetWindowTitle(title)
}

func close() {
	graphicsBackend.Close()
}

func Clear(c color.RGBA) {
	ensureSetupCompletion()
	graphicsBackend.Clear(c)
}

func SetTitle(title string) {
	ensureSetupCompletion()
	windowTitle = title
}

func SetScreenSize(width, height int) {
	ensureSetupCompletion()

	screenWidth = width
	screenHeight = width

	graphicsBackend.SetScreenSize(screenWidth, screenHeight)
}

func ScreenWidth() int {
	return screenWidth
}
func ScreenHeight() int {
	return screenHeight
}

func SetScaleFactor(f int) {
	graphicsBackend.SetScaleFactor(f)
}

func ToggleWireFrame() {
	graphicsBackend.ToggleWireframeMode()
}
