package gl

import (
	"image/color"
	"math"

	"github.com/dfirebaugh/ggez/pkg/graphics"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type ShapeRenderer struct {
	ScreenWidth  int
	ScreenHeight int
	program      graphics.ShaderProgram
}

func NewShapeRenderer(program graphics.ShaderProgram) *ShapeRenderer {
	return &ShapeRenderer{
		program: program,
	}
}

func (s ShapeRenderer) colorToNormalizedRGBA(c color.Color) (r, g, b, a float32) {
	cr, cg, cb, ca := c.RGBA()
	r = float32(cr) / 0xffff
	g = float32(cg) / 0xffff
	b = float32(cb) / 0xffff
	a = float32(ca) / 0xffff
	return r, g, b, a
}

func (s ShapeRenderer) toClipSpace(x int, y int) (float32, float32) {
	return float32(2*x)/float32(s.ScreenWidth) - 1.0, float32(2*y)/float32(s.ScreenHeight) - 1.0
}

func (s *ShapeRenderer) SetScreenSize(width int, height int) {
	s.ScreenHeight = height
	s.ScreenWidth = width
}

func (s ShapeRenderer) DrawLine(x1, y1, x2, y2 int, c color.Color) {
	s.program.Use()
	defer s.program.Delete()

	x1f, y1f := s.toClipSpace(x1, y1)
	x2f, y2f := s.toClipSpace(x2, y2)

	rf, gf, bf, af := s.colorToNormalizedRGBA(c)

	vertices := []float32{
		x1f, y1f, rf, gf, bf, af,
		x2f, y2f, rf, gf, bf, af,
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 6*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 6*4, gl.PtrOffset(2*4))
	gl.EnableVertexAttribArray(1)
	gl.DrawArrays(gl.LINES, 0, 2)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &vbo)
}

func (s ShapeRenderer) FillTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) {
	s.program.Use()
	defer s.program.Delete()

	x1f, y1f := s.toClipSpace(x1, y1)
	x2f, y2f := s.toClipSpace(x2, y2)
	x3f, y3f := s.toClipSpace(x3, y3)

	rf, gf, bf, af := s.colorToNormalizedRGBA(c)

	vertices := []float32{
		x1f, y1f, rf, gf, bf, af,
		x2f, y2f, rf, gf, bf, af,
		x3f, y3f, rf, gf, bf, af,
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 6*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 6*4, gl.PtrOffset(2*4))
	gl.EnableVertexAttribArray(1)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &vbo)
}

func (s ShapeRenderer) DrawTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) {
	s.program.Use()
	defer s.program.Delete()

	x1f, y1f := s.toClipSpace(x1, y1)
	x2f, y2f := s.toClipSpace(x2, y2)
	x3f, y3f := s.toClipSpace(x3, y3)

	rf, gf, bf, af := s.colorToNormalizedRGBA(c)

	vertices := []float32{
		x1f, y1f, rf, gf, bf, af,
		x2f, y2f, rf, gf, bf, af,
		x3f, y3f, rf, gf, bf, af,
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 6*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 6*4, gl.PtrOffset(2*4))
	gl.EnableVertexAttribArray(1)
	gl.DrawArrays(gl.LINE_LOOP, 0, 3)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &vbo)
}

func (s ShapeRenderer) FillPolygon(xPoints, yPoints []int, c color.Color) {}

func (s ShapeRenderer) DrawPolygon(xPoints, yPoints []int, c color.Color) {}

func (s ShapeRenderer) FillRect(x, y, width, height int, c color.Color) {}

func (s ShapeRenderer) DrawRect(x, y, width, height int, c color.Color) {}

func (s ShapeRenderer) simpleFillCirc(xCenter, yCenter, radius int, c color.Color) {
	s.program.Use()
	defer s.program.Delete()

	segments := 100
	angleIncrement := 2.0 * math.Pi / float64(segments)

	rf, gf, bf, af := s.colorToNormalizedRGBA(c)

	vertices := []float32{}

	xCenterF, yCenterF := s.toClipSpace(xCenter, yCenter)
	vertices = append(vertices, xCenterF, yCenterF, rf, gf, bf, af)

	for i := 0; i <= segments; i++ {
		angle := float64(i) * angleIncrement

		x := float32(xCenter) + float32(radius)*float32(math.Cos(angle))
		y := float32(yCenter) + float32(radius)*float32(math.Sin(angle))

		xf, yf := float32(2*x)/float32(s.ScreenWidth)-1.0, float32(2*y)/float32(s.ScreenHeight)-1.0

		vertices = append(vertices, xf, yf, rf, gf, bf, af)
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 6*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 6*4, gl.PtrOffset(2*4))
	gl.EnableVertexAttribArray(1)

	gl.DrawArrays(gl.TRIANGLE_FAN, 0, int32(segments+2))

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &vbo)
}

func (s ShapeRenderer) DrawCirc(x, y, radius int, c color.Color) {
	s.simpleDrawCirc(x, y, radius, c)
}
func (s ShapeRenderer) FillCirc(x, y, radius int, c color.Color) {
	s.simpleFillCirc(x, y, radius, c)
}

func (s ShapeRenderer) simpleDrawCirc(xCenter, yCenter, radius int, c color.Color) {
	s.program.Use()
	defer s.program.Delete()

	segments := 100
	angleIncrement := 2.0 * math.Pi / float64(segments)

	rf, gf, bf, af := s.colorToNormalizedRGBA(c)

	vertices := []float32{}

	for i := 0; i < segments; i++ {
		angle1 := float64(i) * angleIncrement
		angle2 := float64(i+1) * angleIncrement

		x1 := float32(xCenter) + float32(radius)*float32(math.Cos(angle1))
		y1 := float32(yCenter) + float32(radius)*float32(math.Sin(angle1))
		x2 := float32(xCenter) + float32(radius)*float32(math.Cos(angle2))
		y2 := float32(yCenter) + float32(radius)*float32(math.Sin(angle2))

		x1f, y1f := float32(2*x1)/float32(s.ScreenWidth)-1.0, float32(2*y1)/float32(s.ScreenHeight)-1.0
		x2f, y2f := float32(2*x2)/float32(s.ScreenWidth)-1.0, float32(2*y2)/float32(s.ScreenHeight)-1.0

		vertices = append(vertices, x1f, y1f, rf, gf, bf, af)
		vertices = append(vertices, x2f, y2f, rf, gf, bf, af)
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 6*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 6*4, gl.PtrOffset(2*4))
	gl.EnableVertexAttribArray(1)

	gl.DrawArrays(gl.LINES, 0, int32(2*segments))

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &vbo)
}

func (s ShapeRenderer) DrawPoint(x, y int, c color.Color) {
	s.program.Use()
	defer s.program.Delete()

	xf, yf := s.toClipSpace(x, y)

	rf, gf, bf, af := s.colorToNormalizedRGBA(c)

	vertices := []float32{
		xf, yf, rf, gf, bf, af,
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 6*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 6*4, gl.PtrOffset(2*4))
	gl.EnableVertexAttribArray(1)
	gl.DrawArrays(gl.POINTS, 0, 1)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &vbo)
}

func (s ShapeRenderer) DrawCube(x, y, z, size int, c color.Color) {
	s.program.Use()
	defer s.program.Delete()

	// 8 points of the cube
	p1 := [2]int{x, y}
	p2 := [2]int{x + size, y}
	p3 := [2]int{x, y + size}
	p4 := [2]int{x + size, y + size}
	p5 := [2]int{x, y}
	p6 := [2]int{x + size, y}
	p7 := [2]int{x, y + size}
	p8 := [2]int{x + size, y + size}

	// 12 edges of the cube
	edges := [][2][2]int{{p1, p2}, {p1, p3}, {p2, p4}, {p3, p4},
		{p1, p5}, {p2, p6}, {p3, p7}, {p4, p8},
		{p5, p6}, {p5, p7}, {p6, p8}, {p7, p8}}

	for _, edge := range edges {
		s.DrawLine(edge[0][0], edge[0][1], edge[1][0], edge[1][1], c)
	}
}

func (s ShapeRenderer) FillCube(x, y, z, size int, c color.Color) {
	// Since this is an orthographic projection in 2D, we'll just draw 6 rectangles.
	s.FillRect(x, y, size, size, c)        // Bottom
	s.FillRect(x, y+size, size, size, c)   // Top
	s.FillRect(x, y, size, z+size, c)      // Front
	s.FillRect(x+size, y, size, z+size, c) // Back
	s.FillRect(x, y, z+size, size, c)      // Left
	s.FillRect(x, y+size, z+size, size, c) // Right
}
