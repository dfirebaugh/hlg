package draw

import (
	"image/color"

	"github.com/dfirebaugh/hlg/pkg/math/geom"
	"golang.org/x/image/colornames"
	"tinygo.org/x/tinydraw"
)

type Ray struct {
	geom.Ray
}

func (r Ray) Draw(d displayer, clr color.Color) {
	color, ok := clr.(color.RGBA)
	if !ok {
		color = colornames.Black
	}
	tinydraw.Line(d, int16(r.Origin[0]), int16(r.Origin[1]), int16(r.Direction[0]), int16(r.Direction[1]), color)
}
