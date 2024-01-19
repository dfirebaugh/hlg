
# Intro


[![Go Reference](https://pkg.go.dev/badge/github.com/dfirebaugh/hlg.svg)](https://pkg.go.dev/github.com/dfirebaugh/hlg)

`hlg` (High Level Graphics)

This project is a WIP. The goal is to make a high level graphics api for golang primarily using https://github.com/rajveermalviya/go-webgpu.


### Examples
check the `./examples` dir for some basic examples


#### Triangle

```golang
package main

import (
	"github.com/dfirebaugh/hlg"
	"golang.org/x/image/colornames"
)

func main() {
	hlg.SetWindowSize(720, 480)
	hlg.SetScreenSize(240, 160)
	t := hlg.Triangle(0, 160, 120, 0, 240, 160, colornames.Orangered)

	hlg.Update(func() {
		hlg.Clear(colornames.Skyblue)
		t.Render()
	})
}
```

![triangle_example](./assets/triangle_example.png)
