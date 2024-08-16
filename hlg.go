package hlg

import (
	"fmt"
	"image/color"
	"runtime"
	"time"

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

var hlg = &engine{}

func setup() {
	runtime.LockOSThread()
	hlg.inputState = input.NewInputState()
	hlg.uifb = fb.New(int(windowWidth), int(windowHeight))
	var err error
	hlg.graphicsBackend, err = webgpu.NewGraphicsBackend(windowWidth, windowHeight)
	if err != nil {
		panic(err.Error())
	}
	hlg.hasSetupCompleted = true
}

func initWindow() {
	SetWindowSize(windowWidth, windowHeight)
	SetTitle("hlg")
}

func ensureSetupCompletion() {
	if hlg.hasSetupCompleted {
		return
	}
	setup()
	initWindow()
}

// Run is the main update function called to refresh the engine state.
func Run(updateFn func(), renderFn func()) {
	ensureSetupCompletion()
	defer close()
	hlg.fpsCounter = newFPSCounter()

	sw, sh := hlg.graphicsBackend.GetScreenSize()
	hlg.uifb = fb.New(sw, sh)
	hlg.uiTexture, _ = CreateTextureFromImage(hlg.uifb.ToImage())
	defer hlg.uiTexture.Destroy()

	hlg.graphicsBackend.SetInputCallback(func(eventChan chan input.Event) {
		evt := <-eventChan
		handleEvent(evt, hlg.inputState)
	})
	targetFPS := 120.0
	targetFrameDuration := time.Second / time.Duration(targetFPS)

	var lastUpdateTime time.Time
	var accumulator time.Duration

	lastUpdateTime = time.Now()
	for hlg.graphicsBackend.PollEvents() {
		currentTime := time.Now()
		deltaTime := currentTime.Sub(lastUpdateTime)
		lastUpdateTime = currentTime
		accumulator += deltaTime

		frameRendered := false
		for accumulator >= targetFrameDuration {
			updateFn()
			accumulator -= targetFrameDuration
			frameRendered = true
		}

		if frameRendered {
			hlg.uiTexture.UpdateTextureFromImage(hlg.uifb.ToImage())
			hlg.uiTexture.Render()
			hlg.graphicsBackend.Render()
			renderFn()

			calculateFPS()
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

func GetScreenSize() (int, int) {
	ensureSetupCompletion()
	return hlg.graphicsBackend.GetScreenSize()
}

// SetScreenSize sets the size of the screen.
func SetScreenSize(width, height int) {
	ensureSetupCompletion()
	hlg.graphicsBackend.SetScreenSize(width, height)
}

// SetWindowSize sets the size of the window.
func SetWindowSize(width, height int) {
	ensureSetupCompletion()

	hlg.graphicsBackend.SetScreenSize(width, height)
	hlg.graphicsBackend.SetWindowSize(width, height)
}
