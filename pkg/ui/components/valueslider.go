package components

import (
	"fmt"

	"github.com/dfirebaugh/hlg/pkg/input"
	"github.com/dfirebaugh/hlg/pkg/ui"
	"golang.org/x/image/font/basicfont"
)

type ValueSlider struct {
	X, Y, Width, Height int
	Value               float64
	isDragging          bool
	MinValue, MaxValue  float64
}

func (vs *ValueSlider) Update(ctx ui.Context) {
	mouseX, mouseY := ctx.GetCursorPosition()

	if ctx.IsButtonJustPressed(input.MouseButtonLeft) &&
		mouseX > vs.X && mouseX < vs.X+vs.Width &&
		mouseY > vs.Y && mouseY < vs.Y+vs.Height {
		vs.isDragging = true
	}

	if vs.isDragging {
		if ctx.IsButtonPressed(input.MouseButtonLeft) {
			vs.Value = float64(mouseX-vs.X) / float64(vs.Width)
			if vs.Value < 0 {
				vs.Value = 0
			} else if vs.Value > 1 {
				vs.Value = 1
			}
		} else {
			vs.isDragging = false
		}
	}
}

func (vs *ValueSlider) Render(ctx ui.Context) {
	actualValue := vs.MinValue + vs.Value*(vs.MaxValue-vs.MinValue)
	theme := ctx.Theme()

	ctx.FillRoundedRectangle(vs.X, vs.Y, vs.Width, vs.Height, theme.SecondaryColor)

	handleX := vs.X + int(vs.Value*float64(vs.Width))
	handleWidth := vs.Height

	ctx.FillRoundedRectangle(vs.X, vs.Y, handleX-vs.X, vs.Height, theme.PrimaryColor)

	ctx.FillRect(handleX-handleWidth/2, vs.Y, handleWidth, vs.Height, theme.HandleColor)

	valueText := fmt.Sprintf("%.2f", actualValue)
	textColor := theme.TextColor

	textWidth := ctx.TextWidth(valueText)
	face := basicfont.Face7x13
	textHeight := face.Metrics().Ascent.Ceil()
	textX := vs.X + (vs.Width-textWidth)/2
	textY := vs.Y + ((vs.Height - textHeight) / 2)

	ctx.DrawText(textX, textY, valueText, textColor)
}
