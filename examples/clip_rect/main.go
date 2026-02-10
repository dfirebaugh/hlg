package main

import (
	"image"
	"net/http"

	_ "golang.org/x/image/webp" // This is necessary to decode WEBP images

	"github.com/dfirebaugh/hlg"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 320
	screenHeight = 240
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
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetScreenSize(screenWidth, screenHeight)
	hlg.SetTitle("Clip Rect Example")
	t, _ := hlg.CreateTextureFromImage(
		downloadImage(`https://parade.com/.image/c_limit%2Ccs_srgb%2Cq_auto:good%2Cw_620/MTkwNTgxNDg5MjU4ODY1Nzg5/nick-offerman-donkey-thoughts.webp`),
	)

	triangle := hlg.PolygonFromVertices(0, 0, 0, []hlg.Vertex{
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
	hlg.Run(func() {
		// No update logic needed
	}, func() {
		hlg.Clear(colornames.Darkgray)
		hlg.BeginDraw()

		triangle.Render()
		hlg.FilledCircle(80, 60, 40, colornames.Blue)

		hlg.RoundedRectOutline(120, 40, 80, 80, 0, 1, colornames.Cadetblue, colornames.Black)
		hlg.PushClipRect(120, 40, 80, 80)
		// hlg.PushClipRect(140, 60, 40, 40)
		hlg.FilledCircle(160, 80, 50, colornames.Green)
		hlg.FilledCircle(190, 80, 40, colornames.Red)
		// hlg.PopClipRect()
		hlg.FilledRect(105, 110, 30, 15, colornames.Yellow)
		hlg.PopClipRect()

		hlg.FilledCircle(220, 80, 30, colornames.Purple)

		hlg.RoundedRectOutline(40, 140, 240, 60, 0, 2, colornames.White, colornames.Black)
		hlg.PushClipRect(40, 140, 240, 60)
		t.Render()
		triangle.Render()
		hlg.FilledRect(0, 0, screenWidth/2, screenHeight, colornames.Orange)
		hlg.FilledCircle(160, 170, 40, colornames.Purple)
		hlg.PopClipRect()

		hlg.EndDraw()
	})
}
