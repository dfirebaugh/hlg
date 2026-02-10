//go:build js && wasm

// Package gl provides a unified OpenGL/WebGL graphics backend
package gl

import (
	"image"
	"image/color"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/canvas"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/glapi"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/pipelines"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/renderer"
)

// GraphicsBackend provides the WebGL graphics backend
type GraphicsBackend struct {
	*canvas.Canvas
	*renderer.Renderer
	*renderer.RenderQueue
	*renderer.Surface

	ctx *glapi.Context
}

// NewGraphicsBackend creates a new WebGL graphics backend with default canvas ID
func NewGraphicsBackend(width, height int) (*GraphicsBackend, error) {
	return NewGraphicsBackendWithCanvas("hlg-canvas", width, height)
}

// NewGraphicsBackendWithCanvas creates a new WebGL graphics backend with specified canvas
func NewGraphicsBackendWithCanvas(canvasID string, width, height int) (*GraphicsBackend, error) {
	c, err := canvas.NewCanvas(canvasID, width, height)
	if err != nil {
		return nil, err
	}

	// Set up the canvas for aspect ratio preservation from the start
	c.SetWindowSize(width, height)

	ctx := c.GetGL()
	surface := renderer.NewSurface(width, height)
	r := renderer.NewRenderer(ctx, surface)
	rq := r.CreateRenderQueue()

	g := &GraphicsBackend{
		Canvas:      c,
		Renderer:    r,
		RenderQueue: rq,
		Surface:     surface,
		ctx:         ctx,
	}

	return g, nil
}

// Close cleans up resources
func (g *GraphicsBackend) Close() {
	if g.Renderer != nil {
		g.Renderer.Dispose()
	}
	if g.Canvas != nil {
		g.Canvas.Destroy()
	}
}

// PollEvents polls for events (always returns true for web)
func (g *GraphicsBackend) PollEvents() bool {
	return g.Canvas.Poll()
}

// Resize resizes the graphics backend
func (g *GraphicsBackend) Resize(width, height int) {
	g.Surface.Resize(width, height)
}

// GetScreenSize returns the screen size
func (g *GraphicsBackend) GetScreenSize() (int, int) {
	return g.Surface.GetSurfaceSize()
}

// SetScreenSize sets the screen size and locks it so window resizes don't change it
func (g *GraphicsBackend) SetScreenSize(width, height int) {
	g.Surface.SetSurfaceSize(width, height)
	// Also tell the canvas about the target size for aspect ratio preservation
	g.Canvas.SetWindowSize(width, height)
}

// SwapBuffers is a no-op for WebGL (browser handles this)
func (g *GraphicsBackend) SwapBuffers() {
	g.Canvas.SwapBuffers()
}

// Clear clears the screen with a color
func (g *GraphicsBackend) Clear(c color.Color) {
	g.Renderer.Clear(c)
}

// Render renders all queued objects
func (g *GraphicsBackend) Render() {
	g.Renderer.Render(g.Canvas.GetFramebufferSize)
}

// PrepareFrame prepares for a new frame
func (g *GraphicsBackend) PrepareFrame() {
	g.Renderer.PrepareFrame()
}

// CreateTexture creates a new texture from an image
func (g *GraphicsBackend) CreateTexture(img image.Image) graphics.Texture {
	return g.RenderQueue.CreateTexture(img)
}

// CreateMSDFAtlas creates a new MSDF font atlas from an image
func (g *GraphicsBackend) CreateMSDFAtlas(atlasImg image.Image, distanceRange float64) (graphics.MSDFAtlas, error) {
	return pipelines.NewMSDFAtlas(atlasImg, distanceRange)
}

// CompileShader compiles a shader
func (g *GraphicsBackend) CompileShader(code string) graphics.ShaderHandle {
	return g.RenderQueue.CompileShader(code)
}

// SetVSync is a no-op (browser controls refresh rate)
func (g *GraphicsBackend) SetVSync(enabled bool) {
	g.Canvas.SetVSync(enabled)
}

// GetCurrentClipRect returns the current clip rectangle
func (g *GraphicsBackend) GetCurrentClipRect() *[4]int {
	return g.RenderQueue.GetCurrentClipRect()
}

// PushClipRect pushes a clip rectangle
func (g *GraphicsBackend) PushClipRect(x, y, width, height int) {
	g.RenderQueue.PushClipRect(x, y, width, height)
}

// PopClipRect pops a clip rectangle
func (g *GraphicsBackend) PopClipRect() {
	g.RenderQueue.PopClipRect()
}

// GetPrimitiveBuffer returns the primitive buffer
func (g *GraphicsBackend) GetPrimitiveBuffer() *pipelines.PrimitiveBuffer {
	return g.RenderQueue.GetPrimitiveBuffer()
}

// SetOnBeforeAddToQueue sets a callback before adding to queue
func (g *GraphicsBackend) SetOnBeforeAddToQueue(fn func()) {
	g.RenderQueue.SetOnBeforeAddToQueue(fn)
}

// CreateRenderQueue creates a new render queue
func (g *GraphicsBackend) CreateRenderQueue() graphics.RenderQueue {
	return g.Renderer.CreateRenderQueue()
}

// Ensure GraphicsBackend implements graphics.GraphicsBackend
var _ graphics.GraphicsBackend = (*GraphicsBackend)(nil)
