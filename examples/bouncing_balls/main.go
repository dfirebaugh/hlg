package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/gui"
	"github.com/dfirebaugh/hlg/pkg/input"
	"golang.org/x/image/colornames"
)

var (
	balls         []Ball
	numBalls      int
	damping       = float32(0.9)
	enableFade    = true
	enableCollide = true
	sliderValue   float32
)

func main() {
	hlg.SetTitle("Bouncing Balls")
	hlg.SetWindowSize(960, 640)
	hlg.SetScreenSize(960, 640)
	hlg.EnableFPS()
	hlg.SetVSync(true)

	font, err := hlg.LoadDefaultFont()
	if err != nil {
		fmt.Printf("Failed to load font: %v\n", err)
		return
	}
	font.SetAsActiveAtlas()
	hlg.SetDefaultFont(font)

	// Create gui context
	inputCtx := gui.NewDefaultInputContext()
	ctx := gui.NewContext(inputCtx)

	numBalls = 3
	sliderValue = float32(numBalls) / 275 * 100
	balls = make([]Ball, numBalls)
	for i := 0; i < numBalls; i++ {
		balls[i] = NewBall()
	}

	// Panel state (draggable)
	panelState := gui.PanelState{X: 10, Y: 10}
	panelW, panelH := 260, 210

	hlg.Run(func() {
		inputCtx.Update()

		// Handle ball count changes
		maxBalls := 275
		newNumBalls := int(sliderValue / 100 * float32(maxBalls))
		if newNumBalls != numBalls {
			adjustBallCount(newNumBalls)
			numBalls = newNumBalls
		}

		// Update balls
		for i := range balls {
			balls[i].Update()
		}

		// Handle mouse clicks on balls (outside panel)
		if hlg.IsButtonPressed(input.MouseButtonLeft) {
			x, y := hlg.GetCursorPosition()
			// Only affect balls if clicking outside the panel
			if x < panelState.X || x > panelState.X+panelW || y < panelState.Y || y > panelState.Y+panelH {
				for i := range balls {
					if balls[i].containsPoint(float32(x), float32(y)) {
						balls[i].VelocityX += float32(rand.Float64()*2-1) * 5
						balls[i].VelocityY += float32(rand.Float64()*2-1) * 5
					}
				}
			}
		}
	}, func() {
		hlg.Clear(colornames.Black)

		// Render UI (includes drawing calls)
		ctx.Begin()

		// Render balls
		for i := range balls {
			balls[i].Render()
		}

		// Draggable panel
		if ctx.Panel("Ball Control", &panelState, panelW, panelH) {
			// Panel content (only rendered if not collapsed)
			px, py := panelState.X, panelState.Y

			// Slider
			ctx.Slider("ball_count", &sliderValue, 0, 100, px+10, py+40, 240, 12)

			// Fade toggle
			hlg.Text("Enable Color Fade", px+10, py+75, 14, colornames.White)
			ctx.Toggle("fade", &enableFade, px+200, py+70, 40, 22)

			// Collision toggle
			hlg.Text("Enable Collision", px+10, py+105, 14, colornames.White)
			ctx.Toggle("collision", &enableCollide, px+200, py+100, 40, 22)

			// Reset button
			if ctx.Button("Reset Balls", px+10, py+135, 240, 28) {
				for i := range balls {
					balls[i] = NewBall()
				}
			}

			// Ball count label
			hlg.Text(fmt.Sprintf("Number of Balls: %d", numBalls), px+10, py+175, 14, colornames.White)
		}

		ctx.End()
	})

	font.Dispose()
}

type Ball struct {
	X, Y                 float32
	VelocityX, VelocityY float32
	Radius               float32
	Color                color.RGBA
	TargetColor          color.RGBA
	ColorChangeSpeed     float32
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
}

func (b *Ball) updateColor() {
	b.Color.R = uint8(float32(b.Color.R) + (float32(b.TargetColor.R)-float32(b.Color.R))*b.ColorChangeSpeed)
	b.Color.G = uint8(float32(b.Color.G) + (float32(b.TargetColor.G)-float32(b.Color.G))*b.ColorChangeSpeed)
	b.Color.B = uint8(float32(b.Color.B) + (float32(b.TargetColor.B)-float32(b.Color.B))*b.ColorChangeSpeed)
}

func (b *Ball) Render() {
	hlg.FilledCircle(int(b.X), int(b.Y), int(b.Radius), b.Color)
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

	return Ball{
		X:                x,
		Y:                y,
		Radius:           radius,
		VelocityX:        velocityX,
		VelocityY:        velocityY,
		Color:            randomColor,
		TargetColor:      randomColor,
		ColorChangeSpeed: 0.1,
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

	nx := dx / float32(distance)
	ny := dy / float32(distance)

	tx := -ny
	ty := nx

	dpTan1 := b.VelocityX*tx + b.VelocityY*ty
	dpTan2 := other.VelocityX*tx + other.VelocityY*ty

	dpNorm1 := b.VelocityX*nx + b.VelocityY*ny
	dpNorm2 := other.VelocityX*nx + other.VelocityY*ny

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
