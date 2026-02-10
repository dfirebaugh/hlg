// Package graphics provides functionality for 2D graphics rendering,
// including textures, sprites, text, and shapes.
package graphics

import (
	"image"
	"image/color"

	"github.com/dfirebaugh/hlg/pkg/input"
	"github.com/dfirebaugh/hlg/pkg/math/matrix"
)

type GraphicsBackend interface {
	Close()
	GetScreenSize() (int, int)
	SetScreenSize(width int, height int)

	WindowManager
	EventManager
	Renderer
	TextureManager
	ShapeRenderer
	InputManager
	ShaderManager
	FontManager
}

type (
	ShaderHandle  int
	ShaderManager interface {
		CompileShader(shaderCode string) ShaderHandle
	}
)

type Texture interface {
	Handle() uintptr
	UpdateImage(img image.Image) error
	SetShouldBeRendered(souldRender bool)
	Resize(width, height float32)
	Move(x, y float32)
	Rotate(a, pivotX, pivotY float32)
	Scale(x, y float32)
	FlipVertical()
	FlipHorizontal()
	SetFlipHorizontal(shouldFlip bool)
	SetFlipVertical(shouldFlip bool)
	Clip(minX, minY, maxX, maxY float32)
	Render()
	RenderToQueue(rq RenderQueue)
	Dispose()
	IsDisposed() bool
}

type ShaderRenderable interface {
	UpdateUniforms(dataMap map[string][]byte)
	UpdateUniform(name string, data []byte)
	Render()
	Dispose()
	IsDisposed() bool
}

type Renderable interface {
	Render()
	Dispose()
	IsDisposed() bool
}

// ClipRectProvider provides access to the current clip rect
type ClipRectProvider interface {
	GetCurrentClipRect() *[4]int
}

type RenderQueue interface {
	AddToRenderQueue(r Renderable)
	SetPriority(priority int)
	GetCurrentClipRect() *[4]int
	SetOnBeforeAddToQueue(callback func()) // Callback to flush batched primitives before adding renderables
	Present()                              // Present this queue's contents immediately
}

type Renderer interface {
	RenderQueue
	Clear(c color.Color)
	Render()
	CreateRenderQueue() RenderQueue
	SetVSync(enabled bool)
	PushClipRect(x, y, width, height int)
	PopClipRect()
	GetCurrentClipRect() *[4]int
}

type Transformable interface {
	Move(screenX, screenY float32)
	Rotate(angle float32) matrix.Matrix
	Scale(sx, sy float32) matrix.Matrix
}

type WindowManager interface {
	DisableWindowResize()
	SetBorderlessWindowed(v bool)
	SetWindowTitle(title string)
	DestroyWindow()
	SetWindowSize(width int, height int)
	GetWindowSize() (int, int)
	GetFramebufferSize() (int, int)
	GetWindowPosition() (x int, y int)
	IsDisposed() bool
	Renderer
}

type TextureManager interface {
	CreateTextureFromImage(img image.Image) (Texture, error)
	DisposeTexture(h uintptr)
}

type EventManager interface {
	PollEvents() bool
}

type Shape interface {
	Renderable
	Transformable
	SetColor(c color.Color)
}

type Vertex struct {
	Position [3]float32 // x, y, z coordinates
	Color    [4]float32 // RGBA color
}

// OpCode constants for PrimitiveVertex rendering modes
const (
	OpCodeCircle      float32 = 0.0
	OpCodeRoundedRect float32 = 1.0
	OpCodeTriangleSDF float32 = 2.0
	OpCodeMSDF        float32 = 3.0 // MSDF text rendering
	OpCodeSolid       float32 = 4.0 // Simple solid fill (no SDF)
	OpCodeLine        float32 = 5.0 // Line segment SDF
)

// PrimitiveVertex is the vertex format used by the primitive buffer for SDF rendering
// DEPRECATED: Use Primitive instead for the new storage buffer approach
// Note: This struct is uploaded to GPU, so it cannot contain Go pointers.
// ClipRect is tracked separately in the primitive buffer.
type PrimitiveVertex struct {
	Position      [3]float32 // Clip space position
	LocalPosition [2]float32 // Local coordinates for SDF calculation
	OpCode        float32    // Rendering operation code
	Radius        float32    // Corner radius or circle radius
	Color         [4]float32 // RGBA color
	TexCoords     [2]float32 // UV coordinates for MSDF text, or line direction
	HalfSize      [2]float32 // Half width/height of bounding box (for OpenGL SDF)
}

// Primitive is a compact representation for the storage buffer approach.
// One Primitive per shape (vs 6 PrimitiveVertex per shape).
// The vertex shader constructs vertices from this data.
// Memory: 64 bytes per primitive vs 312 bytes (6 x 52) with PrimitiveVertex = ~5x reduction
// IMPORTANT: This struct must match WGSL alignment. vec4<f32> requires 16-byte alignment.
type Primitive struct {
	X, Y, W, H float32    // bytes 0-15: bounding box in screen space
	Color      [4]float32 // bytes 16-31: RGBA color (vec4, 16-byte aligned)
	Radius     float32    // bytes 32-35: corner radius or circle radius
	OpCode     float32    // bytes 36-39: primitive type
	_          [2]float32 // bytes 40-47: padding to align Extra to 16 bytes
	Extra      [4]float32 // bytes 48-63: for MSDF: (u0, v0, u_size, v_size); for shapes: (half_w, half_h, 0, 0)
	ClipRect   *[4]int    // optional clip rect (x, y, width, height) - nil means no clipping
}

type ShapeRenderer interface {
	AddPolygonFromVertices(cx, cy int, width float32, vertices []Vertex) Shape
	AddPolygon(cx, cy int, width float32, c color.Color, sides int) Shape
	AddTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) Shape
	AddRectangle(x, y, width, height int, c color.Color) Shape
	AddRoundedRectangle(x, y, width, height, radius int, c color.Color) Shape
	AddCircle(cx, cy int, radius float32, c color.Color, segments int) Shape
	AddLine(x1, y1, x2, y2 int, width float32, c color.Color) Shape
	AddDynamicRenderable(vertexData []byte, layout VertexBufferLayout, shaderHandle int, uniforms map[string]Uniform, dataMap map[string][]byte) ShaderRenderable
	DrawPrimitiveBuffer(vertices []PrimitiveVertex)                                   // DEPRECATED: Use DrawPrimitives instead
	DrawPrimitiveBufferWithClipRects(vertices []PrimitiveVertex, clipRects []*[4]int) // Draw with per-vertex clip rects
	DrawPrimitives(primitives []Primitive)                                            // New storage buffer approach
	FlushPrimitiveBuffer()                                                            // Force immediate render of pending primitives
	SetMSDFAtlas(atlasImg image.Image, pxRange float64)
	SetMSDFMode(mode int)
	EnableSnapMSDFToPixels(enable bool)
}

type CubeFace int

const (
	CubeFront CubeFace = iota
	CubeBack
	CubeLeft
	CubeRight
	CubeTop
	CubeBottom
)

type InputManager interface {
	SetInputCallback(fn func(eventChan chan input.Event))
}

type Uniform struct {
	Binding uint32
	Size    uint64
}

type VertexAttributeLayout struct {
	ShaderLocation uint32
	Offset         uint64
	Format         string
}

type VertexBufferLayout struct {
	ArrayStride uint64
	Attributes  []VertexAttributeLayout
}

// GlyphQuad represents the UV and position data for rendering a glyph
type GlyphQuad struct {
	// Texture coordinates in atlas (normalized 0-1)
	S0, T0 float64 // Top-left
	S1, T1 float64 // Bottom-right

	// Position coordinates in em units (relative to baseline)
	PL, PB float64 // Bottom-left (plane left, plane bottom)
	PR, PT float64 // Top-right (plane right, plane top)

	// Advance width in em units
	Advance float64
}

// GlyphInfo contains all information needed to render a glyph
type GlyphInfo struct {
	Unicode int
	Quad    GlyphQuad
}

// FontMetrics contains font-level metrics
type FontMetrics struct {
	EmSize     float64
	LineHeight float64
	Ascender   float64
	Descender  float64
}

// MSDFAtlas represents an MSDF font atlas interface
type MSDFAtlas interface {
	AddGlyph(r rune, info *GlyphInfo)
	GetGlyph(r rune) *GlyphInfo
	SetMetrics(metrics FontMetrics)
	GetMetrics() FontMetrics
	Dispose()
	IsDisposed() bool
}

// FontManager provides methods for creating MSDF font atlases
type FontManager interface {
	CreateMSDFAtlas(atlasImg image.Image, distanceRange float64) (MSDFAtlas, error)
}

// GlyphDrawer is an interface for drawing MSDF glyphs.
// This is implemented by gui.DrawContext to allow font rendering without
// the graphics package depending on the gui package.
type GlyphDrawer interface {
	DrawGlyph(x0, y0, x1, y1 float32, u0, v0, u1, v1 float32, c color.Color)
}
