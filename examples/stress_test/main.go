package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/pkg/input"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

type Triangle struct {
	PositionX, PositionY int
	VelocityX, VelocityY float32
	Size                 int
	Color                color.Color
}

var (
	triangles []Triangle
	frames    int
)

func init() {
	triangles = []Triangle{
		newRandomTriangle(),
		newRandomTriangle(),
	}
}

func newRandomTriangle() Triangle {
	return Triangle{
		PositionX: rand.Intn(screenWidth),
		PositionY: rand.Intn(screenHeight),
		VelocityX: randomVelocity(),
		VelocityY: randomVelocity(),
		Size:      rand.Intn(20) + 10,
		Color:     randomColor(),
	}
}

func randomColor() color.Color {
	return color.RGBA{
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		255,
	}
}

func randomVelocity() float32 {
	return (rand.Float32()*2 - 1) * 0.2
}

func update() {
	for i := range triangles {
		triangles[i].PositionX += int(triangles[i].VelocityX * float32(triangles[i].Size))
		triangles[i].PositionY += int(triangles[i].VelocityY * float32(triangles[i].Size))

		if triangles[i].PositionX < 0 || triangles[i].PositionX+triangles[i].Size > screenWidth {
			triangles[i].VelocityX = -triangles[i].VelocityX
		}
		if triangles[i].PositionY < 0 || triangles[i].PositionY+triangles[i].Size > screenHeight {
			triangles[i].VelocityY = -triangles[i].VelocityY
		}
	}

	if hlg.IsKeyJustPressed(input.KeyUp) {
		for i := 0; i < 2000; i++ {
			triangles = append(triangles, newRandomTriangle())
		}
	}

	if hlg.IsKeyJustPressed(input.KeyDown) && len(triangles) > 0 {
		triangles = triangles[:len(triangles)-1]
	}

	frames++
	if frames%60 == 0 {
		fmt.Printf("Rendering %d triangles\n", len(triangles))
	}
}

func render() {
	hlg.Clear(colornames.Skyblue)
	hlg.BeginDraw()

	for _, triangle := range triangles {
		hlg.FilledTriangle(
			triangle.PositionX,
			triangle.PositionY,
			triangle.PositionX+triangle.Size/2,
			triangle.PositionY-triangle.Size,
			triangle.PositionX+triangle.Size,
			triangle.PositionY,
			triangle.Color,
		)
	}

	hlg.EndDraw()
}

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetTitle("Bouncing Triangles")
	hlg.EnableFPS()
	hlg.Run(update, render)
}
