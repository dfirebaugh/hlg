package main

import (
	"image/color"
	"math/rand"

	"github.com/dfirebaugh/ggez"
	"github.com/dfirebaugh/ggez/graphics"
	"github.com/dfirebaugh/ggez/pkg/math/geom"
	"golang.org/x/image/colornames"
)

var balls []Ball

func main() {
	ggez.SetTitle("gpu rendering")
	ggez.SetWindowSize(960, 640)

	numBalls := 400
	balls = make([]Ball, numBalls)
	for i := 0; i < numBalls; i++ {
		balls[i] = NewBall()
	}

	ggez.Update(func() {
		ggez.Clear(colornames.Aliceblue)
		for i := range balls {
			balls[i].Render()
			balls[i].Update()
		}
	})
}

type Ball struct {
	geom.Circle
	Velocity geom.Point
	Color    color.RGBA
	graphics.Shape
}

func (b *Ball) Update() {
	b.X += b.Velocity.X
	b.Y += b.Velocity.Y

	sw, sh := ggez.ScreenSize()
	r := float32(b.Circle.R)

	if b.X-r < 0 || b.X+r > float32(sw) {
		b.Velocity.X = -b.Velocity.X
	}

	if b.Y-r < 0 || b.Y+r > float32(sh) {
		b.Velocity.Y = -b.Velocity.Y
	}
	b.Move(b.Circle.X, b.Circle.Y)
}

func NewBall() Ball {
	sw, sh := ggez.ScreenSize()

	radius := float32(rand.Intn(35) + 10)
	x := radius + float32(rand.Float64()*(float64(sw)-2*float64(radius)))
	y := radius + float32(rand.Float64()*(float64(sh)-2*float64(radius)))

	b := Ball{
		Circle: geom.Circle{
			X: x, Y: y, R: radius,
		},
		Velocity: geom.Point{X: float32(rand.Float64()*4 - 2), Y: float32(rand.Float64()*4 - 2)},
		Color:    color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 255},
	}
	b.Shape = ggez.Circle(int(b.Circle.X), int(b.Circle.Y), b.Circle.R, b.Color)

	return b
}
