package draw

import (
	"image/color"

	"github.com/dfirebaugh/hlg/pkg/math/geom"
)

type Point3D struct {
	X, Y, Z float64
}

type Point geom.Point

func (p Point) Draw(d displayer, c color.Color) {
	d.SetPixel(int16(p.X), int16(p.Y), c.(color.RGBA))
}
