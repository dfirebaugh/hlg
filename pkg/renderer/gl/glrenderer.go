package gl

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/dfirebaugh/ggez/pkg/input"
	"github.com/dfirebaugh/ggez/pkg/renderer"
	"github.com/dfirebaugh/ggez/pkg/shader"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type GLRenderer struct {
	glfwWindow   *glfw.Window
	shaders      []uint32
	mainShader   uint32
	scaleFactor  int
	screenHeight int
	screenWidth  int
	renderer.ShapeRenderer
}

var (
	vertexSource = `
#version 450 core
layout(location = 0) in vec2 position;
layout(location = 1) in vec4 vertexColor;

out vec4 fragmentColor;

void main() {
	vec2 flippedPosition = vec2(position.x, -position.y);
	gl_Position = vec4(flippedPosition, 0.0, 1.0);
	fragmentColor = vertexColor;
}
`

	fragmentSource = `
#version 450 core
in vec4 fragmentColor;

out vec4 color;

void main() {
	color = fragmentColor;
}
`
)

func New() (*GLRenderer, error) {
	var err error
	g := &GLRenderer{
		shaders: []uint32{},
		ShapeRenderer: &GLShapeRenderer{
			ScreenHeight: 160,
			ScreenWidth:  240,
		},
		scaleFactor:  3,
		screenHeight: 160,
		screenWidth:  240,
	}

	if err := glfw.Init(); err != nil {
		log.Fatalf("could not initialize GLFW: %v", err)
	}

	g.glfwWindow, err = glfw.CreateWindow(g.screenWidth, g.screenHeight, "ggez", nil, nil)
	if err != nil {
		log.Fatalf("could not create window: %v", err)
	}

	monitor := glfw.GetPrimaryMonitor()
	vidMode := monitor.GetVideoMode()
	centerX := (vidMode.Width - g.screenWidth*g.scaleFactor) / 2
	centerY := (vidMode.Height - g.screenHeight*g.scaleFactor) / 2
	g.glfwWindow.SetPos(centerX, centerY)

	g.glfwWindow.SetFramebufferSizeCallback(g.resizeCallback)

	g.glfwWindow.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalf("could not initialize OpenGL bindings: %v", err)
	}

	sm, err := shader.NewShaderManager(vertexSource, fragmentSource)
	if err != nil {
		panic(err.Error())
	}
	g.mainShader = sm.Program
	g.shaders = append(g.shaders, g.mainShader)

	return g, nil
}

func (glr *GLRenderer) resizeCallback(window *glfw.Window, width int, height int) {
	windowAspectRatio := float64(width) / float64(height)

	desiredAspectRatio := float64(glr.screenWidth / glr.screenHeight)

	var newWidth, newHeight int
	var offsetX, offsetY int

	if windowAspectRatio > desiredAspectRatio {
		newWidth = int(float64(height) * desiredAspectRatio)
		newHeight = height
		offsetX = (width - newWidth) / 2
	} else {
		newWidth = width
		newHeight = int(float64(width) / desiredAspectRatio)
		offsetY = (height - newHeight) / 2
	}

	gl.Viewport(int32(offsetX), int32(offsetY), int32(newWidth), int32(newHeight))
}

func (glr *GLRenderer) Close() {
	glfw.Terminate()
}

func (glr *GLRenderer) SetScaleFactor(s int) {
	glr.scaleFactor = s
	glr.resizeCallback(glr.glfwWindow, glr.screenWidth, glr.screenHeight)
}

func (glr *GLRenderer) SetWindowTitle(title string) {
	glr.glfwWindow.SetTitle(title)
}

func (glr *GLRenderer) DestroyWindow() {
	glr.glfwWindow.Destroy()
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
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.ClearColor(float32(r)/0xffff, float32(g)/0xffff, float32(b)/0xffff, float32(a)/0xffff)
}

func (glr *GLRenderer) Render() {
	for _, s := range glr.shaders {
		gl.UseProgram(s)
	}
	glr.glfwWindow.SwapBuffers()
}

func (glr *GLRenderer) SetScreenSize(w int, h int) {
	glr.ShapeRenderer.SetScreenSize(w, h)
	glr.SetWindowSize(w, h)
}

func (glr *GLRenderer) SetWindowSize(width, height int) {
	glr.screenWidth = width
	glr.screenHeight = height

	glr.glfwWindow.SetSize(width, height)

	glr.resizeCallback(glr.glfwWindow, width, height)
}

func (glr *GLRenderer) CreateTextureFromImage(img image.Image) (uintptr, error) {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var textureID uint32

	return uintptr(textureID), nil
}

func (glr *GLRenderer) RenderTexture(textureInstance uintptr, x int, y int, w int, h int, angle float64, centerX int, centerY int, flipType int) {
}

func (glr *GLRenderer) DestroyTexture(textureInstance uintptr) {
	textureID := uint32(textureInstance)
	gl.DeleteTextures(1, &textureID)
}

func (glr *GLRenderer) PollEvents(i input.InputDevice) bool {
	defer glfw.PollEvents()
	return !glr.glfwWindow.ShouldClose()
}
