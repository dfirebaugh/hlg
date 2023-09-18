package main

import (
	"math"

	"github.com/dfirebaugh/ggez"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

var (
	camera Camera
)

type Point3D struct {
	X, Y, Z float64
}

type Cube struct {
	Vertices [8]Point3D
	Angle    float64
}

var edges = [12][2]int{
	{0, 1}, {1, 3}, {3, 2}, {2, 0},
	{4, 5}, {5, 7}, {7, 6}, {6, 4},
	{0, 4}, {1, 5}, {2, 6}, {3, 7},
}

type Camera struct {
	X, Y, Z    float64
	Yaw, Pitch float64
}

func (c *Cube) Render(centerX, centerY float64) {
	for _, edge := range edges {
		v1 := c.Vertices[edge[0]]
		v2 := c.Vertices[edge[1]]

		x1 := v1.X*math.Cos(c.Angle) - v1.Z*math.Sin(c.Angle)
		z1 := v1.X*math.Sin(c.Angle) + v1.Z*math.Cos(c.Angle)
		x2 := v2.X*math.Cos(c.Angle) - v2.Z*math.Sin(c.Angle)
		z2 := v2.X*math.Sin(c.Angle) + v2.Z*math.Cos(c.Angle)

		// Apply camera orientation (yaw and pitch)
		x1Rotated := x1*math.Cos(camera.Yaw) - z1*math.Sin(camera.Yaw)
		z1Rotated := x1*math.Sin(camera.Yaw) + z1*math.Cos(camera.Yaw)
		x2Rotated := x2*math.Cos(camera.Yaw) - z2*math.Sin(camera.Yaw)
		z2Rotated := x2*math.Sin(camera.Yaw) + z2*math.Cos(camera.Yaw)

		y1 := v1.Y*math.Cos(camera.Pitch) - z1Rotated*math.Sin(camera.Pitch)
		z1 = v1.Y*math.Sin(camera.Pitch) + z1Rotated*math.Cos(camera.Pitch)
		y2 := v2.Y*math.Cos(camera.Pitch) - z2Rotated*math.Sin(camera.Pitch)
		z2 = v2.Y*math.Sin(camera.Pitch) + z2Rotated*math.Cos(camera.Pitch)

		scale1 := 200.0 / (200.0 + z1)
		scale2 := 200.0 / (200.0 + z2)
		x1Proj := x1Rotated * scale1
		y1Proj := y1 * scale1
		x2Proj := x2Rotated * scale2
		y2Proj := y2 * scale2

		x1Screen := x1Proj*100 + centerX
		y1Screen := y1Proj*100 + centerY
		x2Screen := x2Proj*100 + centerX
		y2Screen := y2Proj*100 + centerY

		ggez.DrawLine(int(x1Screen), int(y1Screen), int(x2Screen), int(y2Screen), colornames.Bisque)
	}
}

func (c *Cube) Update() {
	c.Angle += 0.01
}

func main() {
	ggez.SetTitle("3D Cube")
	ggez.SetScreenSize(screenWidth, screenHeight)

	centerX := float64(screenWidth) / 2
	centerY := float64(screenHeight) / 2

	cubeSize := 2.0

	cube := Cube{
		Vertices: [8]Point3D{
			{cubeSize / 2, cubeSize / 2, cubeSize / 2},
			{cubeSize / 2, cubeSize / 2, -cubeSize / 2},
			{cubeSize / 2, -cubeSize / 2, cubeSize / 2},
			{cubeSize / 2, -cubeSize / 2, -cubeSize / 2},
			{-cubeSize / 2, cubeSize / 2, cubeSize / 2},
			{-cubeSize / 2, cubeSize / 2, -cubeSize / 2},
			{-cubeSize / 2, -cubeSize / 2, cubeSize / 2},
			{-cubeSize / 2, -cubeSize / 2, -cubeSize / 2},
		},
	}

	camera = Camera{
		X:     0,
		Y:     0,
		Z:     5,
		Yaw:   0.3,
		Pitch: 0.2,
	}

	ggez.Update(func() {
		ggez.Clear(colornames.Black)

		cube.Render(centerX, centerY)
		cube.Update()
	})
}
