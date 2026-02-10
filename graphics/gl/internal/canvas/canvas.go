//go:build js && wasm

// Package canvas provides a canvas element wrapper for WebGL rendering
package canvas

import (
	"strconv"
	"syscall/js"

	"github.com/dfirebaugh/hlg/graphics/gl/internal/glapi"
	"github.com/dfirebaugh/hlg/pkg/input"
)

// Canvas wraps an HTML canvas element and its WebGL context
type Canvas struct {
	element   js.Value
	glContext *glapi.Context
	eventChan chan input.Event

	width, height       int     // drawing buffer size (physical pixels)
	cssWidth, cssHeight int     // CSS display size (logical pixels)
	dpr                 float64 // device pixel ratio
	isDisposed          bool

	// Aspect ratio preservation
	targetWidth, targetHeight int     // logical size set by SetWindowSize
	aspectRatio               float64 // targetWidth / targetHeight
	hasTargetSize             bool    // true if SetWindowSize was called

	resizedCallback func(physicalWidth, physicalHeight uint32)
	inputCallback   func(eventChan chan input.Event)
	closeCallback   func()
}

// NewCanvas creates a new Canvas wrapper for the given canvas element ID
func NewCanvas(canvasID string, width, height int) (*Canvas, error) {
	glCtx, err := glapi.NewContextFromCanvas(canvasID)
	if err != nil {
		return nil, err
	}

	// Get device pixel ratio for crisp rendering on high-DPI displays
	dpr := js.Global().Get("devicePixelRatio").Float()
	if dpr < 1 {
		dpr = 1
	}

	element := glCtx.GetCanvas()

	// Get actual CSS display size from the element
	// This allows external CSS to control the canvas size
	cssWidth := element.Get("clientWidth").Int()
	cssHeight := element.Get("clientHeight").Int()

	// Fall back to provided dimensions if element has no size yet
	if cssWidth <= 0 {
		cssWidth = width
	}
	if cssHeight <= 0 {
		cssHeight = height
	}

	canvas := &Canvas{
		element:   element,
		glContext: glCtx,
		eventChan: make(chan input.Event, 100),
		cssWidth:  cssWidth,
		cssHeight: cssHeight,
		dpr:       dpr,
	}

	// Set drawing buffer size (physical pixels)
	canvas.width = int(float64(cssWidth) * dpr)
	canvas.height = int(float64(cssHeight) * dpr)
	canvas.element.Set("width", canvas.width)
	canvas.element.Set("height", canvas.height)

	canvas.setupEventListeners()
	canvas.setupResizeListener()

	return canvas, nil
}

// GetGL returns the WebGL context
func (c *Canvas) GetGL() *glapi.Context {
	return c.glContext
}

// GetSize returns the framebuffer size (same as canvas size for WebGL)
func (c *Canvas) GetSize() (int, int) {
	return c.width, c.height
}

// GetFramebufferSize returns the framebuffer size
func (c *Canvas) GetFramebufferSize() (int, int) {
	return c.width, c.height
}

// GetWindowSize returns the logical window size (the target size set by SetWindowSize)
func (c *Canvas) GetWindowSize() (int, int) {
	if c.hasTargetSize {
		return c.targetWidth, c.targetHeight
	}
	return c.cssWidth, c.cssHeight
}

// SetWindowSize sets the target logical size and aspect ratio
// The canvas will scale to fit the browser window while preserving this aspect ratio
func (c *Canvas) SetWindowSize(width, height int) {
	c.targetWidth = width
	c.targetHeight = height
	c.aspectRatio = float64(width) / float64(height)
	c.hasTargetSize = true

	// Trigger a resize to fit the current window
	c.fitToWindow()

	// If fitToWindow didn't set the size (e.g., window dimensions not available),
	// fall back to target dimensions
	if c.width <= 0 || c.height <= 0 {
		c.cssWidth = width
		c.cssHeight = height
		c.width = int(float64(width) * c.dpr)
		c.height = int(float64(height) * c.dpr)
		c.element.Set("width", c.width)
		c.element.Set("height", c.height)
	}
}

// fitToWindow calculates the largest canvas size that fits the window while preserving aspect ratio
func (c *Canvas) fitToWindow() {
	if !c.hasTargetSize {
		return
	}

	// Get the available space (window inner size)
	windowWidth := js.Global().Get("window").Get("innerWidth").Int()
	windowHeight := js.Global().Get("window").Get("innerHeight").Int()

	if windowWidth <= 0 || windowHeight <= 0 {
		return
	}

	// Calculate the largest size that fits while preserving aspect ratio
	windowAspect := float64(windowWidth) / float64(windowHeight)

	var newCSSWidth, newCSSHeight int
	if windowAspect > c.aspectRatio {
		// Window is wider than target - fit to height
		newCSSHeight = windowHeight
		newCSSWidth = int(float64(windowHeight) * c.aspectRatio)
	} else {
		// Window is taller than target - fit to width
		newCSSWidth = windowWidth
		newCSSHeight = int(float64(windowWidth) / c.aspectRatio)
	}

	c.cssWidth = newCSSWidth
	c.cssHeight = newCSSHeight

	// Scale to physical pixels for crisp rendering
	c.width = int(float64(newCSSWidth) * c.dpr)
	c.height = int(float64(newCSSHeight) * c.dpr)

	// Set drawing buffer size
	c.element.Set("width", c.width)
	c.element.Set("height", c.height)

	// Set CSS display size and center the canvas
	style := c.element.Get("style")
	style.Set("width", strconv.Itoa(newCSSWidth)+"px")
	style.Set("height", strconv.Itoa(newCSSHeight)+"px")
	style.Set("position", "absolute")
	style.Set("left", "50%")
	style.Set("top", "50%")
	style.Set("transform", "translate(-50%, -50%)")

	// Notify callback of new physical size
	if c.resizedCallback != nil {
		c.resizedCallback(uint32(c.width), uint32(c.height))
	}
}

// GetWindowPosition returns the canvas position (always 0, 0 for canvas)
func (c *Canvas) GetWindowPosition() (int, int) {
	return 0, 0
}

// cssToLogical converts CSS canvas coordinates to logical coordinates
// This returns coordinates in the target logical space (set by SetWindowSize)
func (c *Canvas) cssToLogical(cssX, cssY int) (int, int) {
	if !c.hasTargetSize || c.cssWidth <= 0 || c.cssHeight <= 0 {
		return cssX, cssY
	}

	// Scale from CSS size to logical size
	scaleX := float64(c.targetWidth) / float64(c.cssWidth)
	scaleY := float64(c.targetHeight) / float64(c.cssHeight)

	return int(float64(cssX) * scaleX), int(float64(cssY) * scaleY)
}

// cssToLogicalWithRect converts CSS canvas coordinates to logical coordinates
// using the actual rendered size from getBoundingClientRect.
// This ensures correct coordinate translation even immediately after window resize.
func (c *Canvas) cssToLogicalWithRect(relX, relY, rectWidth, rectHeight float64) (int, int) {
	if !c.hasTargetSize || rectWidth <= 0 || rectHeight <= 0 {
		return int(relX), int(relY)
	}

	// Scale from actual rendered size to logical size
	scaleX := float64(c.targetWidth) / rectWidth
	scaleY := float64(c.targetHeight) / rectHeight

	return int(relX * scaleX), int(relY * scaleY)
}

// cssToFramebuffer converts CSS canvas coordinates to framebuffer coordinates
// This returns coordinates that match gl_FragCoord in shaders
func (c *Canvas) cssToFramebuffer(cssX, cssY int) (int, int) {
	// Simply scale by device pixel ratio to get framebuffer coordinates
	return int(float64(cssX) * c.dpr), int(float64(cssY) * c.dpr)
}

// Poll checks for pending events (always returns true for web)
func (c *Canvas) Poll() bool {
	if c.isDisposed {
		return false
	}
	return true
}

// SwapBuffers is a no-op for WebGL (browser handles this)
func (c *Canvas) SwapBuffers() {
	// Browser automatically presents the framebuffer
}

// SetVSync is a no-op (browser controls refresh rate)
func (c *Canvas) SetVSync(enabled bool) {
	// Browser controls this via requestAnimationFrame
}

// SetWindowTitle sets the document title
func (c *Canvas) SetWindowTitle(title string) {
	js.Global().Get("document").Set("title", title)
}

// SetBorderlessWindowed is a no-op for web
func (c *Canvas) SetBorderlessWindowed(v bool) {
	// Not applicable for web
}

// DisableWindowResize is a no-op for web (CSS controls this)
func (c *Canvas) DisableWindowResize() {
	// CSS would control this
}

// SetAspectRatio is handled via CSS for web
func (c *Canvas) SetAspectRatio(numerator, denominator int) {
	// CSS would control this
}

// SetCloseCallback sets a callback for close requests
func (c *Canvas) SetCloseCallback(fn func()) {
	c.closeCallback = fn
	// Listen for beforeunload event
	js.Global().Get("window").Call("addEventListener", "beforeunload", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if c.closeCallback != nil {
			c.closeCallback()
		}
		return nil
	}))
}

// SetCloseRequestedCallback sets a callback for close requests
func (c *Canvas) SetCloseRequestedCallback(fn func()) {
	c.SetCloseCallback(fn)
}

// SetInputCallback sets the input event callback
func (c *Canvas) SetInputCallback(fn func(eventChan chan input.Event)) {
	c.inputCallback = fn
}

// SetResizedCallback sets the resize callback
func (c *Canvas) SetResizedCallback(fn func(physicalWidth, physicalHeight uint32)) {
	c.resizedCallback = fn
}

// Destroy cleans up the canvas resources
func (c *Canvas) Destroy() {
	if c.isDisposed {
		return
	}
	c.isDisposed = true
}

// DestroyWindow is an alias for Destroy
func (c *Canvas) DestroyWindow() {
	c.Destroy()
}

// IsDisposed returns whether the canvas has been disposed
func (c *Canvas) IsDisposed() bool {
	return c.isDisposed
}
