
# Textures

Textures in the `hlg` package can be created from images, allowing for more complex and detailed visuals. Textures are particularly useful for rendering images or sprites.

## Creating a Texture from an Image

```golang
package main

import (
	"image"
	"net/http"
	_ "image/jpeg"
	"github.com/dfirebaugh/hlg"
)

func downloadImage(url string) image.Image {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		panic(err)
	}

	return img
}

func main() {
	hlg.SetWindowSize(154, 240)
	hlg.SetTitle("hlg texture example")

	t, _ := hlg.CreateTextureFromImage(downloadImage(`https://upload.wikimedia.org/wikipedia/commons/thumb/d/da/Nick_Offerman_2013_Headshot_%28cropped%29.jpg/308px-Nick_Offerman_2013_Headshot_%28cropped%29.jpg`))

	hlg.Run(nil, func() {
		t.Render()
	})
}
```

The above code should render this image.

![img](https://upload.wikimedia.org/wikipedia/commons/thumb/d/da/Nick_Offerman_2013_Headshot_%28cropped%29.jpg/308px-Nick_Offerman_2013_Headshot_%28cropped%29.jpg)

> sourced from: https://commons.wikimedia.org/w/index.php?curid=31678974


## Texture Interfaces

```golang
type Renderable interface {
	Render()
	Dispose()
	Hide()
}

type Transformable interface {
	Move(screenX, screenY float32)
	Rotate(angle float32)
	Scale(sx, sy float32)
}

type Texture interface {
	Renderable
	Transformable
	Handle() uintptr
	UpdateImage(img image.Image) error
	Resize(width, height float32)
	FlipVertical()
	FlipHorizontal()
	Clip(minX, minY, maxX, maxY float32)
}
```
