package main

import (
	"image/color"
	"math/rand"

	"github.com/dfirebaugh/ggez"
	"github.com/dfirebaugh/ggez/pkg/draw"
	"github.com/dfirebaugh/ggez/pkg/math/geom"
	"golang.org/x/image/colornames"
)

var balls []Ball

func main() {
	ggez.SetTitle("cpu rendering")
	ggez.SetWindowSize(960, 640)
	screen := NewScreen(960, 640)

	numBalls := 300
	balls = make([]Ball, numBalls)
	for i := 0; i < numBalls; i++ {
		balls[i] = NewBall(screen)
	}

	ggez.Update(func() {
		screen.Clear(colornames.Grey)
		for i := range balls {
			balls[i].Update()
			balls[i].Render()
		}

		screen.Render()
	})
}

type Ball struct {
	screen Screen
	geom.Circle
	Velocity geom.Point
	Color    color.RGBA
}

func (b *Ball) Render() {
	draw.Circle(geom.MakeCircle(b.X, b.Y, b.R)).Fill(b.screen, b.Color)
}

func (b *Ball) Update() {
	// Update the position
	b.X += b.Velocity.X
	b.Y += b.Velocity.Y

	sw, sh := ggez.ScreenSize()

	// Boundary checks - Bounce off the edges
	if b.X < 0 || b.X > float32(sw) {
		b.Velocity.X = -b.Velocity.X
	}

	if b.Y < 0 || b.Y > float32(sh) {
		b.Velocity.Y = -b.Velocity.Y
	}
}

func NewBall(screen Screen) Ball {
	sw, sh := ggez.ScreenSize()
	return Ball{
		screen: screen,
		Circle: geom.Circle{
			X: float32(rand.Float64() * float64(sw)), Y: float32(rand.Float64() * float64(sh)),
			R: float32(rand.Intn(35) + 10),
		},
		Velocity: geom.Point{X: float32(rand.Float64()*4 - 2), Y: float32(rand.Float64()*4 - 2)},
		Color:    color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 255},
	}
}
