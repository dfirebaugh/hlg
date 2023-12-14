package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/dfirebaugh/ggez"
	"github.com/dfirebaugh/ggez/pkg/draw"
	"github.com/dfirebaugh/ggez/pkg/fb"
	"github.com/dfirebaugh/ggez/pkg/input"
	"github.com/dfirebaugh/ggez/pkg/math/geom"
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

type Triangle struct {
	geom.Triangle
	color.Color
	*ggez.Texture
	*fb.ImageFB
}

func (t *Triangle) Render() {
	// using the cpu, we draw the triangle to an image
	draw.Triangle(t.Triangle).Fill(t.ImageFB, t.Color)
	t.Texture.UpdateImage(t.ImageFB.ToImage())
	t.Texture.Render()
}

func main() {
	ggez.SetWindowSize(200, 200)

	triangle := &Triangle{
		Triangle: geom.MakeTriangle([3]geom.Vector{
			geom.MakeVector(0, 200),
			geom.MakeVector(100, 0),
			geom.MakeVector(200, 200),
		}),
		ImageFB: fb.New(200, 200),
		Color:   colornames.Green,
	}

	// upload the image of our triangle to a texture on the gpu
	triangle.Texture, _ = ggez.CreateTextureFromImage(triangle.ToImage())

	ggez.Update(func() {
		ggez.Clear(colornames.Grey)
		if ggez.IsKeyPressed(input.KeySpace) {
			triangle.Color = colors[rand.Intn(len(colors))]
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
	})
}
