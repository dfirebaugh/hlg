package main

import (
	"image/color"
	"math"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/pkg/input"
	"golang.org/x/image/colornames"
)

type Point3D struct {
	X, Y, Z float64
}

type Face struct {
	Vertices [4]int
}

var (
	angleX, angleY float64
	vertices       []Point3D
	faces          []Face
	colors         = []color.Color{
		colornames.Red,
		colornames.Green,
		colornames.Blue,
		colornames.Yellow,
		colornames.Cyan,
		colornames.Magenta,
		colornames.Orange,
		colornames.Purple,
	}
	previousPolygons       []hlg.Shape
	lastMouseX, lastMouseY int
	isMousePressed         bool
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

	// Define the faces of the cube with corrected vertex order
	faces = []Face{
		{Vertices: [4]int{0, 1, 2, 3}}, // Front face
		{Vertices: [4]int{5, 4, 7, 6}}, // Back face
		{Vertices: [4]int{1, 5, 6, 2}}, // Right face
		{Vertices: [4]int{4, 0, 3, 7}}, // Left face
		{Vertices: [4]int{3, 2, 6, 7}}, // Top face
		{Vertices: [4]int{4, 5, 1, 0}}, // Bottom face
	}

	hlg.Run(update, render)
}

func update() {
	handleMouseInput()
}

func render() {
	hlg.Clear(colornames.Black)
	drawCube()
}

func handleMouseInput() {
	mouseX, mouseY := hlg.GetCursorPosition()
	if hlg.IsButtonPressed(input.MouseButtonLeft) {
		if !isMousePressed {
			lastMouseX, lastMouseY = mouseX, mouseY
			isMousePressed = true
		} else {
			deltaX := mouseX - lastMouseX
			deltaY := mouseY - lastMouseY
			angleX += float64(deltaY) * 0.01
			angleY -= float64(deltaX) * 0.01 // Invert the horizontal mouse movement
			lastMouseX = mouseX
			lastMouseY = mouseY
		}
	} else {
		isMousePressed = false
	}
}

func drawCube() {
	// Dispose of previous polygons
	for _, polygon := range previousPolygons {
		polygon.Dispose()
	}
	previousPolygons = []hlg.Shape{}

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

		projected[i] = Point3D{
			X: x*scale + halfWidth,
			Y: y*scale + halfHeight,
			Z: z,
		}
	}

	for i, face := range faces {
		// Calculate face normal
		v0 := projected[face.Vertices[0]]
		v1 := projected[face.Vertices[1]]
		v2 := projected[face.Vertices[2]]
		normal := calculateNormal(v0, v1, v2)

		// Check if the face is facing the camera
		if normal.Z > 0 {
			faceVertices := []hlg.Vertex{
				{Position: [3]float32{float32(projected[face.Vertices[0]].X), float32(projected[face.Vertices[0]].Y), float32(projected[face.Vertices[0]].Z)}, Color: toRGBA(colors[i%len(colors)])},
				{Position: [3]float32{float32(projected[face.Vertices[1]].X), float32(projected[face.Vertices[1]].Y), float32(projected[face.Vertices[1]].Z)}, Color: toRGBA(colors[i%len(colors)])},
				{Position: [3]float32{float32(projected[face.Vertices[2]].X), float32(projected[face.Vertices[2]].Y), float32(projected[face.Vertices[2]].Z)}, Color: toRGBA(colors[i%len(colors)])},
				{Position: [3]float32{float32(projected[face.Vertices[0]].X), float32(projected[face.Vertices[0]].Y), float32(projected[face.Vertices[0]].Z)}, Color: toRGBA(colors[i%len(colors)])},
				{Position: [3]float32{float32(projected[face.Vertices[2]].X), float32(projected[face.Vertices[2]].Y), float32(projected[face.Vertices[2]].Z)}, Color: toRGBA(colors[i%len(colors)])},
				{Position: [3]float32{float32(projected[face.Vertices[3]].X), float32(projected[face.Vertices[3]].Y), float32(projected[face.Vertices[3]].Z)}, Color: toRGBA(colors[i%len(colors)])},
			}
			polygon := hlg.PolygonFromVertices(0, 0, 0, faceVertices)
			polygon.Render()
			previousPolygons = append(previousPolygons, polygon)
		}
	}
}

// calculateNormal calculates the normal vector of a face defined by three vertices.
func calculateNormal(v0, v1, v2 Point3D) Point3D {
	u := Point3D{X: v1.X - v0.X, Y: v1.Y - v0.Y, Z: v1.Z - v0.Z}
	v := Point3D{X: v2.X - v0.X, Y: v2.Y - v0.Y, Z: v2.Z - v0.Z}
	return Point3D{
		X: u.Y*v.Z - u.Z*v.Y,
		Y: u.Z*v.X - u.X*v.Z,
		Z: u.X*v.Y - u.Y*v.X,
	}
}

// toRGBA converts a color.Color to an array of float32 RGBA values.
func toRGBA(c color.Color) [4]float32 {
	r, g, b, a := c.RGBA()
	return [4]float32{
		float32(r) / 0xffff,
		float32(g) / 0xffff,
		float32(b) / 0xffff,
		float32(a) / 0xffff,
	}
}
