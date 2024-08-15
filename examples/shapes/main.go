package main

import (
	"github.com/dfirebaugh/hlg"
	"golang.org/x/image/colornames"
)

func main() {
	t := hlg.Triangle(0, 160, 120, 0, 240, 160, colornames.Green)
	r := hlg.Rectangle(0, 0, 120, 60, colornames.Blue)
	r2 := hlg.Rectangle(50, 50, 120, 60, colornames.Red)
	c := hlg.Circle(120, 80, 20, colornames.Red)
	l := hlg.Line(0, 0, 240, 160, 2, colornames.White)
	hlg.Clear(colornames.Skyblue)
	t.Render()
	r.Render()
	c.Render()
	l.Render()
	r2.Render()

	c.SetColor(colornames.Purple)
	c.Move(0, 0)
	r2.Hide()

	hlg.Run(func() {
	}, func() {})
}
