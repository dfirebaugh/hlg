package main

import (
	"image/color"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/gui"
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
	*Position
	Color   color.Color
	Handles [2]*Handle
}

type Handle struct {
	*Position
	IsDragging bool
	Color      color.Color
}

func NewHandle(p *Position, color color.Color) *Handle {
	return &Handle{
		Position: p,
		Color:    color,
	}
}

func NewRectOverlay(x, y, width, height int, color color.Color) *RectOverlay {
	r := &RectOverlay{
		Position: &Position{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		},
		Color: color,
	}

	r.Handles = [2]*Handle{
		NewHandle(&Position{Width: 10, Height: 10, X: x + width/2, Y: y + height/2}, colornames.Orangered), // Center
		NewHandle(&Position{X: x + width, Y: y + height, Width: 10, Height: 10}, colornames.Orangered),
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
			case 1:
				g.Width = x - g.X
				g.Height = y - g.Y
				handle.UpdatePosition(x, y)
				g.Handles[0].UpdatePosition(x-g.Width/2, y-g.Height/2)
			}
			break // Only one handle can be dragged at a time
		}
	}
}

func (g *RectOverlay) Render(d *gui.Draw) {
	d.DrawRectangle(g.X, g.Y, g.Width, g.Height, &gui.DrawOptions{
		Style: gui.Style{
			FillColor: g.Color,
		},
	})

	for _, handle := range g.Handles {
		d.DrawRectangle(handle.X-handle.Width/2, handle.Y-handle.Height/2, handle.Width, handle.Height, &gui.DrawOptions{
			Style: gui.Style{
				FillColor: handle.Color,
			},
		})
	}
}

func main() {
	hlg.SetWindowSize(640, 480)
	overlay := NewRectOverlay(100, 100, 100, 100, colornames.Green)

	hlg.Run(func() {
		overlay.Update()
	}, func() {
		hlg.Clear(colornames.Skyblue)

		d := gui.Draw{
			ScreenWidth:  640,
			ScreenHeight: 480,
		}

		overlay.Render(&d)
		hlg.SubmitDrawBuffer(d.Encode())
	})
}
