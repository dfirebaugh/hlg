package main

import (
	"image/color"
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

func (p Point3D) subtract(other Point3D) Point3D {
	return Point3D{
		X: p.X - other.X,
		Y: p.Y - other.Y,
		Z: p.Z - other.Z,
	}
}

func (p Point3D) cross(other Point3D) Point3D {
	return Point3D{
		X: p.Y*other.Z - p.Z*other.Y,
		Y: p.Z*other.X - p.X*other.Z,
		Z: p.X*other.Y - p.Y*other.X,
	}
}

func (p Point3D) magnitude() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y + p.Z*p.Z)
}

func (p Point3D) normalize() Point3D {
	mag := p.magnitude()
	if mag == 0 {
		return Point3D{0, 0, 0}
	}
	return Point3D{
		X: p.X / mag,
		Y: p.Y / mag,
		Z: p.Z / mag,
	}
}

type Cube struct {
	Vertices [8]Point3D
	Angle    float64
}

var triangles = [12][3]int{
	{0, 1, 2}, {1, 2, 3}, // Front
	{4, 5, 6}, {5, 6, 7}, // Back
	{0, 1, 4}, {1, 4, 5}, // Top
	{2, 3, 6}, {3, 6, 7}, // Bottom
	{0, 2, 4}, {2, 4, 6}, // Left
	{1, 3, 5}, {3, 5, 7}, // Right
}

type Camera struct {
	X, Y, Z    float64
	Yaw, Pitch float64
}

var faceColors = [6]color.RGBA{
	colornames.Red,
	colornames.Green,
	colornames.Blue,
	colornames.Yellow,
	colornames.Purple,
	colornames.Cyan,
}

func (c *Cube) Render(centerX, centerY float64) {
	for faceIdx, triangle := range triangles {
		coords := [3][2]float64{}
		worldCoords := [3]Point3D{}

		for j, vertexIdx := range triangle {
			v := c.Vertices[vertexIdx]

			x := v.X*math.Cos(c.Angle) - v.Z*math.Sin(c.Angle)
			z := v.X*math.Sin(c.Angle) + v.Z*math.Cos(c.Angle)

			worldCoords[j] = Point3D{x, v.Y, z}

			xRotated := x*math.Cos(camera.Yaw) - z*math.Sin(camera.Yaw)
			zRotated := x*math.Sin(camera.Yaw) + z*math.Cos(camera.Yaw)
			y := v.Y*math.Cos(camera.Pitch) - zRotated*math.Sin(camera.Pitch)
			z = v.Y*math.Sin(camera.Pitch) + zRotated*math.Cos(camera.Pitch)

			scale := 200.0 / (200.0 + z)
			xProj := xRotated * scale
			yProj := y * scale

			coords[j] = [2]float64{xProj*100 + centerX, yProj*100 + centerY}
		}

		v1 := worldCoords[1].subtract(worldCoords[0])
		v2 := worldCoords[2].subtract(worldCoords[0])
		normal := v1.cross(v2).normalize()

		if normal.Z >= 0 {
		}
		color := faceColors[faceIdx/2]

		ggez.DrawTriangle(int(coords[0][0]), int(coords[0][1]), int(coords[1][0]), int(coords[1][1]), int(coords[2][0]), int(coords[2][1]), color)
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
		ggez.PrintAt("hello", 50, 50, colornames.Red)
	})
}
