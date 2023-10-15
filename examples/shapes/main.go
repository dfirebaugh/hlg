package main

import (
	"github.com/dfirebaugh/ggez"
	"golang.org/x/image/colornames"
)

func update() {
	ggez.Clear(colornames.Grey)
	ggez.FillRectangle(20, 20, 60, 60, colornames.Aliceblue)

	ggez.DrawCircle(200, 200, 50, colornames.Red)

	ggez.DrawTriangle(300, 300, 350, 250, 400, 300, colornames.Green)

	ggez.DrawLine(10, 10, 150, 150, colornames.Blue)

	ggez.DrawPoint(300, 10, colornames.Black)

	xPoints := []int{500, 550, 575, 525, 475}
	yPoints := []int{100, 150, 125, 175, 150}
	ggez.FillPolygon(xPoints, yPoints, colornames.Orange)

	ggez.DrawPoint(20, 20, colornames.Red)
	ggez.DrawPoint(10, 20, colornames.Red)
}

func main() {
	ggez.SetScreenSize(960, 640)
	ggez.Update(update)
}
