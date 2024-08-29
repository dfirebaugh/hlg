package components

import (
	"image/color"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/gui"
	"github.com/dfirebaugh/hlg/pkg/input"
)

type Slider struct {
	Node
	x, y, width, height int
	handleRadius        int
	value               float32
	minValue            float32
	maxValue            float32
	isDragging          bool
	isHovered           bool
	normalOutline       color.Color
	hoverOutline        color.Color
	onChange            func(float32)
}

func NewSlider(x, y, width, height int, minValue, maxValue float32, onChange func(float32)) *Slider {
	handleRadius := height / 2

	s := &Slider{
		x:            x,
		y:            y,
		width:        width,
		height:       height,
		handleRadius: handleRadius,
		value:        minValue,
		minValue:     minValue,
		maxValue:     maxValue,
		onChange:     onChange,
	}

	return s
}

func (s *Slider) Render(ctx gui.DrawContext) {
	ctx.DrawRectangle(s.x, s.y, s.width, s.height, &gui.DrawOptions{
		Style: gui.Style{
			FillColor:    ctx.GetTheme().SecondaryColor,
			OutlineColor: s.getTrackOutlineColor(ctx),
			OutlineSize:  2,
			CornerRadius: 5,
		},
	})

	filledWidth := int(float32(s.width) * (s.value - s.minValue) / (s.maxValue - s.minValue))
	ctx.DrawRectangle(s.x, s.y, filledWidth, s.height, &gui.DrawOptions{
		Style: gui.Style{
			FillColor:    ctx.GetTheme().PrimaryColor,
			CornerRadius: 5,
		},
	})

	handleX := s.x + filledWidth
	handleY := s.y + s.height/2
	ctx.DrawCircle(handleX, handleY, s.handleRadius, &gui.DrawOptions{
		Style: gui.Style{
			FillColor:    ctx.GetTheme().TextColor,
			OutlineColor: ctx.GetTheme().HandleColor,
			OutlineSize:  1,
		},
	})
}

func (s *Slider) Update() {
	mouseX, mouseY := hlg.GetCursorPosition()
	if s.isPointWithin(mouseX, mouseY) {
		if !s.isHovered {
			s.isHovered = true
		}
	} else {
		if s.isHovered {
			s.isHovered = false
		}
	}

	if hlg.IsButtonPressed(input.MouseButtonLeft) {
		if s.isDragging || s.isHandlePointWithin(mouseX, mouseY) {
			s.isDragging = true

			newX := mouseX
			if newX < s.x {
				newX = s.x
			} else if newX > s.x+s.width {
				newX = s.x + s.width
			}

			s.value = s.minValue + (float32(newX-s.x)/float32(s.width))*(s.maxValue-s.minValue)

			if s.onChange != nil {
				s.onChange(s.value)
			}
		} else {
			s.isDragging = false
		}
	} else {
		s.isDragging = false
	}
}

func (s *Slider) GetValue() float32 {
	return s.value
}

func (s *Slider) isPointWithin(x, y int) bool {
	return x >= s.x && x <= s.x+s.width && y >= s.y && y <= s.y+s.height
}

func (s *Slider) isHandlePointWithin(x, y int) bool {
	handleX := s.x + int(float32(s.width)*(s.value-s.minValue)/(s.maxValue-s.minValue)) - s.handleRadius
	handleY := s.y + s.height/2 - s.handleRadius
	return x >= handleX && x <= handleX+s.handleRadius*2 && y >= handleY && y <= handleY+s.handleRadius*2
}

func (s *Slider) getTrackOutlineColor(ctx gui.DrawContext) color.Color {
	if s.isHovered {
		return ctx.GetTheme().PrimaryColor
	}
	return ctx.GetTheme().HandleColor
}
