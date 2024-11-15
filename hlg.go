package hlg

import (
	"fmt"
	"image/color"
	"runtime"
	"time"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/webgpu"
	"github.com/dfirebaugh/hlg/pkg/input"
)

type Game interface {
	Update()
	Render()
}

var (
	windowWidth  = 240
	windowHeight = 160
)

type engine struct {
	graphicsBackend   graphics.GraphicsBackend
	inputState        *input.InputState
	windowTitle       string
	fpsCounter        *fpsCounter
	hasSetupCompleted bool
}

var hlg = &engine{}

func setup() {
	runtime.LockOSThread()
	hlg.inputState = input.NewInputState()
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

func run(updateFn func(), renderFn func()) {
	ensureSetupCompletion()
	defer close()
	hlg.fpsCounter = newFPSCounter()

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
			if updateFn != nil {
				updateFn()
			}
			accumulator -= targetFrameDuration
			frameRendered = true

			hlg.inputState.ResetJustPressed()
		}

		if frameRendered {
			hlg.graphicsBackend.Render()
			if renderFn != nil {
				renderFn()
			}

			calculateFPS()
		}
	}
}

func RunGame(game Game) {
	run(game.Update, game.Render)
}

func RunApp(game Game) {
	run(game.Update, game.Render)
}

// Run is the main update function called to refresh the engine state.
func Run(updateFn func(), renderFn func()) {
	run(updateFn, renderFn)
}

func calculateFPS() {
	hlg.fpsCounter.Frame()
	fps := hlg.fpsCounter.GetFPS()
	title := hlg.windowTitle
	if fps != 0 && fpsEnabled {
		title = fmt.Sprintf("%s -- %d\n", title, int(fps))
	}
	if hlg.graphicsBackend.IsDisposed() {
		return
	}
	hlg.graphicsBackend.SetWindowTitle(title)
}

func GetFPS() float64 {
	return hlg.fpsCounter.GetFPS()
}

func Close() {
	close()
}

func close() {
	hlg.graphicsBackend.Close()
}

// Clear clears the screen with the specified color.
func Clear(c color.RGBA) {
	ensureSetupCompletion()
	hlg.graphicsBackend.Clear(c)
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

func GetWindowWidth() int {
	w, _ := hlg.graphicsBackend.GetWindowSize()
	return w
}

func GetWindowHeight() int {
	_, h := hlg.graphicsBackend.GetWindowSize()
	return h
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
	windowWidth = width
	windowHeight = height
}

func CreateRenderQueue() graphics.RenderQueue {
	return hlg.graphicsBackend.CreateRenderQueue()
}

func DisableWindowResize() {
	ensureSetupCompletion()

	hlg.graphicsBackend.DisableWindowResize()

	SetWindowSize(windowWidth, windowHeight)
}

func SetBorderlessWindowed(v bool) {
	hlg.graphicsBackend.SetBorderlessWindowed(v)
}
