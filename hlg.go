package hlg

import (
	"fmt"
	"image/color"
	"runtime"
	"strings"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/webgpu"
	"github.com/dfirebaugh/hlg/pkg/draw"
	"github.com/dfirebaugh/hlg/pkg/fb"
	"github.com/dfirebaugh/hlg/pkg/input"
	"github.com/dfirebaugh/hlg/pkg/math/geom"
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
	fpsCounter        *fpsCounter
	hasSetupCompleted bool
}

var (
	hlg = &engine{}
)

func setup() {
	runtime.LockOSThread()
	hlg.inputState = input.NewInputState()
	hlg.uifb = fb.New(int(windowWidth), int(windowHeight))
	var err error
	hlg.graphicsBackend, err = webgpu.NewGraphicsBackend()
	if err != nil {
		panic(err.Error())
	}
	hlg.hasSetupCompleted = true
}

func initWindow() {
	SetWindowSize(windowWidth, windowHeight)
	SetScaleFactor(3)
	SetTitle("hlg")
}

func ensureSetupCompletion() {
	if hlg.hasSetupCompleted {
		return
	}
	setup()
	initWindow()
}

// Update is the main update function called to refresh the engine state.
func Update(updateFn func()) {
	ensureSetupCompletion()
	defer close()
	hlg.fpsCounter = newFPSCounter()

	hlg.uifb = fb.New(int(windowWidth), int(windowHeight))
	hlg.uiTexture, _ = CreateTextureFromImage(hlg.uifb.ToImage())
	defer hlg.uiTexture.Destroy()

	hlg.graphicsBackend.SetInputCallback(func(eventChan chan input.Event) {
		evt := <-eventChan
		handleEvent(evt, hlg.inputState)
	})

	var err error
	for {
		if !hlg.graphicsBackend.PollEvents() {
			break
		}

		updateFn()
		hlg.uiTexture.UpdateTextureFromImage(hlg.uifb.ToImage())
		hlg.uiTexture.Render()
		hlg.graphicsBackend.Render()

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
	hlg.fpsCounter.Frame()
	fps := hlg.fpsCounter.GetFPS()
	title := hlg.windowTitle
	if fps != 0 && fpsEnabled {
		title = fmt.Sprintf("%s -- FPS: %d\n", title, int(fps))
	}
	if hlg.graphicsBackend.IsDisposed() {
		return
	}
	hlg.graphicsBackend.SetWindowTitle(title)
}

func GetFPS() float64 {
	return hlg.fpsCounter.GetFPS()
}

func close() {
	hlg.graphicsBackend.Close()
}

// Clear clears the screen with the specified color.
func Clear(c color.RGBA) {
	ensureSetupCompletion()
	hlg.graphicsBackend.Clear(c)
	draw.Rect(geom.MakeRect(0, 0, float32(windowWidth), float32(windowHeight))).Fill(hlg.uifb, color.RGBA{0, 0, 0, 0})
}

// SetTitle sets the title of the window.
func SetTitle(title string) {
	ensureSetupCompletion()
	hlg.windowTitle = title
}

// GetWindowSize retrieves the current window size.
func GetWindowSize() (int, int) {
	return hlg.graphicsBackend.GetWindowSize()
}

func GetWindowPosition() (int, int) {
	return hlg.graphicsBackend.GetWindowPosition()
}

// SetScreenSize sets the size of the screen.
func SetScreenSize(width, height int) {
	ensureSetupCompletion()
	hlg.graphicsBackend.SetScreenSize(width, height)
}

// SetWindowSize sets the size of the window.
func SetWindowSize(width, height int) {
	ensureSetupCompletion()

	windowWidth = width
	windowHeight = height

	hlg.graphicsBackend.SetWindowSize(windowWidth, windowHeight)
}

// SetScaleFactor sets the scale factor for the window.
func SetScaleFactor(f int) {
	hlg.graphicsBackend.SetScaleFactor(f)
}

// ScreenSize returns the current screen size.
func ScreenSize() (int, int) {
	return hlg.graphicsBackend.ScreenSize()
}
