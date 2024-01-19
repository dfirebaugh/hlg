// Package graphics provides functionality for 2D graphics rendering,
// including textures, sprites, text, and shapes.
package graphics

import (
	"image"
	"image/color"

	"github.com/dfirebaugh/hlg/pkg/input"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type GraphicsBackend interface {
	Close()
	PrintPlatformAndVersion()
	ScreenSize() (int, int)

	WindowManager
	EventManager
	Renderer
	TextureManager
	ShapeRenderer
	InputManager
}

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
	Clip(minX, minY, maxX, maxY float32)
	RenderPass(pass *wgpu.RenderPassEncoder)
	Render()
	Dispose()
}

type Renderable interface {
	RenderPass(pass *wgpu.RenderPassEncoder)
	Render()
	Dispose()
	Hide()
}

type Renderer interface {
	Clear(c color.Color)
	Render()
}

type Transformable interface {
	Move(screenX, screenY float32)
	Rotate(angle float32)
	Scale(sx, sy float32)
}

type WindowManager interface {
	SetWindowTitle(title string)
	DestroyWindow()
	SetWindowSize(width int, height int)
	SetScreenSize(width int, height int)
	GetWindowSize() (int, int)
	SetScaleFactor(f int)
	IsDisposed() bool
	Renderer
}

type TextureManager interface {
	CreateTextureFromImage(img image.Image) (Texture, error)
	UpdateTextureFromImage(texture Texture, img image.Image)
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

type ShapeRenderer interface {
	AddTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) Shape
	AddRectangle(x, y, width, height int, c color.Color) Shape
	AddCircle(cx, cy int, radius float32, c color.Color, segments int) Shape
	AddLine(x1, y1, x2, y2 int, width float32, c color.Color) Shape
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
