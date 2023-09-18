package renderer

import (
	"image"
	"image/color"

	"github.com/dfirebaugh/ggez/pkg/input"
)

type GraphicsBackend interface {
	Close()
	PrintPlatformAndVersion()
	EventManager
	Renderer
	WindowManager
	TextureManager
	ShapeRenderer
}

type Renderer interface {
	PrintRendererInfo()
	Clear(c color.Color)
	RenderPresent()
}

type WindowManager interface {
	CreateWindow(title string, width int, height int) (uintptr, error)
	SetWindowTitle(title string)
	DestroyWindow()
	SetScreenSize(width int, height int)
	Renderer
}

type TextureManager interface {
	CreateTextureFromImage(img image.Image) (uintptr, error)
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
	PollEvents(DefaultInput input.InputDevice) bool
}

type ShapeRenderer interface {
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
