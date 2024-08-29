package grugui

import (
	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/pkg/input"
)

type inputManager struct {
	surface *Surface
}

func newInputManager(s *Surface) *inputManager {
	return &inputManager{
		surface: s,
	}
}

func (i *inputManager) GetCursorPosition() (int, int) {
	globalX, globalY := hlg.GetCursorPosition()
	surfaceX := i.surface.Renderer.GetX()
	surfaceY := i.surface.Renderer.GetY()

	relX := globalX - surfaceX
	relY := globalY - surfaceY

	if i.surface.isClipped {
		if relX < i.surface.clippedRect.Min.X || relX > i.surface.clippedRect.Max.X ||
			relY < i.surface.clippedRect.Min.Y || relY > i.surface.clippedRect.Max.Y {
			return -1, -1
		}
	}

	return relX, relY
}

func (i *inputManager) GetGlobalCursorPosition() (int, int) {
	return hlg.GetCursorPosition()
}

func (i *inputManager) IsKeyPressed(key input.Key) bool {
	return hlg.IsKeyPressed(key)
}

func (i *inputManager) IsKeyJustPressed(key input.Key) bool {
	return hlg.IsKeyJustPressed(key)
}

func (i *inputManager) IsButtonPressed(button input.MouseButton) bool {
	return hlg.IsButtonPressed(button)
}

func (i *inputManager) IsButtonJustPressed(button input.MouseButton) bool {
	return hlg.IsButtonJustPressed(button)
}
