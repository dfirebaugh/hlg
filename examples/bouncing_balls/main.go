package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/pkg/input"
	"golang.org/x/image/colornames"
)

var (
	balls    []Ball
	numBalls int
	damping  = float32(0.9) // Damping factor to simulate energy loss during collisions
)

func main() {
	hlg.SetTitle("Bouncing Balls")
	hlg.SetWindowSize(960, 640)

	numBalls = 10
	balls = make([]Ball, numBalls)
	for i := 0; i < numBalls; i++ {
		balls[i] = NewBall()
	}

	hlg.Update(func() {
		hlg.Clear(colornames.Aliceblue)
		handleInput()
		for i := range balls {
			balls[i].Update()
			balls[i].Render()
		}
	})
}

type Ball struct {
	X, Y, Radius         float32
	VelocityX, VelocityY float32
	Color                color.RGBA
	Shape                hlg.Shape
}

func (b *Ball) Update() {
	b.X += b.VelocityX
	b.Y += b.VelocityY

	b.handleBorderCollision()

	for i := range balls {
		if &balls[i] != b {
			b.checkCollision(&balls[i])
		}
	}
	b.Shape.Move(b.X, b.Y)
}

func (b *Ball) Render() {
	b.Shape.Render()
}

func NewBall() Ball {
	sw, sh := hlg.GetScreenSize()
	radius := float32(rand.Intn(35) + 10)
	x := radius + float32(rand.Float64()*(float64(sw)-2*float64(radius)))
	y := radius + float32(rand.Float64()*(float64(sh)-2*float64(radius)))
	velocityX := float32(rand.Float64()*4 - 2)
	velocityY := float32(rand.Float64()*4 - 2)

	randomColor := color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}

	shape := hlg.Circle(int(x), int(y), radius, randomColor)

	return Ball{
		X:         x,
		Y:         y,
		Radius:    radius,
		VelocityX: velocityX,
		VelocityY: velocityY,
		Color:     randomColor,
		Shape:     shape,
	}
}

func handleInput() {
	if hlg.IsKeyPressed(input.KeyUp) {
		numBalls++
		balls = append(balls, NewBall())
	}
	if hlg.IsKeyPressed(input.KeyDown) {
		if numBalls > 0 {
			numBalls--
			balls[numBalls].Shape.Dispose()
			balls = balls[:numBalls]
		}
	}
	if hlg.IsKeyPressed(input.KeySpace) {
		for i := range balls {
			balls[i].VelocityX = float32(rand.Float64()*4 - 2)
			balls[i].VelocityY = float32(rand.Float64()*4 - 2)
		}
	}

	if hlg.IsButtonPressed(input.MouseButtonLeft) {
		x, y := hlg.GetCursorPosition()
		for i := range balls {
			if balls[i].containsPoint(float32(x), float32(y)) {
				balls[i].VelocityX += float32(rand.Float64()*2-1) * 5
				balls[i].VelocityY += float32(rand.Float64()*2-1) * 5
			}
		}
	}
}

func (b *Ball) containsPoint(px, py float32) bool {
	dx := b.X - px
	dy := b.Y - py
	distance := math.Sqrt(float64(dx*dx + dy*dy))
	return distance < float64(b.Radius)
}

func (b *Ball) handleBorderCollision() {
	sw, sh := hlg.GetScreenSize()
	if b.X-b.Radius < 0 {
		b.X = b.Radius
		b.VelocityX = -b.VelocityX * damping
	}
	if b.X+b.Radius > float32(sw) {
		b.X = float32(sw) - b.Radius
		b.VelocityX = -b.VelocityX * damping
	}
	if b.Y-b.Radius < 0 {
		b.Y = b.Radius
		b.VelocityY = -b.VelocityY * damping
	}
	if b.Y+b.Radius > float32(sh) {
		b.Y = float32(sh) - b.Radius
		b.VelocityY = -b.VelocityY * damping
	}
}

func (b *Ball) checkCollision(other *Ball) {
	dx := b.X - other.X
	dy := b.Y - other.Y
	distance := math.Sqrt(float64(dx*dx + dy*dy))

	if distance < float64(b.Radius+other.Radius) {
		b.resolveCollision(other)
	}
}

func (b *Ball) resolveCollision(other *Ball) {
	dx := b.X - other.X
	dy := b.Y - other.Y
	distance := math.Sqrt(float64(dx*dx + dy*dy))

	// Normalize the distance
	nx := dx / float32(distance)
	ny := dy / float32(distance)

	// Tangent vector
	tx := -ny
	ty := nx

	// Dot product tangent
	dpTan1 := b.VelocityX*tx + b.VelocityY*ty
	dpTan2 := other.VelocityX*tx + other.VelocityY*ty

	// Dot product normal
	dpNorm1 := b.VelocityX*nx + b.VelocityY*ny
	dpNorm2 := other.VelocityX*nx + other.VelocityY*ny

	// Conservation of momentum in 1D
	m1 := (dpNorm1*(b.Radius-other.Radius) + 2.0*other.Radius*dpNorm2) / (b.Radius + other.Radius)
	m2 := (dpNorm2*(other.Radius-b.Radius) + 2.0*b.Radius*dpNorm1) / (b.Radius + other.Radius)

	b.VelocityX = (tx*dpTan1 + nx*m1) * damping
	b.VelocityY = (ty*dpTan1 + ny*m1) * damping
	other.VelocityX = (tx*dpTan2 + nx*m2) * damping
	other.VelocityY = (ty*dpTan2 + ny*m2) * damping

	// Separate the balls to avoid overlap
	overlap := 0.5 * (float32(distance) - b.Radius - other.Radius)
	b.X -= overlap * nx
	b.Y -= overlap * ny
	other.X += overlap * nx
	other.Y += overlap * ny
}
