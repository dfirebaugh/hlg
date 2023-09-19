package ggez

import (
	"fmt"
	"image/color"
	"runtime"
	"time"

	"github.com/dfirebaugh/ggez/pkg/fb"
	"github.com/dfirebaugh/ggez/pkg/renderer"
	"github.com/dfirebaugh/ggez/pkg/renderer/sdl"
)

type Runner struct {
}

const (
	defaultWidth  = 800
	defaultHeight = 600
)

var (
	graphicsBackend renderer.GraphicsBackend
	windowTitle     string
	uifb            = fb.New(defaultWidth, defaultHeight)
	uiTexture       uintptr
)

func init() {
	runtime.LockOSThread()
	graphicsBackend, _ = sdl.New()

	SetTitle("ggez")
	SetScreenSize(defaultWidth, defaultHeight)
	setupDefaultInput()

	uiTexture, _ = graphicsBackend.CreateTextureFromImage(uifb.ToImage())
}

func Update(updateFn func()) {
	defer close()
	fpsCounter := NewFPSCounter()
	for {
		if !graphicsBackend.PollEvents(DefaultInput) {
			break
		}
		updateFn()

		uiTexture, _ = graphicsBackend.CreateTextureFromImage(uifb.ToImage())
		graphicsBackend.RenderTexture(uiTexture, 0, 0, defaultWidth, defaultHeight, 0, 0, 0, 0)

		graphicsBackend.RenderPresent()

		graphicsBackend.DestroyTexture(uiTexture)

		fpsCounter.Frame()
		fps := fpsCounter.GetFPS()
		title := windowTitle
		if fps != 0 && fpsEnabled {
			title = fmt.Sprintf("%s -- FPS: %d\n", title, int(fps))
		}
		graphicsBackend.SetWindowTitle(title)

		time.Sleep(5 * time.Millisecond)
	}
}

func close() {
	graphicsBackend.Close()
}

func Clear(c color.RGBA) {
	graphicsBackend.Clear(c)
}

func SetTitle(title string) {
	windowTitle = title
}

func SetScreenSize(width, height int) {
	graphicsBackend.SetScreenSize(width, height)
}

func ScreenWidth() int {
	return defaultWidth
}
func ScreenHeight() int {
	return defaultHeight
}
