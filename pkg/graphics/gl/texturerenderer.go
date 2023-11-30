package gl

import (
	"image"
	"image/color"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

func (g *GLRenderer) CreateTextureFromImage(img image.Image) (uintptr, error) {
	imageAspectRatio := float32(img.Bounds().Dx()) / float32(img.Bounds().Dy())

	normalizedWidth := float32(img.Bounds().Dx()) / float32(g.screenWidth)
	normalizedHeight := float32(img.Bounds().Dy()) / float32(g.screenHeight)

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
		return 0, err
	}
	if textures == nil {
		textures = make(map[uintptr]*Texture)
	}

	t.Width = img.Bounds().Dx()
	t.Height = img.Bounds().Dy()
	t.VAO = g.createVAO(vertices, indices)
	textures[uintptr(unsafe.Pointer(t))] = t
	checkGLError()

	return uintptr(unsafe.Pointer(t)), nil
}

func (g *GLRenderer) UpdateTextureFromImage(textureInstance uintptr, img image.Image) {
	textures[textureInstance].UpdateFromImage(img)
}

func (g *GLRenderer) ClearTexture(textureInstance uintptr, c color.Color) {
	textures[textureInstance].Clear(c)
}

func (g *GLRenderer) createVAO(vertices []float32, indices []uint32) uint32 {

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	var EBO uint32
	gl.GenBuffers(1, &EBO)

	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	var stride int32 = 3*4 + 3*4 + 2*4
	var offset int = 0

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, stride, uintptr(offset))
	gl.EnableVertexAttribArray(0)
	offset += 3 * 4

	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, stride, uintptr(offset))
	gl.EnableVertexAttribArray(1)
	offset += 3 * 4

	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, stride, uintptr(offset))
	gl.EnableVertexAttribArray(2)
	offset += 2 * 4

	gl.BindVertexArray(0)

	return VAO
}

func (g *GLRenderer) RenderTexture(textureInstance uintptr, x int, y int, w int, h int, angle float32, centerX int, centerY int, flipType int) {
	g.textureProgram.Use()

	texture0, exists := textures[textureInstance]
	if !exists {
		return
	}

	imgScaleWidth := float32(w) / float32(texture0.Width)
	imgScaleHeight := float32(h) / float32(texture0.Height)

	gl.Uniform2f(g.textureProgram.GetUniformLocation("positionOffset"), float32(x+(w/2)), float32(y+(w/2)))
	gl.Uniform1i(g.textureProgram.GetUniformLocation("windowWidth"), int32(g.screenWidth))
	gl.Uniform1i(g.textureProgram.GetUniformLocation("windowHeight"), int32(g.screenHeight))
	gl.Uniform1f(g.textureProgram.GetUniformLocation("rotationAngle"), float32(angle))
	gl.Uniform1f(g.textureProgram.GetUniformLocation("scaleWidth"), imgScaleWidth)
	gl.Uniform1f(g.textureProgram.GetUniformLocation("scaleHeight"), imgScaleHeight)
	gl.Uniform1f(g.textureProgram.GetUniformLocation("aspectRatioX"), g.aspectRatioX)
	gl.Uniform1f(g.textureProgram.GetUniformLocation("aspectRatioY"), g.aspectRatioY)

	aspectRatio := float32(w) / float32(h)
	location := g.program.GetUniformLocation("desiredAspectRatio")
	gl.Uniform1f(location, aspectRatio)

	gl.ActiveTexture(gl.TEXTURE0)
	texture0.Bind(gl.TEXTURE0)
	defer texture0.UnBind()

	gl.BindVertexArray(texture0.VAO)
	defer gl.BindVertexArray(0)

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, unsafe.Pointer(nil))
}

func (g *GLRenderer) DestroyTexture(textureInstance uintptr) {
	textureID := uint32(textureInstance)
	gl.DeleteTextures(1, &textureID)
}
