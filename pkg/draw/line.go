package draw

import (
	"image/color"

	"github.com/dfirebaugh/ggez/pkg/math/geom"
	"golang.org/x/image/colornames"
	"tinygo.org/x/tinydraw"
)

type Line geom.Line

func (l Line) Draw(d displayer, clr color.Color) {
	color, ok := clr.(color.RGBA)
	if !ok {
		color = colornames.Black
	}
	tinydraw.Line(d, int16(l.Segment.V0[0]), int16(l.Segment.V0[1]), int16(l.Segment.V1[0]), int16(l.Segment.V1[1]), color)
}
