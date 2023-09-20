package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/dfirebaugh/ggez"
	"github.com/dfirebaugh/ggez/pkg/draw"
	"github.com/dfirebaugh/ggez/pkg/fb"
	"github.com/dfirebaugh/ggez/pkg/math/geom"
	"golang.org/x/image/colornames"
)

type Ball struct {
	geom.Circle
	Velocity geom.Point
	Color    color.RGBA
	Texture  *ggez.Texture
}

const (
	screenWidth  = 800
	screenHeight = 600
)

var (
	balls        []Ball
	numBalls     = 500
	circleRadius = 4
)

func (b *Ball) Render() {
	b.Texture.X = int(b.X - b.R)
	b.Texture.Y = int(b.Y - b.R)
	b.Texture.Render()
}

func (b *Ball) Update() {
	b.X += b.Velocity.X
	b.Y += b.Velocity.Y

	if b.X < b.R || b.X > (screenWidth-b.R) {
		b.Velocity.X = -b.Velocity.X
	}

	if b.Y < b.R || b.Y > screenHeight-b.R {
		b.Velocity.Y = -b.Velocity.Y
	}
}

func NewBall() Ball {
	ballColor := color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}

	padding := 2
	f := fb.New((circleRadius*2)+padding, (circleRadius*2)+padding)
	draw.Fill(draw.Circle{X: float64(circleRadius), Y: float64(circleRadius), R: float64(circleRadius)}, f, ballColor)
	texture, err := ggez.CreateTextureFromImage(f.ToImage())
	if err != nil {
		panic(err)
	}

	circle := geom.Circle{
		X: (rand.Float64()*screenWidth - float64(circleRadius)),
		Y: (rand.Float64()*screenHeight - float64(circleRadius)),
		R: float64(circleRadius),
	}

	if circle.X <= circle.R {
		circle.X += circle.R * 2
	}
	if circle.Y <= circle.R {
		circle.Y += circle.R * 2
	}
	return Ball{
		Circle: circle,
		Velocity: geom.Point{
			X: rand.Float64()*4 - 2,
			Y: rand.Float64()*4 - 2,
		},
		Color:   ballColor,
		Texture: texture,
	}
}

func update() {
	ggez.Clear(colornames.Grey)

	for i := range balls {
		balls[i].Update()
		balls[i].Render()
	}

	ggez.PrintAt(fmt.Sprintf("ball count: %d", len(balls)), 100, 100, colornames.Red)
}

func main() {
	balls = make([]Ball, numBalls)
	defer func() {
		// Clean up textures when done
		for _, b := range balls {
			b.Texture.Destroy()
		}
	}()

	for i := 0; i < numBalls; i++ {
		balls[i] = NewBall()
	}

	ggez.SetTitle("ball demo")
	// ggez.DisableFPS()
	ggez.Update(update)
}
