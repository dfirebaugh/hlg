package components

import (
	ui "github.com/dfirebaugh/hlg/pkg/grugui"
	"github.com/dfirebaugh/hlg/pkg/input"
)

type Checkbox struct {
	X, Y, Size int
	IsChecked  bool
	wasPressed bool
}

func (cb *Checkbox) Update(ctx ui.Context) {
	mouseX, mouseY := ctx.GetCursorPosition()

	if ctx.IsButtonJustPressed(input.MouseButtonLeft) &&
		mouseX > cb.X && mouseX < cb.X+cb.Size &&
		mouseY > cb.Y && mouseY < cb.Y+cb.Size {
		if !cb.wasPressed {
			cb.IsChecked = !cb.IsChecked
			cb.wasPressed = true
		}
	} else {
		cb.wasPressed = false
	}
}

func (cb *Checkbox) Render(ctx ui.Context) {
	theme := ctx.Theme()

	ctx.FillRect(cb.X, cb.Y, cb.Size, cb.Size, theme.BackgroundColor)
	ctx.DrawRect(cb.X, cb.Y, cb.Size, cb.Size, theme.PrimaryColor)

	if cb.IsChecked {
		ctx.DrawLine(cb.X+4, cb.Y+cb.Size/2, cb.X+cb.Size/2, cb.Y+cb.Size-4, theme.TextColor)
		ctx.DrawLine(cb.X+cb.Size/2, cb.Y+cb.Size-4, cb.X+cb.Size-4, cb.Y+4, theme.TextColor)
	}
}
