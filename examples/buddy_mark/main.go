package main

import (
	"bytes"
	"fmt"
	"image"
	"math/rand"
	"net/http"
	"time"

	_ "image/png"

	_ "golang.org/x/image/webp" // This is necessary to decode WEBP images

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/assets"
	"github.com/dfirebaugh/hlg/pkg/input"
	"golang.org/x/image/colornames"
)

const (
	gravity      = 0.1
	damping      = 0.9
	buddyWidth   = 32
	buddyHeight  = 32
	screenWidth  = 800
	screenHeight = 600
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

type Buddy struct {
	X, Y                 float32
	VelocityX, VelocityY float32
	Sprite               *hlg.Sprite
}

var (
	buddies []*Buddy
	img     image.Image
)

func main() {
	var err error
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetTitle("BuddyMark Stress Test")
	hlg.EnableFPS()

	reader := bytes.NewReader(assets.BuddyDanceSpriteSheet)
	img, _, err = image.Decode(reader)
	if err != nil {
		panic(err)
	}

	lastFrameTime := time.Now()
	frameDuration := time.Millisecond * 100
	backgroundImg := downloadImage(`https://www.gstatic.com/webp/gallery/1.webp`)

	background, _ := hlg.CreateTextureFromImage(
		backgroundImg,
	)
	rq := hlg.CreateRenderQueue()
	background.RenderToQueue(rq)

	background.Resize(screenWidth, screenHeight)
	background.Move(0, 0)

	hlg.Run(func() {
		if time.Since(lastFrameTime) >= frameDuration {
			lastFrameTime = time.Now()
			for _, buddy := range buddies {
				buddy.Sprite.NextFrame()
			}
		}

		handleInput()

		for _, buddy := range buddies {
			buddy.Update()
		}
	}, func() {
		hlg.Clear(colornames.Skyblue)
		rq.Present() // Present background layer first (behind everything)
		for _, buddy := range buddies {
			buddy.Render()
		}
		hlg.PrintAt(fmt.Sprintf("buddies: %d", len(buddies)), 20, 20, colornames.Red)
	})
}

func handleInput() {
	if hlg.IsButtonPressed(input.MouseButtonLeft) {
		if hlg.GetFPS() < 40 {
			return
		}
		x, y := hlg.GetCursorPosition()
		for range 5 {
			buddy := NewBuddy(float32(x), float32(y))
			buddies = append(buddies, buddy)
		}
	}

	if hlg.IsButtonPressed(input.MouseButtonRight) {
		for _, b := range buddies {
			b.Sprite.Dispose()
		}
		buddies = []*Buddy{}
	}
}

func NewBuddy(x, y float32) *Buddy {
	frameSize := image.Point{X: buddyWidth, Y: buddyHeight}
	sheetSize := image.Point{X: 4, Y: 1}

	sprite := hlg.NewSprite(img, frameSize, sheetSize)
	sprite.Scale(2, 2)

	return &Buddy{
		X:         x,
		Y:         y,
		VelocityX: rand.Float32()*10 - 5,
		VelocityY: rand.Float32()*10 - 5,
		Sprite:    sprite,
	}
}

func (b *Buddy) Update() {
	b.VelocityY += gravity

	b.X += b.VelocityX
	b.Y += b.VelocityY

	if b.X < 0 {
		b.X = 0
		b.VelocityX *= -damping
	} else if b.X > screenWidth-buddyWidth {
		b.X = screenWidth - buddyWidth
		b.VelocityX *= -damping
	}

	if b.Y > screenHeight-buddyHeight {
		b.Y = screenHeight - buddyHeight
		b.VelocityY *= -damping
	}

	b.Sprite.Move(b.X, b.Y)
}

func (b *Buddy) Render() {
	b.Sprite.Render()
}
