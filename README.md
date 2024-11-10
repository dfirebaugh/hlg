# hlg  (High Level Graphics)
This project is a WIP. The goal is to make a high level graphics api for golang.

[![Go Reference](https://pkg.go.dev/badge/github.com/dfirebaugh/hlg.svg)](https://pkg.go.dev/github.com/dfirebaugh/hlg)
Documentation: https://dfirebaugh.github.io/hlg/

### Examples
check the `./examples` dir for some basic examples


#### Triangle

```golang
package main

import (
	"github.com/dfirebaugh/hlg"
	"golang.org/x/image/colornames"
)

var triangle hlg.Shape

// update operation need to happen less frequently than render operations
func update() {
}

func render() {
	hlg.Clear(colornames.Skyblue)
	triangle.Render()
}

func main() {
	hlg.SetWindowSize(720, 480)
	hlg.SetScreenSize(240, 160)
	triangle = hlg.Triangle(0, 160, 120, 0, 240, 160, colornames.Orangered)

	hlg.Run(update, render)
}
```

![triangle_example](./assets/images/triangle_example.png)

#### Colored Triangle

```golang
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

// update operations happen less frequently than render operations
func update() {
}

func render() {
	hlg.Clear(colornames.Skyblue)
	triangle.Render()
}

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetTitle("color triangle")

	triangle = hlg.PolygonFromVertices(0, 0, 0, []hlg.Vertex{
		{
			Position: [3]float32{0, screenHeight, 0},
			Color:    colornames.Red,
		},
		{
			Position: [3]float32{screenWidth / 2, 0, 0},
			Color:    colornames.Green,
		},
		{
			Position: [3]float32{screenWidth, screenHeight, 0},
			Color:    colornames.Blue,
		},
	})

	hlg.Run(update, render)
}
```

![color_triangle_example](./assets/images/color_triangle_example.png)

