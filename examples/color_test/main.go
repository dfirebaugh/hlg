package main

import (
	"image/color"

	"github.com/dfirebaugh/hlg"
)

const (
	screenWidth  = 400
	screenHeight = 300
)

// Test colors with their expected sRGB values
var testColors = []struct {
	name  string
	color color.RGBA
}{
	// Primaries
	{"Red", color.RGBA{255, 0, 0, 255}},
	{"Green", color.RGBA{0, 255, 0, 255}},
	{"Blue", color.RGBA{0, 0, 255, 255}},

	// Secondaries
	{"Yellow", color.RGBA{255, 255, 0, 255}},
	{"Cyan", color.RGBA{0, 255, 255, 255}},
	{"Magenta", color.RGBA{255, 0, 255, 255}},

	// Grays - key test for gamma correctness
	// 50% gray (128) should look perceptually mid-gray
	{"White", color.RGBA{255, 255, 255, 255}},
	{"Gray 75%", color.RGBA{191, 191, 191, 255}},
	{"Gray 50%", color.RGBA{128, 128, 128, 255}}, // Critical test!
	{"Gray 25%", color.RGBA{64, 64, 64, 255}},
	{"Black", color.RGBA{0, 0, 0, 255}},

	// Mid-tones
	{"Orange", color.RGBA{255, 128, 0, 255}},
	{"Pink", color.RGBA{255, 128, 128, 255}},
	{"Purple", color.RGBA{128, 0, 255, 255}},
	{"Teal", color.RGBA{0, 128, 128, 255}},

	// Semi-transparent (test alpha handling)
	{"Red 50%", color.RGBA{255, 0, 0, 128}},
	{"Green 50%", color.RGBA{0, 255, 0, 128}},
	{"Blue 50%", color.RGBA{0, 0, 255, 128}},
}

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetTitle("Color Accuracy Test")

	hlg.Run(func() {}, func() {
		hlg.Clear(color.RGBA{40, 40, 40, 255})
		hlg.BeginDraw()

		// Draw color swatches in a grid
		cols := 6
		swatchW := 50
		swatchH := 50
		padding := 10
		startX := 20
		startY := 20

		for i, tc := range testColors {
			col := i % cols
			row := i / cols

			x := startX + col*(swatchW+padding)
			y := startY + row*(swatchH+padding)

			// Draw outline (white border)
			hlg.FilledRect(x-1, y-1, swatchW+2, swatchH+2, color.RGBA{255, 255, 255, 255})
			// Draw swatch
			hlg.FilledRect(x, y, swatchW, swatchH, tc.color)
		}

		hlg.EndDraw()
	})
}
