package main

import (
	"image"
	"net/http"

	"golang.org/x/image/colornames"
	_ "golang.org/x/image/webp" // This is necessary to decode WEBP images

	"github.com/dfirebaugh/ggez"
)

var tex *ggez.Texture

func update() {
	ggez.Clear(colornames.Aliceblue)
	tex.Render()
}

// downloadImage fetches the image from the given URL and returns it as an image.Image
func downloadImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func main() {
	ggez.SetRenderer(ggez.GLRenderer)
	imgURL := "https://parade.com/.image/c_limit%2Ccs_srgb%2Cq_auto:good%2Cw_620/MTkwNTgxNDg5MjU4ODY1Nzg5/nick-offerman-donkey-thoughts.webp"
	img, err := downloadImage(imgURL)
	if err != nil {
		panic(err)
	}
	ggez.SetScreenSize(img.Bounds().Max.X, img.Bounds().Max.Y)

	tex, err = ggez.CreateTextureFromImage(img)
	if err != nil {
		panic(err)
	}

	ggez.Update(update)
}
