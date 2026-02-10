package main

import (
	"github.com/dfirebaugh/hlg"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 240
	screenHeight = 160
)

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetScreenSize(screenWidth, screenHeight)
	hlg.SetTitle("Shapes Example")

	hlg.Run(func() {
		// No update logic needed for this example
	}, func() {
		hlg.Clear(colornames.Skyblue)
		hlg.BeginDraw()

		hlg.FilledRect(20, 20, 20, 20, colornames.Green)
		hlg.RoundedRect(50, 50, 70, 70, 2, colornames.Green)
		hlg.RoundedRectOutline(140, 50, 40, 40, 2, 4, colornames.Green, colornames.Purple)

		hlg.FilledCircle(60, 100, 17, colornames.Black)
		hlg.FilledCircle(60, 100, 15, colornames.Red)
		hlg.FilledCircle(100, 100, 24, colornames.White)
		hlg.FilledCircle(100, 100, 20, colornames.Blue)

		hlg.FilledTriangle(160, 140, 200, 140, 180, 110, colornames.Orange)

		hlg.Segment(20, 150, 60, 130, 1, colornames.Red)
		hlg.Segment(60, 130, 100, 150, 1, colornames.Red)
		hlg.Segment(100, 150, 140, 130, 1, colornames.Red)
		hlg.Segment(140, 130, 180, 150, 1, colornames.Red)

		hlg.EndDraw()
	})
}
