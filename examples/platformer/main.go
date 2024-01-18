package main

import (
	"bytes"
	"fmt"
	"image"
	"time"

	_ "image/png"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/assets"
	"github.com/dfirebaugh/hlg/pkg/input"
	"golang.org/x/image/colornames"
)

const (
	windowWidth  = 800
	windowHeight = 600
	playerSpeed  = 5
	gravity      = 0.5
	jumpSpeed    = 10
)

type Player struct {
	X         float64
	Y         float64
	W         float64
	H         float64
	VelY      float64
	Ground    bool
	Sprite    *hlg.Sprite
	LastFrame time.Time
}

func (p *Player) Update() {
	// Basic gravity
	p.VelY += gravity
	p.Y += p.VelY

	// Collision with ground
	if p.Y > float64(windowHeight)-p.H {
		p.Y = float64(windowHeight) - p.H
		p.VelY = 0
		p.Ground = true
	}

	// Movement
	if hlg.IsKeyPressed(input.KeyA) || hlg.IsKeyPressed(input.KeyLeft) {
		p.X -= playerSpeed
	}
	if hlg.IsKeyPressed(input.KeyD) || hlg.IsKeyPressed(input.KeyRight) {
		p.X += playerSpeed
	}

	// Jumping
	if hlg.IsKeyPressed(input.KeySpace) && p.Ground {
		p.VelY = -jumpSpeed
		p.Ground = false
	}

	// Update sprite frame
	if time.Since(p.LastFrame) >= time.Millisecond*200 {
		p.LastFrame = time.Now()
		p.Sprite.NextFrame()
	}

	p.Sprite.Move(float32(p.X), float32(p.Y))
}

func (p *Player) Render() {
	p.Sprite.Render()
}

func main() {
	hlg.SetWindowSize(windowWidth, windowHeight)
	hlg.SetScreenSize(windowWidth, windowHeight)

	// Load sprite sheet (replace assets.BuddyDanceSpriteSheet with your sprite sheet)
	reader := bytes.NewReader(assets.BuddyDanceSpriteSheet)
	img, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}
	frameSize := image.Point{X: 32, Y: 32}
	sheetSize := image.Point{X: 4, Y: 1}
	sprite := hlg.NewSprite(img, frameSize, sheetSize)

	player := &Player{
		X:         100,
		Y:         float64(windowHeight) - 100,
		W:         50,
		H:         50,
		Sprite:    sprite,
		LastFrame: time.Now(),
	}
	sprite.Resize(float32(player.W), float32(player.H))

	hlg.Update(func() {
		hlg.Clear(colornames.White)

		player.Update()
		player.Render()

		// Display player coordinates
		hlg.PrintAt(fmt.Sprintf("Player X: %d Y: %d", int(player.X), int(player.Y)),
			10, windowHeight-20, colornames.Black)
	})
}
