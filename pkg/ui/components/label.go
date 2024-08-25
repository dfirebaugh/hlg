package components

import (
	"image/color"

	"github.com/dfirebaugh/hlg/pkg/ui"
)

type Label struct {
	X, Y  int
	Text  string
	Color color.Color
}

func NewLabel(x, y int, text string, color color.Color) *Label {
	return &Label{
		X:     x,
		Y:     y,
		Text:  text,
		Color: color,
	}
}

func (l *Label) Update(ctx ui.Context) {
}

func (l *Label) Render(ctx ui.Context) {
	if l.Text == "" {
		return
	}

	if l.Color == (color.Color)(nil) {
		l.Color = color.RGBA{255, 255, 255, 255}
	}

	ctx.DrawText(l.X, l.Y, l.Text, l.Color)
}
