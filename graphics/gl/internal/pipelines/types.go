package pipelines

import "github.com/dfirebaugh/hlg/graphics"

// Surface provides screen size information
type Surface interface {
	GetSurfaceSize() (int, int)
}

// RenderQueue interface for adding renderables
type RenderQueue interface {
	AddToRenderQueue(r graphics.Renderable)
	GetCurrentClipRect() *[4]int
	GetSurfaceSize() (int, int)
}

// PrimitiveBufferProvider provides access to the primitive buffer for rendering
type PrimitiveBufferProvider interface {
	GetPrimitiveBuffer() *PrimitiveBuffer
}

// ClipRectRun tracks a contiguous range of vertices with the same clip rect
type ClipRectRun struct {
	ClipRect *[4]int
	StartIdx int
	Count    int
}
