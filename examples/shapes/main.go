package main

import (
	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/pkg/primitives"
	"golang.org/x/image/colornames"
)

func main() {
	l := primitives.NewLine([]primitives.Position{
		{X: 0, Y: 0},
		{X: 30, Y: 40},
		{X: 220, Y: 10},
		{X: 140, Y: 200},
		{X: 80, Y: 90},
	}, 4, colornames.Tomato)
	c := primitives.NewCircle(40, 40, 20, colornames.Steelblue, 2)
	c.SetColor(colornames.Purple)
	c.SetOutlineColor(colornames.Purple)

	r := primitives.NewRect(140, 40, 60, 60, 5, colornames.Midnightblue)
	r.SetOutlineWidth(3)

	hlg.Run(func() {
	}, func() {
		hlg.Clear(colornames.Skyblue)
		c.Render()
		l.Render()
		r.Render()
	})
}
