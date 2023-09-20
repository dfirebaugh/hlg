package renderer

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"strings"

	"github.com/dfirebaugh/ggez/pkg/input"
	"github.com/dfirebaugh/ggez/pkg/renderer/sdl"
	"github.com/go-gl/gl/v4.5-core/gl"
)

type GLRenderer struct {
	window     uintptr
	glContext  uintptr
	shaders    []uint32
	lineShader uint32
	lineVAO    uint32
	lineVBO    uint32
}

var (
	width, height int
)

var (
	vertexLineShaderSource = `
	#version 450 core
	
	layout(location = 0) in vec2 inPosition;
	
	void main() {
		gl_Position = vec4(inPosition, 0.0, 1.0);
	}
	`

	fragmentLineShaderSource = `
	#version 450 core
	
	uniform vec4 lineColor;  // Add uniform for the line color
	out vec4 outColor;
	
	void main() {
		outColor = lineColor;
	}
	`
)

func NewGLRenderer() (*GLRenderer, error) {
	var err error
	g := &GLRenderer{
		shaders: []uint32{},
	}

	sdl.SDL_Init(sdl.SDL_INIT_VIDEO)
	sdl.SDL_GL_SetAttribute(sdl.SDL_GL_CONTEXT_MAJOR_VERSION, 4)
	sdl.SDL_GL_SetAttribute(sdl.SDL_GL_CONTEXT_MINOR_VERSION, 5)
	sdl.SDL_GL_SetAttribute(sdl.SDL_GL_CONTEXT_PROFILE_MASK, int(sdl.SDL_GL_CONTEXT_PROFILE_ES))

	g.window, err = g.CreateWindow("ggez", 800, 600)
	if err != nil {
		log.Fatalf("could not create OpenGL context: %s", sdl.SDL_GetError())
	}
	g.glContext = sdl.SDL_GL_CreateContext(g.window)
	if g.glContext == 0 {
		log.Fatalf("could not create OpenGL context: %s", sdl.SDL_GetError())
	}

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}
	gl.Viewport(0, 0, 800, 600)
	gl.ClearColor(0, 0, 0, 1)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	g.lineShader = createShader(vertexLineShaderSource, fragmentLineShaderSource)
	g.shaders = append(g.shaders, g.lineShader)
	g.lineVAO, g.lineVBO = createVBOandVAO([]float32{
		-0.5, 0.0, // Start point
		0.5, 0.0, // End point
	})

	if err := gl.GetError(); err != gl.NO_ERROR {
		log.Printf("OpenGL error: %v", err)
	}

	for _, s := range g.shaders {
		var status int32
		gl.GetShaderiv(s, gl.COMPILE_STATUS, &status)
		if status == gl.FALSE {
			var logLength int32
			gl.GetShaderiv(s, gl.INFO_LOG_LENGTH, &logLength)

			log := strings.Repeat("\x00", int(logLength+1))
			gl.GetShaderInfoLog(s, logLength, nil, gl.Str(log))

			fmt.Print("failed to compile vertex shader:", log)
		}
	}

	for _, s := range g.shaders {
		checkShaderCompileStatus(s)
		checkProgramLinkStatus(s)
	}
	return g, nil
}

func (glr *GLRenderer) Close() {
	sdl.SDL_DestroyWindow(glr.window)
	sdl.SDL_GL_DeleteContext(glr.glContext)
}

func checkShaderCompileStatus(shader uint32) {
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		logStr := make([]byte, logLength+1)

		gl.GetShaderInfoLog(shader, logLength, nil, &logStr[0])

		fmt.Println("Shader compilation failed:", string(logStr))
	}
}

func checkProgramLinkStatus(program uint32) {
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		logStr := make([]byte, logLength)
		gl.GetProgramInfoLog(program, logLength, nil, &logStr[0])

		fmt.Println("Program link failed:", string(logStr))
	}
}

func createVBOandVAO(vertices []float32) (uint32, uint32) {
	var VBO, VAO uint32

	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)

	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return VBO, VAO
}

func createShader(vertexShaderSource string, fragmentShaderSource string) uint32 {
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	cstr, free := gl.Strs(vertexShaderSource)
	gl.ShaderSource(vertexShader, 1, cstr, nil)
	free()
	gl.CompileShader(vertexShader)

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	cstr, free = gl.Strs(fragmentShaderSource)
	gl.ShaderSource(fragmentShader, 1, cstr, nil)
	free()
	gl.CompileShader(fragmentShader)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program
}

func (glr *GLRenderer) CreateWindow(title string, width int, height int) (uintptr, error) {
	window := sdl.SDL_CreateWindow(title, 0, 0, width, height, sdl.SDL_WINDOW_OPENGL)
	if window == 0 {
		return 0, fmt.Errorf("could not create window")
	}
	glr.window = window
	return window, nil
}

func (glr *GLRenderer) SetWindowTitle(title string) {
	sdl.SDL_SetWindowTitle(glr.window, title)
}

func (glr *GLRenderer) DestroyWindow() {
	sdl.SDL_DestroyWindow(glr.window)
}

func (glr *GLRenderer) CreateTextureFromImage(img image.Image) (uintptr, error) {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	gl.GenerateMipmap(gl.TEXTURE_2D)
	return uintptr(textureID), nil
}

func (glr *GLRenderer) RenderTexture(textureInstance uintptr, x int, y int, w int, h int, angle float64, centerX int, centerY int, flipType int) {
}

func (glr *GLRenderer) DrawLine(x1, y1, x2, y2 int, c color.Color) {
}

func (glr *GLRenderer) PrintPlatformAndVersion() {
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version:", version)
}

func (glr *GLRenderer) PrintRendererInfo() {
	renderer := gl.GoStr(gl.GetString(gl.RENDERER))
	fmt.Println("Renderer:", renderer)
}

func (glr *GLRenderer) Clear(c color.Color) {
	r, g, b, a := c.RGBA()
	gl.ClearColor(float32(r)/0xffff, float32(g)/0xffff, float32(b)/0xffff, float32(a)/0xffff)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (glr *GLRenderer) RenderPresent() {
	for _, s := range glr.shaders {
		gl.UseProgram(s)
	}
	sdl.SDL_GL_SwapWindow(glr.window)
}

func (glr *GLRenderer) SetScreenSize(w int, h int) {
	width = w
	height = h
	gl.Viewport(0, 0, int32(width), int32(height))
}

func (glr *GLRenderer) DestroyTexture(textureInstance uintptr) {
	textureID := uint32(textureInstance)
	gl.DeleteTextures(1, &textureID)
}

func (glr *GLRenderer) FillTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) {
}

func (glr *GLRenderer) DrawTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) {
}

func (glr *GLRenderer) FillPolygon(xPoints, yPoints []int, c color.Color) {
}

func (glr *GLRenderer) DrawPolygon(xPoints, yPoints []int, c color.Color) {
}

func (glr *GLRenderer) FillRect(x, y, width, height int, c color.Color) {
}

func (glr *GLRenderer) DrawRect(x, y, width, height int, c color.Color) {
}

func (glr *GLRenderer) FillCirc(x, y, radius int, c color.Color) {
}

func (glr *GLRenderer) DrawCirc(xCenter, yCenter, radius int, c color.Color) {
}

func (glr *GLRenderer) DrawPoint(x, y int, c color.Color) {
}

func (glr *GLRenderer) PollEvents(i input.InputDevice) bool {
	return sdl.PollEvents(i)
}
