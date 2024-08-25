package components

import (
	"github.com/dfirebaugh/hlg/pkg/input"
	"github.com/dfirebaugh/hlg/pkg/ui"
)

type Slider struct {
	X, Y, Width, Height int
	Value               float64
	isDragging          bool
}

func (s *Slider) Update(ctx ui.Context) {
	mouseX, mouseY := ctx.GetCursorPosition()

	if ctx.IsButtonJustPressed(input.MouseButtonLeft) &&
		mouseX > s.X && mouseX < s.X+s.Width &&
		mouseY > s.Y && mouseY < s.Y+s.Height {
		s.isDragging = true
	}

	if s.isDragging {
		if ctx.IsButtonPressed(input.MouseButtonLeft) {
			s.Value = float64(mouseX-s.X) / float64(s.Width)
			if s.Value < 0 {
				s.Value = 0
			} else if s.Value > 1 {
				s.Value = 1
			}
		} else {
			s.isDragging = false
		}
	}
}

func (s *Slider) Render(ctx ui.Context) {
	theme := ctx.Theme()

	ctx.FillRoundedRectangle(s.X, s.Y, s.Width, s.Height, theme.SecondaryColor)

	handleX := s.X + int(s.Value*float64(s.Width))
	handleRadius := s.Height / 2

	ctx.FillRoundedRectangle(s.X, s.Y, handleX-s.X, s.Height, theme.PrimaryColor)

	ctx.FillCircle(handleX, s.Y+handleRadius, handleRadius, theme.HandleColor)
}
