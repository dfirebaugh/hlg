package graphics

import (
	"image"
	"image/color"

	"github.com/dfirebaugh/ggez/pkg/math/geom"
)

type GraphicsBackend interface {
	Close()
	PrintPlatformAndVersion()
	WindowManager
	EventManager
	Renderer
	TextureManager
	ShapeRenderer
	ModelRenderer
	InputManager
}

type Renderer interface {
	PrintRendererInfo()
	Clear(c color.Color)
	Render()
	ToggleWireframeMode()
}

type WindowManager interface {
	SetWindowTitle(title string)
	DestroyWindow()
	SetScreenSize(width int, height int)
	SetScaleFactor(f int)
	Renderer
}

type TextureManager interface {
	CreateTextureFromImage(img image.Image) (uintptr, error)
	UpdateTextureFromImage(texture uintptr, img image.Image)
	RenderTexture(
		textureInstance uintptr,
		x int,
		y int,
		w int,
		h int,
		angle float64,
		centerX int,
		centerY int,
		flipType int,
	)
	DestroyTexture(textureInstance uintptr)
}

type EventManager interface {
	PollEvents() bool
}

type ShapeRenderer interface {
	// SetScreenSize(width int, height int)
	DrawLine(x1, y1, x2, y2 int, c color.Color)
	FillTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color)
	DrawTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color)
	FillPolygon(xPoints, yPoints []int, c color.Color)
	DrawPolygon(xPoints, yPoints []int, c color.Color)
	FillRect(x, y, width, height int, c color.Color)
	DrawRect(x, y, width, height int, c color.Color)
	FillCirc(x, y, radius int, c color.Color)
	DrawCirc(xCenter, yCenter, radius int, c color.Color)
	DrawPoint(x, y int, c color.Color)
}

type Model interface {
	Rotate(angle float32, axis geom.Vector3)
	Scale(factor float32)
	SetPosition(v geom.Vector3)
	Translate(v geom.Vector3)
	GetMeshes() []*geom.Mesh
	GetScaleFactor() float32
	GetPosition() geom.Vector3
	GetRotation() geom.Matrix4
}

type ModelRenderer interface {
	RenderModel(m Model, t Texture)
}

type Texture interface {
	Handle() uintptr
}

type Shader interface {
	Handle() uint32
	Delete()
}

type ShaderProgram interface {
	Handle() uint32
	Attach(shaders ...Shader)
	Delete()
	GetUniformLocation(name string) int32
	Link() error
	Use()
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
