package main

import (
	_ "embed"
	"image/color"
	"math"
	"unsafe"

	"github.com/dfirebaugh/hlg"
	"golang.org/x/image/colornames"
)

var mousePosition [2]float32

//go:embed shader.wgsl
var shaderCode string

var (
	screenWidth  float32 = 800
	screenHeight float32 = 600
)

const (
	DrawOpCircle = iota
	DrawOpRoundedRect
	DrawOpTriangle
	DrawOpOrientedBox
	DrawOpSegment
	DrawOpBezier
)

type Vertex struct {
	Position      [3]float32
	LocalPosition [2]float32
	OpCode        float32
	Radius        float32
	Color         [4]float32
}

func VerticesToBytes(vertices []Vertex) []byte {
	size := len(vertices) * int(unsafe.Sizeof(Vertex{}))
	data := make([]byte, size)
	copy(data, unsafe.Slice((*byte)(unsafe.Pointer(&vertices[0])), size))
	return data
}

func main() {
	hlg.SetWindowSize(int(screenWidth), int(screenHeight))
	hlg.SetTitle("SDF primitive buffer")

	shader := hlg.CompileShader(shaderCode)

	vertexLayout := hlg.VertexBufferLayout{
		ArrayStride: 4*3 + 4*2 + 4 + 4 + 4*4,
		Attributes: []hlg.VertexAttributeLayout{
			{
				ShaderLocation: 0,
				Offset:         0,
				Format:         "float32x3", // Position (Clip space)
			},
			{
				ShaderLocation: 1,
				Offset:         3 * 4,
				Format:         "float32x2", // Local Position
			},
			{
				ShaderLocation: 2,
				Offset:         (3 + 2) * 4,
				Format:         "float32", // OpCode
			},
			{
				ShaderLocation: 3,
				Offset:         (3 + 2 + 1) * 4,
				Format:         "float32", // Radius
			},
			{
				ShaderLocation: 4,
				Offset:         (3 + 2 + 1 + 1) * 4,
				Format:         "float32x4", // Color
			},
		},
	}

	var quad hlg.Renderable

	hlg.Run(func() {
		screenWidth = float32(hlg.GetWindowWidth())
		screenHeight = float32(hlg.GetWindowHeight())
	}, func() {
		x, y := hlg.GetCursorPosition()

		mousePosition[0] = float32(x)
		mousePosition[1] = float32(y)

		quadVertices := makeVerticesForShapes(mousePosition)
		quadVertexData := VerticesToBytes(quadVertices)

		if quad != nil {
			quad = nil
		}

		quad = hlg.CreateRenderable(shader, quadVertexData, vertexLayout, nil, nil)
		if quad == nil {
			panic("Failed to create renderable")
		}

		hlg.Clear(colornames.Skyblue)
		quad.Render()
	})
}

func makeVerticesForShapes(mousePos [2]float32) []Vertex {
	circleColor := colornames.Tomato
	rectColor := colornames.Purple
	triangleColor := colornames.Steelblue

	circleVertices := makeCircleVertices(mousePos[0]-60, mousePos[1], 40, circleColor)
	rectVertices := makeRoundedRectVertices(mousePos[0], mousePos[1], 240, 180, rectColor)
	triangleVertices := makeTriangleVertices(mousePos[0]+60, mousePos[1], 100, triangleColor)

	vertices := append(rectVertices, circleVertices...)
	vertices = append(vertices, triangleVertices...)

	return vertices
}

func makeCircleVertices(centerX, centerY, radius float32, color color.Color) []Vertex {
	ndcCenterX := (centerX/screenWidth)*2.0 - 1.0
	ndcCenterY := 1.0 - (centerY/screenHeight)*2.0

	left := ndcCenterX - radius/screenWidth*2.0
	right := ndcCenterX + radius/screenWidth*2.0
	bottom := ndcCenterY - radius/screenHeight*2.0
	top := ndcCenterY + radius/screenHeight*2.0

	colorVec := colorToFloat32(color)

	return []Vertex{
		{Position: [3]float32{left, bottom, 0.0}, LocalPosition: [2]float32{-1.0, -1.0}, OpCode: DrawOpCircle, Radius: radius, Color: colorVec},
		{Position: [3]float32{right, bottom, 0.0}, LocalPosition: [2]float32{1.0, -1.0}, OpCode: DrawOpCircle, Radius: radius, Color: colorVec},
		{Position: [3]float32{left, top, 0.0}, LocalPosition: [2]float32{-1.0, 1.0}, OpCode: DrawOpCircle, Radius: radius, Color: colorVec},

		{Position: [3]float32{left, top, 0.0}, LocalPosition: [2]float32{-1.0, 1.0}, OpCode: DrawOpCircle, Radius: radius, Color: colorVec},
		{Position: [3]float32{right, bottom, 0.0}, LocalPosition: [2]float32{1.0, -1.0}, OpCode: DrawOpCircle, Radius: radius, Color: colorVec},
		{Position: [3]float32{right, top, 0.0}, LocalPosition: [2]float32{1.0, 1.0}, OpCode: DrawOpCircle, Radius: radius, Color: colorVec},
	}
}

func makeRoundedRectVertices(centerX, centerY, width, height float32, color color.Color) []Vertex {
	ndcCenterX := (centerX/screenWidth)*2.0 - 1.0
	ndcCenterY := 1.0 - (centerY/screenHeight)*2.0
	ndcWidth := (width / screenWidth) * 2.0
	ndcHeight := (height / screenHeight) * 2.0

	left := ndcCenterX - ndcWidth/2
	right := ndcCenterX + ndcWidth/2
	bottom := ndcCenterY - ndcHeight/2
	top := ndcCenterY + ndcHeight/2

	colorVec := colorToFloat32(color)

	return []Vertex{
		{Position: [3]float32{left, bottom, 0.0}, LocalPosition: [2]float32{-1.0, -1.0}, OpCode: DrawOpRoundedRect, Radius: 0.1, Color: colorVec},
		{Position: [3]float32{right, bottom, 0.0}, LocalPosition: [2]float32{1.0, -1.0}, OpCode: DrawOpRoundedRect, Radius: 0.1, Color: colorVec},
		{Position: [3]float32{left, top, 0.0}, LocalPosition: [2]float32{-1.0, 1.0}, OpCode: DrawOpRoundedRect, Radius: 0.1, Color: colorVec},

		{Position: [3]float32{left, top, 0.0}, LocalPosition: [2]float32{-1.0, 1.0}, OpCode: DrawOpRoundedRect, Radius: 0.1, Color: colorVec},
		{Position: [3]float32{right, bottom, 0.0}, LocalPosition: [2]float32{1.0, -1.0}, OpCode: DrawOpRoundedRect, Radius: 0.1, Color: colorVec},
		{Position: [3]float32{right, top, 0.0}, LocalPosition: [2]float32{1.0, 1.0}, OpCode: DrawOpRoundedRect, Radius: 0.1, Color: colorVec},
	}
}

func makeTriangleVertices(centerX, centerY, size float32, color color.Color) []Vertex {
	ndcCenterX := (centerX/screenWidth)*2.0 - 1.0
	ndcCenterY := 1.0 - (centerY/screenHeight)*2.0

	height := size * float32(math.Sqrt(3.0)) / 2.0

	v0 := [2]float32{ndcCenterX, ndcCenterY + height/2/screenHeight*2.0}
	v1 := [2]float32{ndcCenterX - size/2/screenWidth*2.0, ndcCenterY - height/2/screenHeight*2.0}
	v2 := [2]float32{ndcCenterX + size/2/screenWidth*2.0, ndcCenterY - height/2/screenHeight*2.0}

	colorVec := colorToFloat32(color)

	return []Vertex{
		{Position: [3]float32{v0[0], v0[1], 0.0}, LocalPosition: [2]float32{0.0, 1.0}, OpCode: DrawOpTriangle, Radius: 0.0, Color: colorVec},
		{Position: [3]float32{v1[0], v1[1], 0.0}, LocalPosition: [2]float32{-1.0, -1.0}, OpCode: DrawOpTriangle, Radius: 0.0, Color: colorVec},
		{Position: [3]float32{v2[0], v2[1], 0.0}, LocalPosition: [2]float32{1.0, -1.0}, OpCode: DrawOpTriangle, Radius: 0.0, Color: colorVec},
	}
}

func colorToFloat32(c color.Color) [4]float32 {
	r, g, b, a := c.RGBA()
	return [4]float32{
		float32(r) / 0xffff,
		float32(g) / 0xffff,
		float32(b) / 0xffff,
		float32(a) / 0xffff,
	}
}
