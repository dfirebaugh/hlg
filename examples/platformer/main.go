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
	windowWidth        = 800
	windowHeight       = 600
	playerSpeed        = 5
	gravity            = 0.2
	jumpSpeed          = 10
	coyoteTimeDuration = 400 // milliseconds
	debug              = false
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

	CoyoteTimeLeft int
	hlg.Shape
}

func (p *Player) handleMovement() {
	if hlg.IsKeyPressed(input.KeyA) || hlg.IsKeyPressed(input.KeyLeft) {
		p.X -= playerSpeed
	}
	if hlg.IsKeyPressed(input.KeyD) || hlg.IsKeyPressed(input.KeyRight) {
		p.X += playerSpeed
	}
}

func (p *Player) handleCoyoteTime() {
	if p.Ground {
		p.CoyoteTimeLeft = coyoteTimeDuration
	} else {
		p.CoyoteTimeLeft -= 17
		if p.CoyoteTimeLeft < 0 {
			p.CoyoteTimeLeft = 0
		}
	}
}

func (p *Player) handlePlatformCollision() {
	p.Ground = false

	playerBottomCenterX := p.X + p.W/2
	playerBottomCenterY := p.Y + p.H - 8

	for _, pl := range p.platforms {
		if playerBottomCenterX > pl.X && playerBottomCenterX < (pl.X+pl.W) {
			if playerBottomCenterY >= pl.Y && playerBottomCenterY <= (pl.Y+pl.H) && p.VelY >= 0 {
				p.Y = pl.Y - p.H + (p.H / 2)
				p.VelY = 0
				p.Ground = true
				break
			}
		}
	}
}

func (p *Player) handleGroundCollision() {
	if p.Y > float64(windowHeight)-p.H+(p.H/2) {
		p.Y = float64(windowHeight) - p.H + (p.H / 2)
		p.VelY = 0
		p.Ground = true
	}
}

func (p *Player) updateVelocity() {
	if p.CoyoteTimeLeft <= 0 {
		p.VelY += gravity
	}
	p.Y += p.VelY
}

func (p *Player) updateSpriteFrame() {
	if time.Since(p.LastFrame) >= time.Millisecond*200 {
		p.LastFrame = time.Now()
		p.Sprite.NextFrame()
	}
}

func (p *Player) handleJump() {
	if hlg.IsKeyPressed(input.KeySpace) && (p.Ground || p.CoyoteTimeLeft > 0) {
		p.VelY = -jumpSpeed
		p.Ground = false
		p.CoyoteTimeLeft = 0
	}
}

func (p *Player) Update() {
	p.updateVelocity()
	p.handleGroundCollision()
	p.handleMovement()
	p.handleCoyoteTime()

	p.handleJump()
	p.updateSpriteFrame()

	p.Sprite.Move(float32(p.X), float32(p.Y))
	p.handlePlatformCollision()
	p.Shape.Move(float32(p.X), float32(p.Y))
}

func (p *Player) Render() {
	p.Sprite.Render()

	if debug {
		p.Shape.Render()
	}
}

func main() {
	hlg.SetWindowSize(windowWidth, windowHeight)
	hlg.SetScreenSize(windowWidth, windowHeight)
	if debug {
		hlg.EnableFPS()
	}

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
		W:         32,
		H:         32,
		Sprite:    sprite,
		LastFrame: time.Now(),
		platforms: platforms,
	}
	player.Shape = hlg.Rectangle(int(player.X), int(player.Y), int(player.W), int(player.H), colornames.Mediumpurple)
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
		hlg.PrintAt(fmt.Sprintf("VelY: %.2f, Ground: %t, CoyoteTimeLeft: %d", player.VelY, player.Ground, player.CoyoteTimeLeft), 10, windowHeight-40, colornames.Black)
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
