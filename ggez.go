package ggez

import (
	"fmt"
	"image/color"
	"runtime"
	"strings"

	"github.com/dfirebaugh/ggez/graphics"
	"github.com/dfirebaugh/ggez/graphics/webgpu"
	"github.com/dfirebaugh/ggez/pkg/draw"
	"github.com/dfirebaugh/ggez/pkg/fb"
	"github.com/dfirebaugh/ggez/pkg/input"
	"github.com/dfirebaugh/ggez/pkg/math/geom"
)

var (
	windowWidth  = 240
	windowHeight = 160
)

type engine struct {
	graphicsBackend   graphics.GraphicsBackend
	inputState        *input.InputState
	windowTitle       string
	uifb              *fb.ImageFB
	uiTexture         *Texture
	fpsCounter        *FPSCounter
	hasSetupCompleted bool
}

var (
	ggez = &engine{}
)

func Setup() {
	runtime.LockOSThread()
	ggez.inputState = input.NewInputState()
	ggez.uifb = fb.New(int(windowWidth), int(windowHeight))
	var err error
	ggez.graphicsBackend, err = webgpu.NewGraphicsBackend()
	if err != nil {
		panic(err.Error())
	}
	ggez.hasSetupCompleted = true
}

func initWindow() {
	SetWindowSize(windowWidth, windowHeight)
	SetScaleFactor(3)
	SetTitle("ggez")
}

func ensureSetupCompletion() {
	if ggez.hasSetupCompleted {
		return
	}
	Setup()
	initWindow()
}

func Update(updateFn func()) {
	ensureSetupCompletion()
	defer close()
	ggez.fpsCounter = NewFPSCounter()

	ggez.uifb = fb.New(int(windowWidth), int(windowHeight))
	ggez.uiTexture, _ = CreateTextureFromImage(ggez.uifb.ToImage())
	defer ggez.uiTexture.Destroy()

	ggez.graphicsBackend.SetInputCallback(func(eventChan chan input.Event) {
		evt := <-eventChan
		handleEvent(evt, ggez.inputState)
	})

	var err error
	for {
		if !ggez.graphicsBackend.PollEvents() {
			break
		}

		updateFn()
		ggez.uiTexture.UpdateTextureFromImage(ggez.uifb.ToImage())
		ggez.uiTexture.Render()
		ggez.graphicsBackend.Render()

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
}

func calculateFPS() {
	ggez.fpsCounter.Frame()
	fps := ggez.fpsCounter.GetFPS()
	title := ggez.windowTitle
	if fps != 0 && fpsEnabled {
		title = fmt.Sprintf("%s -- FPS: %d\n", title, int(fps))
	}
	if ggez.graphicsBackend.IsDisposed() {
		return
	}
	ggez.graphicsBackend.SetWindowTitle(title)
}

func GetFPS() float64 {
	return ggez.fpsCounter.GetFPS()
}

func close() {
	ggez.graphicsBackend.Close()
}

func Clear(c color.RGBA) {
	ensureSetupCompletion()
	ggez.graphicsBackend.Clear(c)
	draw.Rect(geom.MakeRect(0, 0, float32(windowWidth), float32(windowHeight))).Fill(ggez.uifb, color.RGBA{0, 0, 0, 0})
}

func SetTitle(title string) {
	ensureSetupCompletion()
	ggez.windowTitle = title
}

func GetWindowSize() (int, int) {
	return ggez.graphicsBackend.GetWindowSize()
}

func SetScreenSize(width, height int) {
	ensureSetupCompletion()
	ggez.graphicsBackend.SetScreenSize(width, height)
}

func SetWindowSize(width, height int) {
	ensureSetupCompletion()

	windowWidth = width
	windowHeight = height

	ggez.graphicsBackend.SetWindowSize(windowWidth, windowHeight)
}

func SetScaleFactor(f int) {
	ggez.graphicsBackend.SetScaleFactor(f)
}

func ToggleWireFrame() {
	// graphicsBackend.ToggleWireframeMode()
}

func ScreenSize() (int, int) {
	return ggez.graphicsBackend.ScreenSize()
}
