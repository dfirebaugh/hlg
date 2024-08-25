package ui

import (
	"image/color"

	"github.com/dfirebaugh/hlg"
)

type App struct {
	surfaces   []*Surface
	clearColor color.RGBA
}

func NewApp(screenWidth, screenHeight int) *App {
	return &App{
		surfaces:   []*Surface{},
		clearColor: color.RGBA{0, 0, 0, 0},
	}
}

func (a *App) CreateSurface(width, height int) *Surface {
	surface := NewSurface(width, height)
	a.surfaces = append(a.surfaces, surface)
	return surface
}

func (a *App) Render() {
	hlg.Clear(a.clearColor)
	for _, surface := range a.surfaces {
		surface.Render()
	}
}

func (a *App) Update() {
	for _, surface := range a.surfaces {
		surface.Update()
	}
}
