package main

import (
	"image/color"

	"github.com/dfirebaugh/hlg"
	"golang.org/x/image/colornames"
)

var triangle hlg.Shape

const (
	screenWidth  = 240
	screenHeight = 160
)

// update operations happen less frequently than render operations
func update() {
}

func render() {
	hlg.Clear(colornames.Skyblue)
	triangle.Render()
}

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetTitle("color triangle")

	triangle = hlg.PolygonFromVertices(0, 0, 0, []hlg.Vertex{
		{
			Position: [3]float32{0, screenHeight, 0},
			Color:    toRGBA(colornames.Red),
		},
		{
			Position: [3]float32{screenWidth / 2, 0, 0},
			Color:    toRGBA(colornames.Green),
		},
		{
			Position: [3]float32{screenWidth, screenHeight, 0},
			Color:    toRGBA(colornames.Blue),
		},
	})

	hlg.Run(update, render)
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
