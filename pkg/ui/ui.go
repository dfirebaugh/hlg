package ui

import (
	"image"
	"image/color"

	"github.com/dfirebaugh/hlg/pkg/input"
	"github.com/dfirebaugh/hlg/pkg/math/geom"
)

type View interface {
	Update(ctx Context)
	Render(ctx Context)
}

type Context interface {
	Renderer
	Theme() *Theme
	Input
	RemoveChildSurface(child *Surface)
	AddChildSurface(width, height int) *Surface
}

type Renderer interface {
	SetVisibility(visible bool)
	DrawRect(x, y, width, height int, c color.Color)
	FillRect(x, y, width, height int, c color.Color)
	DrawCircle(x, y, radius int, c color.Color)
	FillCircle(x, y, radius int, c color.Color)
	DrawTriangle(points [3]geom.Point, c color.Color)
	FillTriangle(points [3]geom.Point, c color.Color)
	DrawText(x, y int, text string, c color.Color)
	DrawLine(x1, y1, x2, y2 int, c color.Color)
	DrawRoundedRectangle(x, y, width, height int, c color.Color)
	FillRoundedRectangle(x, y, width, height int, c color.Color)
	TextWidth(text string) int
	TextHeight(text string) int
	UpdateImage(image.Image) error
	ToImage() *image.RGBA
	Clear(color.Color)
	Render()
	Move(int, int)
	Resize(int, int)
	Clip(x, y, width, height int)
	GetX() int
	GetY() int
	Height() int
	Width() int
}

type Input interface {
	GetCursorPosition() (int, int)
	GetGlobalCursorPosition() (int, int)
	IsKeyPressed(key input.Key) bool
	IsKeyJustPressed(key input.Key) bool
	IsButtonPressed(button input.MouseButton) bool
	IsButtonJustPressed(button input.MouseButton) bool
}
