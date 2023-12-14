package main

import (
	"image/color"

	"github.com/dfirebaugh/ggez"
	"github.com/dfirebaugh/ggez/pkg/draw"
	"github.com/dfirebaugh/ggez/pkg/fb"
	"github.com/dfirebaugh/ggez/pkg/math/geom"
	"golang.org/x/image/colornames"
)

type Triangle struct {
	geom.Triangle
	Color color.RGBA
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
		triangle.Render()
	})
}
