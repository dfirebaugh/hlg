package main

import (
	"image/color"
	"math"

	"github.com/dfirebaugh/hlg"
	"golang.org/x/image/colornames"
)

type Point3D struct {
	X, Y, Z float64
}

type Edge struct {
	Start, End int
}

var (
	angleX, angleY, angleZ float64
	vertices               []Point3D
	edges                  []Edge
	colors                 = []color.Color{
		colornames.Red,
		colornames.Green,
		colornames.Blue,
		colornames.Yellow,
		colornames.Cyan,
		colornames.Magenta,
		colornames.Orange,
		colornames.Purple,
	}
	previousLines []hlg.Shape
)

func main() {
	hlg.SetTitle("Rotating Cube Example")
	hlg.SetWindowSize(800, 600)
	hlg.EnableFPS()

	// Define the vertices of a cube
	vertices = []Point3D{
		{-1, -1, -1},
		{1, -1, -1},
		{1, 1, -1},
		{-1, 1, -1},
		{-1, -1, 1},
		{1, -1, 1},
		{1, 1, 1},
		{-1, 1, 1},
	}

	// Define the edges of the cube connecting the vertices
	edges = []Edge{
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 0},
		{4, 5},
		{5, 6},
		{6, 7},
		{7, 4},
		{0, 4},
		{1, 5},
		{2, 6},
		{3, 7},
	}

	hlg.Run(func() {
		angleX += 0.01
		angleY += 0.01
		angleZ += 0.01
	}, func() {
		hlg.Clear(colornames.Black)
		drawCube()
	})
}

func drawCube() {
	// Dispose of previous lines
	for _, line := range previousLines {
		line.Dispose()
	}
	previousLines = []hlg.Shape{}

	screenWidth, screenHeight := hlg.GetScreenSize()
	scale := float64(screenWidth) / 4
	halfWidth := float64(screenWidth) / 2
	halfHeight := float64(screenHeight) / 2

	projected := make([]Point3D, len(vertices))
	for i, v := range vertices {
		x := v.X
		y := v.Y
		z := v.Z

		// Rotate around X axis
		xy := y*math.Cos(angleX) - z*math.Sin(angleX)
		yz := y*math.Sin(angleX) + z*math.Cos(angleX)
		y = xy
		z = yz

		// Rotate around Y axis
		xz := z*math.Cos(angleY) - x*math.Sin(angleY)
		yz = z*math.Sin(angleY) + x*math.Cos(angleY)
		x = xz
		z = yz

		// Rotate around Z axis
		xy = x*math.Cos(angleZ) - y*math.Sin(angleZ)
		yz = x*math.Sin(angleZ) + y*math.Cos(angleZ)
		x = xy
		y = yz

		projected[i] = Point3D{
			X: x*scale + halfWidth,
			Y: y*scale + halfHeight,
			Z: z,
		}
	}

	for i, edge := range edges {
		start := projected[edge.Start]
		end := projected[edge.End]
		colorIndex := i % len(colors)
		line := hlg.Line(int(start.X), int(start.Y), int(end.X), int(end.Y), 2, colors[colorIndex])
		line.Render()
		previousLines = append(previousLines, line)
	}
}
