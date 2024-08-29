package components

import (
	ui "github.com/dfirebaugh/hlg/pkg/grugui"
	"github.com/dfirebaugh/hlg/pkg/input"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

type Button struct {
	X, Y, Width, Height int
	Text                string
	IsPressed           bool
}

func (b *Button) Update(ctx ui.Context) {
	mouseX, mouseY := ctx.GetCursorPosition()

	if ctx.IsButtonPressed(input.MouseButtonLeft) &&
		mouseX > b.X && mouseX < b.X+b.Width &&
		mouseY > b.Y && mouseY < b.Y+b.Height {
		b.IsPressed = true
	} else {
		b.IsPressed = false
	}
}

func (b *Button) Render(ctx ui.Context) {
	theme := ctx.Theme()
	rectColor := theme.PrimaryColor
	if b.IsPressed {
		rectColor = theme.SecondaryColor
	}
	ctx.FillRoundedRectangle(b.X, b.Y, b.Width, b.Height, rectColor)

	face := basicfont.Face7x13
	textWidth := calculateTextWidth(b.Text, face)
	textHeight := face.Metrics().Ascent.Ceil() + face.Metrics().Descent.Ceil()
	textX := b.X + (b.Width-textWidth)/2
	textY := b.Y + (b.Height-textHeight)/2

	ctx.DrawText(textX, textY, b.Text, theme.TextColor)
}

func calculateTextWidth(text string, face font.Face) int {
	width := 0
	for _, char := range text {
		advance, _ := face.GlyphAdvance(rune(char))
		width += int(advance.Round())
	}
	return width
}
