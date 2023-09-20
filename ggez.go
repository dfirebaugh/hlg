package ggez

import (
	"fmt"
	"image/color"
	"runtime"
	"time"

	"github.com/dfirebaugh/ggez/pkg/fb"
	"github.com/dfirebaugh/ggez/pkg/renderer"
	"github.com/dfirebaugh/ggez/pkg/renderer/gl"
	"github.com/dfirebaugh/ggez/pkg/renderer/sdl"
)

type Runner struct {
}

var (
	screenWidth  = 240
	screenHeight = 160
)

var (
	graphicsBackend renderer.GraphicsBackend

	windowTitle string
	uifb        = fb.New(screenWidth, screenHeight)
	uiTexture   uintptr

	ConfiguredRenderer RendererType

	hasSetupCompleted = false
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

	switch ConfiguredRenderer {
	case SDLAutoRenderer:
		graphicsBackend, _ = sdl.New()
	case GLRenderer:
		graphicsBackend, _ = gl.New()
	}
	hasSetupCompleted = true

	SetTitle("ggez")
	SetScreenSize(screenWidth, screenHeight)
	setupDefaultInput()
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
	fpsCounter := NewFPSCounter()
	for {
		if !graphicsBackend.PollEvents(DefaultInput) {
			break
		}
		updateFn()

		uiTexture, _ = graphicsBackend.CreateTextureFromImage(uifb.ToImage())
		graphicsBackend.RenderTexture(uiTexture, 0, 0, screenWidth, screenHeight, 0, 0, 0, 0)

		graphicsBackend.Render()

		graphicsBackend.DestroyTexture(uiTexture)

		fpsCounter.Frame()
		fps := fpsCounter.GetFPS()
		title := windowTitle
		if fps != 0 && fpsEnabled {
			title = fmt.Sprintf("%s -- FPS: %d\n", title, int(fps))
		}
		graphicsBackend.SetWindowTitle(title)

		if ConfiguredRenderer == SDLAutoRenderer {
			time.Sleep(5 * time.Millisecond)
		}
	}
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
