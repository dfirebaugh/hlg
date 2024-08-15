package draw

import (
	"image/color"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/proggy"
)

func WriteLine(d displayer, s string, x int, y int, c color.Color) {
	tinyfont.WriteLine(d, &proggy.TinySZ8pt7b, int16(x), int16(y), s, c.(color.RGBA))
}
