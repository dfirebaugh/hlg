package main

import (
	"github.com/dfirebaugh/ggez"
	"golang.org/x/image/colornames"
)

var (
	x, y = 100, 100
)

func update() {
	ggez.Clear(colornames.Grey)
	if ggez.DefaultInput.IsKeyPressed(ggez.KeyA) || ggez.DefaultInput.IsKeyPressed(ggez.KeyLEFT) {
		x -= 1
	}
	if ggez.DefaultInput.IsKeyPressed(ggez.KeyD) || ggez.DefaultInput.IsKeyPressed(ggez.KeyRIGHT) {
		x += 1
	}
	if ggez.DefaultInput.IsKeyPressed(ggez.KeyW) || ggez.DefaultInput.IsKeyPressed(ggez.KeyUP) {
		y -= 1
	}
	if ggez.DefaultInput.IsKeyPressed(ggez.KeyS) || ggez.DefaultInput.IsKeyPressed(ggez.KeyDOWN) {
		y += 1
	}
	ggez.FillRectangle(x, y, 20, 20, colornames.Aliceblue)
}

func main() {
	ggez.Update(update)
}
