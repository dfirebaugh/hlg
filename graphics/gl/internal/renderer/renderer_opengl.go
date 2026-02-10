//go:build !js

package renderer

import (
	"image/color"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/glapi"
)

// Renderer manages OpenGL rendering state
type Renderer struct {
	ctx           *glapi.Context
	surface       *Surface
	renderQueues  []*RenderQueue
	clipRectStack [][4]int
}

// NewRenderer creates a new renderer
func NewRenderer(ctx *glapi.Context, surface *Surface) *Renderer {
	r := &Renderer{
		ctx:           ctx,
		surface:       surface,
		renderQueues:  make([]*RenderQueue, 0),
		clipRectStack: make([][4]int, 0),
	}

	// Enable blending
	ctx.Enable(glapi.BLEND)
	ctx.BlendFunc(glapi.SRC_ALPHA, glapi.ONE_MINUS_SRC_ALPHA)

	return r
}

// GetContext returns the GL context
func (r *Renderer) GetContext() *glapi.Context {
	return r.ctx
}

// GetSurface returns the surface
func (r *Renderer) GetSurface() *Surface {
	return r.surface
}

// AddRenderQueue adds a render queue to be rendered
func (r *Renderer) AddRenderQueue(rq *RenderQueue) {
	r.renderQueues = append(r.renderQueues, rq)
}

// CreateRenderQueue creates and registers a new render queue
func (r *Renderer) CreateRenderQueue() *RenderQueue {
	rq := NewRenderQueue(r.ctx, r.surface)
	r.AddRenderQueue(rq)
	return rq
}

// Clear clears the screen with the given color
func (r *Renderer) Clear(c color.Color) {
	cr, cg, cb, ca := c.RGBA()
	r.ctx.ClearColor(
		float32(cr)/0xffff,
		float32(cg)/0xffff,
		float32(cb)/0xffff,
		float32(ca)/0xffff,
	)
	r.ctx.Clear(glapi.COLOR_BUFFER_BIT | glapi.DEPTH_BUFFER_BIT)
}

var renderDebugOnce bool

// Render renders all queued objects
func (r *Renderer) Render(getFramebufferSize func() (int, int)) {
	fbWidth, fbHeight := getFramebufferSize()
	surfW, surfH := r.surface.GetSurfaceSize()

	if !renderDebugOnce {
		println("[HLG Debug] framebuffer:", fbWidth, "x", fbHeight, "surface:", surfW, "x", surfH)
		renderDebugOnce = true
	}

	r.ctx.Viewport(0, 0, fbWidth, fbHeight)

	r.ctx.Enable(glapi.BLEND)
	r.ctx.BlendFunc(glapi.SRC_ALPHA, glapi.ONE_MINUS_SRC_ALPHA)

	for _, rq := range r.renderQueues {
		rq.RenderFrame()
	}
}

// PrepareFrame prepares all render queues for a new frame
func (r *Renderer) PrepareFrame() {
	for _, rq := range r.renderQueues {
		rq.PrepareFrame()
	}
}

// Resize updates the surface size
func (r *Renderer) Resize(width, height int) {
	r.surface.Resize(width, height)
}

// Dispose cleans up renderer resources
func (r *Renderer) Dispose() {
	for _, rq := range r.renderQueues {
		rq.Dispose()
	}
	r.renderQueues = nil
}

// PushClipRect pushes a clip rectangle onto the stack
func (r *Renderer) PushClipRect(x, y, width, height int) {
	r.clipRectStack = append(r.clipRectStack, [4]int{x, y, width, height})
}

// PopClipRect pops a clip rectangle from the stack
func (r *Renderer) PopClipRect() {
	if len(r.clipRectStack) > 0 {
		r.clipRectStack = r.clipRectStack[:len(r.clipRectStack)-1]
	}
}

// GetCurrentClipRect returns the current clip rectangle or nil if none
func (r *Renderer) GetCurrentClipRect() *[4]int {
	if len(r.clipRectStack) == 0 {
		return nil
	}
	rect := r.clipRectStack[len(r.clipRectStack)-1]
	return &rect
}

// CreateMSDFAtlas creates a new MSDF atlas
func (r *Renderer) CreateMSDFAtlas() graphics.MSDFAtlas {
	// Will be implemented via RenderQueue
	return nil
}
