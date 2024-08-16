package main

import (
	"bytes"
	"fmt"
	"image"
	"math/rand"
	"time"

	_ "image/png"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/assets"
	"github.com/dfirebaugh/hlg/pkg/input"
	"golang.org/x/image/colornames"
)

const (
	gravity     = 0.1
	damping     = 0.9
	buddyWidth  = 32
	buddyHeight = 32
)

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
	hlg.SetWindowSize(800, 600)
	hlg.SetTitle("BuddyMark Stress Test")
	hlg.EnableFPS()

	reader := bytes.NewReader(assets.BuddyDanceSpriteSheet)
	img, _, err = image.Decode(reader)
	if err != nil {
		panic(err)
	}

	lastFrameTime := time.Now()
	frameDuration := time.Millisecond * 100

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
		for _, buddy := range buddies {
			buddy.Render()
		}
		hlg.PrintAt(fmt.Sprintf("buddies: %d", len(buddies)), 20, 20, colornames.Red)
	})
}

func handleInput() {
	if hlg.IsButtonPressed(input.MouseButtonLeft) {
		x, y := hlg.GetCursorPosition()
		for i := 0; i < 5; i++ {
			buddy := NewBuddy(float32(x), float32(y))
			buddies = append(buddies, buddy)
		}
	}

	if hlg.IsButtonPressed(input.MouseButtonRight) {
		buddies = []*Buddy{}
	}
}

func NewBuddy(x, y float32) *Buddy {
	frameSize := image.Point{X: buddyWidth, Y: buddyHeight}
	sheetSize := image.Point{X: 4, Y: 1}

	sprite := hlg.NewSprite(img, frameSize, sheetSize)
	// sprite.Scale(0.09, 0.09)
	// sprite.Resize(32, 32)
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
	} else if b.X > 800-buddyWidth {
		b.X = 800 - buddyWidth
		b.VelocityX *= -damping
	}

	if b.Y > 600-buddyHeight {
		b.Y = 600 - buddyHeight
		b.VelocityY *= -damping
	}

	b.Sprite.Move(b.X, b.Y)
}

func (b *Buddy) Render() {
	b.Sprite.Render()
}
