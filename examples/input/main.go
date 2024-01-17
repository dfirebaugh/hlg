package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/dfirebaugh/ggez"
	"github.com/dfirebaugh/ggez/pkg/input"
	"golang.org/x/image/colornames"
)

var (
	colors = []color.Color{
		colornames.Red,
		colornames.Cadetblue,
		colornames.Black,
		colornames.Violet,
		colornames.Blue,
		colornames.Blueviolet,
		colornames.Orange,
	}
)

func main() {
	ggez.SetWindowSize(200, 200)

	triangle := ggez.Triangle(0, 200, 100, 0, 200, 200, colornames.Green)

	circle := ggez.Circle(10, 10, 10, colornames.Red)

	ggez.Update(func() {
		ggez.Clear(colornames.Grey)
		if ggez.IsKeyPressed(input.KeySpace) {
			triangle.SetColor(colors[rand.Intn(len(colors))])
		}
		if ggez.IsButtonPressed(input.MouseButtonLeft) {
			ggez.PrintAt("left mouse button pressed", 10, 25, colornames.Red)
		}
		if ggez.IsButtonPressed(input.MouseButtonRight) {
			ggez.PrintAt("right mouse button pressed", 10, 35, colornames.Red)
		}
		if ggez.IsButtonPressed(input.MouseButtonMiddle) {
			ggez.PrintAt("middle mouse button pressed", 10, 45, colornames.Red)
		}

		x, y := ggez.GetCursorPosition()
		ggez.PrintAt(fmt.Sprintf("%d:%d", x, y), 10, 10, colornames.Red)

		ggez.PrintAt("press space", 65, 180, colornames.White)

		triangle.Render()
		circle.Move(float32(x), float32(y))
		circle.Render()
	})
}
