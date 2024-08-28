package main

import (
	"bytes"
	"image"
	"time"

	_ "image/png"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/assets"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetTitle("hlg sprite example")

	reader := bytes.NewReader(assets.BuddyDanceSpriteSheet)
	img, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}

	frameSize := image.Point{X: 32, Y: 32}
	sheetSize := image.Point{X: 4, Y: 1}

	sprite := hlg.NewSprite(img, frameSize, sheetSize)
	sprite.Resize(512, 512)
	sprite.Move(screenWidth/2-256, screenHeight/2-256)
	// sprite.Scale(4, 4)

	lastFrameTime := time.Now()
	frameDuration := time.Millisecond * 200

	hlg.Run(func() {
		if time.Since(lastFrameTime) >= frameDuration {
			lastFrameTime = time.Now()
			sprite.NextFrame()
		}
	}, func() {
		hlg.Clear(colornames.Skyblue)
		sprite.Render()
	})
}
