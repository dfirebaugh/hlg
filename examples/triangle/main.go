package main

import (
	"github.com/dfirebaugh/hlg"
	"golang.org/x/image/colornames"
)

func main() {
	hlg.SetWindowSize(720, 480)
	hlg.SetScreenSize(240, 160)
	t := hlg.Triangle(0, 160, 120, 0, 240, 160, colornames.Orangered)

	hlg.Update(func() {
		hlg.Clear(colornames.Skyblue)
		t.Render()
	})
}
