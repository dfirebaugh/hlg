package main

import (
	"image"
	"net/http"

	"golang.org/x/image/colornames"
	_ "golang.org/x/image/webp" // This is necessary to decode WEBP images

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
	hlg.SetWindowSize(600, 412)
	hlg.SetTitle("hlg texture example")

	t, _ := hlg.CreateTextureFromImage(
		downloadImage(`https://parade.com/.image/c_limit%2Ccs_srgb%2Cq_auto:good%2Cw_620/MTkwNTgxNDg5MjU4ODY1Nzg5/nick-offerman-donkey-thoughts.webp`),
	)

	mountain, _ := hlg.CreateTextureFromImage(
		downloadImage(`https://www.gstatic.com/webp/gallery/1.webp`),
	)
	mountain.Resize(300, 206)
	// mountain.Move(100, 100)

	t.Scale(.4, .4)
	t.Move(0, 0)
	// t.FlipVertical()
	// t.FlipHorizontal()

	hlg.Run(func() {
	}, func() {
		hlg.Clear(colornames.Aliceblue)
		mountain.Render()
		t.Render()
	})
}
