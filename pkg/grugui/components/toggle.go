package components

import (
	ui "github.com/dfirebaugh/hlg/pkg/grugui"
	"github.com/dfirebaugh/hlg/pkg/input"
)

type Toggle struct {
	X, Y, Width, Height int
	IsOn                bool
	wasPressed          bool
}

func (t *Toggle) Update(ctx ui.Context) {
	mouseX, mouseY := ctx.GetCursorPosition()

	if ctx.IsButtonJustPressed(input.MouseButtonLeft) &&
		mouseX > t.X && mouseX < t.X+t.Width &&
		mouseY > t.Y && mouseY < t.Y+t.Height {
		if !t.wasPressed {
			t.IsOn = !t.IsOn
			t.wasPressed = true
		}
	} else {
		t.wasPressed = false
	}
}

func (t *Toggle) Render(ctx ui.Context) {
	theme := ctx.Theme()

	bgColor := theme.SecondaryColor
	if t.IsOn {
		bgColor = theme.PrimaryColor
	}

	ctx.FillRoundedRectangle(t.X, t.Y, t.Width, t.Height, bgColor)

	sliderRadius := t.Height / 2
	sliderX := t.X
	if t.IsOn {
		sliderX = t.X + t.Width - sliderRadius*2
	}

	ctx.FillCircle(sliderX+sliderRadius, t.Y+sliderRadius, sliderRadius, theme.HandleColor)
}
