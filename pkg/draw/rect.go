package draw

import (
	"image/color"

	"github.com/dfirebaugh/hlg/pkg/math/geom"
	"golang.org/x/image/colornames"
	"tinygo.org/x/tinydraw"
)

type Rect geom.Rect

func (r Rect) Draw(d displayer, clr color.Color) {
	color, ok := clr.(color.RGBA)
	if !ok {
		color = colornames.Black
	}
	tinydraw.Rectangle(d, int16(r[0]), int16(r[1]), int16(r[2]), int16(r[3]), color)
}

func (r Rect) Fill(d displayer, clr color.Color) {
	color, ok := clr.(color.RGBA)
	if !ok {
		color = colornames.Black
	}
	tinydraw.FilledRectangle(d, int16(r[0]), int16(r[1]), int16(r[2]), int16(r[3]), color)
}
