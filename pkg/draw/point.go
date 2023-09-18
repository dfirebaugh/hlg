package draw

import (
	"image/color"

	"github.com/dfirebaugh/ggez/pkg/math/geom"
)

type Point geom.Point

func (p Point) Draw(d displayer, c color.Color) {
	d.SetPixel(int16(p.X), int16(p.Y), c.(color.RGBA))
}
