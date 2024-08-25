package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/pkg/input"
	"github.com/dfirebaugh/hlg/pkg/ui"
	"github.com/dfirebaugh/hlg/pkg/ui/components"
	"golang.org/x/image/colornames"
)

var (
	balls          []Ball
	numBalls       int
	damping        = float32(0.9) // Damping factor to simulate energy loss during collisions
	enableFade     = true
	enableCollide  = true
	slider         *components.Slider
	fadeToggle     *components.Toggle
	collideToggle  *components.Toggle
	resetButton    *components.Button
	ballCountLabel *components.Label
)

func main() {
	hlg.SetTitle("Bouncing Balls")
	hlg.SetWindowSize(960, 640)
	hlg.EnableFPS()

	surface := setupUI()

	numBalls = 3
	balls = make([]Ball, numBalls)
	for i := 0; i < numBalls; i++ {
		balls[i] = NewBall()
	}

	hlg.Run(func() {
		handleInput()
		for i := range balls {
			balls[i].Update()
		}
		surface.Update()
	}, func() {
		hlg.Clear(colornames.Black)
		for i := range balls {
			balls[i].Render()
		}
		surface.Render()
	})
}

type Ball struct {
	X, Y                 float32
	VelocityX, VelocityY float32
	Radius               float32
	Color                color.RGBA
	TargetColor          color.RGBA
	ColorChangeSpeed     float32
	Shape                hlg.Shape
}

func (b *Ball) Update() {
	b.X += b.VelocityX
	b.Y += b.VelocityY

	b.handleBorderCollision()

	if enableCollide {
		for i := range balls {
			if &balls[i] != b {
				b.checkCollision(&balls[i])
			}
		}
	}

	if enableFade {
		b.updateColor()
	}

	b.Shape.Move(b.X, b.Y)
}

func (b *Ball) updateColor() {
	b.Color.R = uint8(float32(b.Color.R) + (float32(b.TargetColor.R)-float32(b.Color.R))*b.ColorChangeSpeed)
	b.Color.G = uint8(float32(b.Color.G) + (float32(b.TargetColor.G)-float32(b.Color.G))*b.ColorChangeSpeed)
	b.Color.B = uint8(float32(b.Color.B) + (float32(b.TargetColor.B)-float32(b.Color.B))*b.ColorChangeSpeed)

	b.Shape.SetColor(b.Color)
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
		X:                x,
		Y:                y,
		Radius:           radius,
		VelocityX:        velocityX,
		VelocityY:        velocityY,
		Color:            randomColor,
		TargetColor:      randomColor,
		ColorChangeSpeed: 0.1,
		Shape:            shape,
	}
}

func handleInput() {
	maxBalls := 275
	newNumBalls := int(slider.Value * float64(maxBalls))

	if newNumBalls != numBalls {
		adjustBallCount(newNumBalls)
		numBalls = newNumBalls
		if ballCountLabel != nil {
			ballCountLabel.Text = fmt.Sprintf("Number of Balls: %d", numBalls)
		}
	}

	if fadeToggle != nil {
		enableFade = fadeToggle.IsOn
	}

	if collideToggle != nil {
		enableCollide = collideToggle.IsOn
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

	if resetButton != nil && resetButton.IsPressed {
		for i := range balls {
			balls[i] = NewBall()
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

	overlap := 0.5 * (float32(distance) - b.Radius - other.Radius)
	b.X -= overlap * nx
	b.Y -= overlap * ny
	other.X += overlap * nx
	other.Y += overlap * ny

	if enableFade {
		b.TargetColor = randomColor()
		other.TargetColor = randomColor()
	}
}

func randomColor() color.RGBA {
	return color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}
}

func adjustBallCount(targetCount int) {
	currentCount := len(balls)

	if currentCount < targetCount {
		for i := 0; i < targetCount-currentCount; i++ {
			balls = append(balls, NewBall())
		}
	} else if currentCount > targetCount {
		balls = balls[:targetCount]
	}
}

func setupUI() *ui.Surface {
	surface := ui.NewSurface(260, 200)
	surface.Move(10, 10)

	surface.Add(&components.SurfaceHandle{
		Width:  260,
		Height: 30,
		Label:  "Ball Control",
	})

	slider = &components.Slider{
		X:      10,
		Y:      40,
		Width:  240,
		Height: 20,
		Value:  float64(numBalls) / 275,
	}
	surface.Add(slider)

	fadeToggle = &components.Toggle{
		X:      200,
		Y:      70,
		Width:  30,
		Height: 20,
		IsOn:   true,
	}
	surface.Add(fadeToggle)

	surface.Add(&components.Label{
		X:    10,
		Y:    70,
		Text: "Enable Color Fade",
	})

	collideToggle = &components.Toggle{
		X:      200,
		Y:      100,
		Width:  30,
		Height: 20,
		IsOn:   true,
	}
	surface.Add(collideToggle)

	surface.Add(&components.Label{
		X:    10,
		Y:    100,
		Text: "Enable Collision",
	})

	resetButton = &components.Button{
		X:      10,
		Y:      130,
		Width:  240,
		Height: 20,
		Text:   "Reset Balls",
	}
	surface.Add(resetButton)

	ballCountLabel = &components.Label{
		X:    10,
		Y:    160,
		Text: fmt.Sprintf("Number of Balls: %d", numBalls),
	}
	surface.Add(ballCountLabel)

	return surface
}
