package gl

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"unsafe"

	"github.com/dfirebaugh/ggez/pkg/graphics"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type GLRenderer struct {
	window           *glfw.Window
	programs         []*Program
	modelProgram     *Program
	textureProgram   *Program
	scaleFactor      int
	screenHeight     int
	screenWidth      int
	aspectRatioX     float32
	aspectRatioY     float32
	VAO              uint32
	wireframeEnabled bool
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

	g.ModelRenderer = NewModelRenderer(g.modelProgram)
	g.ShapeRenderer = NewShapeRenderer()
	g.InputManager = NewInputDeviceGlfw(g.window)

	checkGLError()

	return g, nil
}

func (g *GLRenderer) setupShaders() error {
	textureVertexShader, err := NewShader(BasicVert, gl.VERTEX_SHADER)
	if err != nil {
		return err
	}

	textureFragmentShader, err := NewShader(BasicFrag, gl.FRAGMENT_SHADER)
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
	if g.textureProgram != nil {
		g.textureProgram.Delete()
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
	// checkGLError()
}

func (g *GLRenderer) PrintPlatformAndVersion() {
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version:", version)
	// checkGLError()
}

func (g *GLRenderer) PrintRendererInfo() {
	renderer := gl.GoStr(gl.GetString(gl.RENDERER))
	fmt.Println("Renderer:", renderer)
	// checkGLError()
}

func (g *GLRenderer) Clear(c color.Color) {
	r, green, b, a := c.RGBA()
	gl.ClearColor(float32(r)/0xffff, float32(green)/0xffff, float32(b)/0xffff, float32(a)/0xffff)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	// checkGLError()
}

func (g *GLRenderer) Render() {
	g.window.SwapBuffers()
	// checkGLError()
}

func (g *GLRenderer) SetScreenSize(w int, h int) {
	// Shapegraphics.SetScreenSize(w, h)
	g.SetWindowSize(w, h)
	// checkGLError()
}

func (g *GLRenderer) SetWindowSize(width, height int) {
	g.screenWidth = width
	g.screenHeight = height
	windowWidth = width
	windowHeight = height

	g.window.SetSize(width, height)

	g.resizeCallback(g.window, width, height)
	// checkGLError()
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func (g *GLRenderer) CreateTextureFromImage(img image.Image) (uintptr, error) {
	// Compute the aspect ratio of the image
	imageAspectRatio := float32(img.Bounds().Dx()) / float32(img.Bounds().Dy())

	// Normalize the image dimensions to fit within the OpenGL coordinate system
	normalizedWidth := float32(img.Bounds().Dx()) / float32(g.screenWidth)
	normalizedHeight := float32(img.Bounds().Dy()) / float32(g.screenHeight)

	// Adjust normalization to maintain aspect ratio
	scaleFactor := min(normalizedWidth/imageAspectRatio, normalizedHeight)
	normalizedWidth = imageAspectRatio * scaleFactor
	normalizedHeight = scaleFactor

	vertices := []float32{
		-1, 1, 0.0, 1.0, 1.0, 1.0, 0.0, 0.0, // top left
		-1 + 2*normalizedWidth, 1, 0.0, 1.0, 1.0, 1.0, 1.0, 0.0, // top right
		-1 + 2*normalizedWidth, 1 - 2*normalizedHeight, 0.0, 1.0, 1.0, 1.0, 1.0, 1.0, // bottom right
		-1, 1 - 2*normalizedHeight, 0.0, 1.0, 1.0, 1.0, 0.0, 1.0, // bottom left
	}

	indices := []uint32{
		// rectangle
		0, 1, 2, // top triangle
		0, 2, 3, // bottom triangle
	}

	t, err := NewTexture(img)
	if err != nil {
		return 0, err // Return an error instead of panicking
	}
	if textures == nil {
		textures = make(map[uintptr]*Texture)
	}

	t.Width = img.Bounds().Dx()
	t.Height = img.Bounds().Dy()
	t.VAO = createVAO(vertices, indices)
	textures[uintptr(unsafe.Pointer(t))] = t
	checkGLError()

	// Return the texture instance as a uintptr
	return uintptr(unsafe.Pointer(t)), nil
}

func (g *GLRenderer) UpdateTextureFromImage(textureInstance uintptr, img image.Image) {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	textureID := uint32(textureInstance)
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.TexSubImage2D(
		gl.TEXTURE_2D, 0, 0, 0,
		int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y),
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix),
	)

	// Check for OpenGL errors
	if err := gl.GetError(); err != gl.NO_ERROR {
		log.Printf("OpenGL error: %s\n", err)
	}

	// Unbind the texture
	gl.BindTexture(gl.TEXTURE_2D, 0)

	gl.GenerateMipmap(gl.TEXTURE_2D)

	// Check for OpenGL errors after generating mipmaps
	if err := gl.GetError(); err != gl.NO_ERROR {
		log.Printf("OpenGL error after generating mipmaps: %s\n", err)
	}
}

func createVAO(vertices []float32, indices []uint32) uint32 {

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	var EBO uint32
	gl.GenBuffers(1, &EBO)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(VAO)

	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// copy indices into element buffer
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// size of one whole vertex (sum of attrib sizes)
	var stride int32 = 3*4 + 3*4 + 2*4
	var offset int = 0

	// position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(0)
	offset += 3 * 4

	// color
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(1)
	offset += 3 * 4

	// texture position
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(2)
	offset += 2 * 4

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray(0)

	return VAO
}
func (g *GLRenderer) RenderTexture(textureInstance uintptr, x int, y int, w int, h int, angle float32, centerX int, centerY int, flipType int) {
	g.textureProgram.Use()
	defer g.textureProgram.Delete()

	// Retrieve the texture associated with the provided textureInstance
	texture0, exists := textures[textureInstance]
	if !exists {
		return // Texture not found, exit early
	}

	// Calculate the scale based on the image's width and height
	imgScaleWidth := float32(w) / float32(texture0.Width)
	imgScaleHeight := float32(h) / float32(texture0.Height)

	// Set uniform variables in the shader program
	gl.Uniform2f(g.textureProgram.GetUniformLocation("positionOffset"), float32(x+(w/2)), float32(y+(w/2)))
	gl.Uniform1i(g.textureProgram.GetUniformLocation("windowWidth"), int32(g.screenWidth))
	gl.Uniform1i(g.textureProgram.GetUniformLocation("windowHeight"), int32(g.screenHeight))
	gl.Uniform1f(g.textureProgram.GetUniformLocation("rotationAngle"), float32(angle))
	gl.Uniform1f(g.textureProgram.GetUniformLocation("scaleWidth"), imgScaleWidth)
	gl.Uniform1f(g.textureProgram.GetUniformLocation("scaleHeight"), imgScaleHeight)
	gl.Uniform1f(g.textureProgram.GetUniformLocation("aspectRatioX"), g.aspectRatioX)
	gl.Uniform1f(g.textureProgram.GetUniformLocation("aspectRatioY"), g.aspectRatioY)

	aspectRatio := float32(w) / float32(h)
	location := g.textureProgram.GetUniformLocation("desiredAspectRatio")
	gl.Uniform1f(location, aspectRatio)

	// Activate the texture unit and bind the texture
	gl.ActiveTexture(gl.TEXTURE0)
	texture0.Bind(gl.TEXTURE0)
	defer texture0.UnBind()

	// Bind the VAO associated with the texture
	gl.BindVertexArray(texture0.VAO)
	defer gl.BindVertexArray(0) // Unbind the VAO after rendering

	// Draw the texture using triangles
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, unsafe.Pointer(nil))
}

func (g *GLRenderer) DestroyTexture(textureInstance uintptr) {
	textureID := uint32(textureInstance)
	gl.DeleteTextures(1, &textureID)
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
