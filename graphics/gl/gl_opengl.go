//go:build !js

// Package gl provides a unified OpenGL/WebGL graphics backend
package gl

import (
	"image"
	"image/color"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/glapi"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/pipelines"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/renderer"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/window"
	"github.com/go-gl/gl/v4.1-core/gl"
)

// GraphicsBackend provides the OpenGL graphics backend
type GraphicsBackend struct {
	*window.Window
	*renderer.Renderer
	*renderer.RenderQueue
	*renderer.Surface

	ctx *glapi.Context
}

// NewGraphicsBackend creates a new OpenGL graphics backend
func NewGraphicsBackend(width, height int) (*GraphicsBackend, error) {
	win, err := window.NewWindow(width, height)
	if err != nil {
		return nil, err
	}

	// Initialize OpenGL
	if err := gl.Init(); err != nil {
		win.Destroy()
		return nil, err
	}

	ctx := glapi.NewContext()
	surface := renderer.NewSurface(width, height)
	r := renderer.NewRenderer(ctx, surface)
	rq := r.CreateRenderQueue()

	g := &GraphicsBackend{
		Window:      win,
		Renderer:    r,
		RenderQueue: rq,
		Surface:     surface,
		ctx:         ctx,
	}

	// Set up resize callback to update surface when window is resized
	win.SetResizedCallback(func(physicalWidth, physicalHeight uint32) {
		g.Surface.Resize(int(physicalWidth), int(physicalHeight))
	})

	return g, nil
}

// Close cleans up resources
func (g *GraphicsBackend) Close() {
	if g.Renderer != nil {
		g.Renderer.Dispose()
	}
	if g.Window != nil {
		g.Window.Destroy()
	}
}

// PollEvents polls for window events
func (g *GraphicsBackend) PollEvents() bool {
	return g.Window.Poll()
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
}

// SwapBuffers swaps the front and back buffers
func (g *GraphicsBackend) SwapBuffers() {
	g.Window.SwapBuffers()
}

// Clear clears the screen with a color
func (g *GraphicsBackend) Clear(c color.Color) {
	g.Renderer.Clear(c)
}

// Render renders all queued objects
func (g *GraphicsBackend) Render() {
	g.Renderer.Render(g.Window.GetFramebufferSize)
	g.Window.SwapBuffers()
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

// SetVSync enables or disables vsync
func (g *GraphicsBackend) SetVSync(enabled bool) {
	g.Window.SetVSync(enabled)
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
