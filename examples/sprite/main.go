package main

import (
	"bytes"
	"image"
	"time"

	_ "image/png"

	"github.com/dfirebaugh/ggez"
	"github.com/dfirebaugh/ggez/assets"
)

func main() {
	ggez.SetWindowSize(600, 600)
	ggez.SetTitle("ggez sprite example")

	reader := bytes.NewReader(assets.BuddyDanceSpriteSheet)
	img, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}

	frameSize := image.Point{X: 32, Y: 32}
	sheetSize := image.Point{X: 4, Y: 1}

	sprite := ggez.NewSprite(img, frameSize, sheetSize)
	sprite.Scale(.5, .5)

	lastFrameTime := time.Now()
	frameDuration := time.Millisecond * 200

	ggez.Update(func() {
		if time.Since(lastFrameTime) >= frameDuration {
			lastFrameTime = time.Now()
			sprite.NextFrame()
		}
	})
}
