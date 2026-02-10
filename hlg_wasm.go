//go:build js && wasm

package hlg

import (
	"fmt"
	"image/color"
	"syscall/js"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/gl"
	"github.com/dfirebaugh/hlg/pkg/input"
)

// Backend represents the graphics backend type
type Backend int

const (
	// BackendWebGL is the only available backend for WASM
	BackendWebGL Backend = iota
	// BackendOpenGL is defined for API compatibility but not available in WASM
	BackendOpenGL
	// BackendWebGPU is defined for API compatibility but not available in WASM
	BackendWebGPU
)

var selectedBackend = BackendWebGL

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
	lastDisplayedFPS  int
}

var hlg = &engine{}

// SetBackend is a no-op for WASM builds (always uses WebGL)
func SetBackend(backend Backend) {
	// WebGL is the only backend available for WASM
	selectedBackend = BackendWebGL
}

// GetBackend returns the currently selected graphics backend
func GetBackend() Backend {
	return selectedBackend
}

func setup() {
	hlg.inputState = input.NewInputState()
	var err error

	hlg.graphicsBackend, err = gl.NewGraphicsBackend(windowWidth, windowHeight)
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
	hlg.fpsCounter = newFPSCounter()

	hlg.graphicsBackend.SetInputCallback(func(eventChan chan input.Event) {
		evt := <-eventChan
		handleEvent(evt, hlg.inputState)
	})

	// Use requestAnimationFrame for the game loop
	var frameFunc js.Func
	frameFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if hlg.graphicsBackend.IsDisposed() {
			frameFunc.Release()
			return nil
		}

		// Update
		if updateFn != nil {
			updateFn()
		}

		// Render
		if renderFn != nil {
			renderFn()
		}
		hlg.graphicsBackend.Render()

		calculateFPS()
		hlg.inputState.ResetJustPressed()

		// Schedule next frame
		js.Global().Call("requestAnimationFrame", frameFunc)
		return nil
	})

	// Start the animation loop
	js.Global().Call("requestAnimationFrame", frameFunc)

	// Keep the Go program alive
	select {}
}

func RunGame(game Game) {
	run(game.Update, game.Render)
}

// RunApp is an alias for RunGame for backwards compatibility
// Deprecated: Use RunGame instead
func RunApp(game Game) {
	RunGame(game)
}

// Run is the main update function called to refresh the engine state
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

// Clear clears the screen with the specified color
func Clear(c color.RGBA) {
	ensureSetupCompletion()
	hlg.graphicsBackend.Clear(c)
}

// SetTitle sets the title of the window (document title for web)
func SetTitle(title string) {
	ensureSetupCompletion()
	hlg.windowTitle = title
}

// GetWindowSize retrieves the current window size
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
// On HiDPI displays, this may be larger than the window size.
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

// SetScreenSize sets the size of the screen
func SetScreenSize(width, height int) {
	// Store dimensions before setup so the backend is created with correct size
	windowWidth = width
	windowHeight = height
	ensureSetupCompletion()
	hlg.graphicsBackend.SetScreenSize(width, height)
}

// SetWindowSize sets the size of the window
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

// SetVSync is a no-op for WASM (browser controls refresh rate)
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
