package main

import (
	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/gui"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 240
	screenHeight = 160
)

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetTitle("Shapes Example")

	hlg.Run(func() {
		// No update logic needed for this example
	}, func() {
		hlg.Clear(colornames.Skyblue)
		d := gui.Draw{
			ScreenWidth:  screenWidth,
			ScreenHeight: screenHeight,
		}
		d.DrawRectangle(20, 20, 20, 20, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Green,
				OutlineColor: colornames.Purple,
			},
		})
		d.DrawRectangle(50, 50, 70, 70, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Green,
				OutlineColor: colornames.Purple,
				CornerRadius: 2,
			},
		})
		d.DrawRectangle(140, 50, 40, 40, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Green,
				OutlineColor: colornames.Purple,
				CornerRadius: 2,
				OutlineSize:  4,
			},
		})

		d.DrawCircle(60, 100, 15, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Red,
				OutlineColor: colornames.Black,
				OutlineSize:  2,
			},
		})
		d.DrawCircle(100, 100, 20, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Blue,
				OutlineColor: colornames.White,
				OutlineSize:  4,
			},
		})

		d.DrawTriangle(160, 120, 180, 140, 140, 140, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Orange,
				OutlineColor: colornames.Black,
				OutlineSize:  1,
			},
		})
		d.DrawTriangle(180, 60, 200, 80, 160, 80, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Yellow,
				OutlineColor: colornames.Black,
				OutlineSize:  1,
			},
		})

		d.DrawLine([]gui.Position{
			{X: 20, Y: 150},
			{X: 60, Y: 130},
			{X: 100, Y: 150},
			{X: 140, Y: 130},
			{X: 180, Y: 150},
		}, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Red,
				OutlineColor: colornames.Black,
				OutlineSize:  1,
			},
		})
		hlg.SubmitDrawBuffer(d.Encode())
	})
}
