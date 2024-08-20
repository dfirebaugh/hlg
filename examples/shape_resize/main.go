package main

import (
	"image/color"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/pkg/input"
	"golang.org/x/image/colornames"
)

type Position struct {
	X, Y, Width, Height int
}

func (p Position) IsWithin(x, y int) bool {
	return x >= p.X && x <= p.X+p.Width && y >= p.Y && y <= p.Y+p.Height
}

type RectOverlay struct {
	hlg.Shape
	*Position
	Color   color.Color
	Handles [2]*Handle
}

type Handle struct {
	hlg.Shape
	*Position
	IsDragging bool
}

func NewHandle(p *Position) *Handle {
	return &Handle{
		Position: p,
		Shape:    hlg.Rectangle(p.X-p.Width/2, p.Y-p.Height/2, p.Width, p.Height, colornames.Orangered),
	}
}

func NewRectOverlay(x, y, width, height int, color color.Color) *RectOverlay {
	r := &RectOverlay{
		Shape: hlg.Rectangle(x, y, width, height, color),
		Position: &Position{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		},
		Color: color,
	}

	r.Handles = [2]*Handle{
		NewHandle(&Position{Width: 10, Height: 10, X: x + width/2, Y: y + width/2}), // Center
		NewHandle(&Position{X: x + width, Y: y + height, Width: 10, Height: 10}),
	}

	return r
}

func (h *Handle) IsWithin(x, y int) bool {
	return x >= h.Position.X-(h.Position.Width/2) && x <= (h.Position.X+(h.Position.Width/2)) &&
		y >= h.Position.Y-(h.Position.Height/2) && y <= (h.Position.Y+(h.Position.Height/2))
}

func (h *Handle) UpdatePosition(x, y int) {
	h.Position.X = x
	h.Position.Y = y
	h.Shape.Move(float32(h.Position.X), float32(h.Position.Y))
}

func (h *Handle) Render() {
	h.Shape.Render()
}

func (g *RectOverlay) Update() {
	x, y := hlg.GetCursorPosition()

	for _, handle := range g.Handles {
		if hlg.IsButtonPressed(input.MouseButtonLeft) && handle.IsWithin(x, y) {
			handle.IsDragging = true
		} else if !hlg.IsButtonPressed(input.MouseButtonLeft) {
			handle.IsDragging = false
		}
	}

	for index, handle := range g.Handles {
		if handle.IsDragging {
			switch index {
			case 0:
				g.X = x - g.Width/2
				g.Y = y - g.Height/2
				handle.UpdatePosition(x, y)
				g.Handles[1].UpdatePosition(x+g.Width/2, y+g.Height/2)
				g.newRect()
			case 1:
				g.Width = x - g.X
				g.Height = y - g.Y
				handle.UpdatePosition(x, y)
				g.Handles[0].UpdatePosition(x-g.Width/2, y-g.Height/2)
				g.newRect()
			}

			break // Only one handle can be dragged at a time
		}
	}
}

func (g *RectOverlay) newRect() {
	g.Shape = hlg.Rectangle(g.X, g.Y, g.Width, g.Height, g.Color)
}

func (g *RectOverlay) Render() {
	if g.Shape == nil {
		return
	}
	g.Shape.Render()
	for _, handle := range g.Handles {
		handle.Render()
	}
}

func main() {
	hlg.SetWindowSize(640, 480)
	overlay := NewRectOverlay(100, 100, 100, 100, colornames.Green)
	hlg.Run(func() {
		overlay.Update()
	}, func() {
		hlg.Clear(colornames.Skyblue)
		overlay.Render()
	})
}
