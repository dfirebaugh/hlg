package main

import (
	"github.com/dfirebaugh/ggez"
	"golang.org/x/image/colornames"
)

func update() {
	ggez.Clear(colornames.Grey)
	ggez.PrintAt("hello, world", 10, 30, colornames.Red)
}

func main() {
	ggez.SetScreenSize(200, 200)
	ggez.Update(update)
}
