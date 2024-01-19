
# Textures

Textures can be created from images using `hlg.CreateTextureFromImage(img image.Image)`.

```golang
package main

import (
	"image"
	"net/http"

	_ "image/jpeg" // This is necessary to decode jpeg images

	"github.com/dfirebaugh/hlg"
)

// downloadImage fetches the image from the given URL and returns it as an image.Image
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

	t, _ := hlg.CreateTextureFromImage(
		downloadImage(`https://upload.wikimedia.org/wikipedia/commons/thumb/d/da/Nick_Offerman_2013_Headshot_%28cropped%29.jpg/308px-Nick_Offerman_2013_Headshot_%28cropped%29.jpg`),
	)

	hlg.Update(func() {
		t.Render()
	})
}
```

The above code should render this image.

![img](https://upload.wikimedia.org/wikipedia/commons/thumb/d/da/Nick_Offerman_2013_Headshot_%28cropped%29.jpg/308px-Nick_Offerman_2013_Headshot_%28cropped%29.jpg)
> sourced from: https://commons.wikimedia.org/w/index.php?curid=31678974

Textures implement the `Transformable` interface.
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
