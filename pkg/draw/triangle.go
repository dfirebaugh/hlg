package draw

import (
	"image/color"

	"github.com/dfirebaugh/ggez/pkg/math/geom"
	"golang.org/x/image/colornames"
	"tinygo.org/x/tinydraw"
)

type Triangle geom.Triangle

func (t Triangle) Draw(d displayer, clr color.Color) {
	color, ok := clr.(color.RGBA)
	if !ok {
		color = colornames.Black
	}
	tinydraw.Triangle(d, int16(t[0][0]), int16(t[0][1]), int16(t[1][0]), int16(t[1][1]), int16(t[2][0]), int16(t[2][1]), color)
}

func (t Triangle) Fill(d displayer, clr color.Color) {
	color, ok := clr.(color.RGBA)
	if !ok {
		color = colornames.Black
	}
	tinydraw.FilledTriangle(d, int16(t[0][0]), int16(t[0][1]), int16(t[1][0]), int16(t[1][1]), int16(t[2][0]), int16(t[2][1]), color)
}
