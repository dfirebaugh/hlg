package main

import (
	"github.com/dfirebaugh/hlg"
	"golang.org/x/image/colornames"
)

func update() {
}

func render() {
	hlg.Clear(colornames.Grey)
	hlg.PrintAt("hello, world", 10, 30, colornames.Red)
}

func main() {
	hlg.SetWindowSize(200, 200)
	hlg.Run(update, render)
}
