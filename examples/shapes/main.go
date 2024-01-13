package main

import (
	"github.com/dfirebaugh/ggez"
	"golang.org/x/image/colornames"
)

func main() {
	t := ggez.Triangle(0, 160, 120, 0, 240, 160, colornames.Green)
	r := ggez.Rectangle(0, 0, 120, 60, colornames.Blue)
	c := ggez.Circle(120, 80, 20, colornames.Red)
	l := ggez.Line(0, 0, 240, 160, 2, colornames.White)
	ggez.Clear(colornames.Skyblue)
	t.Render()
	r.Render()
	c.Render()
	l.Render()

	ggez.Update(func() {})
}
