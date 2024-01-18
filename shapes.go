package ggez

import (
	"image/color"

	"github.com/dfirebaugh/ggez/graphics"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/proggy"
)

func Triangle(x1, y1, x2, y2, x3, y3 int, c color.Color) graphics.Shape {
	ensureSetupCompletion()
	return ggez.graphicsBackend.AddTriangle(x1, y1, x2, y2, x3, y3, c)
}

func Rectangle(x, y, width, height int, c color.Color) graphics.Shape {
	ensureSetupCompletion()
	return ggez.graphicsBackend.AddRectangle(x, y, width, height, c)
}

func Polygon(x, y int, width float32, sides int, c color.Color) graphics.Shape {
	ensureSetupCompletion()
	return ggez.graphicsBackend.AddCircle(x, y, width/2, c, sides)
}

func Circle(x, y int, radius float32, c color.Color) graphics.Shape {
	ensureSetupCompletion()
	return ggez.graphicsBackend.AddCircle(x, y, radius, c, 32)
}

func Line(x1, y1, x2, y2 int, width float32, c color.Color) graphics.Shape {
	ensureSetupCompletion()
	return ggez.graphicsBackend.AddLine(x1, y1, x2, y2, width, c)
}

func DrawTexture(t Texture) {
	ensureSetupCompletion()
	t.Render()
}

func PrintAt(s string, x int, y int, c color.Color) {
	ensureSetupCompletion()
	tinyfont.WriteLine(ggez.uifb, &proggy.TinySZ8pt7b, int16(x), int16(y), s, c.(color.RGBA))
}

func DrawModel(m graphics.Model, t graphics.Texture) {
	// graphicsBackend.RenderModel(m, t)
}
