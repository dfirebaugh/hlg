package gl

import (
	"fmt"
	"image/color"
	"log"

	"github.com/dfirebaugh/ggez/pkg/graphics"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type GLRenderer struct {
	window                  *glfw.Window
	programs                []*Program
	modelProgram            *Program
	nonTexturedModelProgram *Program
	textureProgram          *Program
	scaleFactor             int
	screenHeight            int
	screenWidth             int
	aspectRatioX            float32
	aspectRatioY            float32
	VAO                     uint32
	wireframeEnabled        bool
	graphics.ShapeRenderer
	graphics.InputManager
	ModelRenderer
}

var (
	windowWidth, windowHeight int
)

var textures map[uintptr]*Texture

func New() (*GLRenderer, error) {
	var err error
	g := &GLRenderer{
		scaleFactor:  3,
		screenHeight: 160,
		screenWidth:  240,
	}

	if err := glfw.Init(); err != nil {
		log.Fatalf("could not initialize GLFW: %v", err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	g.window, err = glfw.CreateWindow(g.screenWidth, g.screenHeight, "ggez", nil, nil)
	if err != nil {
		log.Fatalf("could not create window: %v", err)
	}

	g.window.MakeContextCurrent()
	g.window.SetKeyCallback(g.keyCallback)

	monitor := glfw.GetPrimaryMonitor()
	vidMode := monitor.GetVideoMode()
	centerX := (vidMode.Width - g.screenWidth*g.scaleFactor) / 2
	centerY := (vidMode.Height - g.screenHeight*g.scaleFactor) / 2
	g.window.SetPos(centerX, centerY)

	g.window.SetFramebufferSizeCallback(g.resizeCallback)

	if err := gl.Init(); err != nil {
		log.Fatalf("could not initialize OpenGL bindings: %v", err)
	}

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	g.setupShaders()
	g.PrintPlatformAndVersion()
	g.PrintRendererInfo()

	g.ModelRenderer = NewModelRenderer(g.modelProgram, g.nonTexturedModelProgram)
	g.ShapeRenderer = NewShapeRenderer()
	g.InputManager = NewInputDeviceGlfw(g.window)

	checkGLError()

	return g, nil
}

func (g *GLRenderer) setupShaders() error {
	textureVertexShader, err := NewShader(TextureVert, gl.VERTEX_SHADER)
	if err != nil {
		return err
	}

	textureFragmentShader, err := NewShader(TextureFrag, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	g.textureProgram, err = NewProgram(textureVertexShader, textureFragmentShader)
	if err != nil {
		return err
	}
	g.programs = append(g.programs, g.textureProgram)

	modelVertShader, err := NewShader(ModelVert, gl.VERTEX_SHADER)
	if err != nil {
		return err
	}

	modelFragShader, err := NewShader(ModelFrag, gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}

	g.modelProgram, err = NewProgram(modelVertShader, modelFragShader)
	if err != nil {
		return err
	}
	g.programs = append(g.programs, g.modelProgram)

	nonTexturedModelVertShader, err := NewShader(NonTexturedModelVert, gl.VERTEX_SHADER)
	if err != nil {
		return err
	}

	nonTexturedModelFragShader, err := NewShader(NonTexturedModelFrag, gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}

	g.nonTexturedModelProgram, err = NewProgram(nonTexturedModelVertShader, nonTexturedModelFragShader)
	if err != nil {
		return err
	}
	g.programs = append(g.programs, g.nonTexturedModelProgram)

	return nil
}

func (g *GLRenderer) resizeCallback(window *glfw.Window, width int, height int) {
	if width == 0 || height == 0 {
		return
	}

	desiredAspectRatio := float64(g.screenWidth) / float64(g.screenHeight)
	viewportAspectRatio := float64(width) / float64(height)

	if desiredAspectRatio > viewportAspectRatio {
		g.aspectRatioX = 1.0
		g.aspectRatioY = float32(viewportAspectRatio / desiredAspectRatio)
	} else {
		g.aspectRatioX = float32(desiredAspectRatio / viewportAspectRatio)
		g.aspectRatioY = 1.0
	}

	gl.Viewport(0, 0, int32(width), int32(height))
	checkGLError()
}

func (g *GLRenderer) Close() {
	for _, p := range g.programs {
		if p != nil {
			p.Delete()
		}
	}

	checkGLError()
}

func (g *GLRenderer) SetScaleFactor(s int) {
	g.scaleFactor = s
	g.resizeCallback(g.window, g.screenWidth, g.screenHeight)
	checkGLError()
}

func (g *GLRenderer) SetWindowTitle(title string) {
	g.window.SetTitle(title)
	checkGLError()
}

func (g *GLRenderer) DestroyWindow() {
	g.window.Destroy()
	checkGLError()
}

func (g *GLRenderer) PrintPlatformAndVersion() {
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version:", version)
	checkGLError()
}

func (g *GLRenderer) PrintRendererInfo() {
	renderer := gl.GoStr(gl.GetString(gl.RENDERER))
	fmt.Println("Renderer:", renderer)
	checkGLError()
}

func (g *GLRenderer) Clear(c color.Color) {
	r, green, b, a := c.RGBA()
	gl.ClearColor(float32(r)/0xffff, float32(green)/0xffff, float32(b)/0xffff, float32(a)/0xffff)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	checkGLError()
}

func (g *GLRenderer) Render() {
	g.window.SwapBuffers()
	checkGLError()
}

func (g *GLRenderer) SetScreenSize(w int, h int) {
	// Shapegraphics.SetScreenSize(w, h)
	g.SetWindowSize(w, h)
	checkGLError()
}

func (g *GLRenderer) SetWindowSize(width, height int) {
	g.screenWidth = width
	g.screenHeight = height
	windowWidth = width
	windowHeight = height

	g.window.SetSize(width, height)

	g.resizeCallback(g.window, width, height)
	checkGLError()
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func (g *GLRenderer) PollEvents() bool {
	defer glfw.PollEvents()
	g.InputManager.Update()
	return !g.window.ShouldClose()
}

func (g *GLRenderer) ToggleWireframeMode() {
	g.wireframeEnabled = !g.wireframeEnabled
	if g.wireframeEnabled {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		return
	}
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
}

func checkGLError() {
	if errCode := gl.GetError(); errCode != gl.NO_ERROR {
		var errString string
		switch errCode {
		case gl.INVALID_ENUM:
			errString = "INVALID_ENUM"
		case gl.INVALID_VALUE:
			errString = "INVALID_VALUE"
		case gl.INVALID_OPERATION:
			errString = "INVALID_OPERATION"
		// Add other error codes as needed
		default:
			errString = "UNKNOWN_ERROR"
		}
		fmt.Printf("OpenGL error: %v (%s)\n", errCode, errString)
	}
}
