package grugui

import (
	"image"

	"github.com/dfirebaugh/hlg/pkg/grugui/renderer"
)

type Surface struct {
	Renderer
	inputManager
	theme         *Theme
	views         []View
	childSurfaces []*Surface
	clippedRect   image.Rectangle
	isClipped     bool
}

func NewSurface(width, height int) *Surface {
	s := &Surface{
		Renderer:      renderer.NewGrugRenderer(width, height),
		theme:         DefaultTheme(),
		views:         []View{},
		childSurfaces: []*Surface{},
		clippedRect:   image.Rect(0, 0, width, height),
	}

	s.inputManager = *newInputManager(s)

	return s
}

func (s *Surface) Update() {
	if s.isSurfaceInFocus() {
		for _, view := range s.views {
			view.Update(s)
		}
	}

	for _, child := range s.childSurfaces {
		child.Update()
	}
}

func (s *Surface) Render() {
	s.Renderer.Clear(s.theme.BackgroundColor)
	for _, view := range s.views {
		view.Render(s)
	}

	for _, child := range s.childSurfaces {
		child.Render()
	}
	s.Renderer.Render()
}

func (s *Surface) Add(v View) {
	s.views = append(s.views, v)
}

func (s *Surface) Theme() *Theme {
	return s.theme
}

func (s *Surface) isSurfaceInFocus() bool {
	cursorX, cursorY := s.inputManager.GetGlobalCursorPosition()
	surfaceX, surfaceY := s.Renderer.GetX(), s.Renderer.GetY()
	width, height := s.Renderer.Width(), s.Renderer.Height()

	if cursorX < surfaceX || cursorX > surfaceX+width || cursorY < surfaceY || cursorY > surfaceY+height {
		return false
	}

	relX := cursorX - surfaceX
	relY := cursorY - surfaceY
	if s.isClipped {
		if relX < s.clippedRect.Min.X || relX > s.clippedRect.Max.X ||
			relY < s.clippedRect.Min.Y || relY > s.clippedRect.Max.Y {
			return false
		}
	}

	return true
}

func (s *Surface) AddChildSurface(width, height int) *Surface {
	child := NewSurface(width, height)
	s.childSurfaces = append(s.childSurfaces, child)
	return child
}

func (s *Surface) RemoveChildSurface(child *Surface) {
	for i, surface := range s.childSurfaces {
		if surface == child {
			s.childSurfaces = append(s.childSurfaces[:i], s.childSurfaces[i+1:]...)
			break
		}
	}
}

func (s *Surface) Clip(x, y, width, height int) {
	clipRect := image.Rect(x, y, x+width, y+height)
	s.ClipToRect(clipRect)
	s.Renderer.Clip(x, y, width, height)
}

func (s *Surface) ClipToRect(rect image.Rectangle) {
	s.clippedRect = rect
	s.isClipped = true
}
