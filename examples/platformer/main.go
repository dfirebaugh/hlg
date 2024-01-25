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
	gravity      = 0.2
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
	platforms []*Platform
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

	p.Ground = false

	playerBottomCenterX := p.X + p.W/2
	playerBottomCenterY := p.Y + p.H

	for _, pl := range p.platforms {
		if playerBottomCenterX > pl.X && playerBottomCenterX < (pl.X+pl.W) {
			if playerBottomCenterY >= pl.Y && playerBottomCenterY <= (pl.Y+pl.H) && p.VelY >= 0 {
				p.Y = pl.Y - p.H
				p.VelY = 0
				p.Ground = true
				break
			}
		}
	}
}

func (p *Player) Render() {
	p.Sprite.Render()
}

func main() {
	hlg.SetWindowSize(windowWidth, windowHeight)
	hlg.SetScreenSize(windowWidth, windowHeight)

	reader := bytes.NewReader(assets.BuddyDanceSpriteSheet)
	img, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}
	frameSize := image.Point{X: 32, Y: 32}
	sheetSize := image.Point{X: 4, Y: 1}
	sprite := hlg.NewSprite(img, frameSize, sheetSize)

	platforms := []*Platform{
		NewPlatform(200, 400, 100, 20),
		NewPlatform(400, 300, 150, 20),
	}

	player := &Player{
		X:         100,
		Y:         float64(windowHeight) - 100,
		W:         16,
		H:         16,
		Sprite:    sprite,
		LastFrame: time.Now(),
		platforms: platforms,
	}
	sprite.Resize(float32(player.W), float32(player.H))

	hlg.Update(func() {
		hlg.Clear(colornames.White)

		player.Update()
		player.Render()

		for _, pl := range platforms {
			pl.Render()
		}

		hlg.PrintAt(fmt.Sprintf("Player X: %d Y: %d", int(player.X), int(player.Y)),
			10, windowHeight-20, colornames.Black)
	})
}

type Platform struct {
	X float64
	Y float64
	W float64
	H float64
	hlg.Shape
}

func NewPlatform(x, y, w, h float64) *Platform {
	shape := hlg.Rectangle(int(x), int(y), int(w), int(h), colornames.Royalblue)
	return &Platform{X: x, Y: y, W: w, H: h, Shape: shape}
}

func (pl *Platform) Render() {
	pl.Shape.Render()
}
