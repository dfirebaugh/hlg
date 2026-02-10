//go:build !js

package hlg

import (
	"fmt"
	"image/color"
	"runtime"
	"time"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/gl"
	"github.com/dfirebaugh/hlg/graphics/webgpu"
	"github.com/dfirebaugh/hlg/pkg/input"
)

// Backend represents the graphics backend type
type Backend int

const (
	BackendOpenGL Backend = iota
	BackendWebGPU
)

var selectedBackend = BackendOpenGL

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
	lastDisplayedFPS  int // cached FPS to avoid string allocation every frame
}

var hlg = &engine{}

// SetBackend sets the graphics backend to use. Must be called before Run() or any
// graphics operations. If not called, defaults to BackendWebGPU.
func SetBackend(backend Backend) {
	if hlg.hasSetupCompleted {
		panic("SetBackend must be called before Run() or any graphics operations")
	}
	selectedBackend = backend
}

// GetBackend returns the currently selected graphics backend.
func GetBackend() Backend {
	return selectedBackend
}

func setup() {
	runtime.LockOSThread()
	hlg.inputState = input.NewInputState()
	var err error

	switch selectedBackend {
	case BackendOpenGL:
		hlg.graphicsBackend, err = gl.NewGraphicsBackend(windowWidth, windowHeight)
	case BackendWebGPU:
		fallthrough
	default:
		hlg.graphicsBackend, err = webgpu.NewGraphicsBackend(windowWidth, windowHeight)
	}

	if err != nil {
		panic(err.Error())
	}

	// Register callback to flush batched primitives when a Shape is rendered
	// This preserves draw order between batched primitives and Shapes
	hlg.graphicsBackend.SetOnBeforeAddToQueue(flushBatch)

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
		}

		if frameRendered {
			if renderFn != nil {
				renderFn()
			}
			hlg.graphicsBackend.Render()

			calculateFPS()

			// Reset input state after render so gui widgets can see JustPressed events
			hlg.inputState.ResetJustPressed()
		}
	}
}

func RunGame(game Game) {
	run(game.Update, game.Render)
}

// RunApp is an alias for RunGame for backwards compatibility.
// Deprecated: Use RunGame instead.
func RunApp(game Game) {
	RunGame(game)
}

// Run is the main update function called to refresh the engine state.
func Run(updateFn func(), renderFn func()) {
	run(updateFn, renderFn)
}

func calculateFPS() {
	hlg.fpsCounter.Frame()
	if hlg.graphicsBackend.IsDisposed() {
		return
	}

	if !fpsEnabled {
		return
	}

	// Only update title when FPS value changes to avoid string allocation every frame
	currentFPS := int(hlg.fpsCounter.GetFPS())
	if currentFPS == hlg.lastDisplayedFPS {
		return
	}
	hlg.lastDisplayedFPS = currentFPS

	title := fmt.Sprintf("%s -- %d", hlg.windowTitle, currentFPS)
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
	hlg.graphicsBackend.SetWindowTitle(title)
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

// GetFramebufferSize returns the framebuffer size in pixels.
// On HiDPI/Retina displays, this may be larger than the window size.
func GetFramebufferSize() (int, int) {
	return hlg.graphicsBackend.GetFramebufferSize()
}

// GetPixelScale returns the ratio of framebuffer size to window size.
// This is useful for scaling mouse coordinates to match gl_FragCoord in shaders.
// Returns 1.0 on standard displays, 2.0 on Retina displays, etc.
func GetPixelScale() float32 {
	winW, _ := hlg.graphicsBackend.GetWindowSize()
	fbW, _ := hlg.graphicsBackend.GetFramebufferSize()
	if winW == 0 {
		return 1.0
	}
	return float32(fbW) / float32(winW)
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
	// Store dimensions before setup so the backend is created with correct size
	windowWidth = width
	windowHeight = height
	ensureSetupCompletion()
	hlg.graphicsBackend.SetScreenSize(width, height)
}

// SetWindowSize sets the size of the window.
func SetWindowSize(width, height int) {
	// Store dimensions before setup so the backend is created with correct size
	windowWidth = width
	windowHeight = height
	ensureSetupCompletion()
	hlg.graphicsBackend.SetWindowSize(width, height)
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

// SetVSync enables or disables vertical sync.
// When enabled, the frame rate is limited to the display's refresh rate.
func SetVSync(enabled bool) {
	ensureSetupCompletion()
	hlg.graphicsBackend.SetVSync(enabled)
}

// PushClipRect pushes a clip rectangle onto the stack.
// All subsequent rendering will be clipped to this rectangle.
// Clip rectangles can be nested - each push further restricts the clip region.
func PushClipRect(x, y, width, height int) {
	ensureSetupCompletion()
	pushClipRectToStack(x, y, width, height)
	hlg.graphicsBackend.PushClipRect(x, y, width, height)
}

// PopClipRect pops the top clip rectangle from the stack.
// Rendering will be clipped to the previous rectangle, or unclipped if the stack is empty.
func PopClipRect() {
	ensureSetupCompletion()
	popClipRectFromStack()
	hlg.graphicsBackend.PopClipRect()
}
