// Package graphics provides functionality for 2D graphics rendering,
// including textures, sprites, text, and shapes.
package graphics

import (
	"image"
	"image/color"

	"github.com/dfirebaugh/hlg/pkg/input"
	"github.com/dfirebaugh/hlg/pkg/math/matrix"
	"github.com/rajveermalviya/go-webgpu/wgpu"
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
	RenderPass(pass *wgpu.RenderPassEncoder)
	Render()
	RenderToQueue(rq RenderQueue)
	Dispose()
	IsDisposed() bool
}

type ShaderRenderable interface {
	UpdateUniforms(dataMap map[string][]byte)
	UpdateUniform(name string, data []byte)
	RenderPass(pass *wgpu.RenderPassEncoder)
	Render()
	Dispose()
	IsDisposed() bool
	Hide()
}

type Renderable interface {
	RenderPass(pass *wgpu.RenderPassEncoder)
	Render()
	Dispose()
	IsDisposed() bool
	Hide()
}
type RenderQueue interface {
	AddToRenderQueue(r Renderable)
	SetPriority(priority int)
}

type Renderer interface {
	RenderQueue
	Clear(c color.Color)
	Render()
	CreateRenderQueue() RenderQueue
}

type Transformable interface {
	Move(screenX, screenY float32)
	Rotate(angle float32) matrix.Matrix
	Scale(sx, sy float32) matrix.Matrix
}

type WindowManager interface {
	DisableWindowResize()
	SetWindowTitle(title string)
	DestroyWindow()
	SetWindowSize(width int, height int)
	GetWindowSize() (int, int)
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

type ShapeRenderer interface {
	AddPolygonFromVertices(cx, cy int, width float32, vertices []Vertex) Shape
	AddPolygon(cx, cy int, width float32, c color.Color, sides int) Shape
	AddTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) Shape
	AddRectangle(x, y, width, height int, c color.Color) Shape
	AddCircle(cx, cy int, radius float32, c color.Color, segments int) Shape
	AddLine(x1, y1, x2, y2 int, width float32, c color.Color) Shape
	AddDynamicRenderable(vertices []Vertex, shaderHandle int, uniforms map[string]Uniform, dataMap map[string][]byte) ShaderRenderable
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
