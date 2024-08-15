package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/pkg/input"
	"golang.org/x/image/colornames"
)

var colors = []color.Color{
	colornames.Red,
	colornames.Cadetblue,
	colornames.Black,
	colornames.Violet,
	colornames.Blue,
	colornames.Blueviolet,
	colornames.Orange,
}

func main() {
	hlg.SetWindowSize(200, 200)

	triangle := hlg.Triangle(0, 200, 100, 0, 200, 200, colornames.Green)

	circle := hlg.Circle(10, 10, 10, colornames.Red)

	hlg.Run(func() {
	}, func() {
		hlg.Clear(colornames.Grey)
		if hlg.IsKeyPressed(input.KeySpace) {
			triangle.SetColor(colors[rand.Intn(len(colors))])
		}
		if hlg.IsButtonPressed(input.MouseButtonLeft) {
			hlg.PrintAt("left mouse button pressed", 10, 25, colornames.Red)
		}
		if hlg.IsButtonPressed(input.MouseButtonRight) {
			hlg.PrintAt("right mouse button pressed", 10, 35, colornames.Red)
		}
		if hlg.IsButtonPressed(input.MouseButtonMiddle) {
			hlg.PrintAt("middle mouse button pressed", 10, 45, colornames.Red)
		}

		x, y := hlg.GetCursorPosition()
		hlg.PrintAt(fmt.Sprintf("%d:%d", x, y), 10, 10, colornames.Red)

		hlg.PrintAt("press space", 65, 180, colornames.White)

		triangle.Render()
		circle.Move(float32(x), float32(y))
		circle.Render()
	})
}
