package gl

import (
	"image/color"
	"math"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type ShapeRenderer struct {
	ScreenWidth  int
	ScreenHeight int
	shapeProgram *Program
	// polygonProgram *Program
}

func NewShapeRenderer() *ShapeRenderer {
	shapeVertShader, err := NewShader(ShapeVert, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	shapeFragShader, err := NewShader(ShapeFrag, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	shapeProgram, err := NewProgram(shapeVertShader, shapeFragShader)
	if err != nil {
		panic(err)
	}

	return &ShapeRenderer{
		shapeProgram: shapeProgram,
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

func (s ShapeRenderer) createTriangleVAO(x1, y1, x2, y2, x3, y3 int) uint32 {
	vertices := []float32{
		float32(x1)/float32(windowWidth)*2 - 1.0, -float32(y1)/float32(windowHeight)*2 + 1.0, 0.0,
		float32(x2)/float32(windowWidth)*2 - 1.0, -float32(y2)/float32(windowHeight)*2 + 1.0, 0.0,
		float32(x3)/float32(windowWidth)*2 - 1.0, -float32(y3)/float32(windowHeight)*2 + 1.0, 0.0,
	}

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteBuffers(1, &VBO)

	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, uintptr(0))
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)

	return VAO
}

func (s ShapeRenderer) FillTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) {
	VAO := s.createTriangleVAO(x1, y1, x2, y2, x3, y3)
	s.shapeProgram.Use()

	gl.BindVertexArray(VAO)

	// Set the uniform color variable in the fragment shader
	triangleColor := s.shapeProgram.GetUniformLocation("shapeColor")

	r, g, b, a := s.colorToNormalizedRGBA(c)
	gl.Uniform4f(triangleColor, r, g, b, a)

	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.BindVertexArray(0)
}

func (s ShapeRenderer) DrawTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) {
	s.DrawLine(x1, y1, x2, y2, c)
	s.DrawLine(x2, y2, x3, y3, c)
	s.DrawLine(x3, y3, x1, y1, c)
}

func (s ShapeRenderer) DrawLine(x1, y1, x2, y2 int, c color.Color) {
	s.shapeProgram.Use()

	// Convert coordinates to OpenGL's normalized device coordinates (-1 to 1)
	nx1 := float32(x1)/float32(windowWidth)*2 - 1.0
	ny1 := -float32(y1)/float32(windowHeight)*2 + 1.0
	nx2 := float32(x2)/float32(windowWidth)*2 - 1.0
	ny2 := -float32(y2)/float32(windowHeight)*2 + 1.0

	vertices := []float32{
		nx1, ny1, 0.0,
		nx2, ny2, 0.0,
	}

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)
	gl.BindVertexArray(VAO)

	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteBuffers(1, &VBO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, uintptr(0))
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)

	s.shapeProgram.Use()
	gl.BindVertexArray(VAO)

	// Set the uniform color variable in the fragment shader
	lineColor := s.shapeProgram.GetUniformLocation("shapeColor")
	r, g, b, a := s.colorToNormalizedRGBA(c)
	gl.Uniform4f(lineColor, r, g, b, a)

	gl.DrawArrays(gl.LINES, 0, 2)
	gl.BindVertexArray(0)
}

// Fill a polygon with a specified color
func (s ShapeRenderer) FillPolygon(xPoints, yPoints []int, c color.Color) {
	s.shapeProgram.Use()

	if len(xPoints) != len(yPoints) || len(xPoints) < 3 {
		// Invalid input, can't draw a polygon
		return
	}

	shapeColor := s.shapeProgram.GetUniformLocation("shapeColor")
	r, g, b, a := s.colorToNormalizedRGBA(c)
	gl.Uniform4f(shapeColor, r, g, b, a)

	vertices := make([]float32, 0, len(xPoints)*2)
	for i := 0; i < len(xPoints); i++ {
		x := float32(xPoints[i])/float32(windowWidth)*2 - 1.0
		y := -float32(yPoints[i])/float32(windowHeight)*2 + 1.0
		vertices = append(vertices, x, y)
	}

	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)

	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteBuffers(1, &VBO)

	gl.BindVertexArray(VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)

	gl.DrawArrays(gl.TRIANGLE_FAN, 0, int32(len(xPoints)))
	gl.BindVertexArray(0)
}

// DrawPolygon draws the outline of a polygon with a specified color
func (s ShapeRenderer) DrawPolygon(xPoints, yPoints []int, c color.Color) {
	if len(xPoints) != len(yPoints) {
		return
	}

	s.shapeProgram.Use()

	// Set the uniform color variable in the fragment shader
	polygonColor := s.shapeProgram.GetUniformLocation("shapeColor")
	r, g, b, a := s.colorToNormalizedRGBA(c)
	gl.Uniform4f(polygonColor, r, g, b, a)

	// Create vertex and index slices for the polygon
	vertices := make([]float32, 0, len(xPoints)*2)
	indices := make([]uint32, 0, len(xPoints))

	for i := 0; i < len(xPoints); i++ {
		x := float32(xPoints[i])/float32(windowWidth)*2 - 1.0
		y := -float32(yPoints[i])/float32(windowHeight)*2 + 1.0
		vertices = append(vertices, x, y)
		indices = append(indices, uint32(i))
	}

	// Create a VAO and VBO for the polygon
	var VAO, VBO, EBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.GenBuffers(1, &EBO)

	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteBuffers(1, &VBO)

	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// Specify the layout of the vertex data
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)

	// Draw the polygon using indexed drawing
	gl.DrawElements(gl.LINE_LOOP, int32(len(indices)), gl.UNSIGNED_INT, nil)

	gl.BindVertexArray(0)
}

// FillRect fills a rectangle with a specified color
func (s ShapeRenderer) FillRect(x, y, width, height int, c color.Color) {
	xPoints := []int{x, x + width, x + width, x}
	yPoints := []int{y, y, y + height, y + height}
	s.FillPolygon(xPoints, yPoints, c)
}

// DrawRect draws the outline of a rectangle with a specified color
func (s ShapeRenderer) DrawRect(x, y, width, height int, c color.Color) {
	xPoints := []int{x, x + width, x + width, x}
	yPoints := []int{y, y, y + height, y + height}
	s.DrawPolygon(xPoints, yPoints, c)
}

func (s ShapeRenderer) simpleFillCirc(xCenter, yCenter, radius int, c color.Color) {
	const (
		circSegments = 100
	)
	xPoints := make([]int, circSegments)
	yPoints := make([]int, circSegments)

	for i := 0; i < circSegments; i++ {
		angle := 2.0 * math.Pi * float64(i) / float64(circSegments)
		x := float64(xCenter) + float64(radius)*math.Cos(angle)
		y := float64(yCenter) + float64(radius)*math.Sin(angle)
		xPoints[i] = int(x)
		yPoints[i] = int(y)
	}

	s.FillPolygon(xPoints, yPoints, c)
}

func (s ShapeRenderer) simpleDrawCirc(xCenter, yCenter, radius int, c color.Color) {
	const (
		circSegments = 100
	)
	xPoints := make([]int, circSegments)
	yPoints := make([]int, circSegments)

	for i := 0; i < circSegments; i++ {
		angle := 2.0 * math.Pi * float64(i) / float64(circSegments)
		x := float64(xCenter) + float64(radius)*math.Cos(angle)
		y := float64(yCenter) + float64(radius)*math.Sin(angle)
		xPoints[i] = int(x)
		yPoints[i] = int(y)
	}

	s.DrawPolygon(xPoints, yPoints, c)
}

func (s ShapeRenderer) DrawPoint(x, y int, c color.Color) {
	const pointSize = 5
	s.DrawRect(x, y, pointSize, pointSize, c)
}

func (s ShapeRenderer) DrawCircle(xCenter, yCenter, radius int, c color.Color) {
	s.simpleDrawCirc(xCenter, yCenter, radius, c)
}

func (s ShapeRenderer) FillCircle(xCenter, yCenter, radius int, c color.Color) {
	s.simpleFillCirc(xCenter, yCenter, radius, c)
}
