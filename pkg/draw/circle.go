package draw

import (
	"image/color"

	"github.com/dfirebaugh/ggez/pkg/math/geom"
	"golang.org/x/image/colornames"
	"tinygo.org/x/tinydraw"
)

type Circle geom.Circle

func (c Circle) Draw(d displayer, clr color.Color) {
	color, ok := clr.(color.RGBA)
	if !ok {
		color = colornames.Black
	}
	tinydraw.Circle(d, int16(c.X), int16(c.Y), int16(c.R), color)
}

func (c Circle) Fill(d displayer, clr color.Color) {
	color, ok := clr.(color.RGBA)
	if !ok {
		color = colornames.Black
	}
	tinydraw.FilledCircle(d, int16(c.X), int16(c.Y), int16(c.R), color)
}
