package main

import (
	"github.com/dfirebaugh/hlg"
	"golang.org/x/image/colornames"
)

var triangle hlg.Shape

const (
	screenWidth  = 240
	screenHeight = 160
)

// update operation need to happen less frequently than render operations
func update() {
}

func render() {
	hlg.Clear(colornames.Skyblue)
	triangle.Render()
}

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	triangle = hlg.Triangle(0, screenHeight, screenWidth/2, 0, screenWidth, screenHeight, colornames.Orangered)

	hlg.Run(update, render)
}
