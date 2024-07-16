package main

import (
	"image/color"

	"github.com/dfirebaugh/hlg"
	"golang.org/x/image/colornames"
)

func main() {
	hlg.SetWindowSize(720, 480)
	hlg.SetScreenSize(240, 160)

	t := hlg.PolygonFromVertices(0, 0, 0, []hlg.Vertex{
		{
			Position: [3]float32{0, 160, 0},
			Color:    toRGBA(colornames.Red),
		},
		{
			Position: [3]float32{120, 0, 0},
			Color:    toRGBA(colornames.Green),
		},
		{
			Position: [3]float32{240, 160, 0},
			Color:    toRGBA(colornames.Blue),
		},
	})

	hlg.Update(func() {
		hlg.Clear(colornames.Skyblue)
		t.Render()
	})
}

func toRGBA(c color.Color) [4]float32 {
	r, g, b, a := c.RGBA()
	return [4]float32{
		float32(r) / 0xffff,
		float32(g) / 0xffff,
		float32(b) / 0xffff,
		float32(a) / 0xffff,
	}
}
