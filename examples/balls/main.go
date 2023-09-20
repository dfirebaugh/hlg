package main

import (
	"image/color"
	"math/rand"

	"github.com/dfirebaugh/ggez"
	"github.com/dfirebaugh/ggez/pkg/math/geom"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

var balls []Ball

func update() {
	ggez.Clear(colornames.Grey)

	for i := range balls {
		balls[i].Update(screenWidth, screenHeight)
		balls[i].Render()
	}
}

func main() {
	ggez.Setup(ggez.GLRenderer)

	numBalls := 100
	balls = make([]Ball, numBalls)
	for i := 0; i < numBalls; i++ {
		balls[i] = NewBall(screenWidth, screenHeight)
	}

	ggez.Update(update)
}

type Ball struct {
	geom.Circle
	Velocity geom.Point
	Color    color.RGBA
}

func (b *Ball) Render() {
	ggez.FillCircle(int(b.X), int(b.Y), int(b.R), b.Color)
}

func (b *Ball) Update(screenWidth, screenHeight float64) {
	// Update the position
	b.X += b.Velocity.X
	b.Y += b.Velocity.Y

	// Boundary checks - Bounce off the edges
	if b.X < 0 || b.X > screenWidth {
		b.Velocity.X = -b.Velocity.X
	}

	if b.Y < 0 || b.Y > screenHeight {
		b.Velocity.Y = -b.Velocity.Y
	}
}

func NewBall(screenWidth, screenHeight float64) Ball {
	return Ball{
		Circle: geom.Circle{
			X: rand.Float64() * screenWidth, Y: rand.Float64() * screenHeight,
			R: 50,
		},
		Velocity: geom.Point{X: rand.Float64()*4 - 2, Y: rand.Float64()*4 - 2}, // random velocity between -2 to 2
		Color:    color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 255},
	}
}
