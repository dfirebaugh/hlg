package main

import (
	"image/color"
	"math/rand"

	"github.com/dfirebaugh/ggez"
	"github.com/dfirebaugh/ggez/pkg/draw"
	"github.com/dfirebaugh/ggez/pkg/fb"
	"github.com/dfirebaugh/ggez/pkg/math/geom"
	"golang.org/x/image/colornames"
)

var balls []Ball

func update() {
	ggez.Clear(colornames.Grey)

	for i := range balls {
		balls[i].Update()
		balls[i].Render()
	}
}

func main() {
	ggez.SetWindowSize(960, 640)
	ggez.SetTitle("texture rendering")

	numBalls := 900
	balls = make([]Ball, numBalls)
	for i := 0; i < numBalls; i++ {
		balls[i] = NewBall()
	}

	ggez.Update(update)
}

type Ball struct {
	geom.Circle
	Velocity geom.Point
	Color    color.RGBA
	*ggez.Texture
}

func (b *Ball) Render() {
	b.Texture.Render()
}

func (b *Ball) Update() {
	b.X += b.Velocity.X
	b.Y += b.Velocity.Y

	sw, sh := ggez.ScreenSize()

	if b.X < 0 || b.X > float32(sw) {
		b.Velocity.X = -b.Velocity.X
	}

	if b.Y < 0 || b.Y > float32(sh) {
		b.Velocity.Y = -b.Velocity.Y
	}
	b.Texture.Move(b.X, b.Y)
}

func NewBall() Ball {
	circleRadius := 80
	padding := 0
	ballColor := color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 255}
	f := fb.New((circleRadius*2)+padding, (circleRadius*2)+padding)
	draw.Circle(geom.MakeCircle(float32(circleRadius), float32(circleRadius), float32(rand.Intn(circleRadius)+10))).Fill(f, ballColor)
	texture, err := ggez.CreateTextureFromImage(f.ToImage())
	texture.Resize(float32(circleRadius/4), float32(circleRadius/4))
	if err != nil {
		panic(err)
	}

	sw, sh := ggez.ScreenSize()

	return Ball{
		Circle: geom.Circle{
			X: float32(rand.Float64() * float64(sw)), Y: float32(rand.Float64() * float64(sh)),
			R: 25,
		},
		Velocity: geom.Point{X: float32(rand.Float64()*4 - 2), Y: float32(rand.Float64()*4 - 2)},
		// Color:    color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 255},
		Texture: texture,
	}
}
