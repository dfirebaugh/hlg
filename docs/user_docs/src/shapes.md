
# Shapes

The `hlg` package provides functionality to create and render various shapes. These shapes include triangles, rectangles, circles, and lines. Each shape can be manipulated through transformations such as move, rotate, and scale.

```golang
func Triangle(x1, y1, x2, y2, x3, y3 int, c color.Color) graphics.Shape 
func Rectangle(x, y, width, height int, c color.Color) graphics.Shape 
func Polygon(x, y int, width float32, sides int, c color.Color) graphics.Shape 
func Circle(x, y int, radius float32, c color.Color) graphics.Shape 
func Line(x1, y1, x2, y2 int, width float32, c color.Color) graphics.Shape 
```

## Example Usage

```golang
package main

import (
	"github.com/dfirebaugh/hlg"
	"golang.org/x/image/colornames"
)

func main() {
	t := hlg.Triangle(0, 160, 120, 0, 240, 160, colornames.Green)
	r := hlg.Rectangle(0, 0, 120, 60, colornames.Blue)
	r2 := hlg.Rectangle(50, 50, 120, 60, colornames.Red)
	c := hlg.Circle(120, 80, 20, colornames.Red)
	l := hlg.Line(0, 0, 240, 160, 2, colornames.White)

	c.SetColor(colornames.Purple)
	c.Move(0, 0)

	hlg.Run(nil, func() {
		hlg.Clear(colornames.Skyblue)
		t.Render()
		r.Render()
		c.Render()
		l.Render()
		r2.Render()
		r2.Hide()
	})
}
```

## Shape Interfaces

```golang
type Transformable interface {
	Move(screenX, screenY float32)
	Rotate(angle float32)
	Scale(sx, sy float32)
}

type Renderable interface {
	Render()
	Dispose()
	Hide()
}

type Shape interface {
	Renderable
	Transformable
	SetColor(c color.Color)
}
```
