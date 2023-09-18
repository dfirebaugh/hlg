package ggez

import (
	"image/color"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/proggy"
)

func DrawLine(x1, y1, x2, y2 int, c color.Color) {
	graphicsBackend.DrawLine(x1, y1, x2, y2, c)
}
func FillTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) {
	graphicsBackend.FillTriangle(x1, y1, x2, y2, x3, y3, c)
}
func DrawTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) {
	graphicsBackend.DrawTriangle(x1, y1, x2, y2, x3, y3, c)
}
func FillPolygon(xPoints, yPoints []int, c color.Color) {
	graphicsBackend.FillPolygon(xPoints, yPoints, c)
}
func FillRectangle(x, y, width, height int, c color.Color) {
	graphicsBackend.FillRect(x, y, width, height, c)
}
func DrawRectangle(x, y, width, height int, c color.Color) {
	graphicsBackend.DrawRect(x, y, width, height, c)
}
func FillCircle(x, y, radius int, c color.Color) {
	graphicsBackend.FillCirc(x, y, radius, c)
}
func DrawCircle(x, y, radius int, c color.Color) {
	graphicsBackend.DrawCirc(x, y, radius, c)
}

func DrawPoint(x, y int, c color.Color) {
	graphicsBackend.DrawPoint(x, y, c)
}

func DrawTexture(t Texture) {
	t.Render()
}

func PrintAt(s string, x int, y int, c color.Color) {
	tinyfont.WriteLine(uifb, &proggy.TinySZ8pt7b, int16(x), int16(y), s, c.(color.RGBA))
}
