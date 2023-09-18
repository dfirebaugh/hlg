package main

import (
	"github.com/dfirebaugh/ggez"
	"github.com/dfirebaugh/ggez/pkg/draw"
	"github.com/dfirebaugh/ggez/pkg/fb"
	"golang.org/x/image/colornames"
)

var (
	ui   *ggez.Texture
	uiFB fb.Displayer
)

func update() {
	// ggez.Clear(colornames.Grey)
	ui.Render()
}

func main() {
	ggez.SetScreenSize(200, 200)
	f := fb.New(ggez.ScreenWidth(), ggez.ScreenHeight())
	uiFB = f
	draw.WriteLine(uiFB, "hello, world", 20, 20, colornames.Red)

	ui, _ = ggez.CreateTextureFromImage(f.ToImage())

	ggez.Update(update)
}
