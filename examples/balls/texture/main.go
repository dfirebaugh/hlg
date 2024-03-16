package main

import (
	"image/color"
	"math/rand"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/pkg/draw"
	"github.com/dfirebaugh/hlg/pkg/fb"
	"github.com/dfirebaugh/hlg/pkg/math/geom"
	"golang.org/x/image/colornames"
)

var balls []Ball

func update() {
	hlg.Clear(colornames.Grey)

	for i := range balls {
		balls[i].Update()
		balls[i].Render()
	}
}

func main() {
	hlg.SetWindowSize(960, 640)
	hlg.SetTitle("texture rendering")

	numBalls := 900
	balls = make([]Ball, numBalls)
	for i := 0; i < numBalls; i++ {
		balls[i] = NewBall()
	}

	hlg.Update(update)
}

type Ball struct {
	geom.Circle
	Velocity geom.Point
	Color    color.RGBA
	*hlg.Texture
}

func (b *Ball) Render() {
	b.Texture.Render()
}

func (b *Ball) Update() {
	b.X += b.Velocity.X
	b.Y += b.Velocity.Y

	sw, sh := hlg.GetScreenSize()
	r := float32(b.Circle.R)

	if b.X-r < 0 || b.X+r > float32(sw) {
		b.Velocity.X = -b.Velocity.X
	}

	if b.Y-r < 0 || b.Y+r > float32(sh) {
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
	texture, err := hlg.CreateTextureFromImage(f.ToImage())
	texture.Resize(float32(circleRadius/4), float32(circleRadius/4))
	if err != nil {
		panic(err)
	}

	sw, sh := hlg.GetScreenSize()

	radius := float32(rand.Intn(35) + 10)
	x := radius + float32(rand.Float64()*(float64(sw)-2*float64(radius)))
	y := radius + float32(rand.Float64()*(float64(sh)-2*float64(radius)))

	return Ball{
		Circle: geom.Circle{
			X: x, Y: y, R: radius,
		},
		Velocity: geom.Point{X: float32(rand.Float64()*4 - 2), Y: float32(rand.Float64()*4 - 2)},
		Texture:  texture,
	}
}
