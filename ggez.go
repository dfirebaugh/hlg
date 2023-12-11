package ggez

import (
	"fmt"
	"image/color"
	"runtime"
	"strings"

	"github.com/dfirebaugh/ggez/graphics"
	"github.com/dfirebaugh/ggez/graphics/webgpu"
	"github.com/dfirebaugh/ggez/pkg/fb"
)

type Runner struct {
}

var (
	windowWidth  = 240
	windowHeight = 160
)

var (
	graphicsBackend graphics.GraphicsBackend

	windowTitle string
	uifb        = fb.New(int(windowWidth), int(windowHeight))
	uiTexture   *Texture

	hasSetupCompleted = false
	fpsCounter        *FPSCounter
)

func Setup() {
	runtime.LockOSThread()

	graphicsBackend, _ = webgpu.NewGraphicsBackend()
	hasSetupCompleted = true
}

func initWindow() {
	SetWindowSize(windowWidth, windowHeight)
	SetScaleFactor(3)
	SetTitle("ggez")
}

func ensureSetupCompletion() {
	if hasSetupCompleted {
		return
	}
	Setup()
	initWindow()
}

func Update(updateFn func()) {
	ensureSetupCompletion()
	defer close()
	fpsCounter = NewFPSCounter()

	uifb = fb.New(int(windowWidth), int(windowHeight))
	uiTexture, _ = CreateTextureFromImage(uifb.ToImage())

	var err error
	for {
		if !graphicsBackend.PollEvents() {
			break
		}

		uiTexture.UpdateTextureFromImage(uifb.ToImage())
		uiTexture.Render()
		updateFn()
		graphicsBackend.Render()

		calculateFPS()
		if err != nil {
			fmt.Println("error occured while rendering:", err)

			errstr := err.Error()
			switch {
			case strings.Contains(errstr, "Surface timed out"): // do nothing
			case strings.Contains(errstr, "Surface is outdated"): // do nothing
			case strings.Contains(errstr, "Surface was lost"): // do nothing
			default:
				panic(err)
			}
		}
	}
	uiTexture.Destroy()
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

func GetFPS() float64 {
	return fpsCounter.GetFPS()
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

func GetWindowSize() (int, int) {
	return graphicsBackend.GetWindowSize()
}

func SetScreenSize(width, height int) {
	ensureSetupCompletion()
	graphicsBackend.SetScreenSize(width, height)
}

func SetWindowSize(width, height int) {
	ensureSetupCompletion()

	windowWidth = width
	windowHeight = height

	graphicsBackend.SetWindowSize(windowWidth, windowHeight)
}

func SetScaleFactor(f int) {
	graphicsBackend.SetScaleFactor(f)
}

func ToggleWireFrame() {
	// graphicsBackend.ToggleWireframeMode()
}

func ScreenHeight() int {
	return graphicsBackend.ScreenHeight()
}
func ScreenWidth() int {
	return graphicsBackend.ScreenWidth()
}
