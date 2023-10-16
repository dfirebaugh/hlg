package gl

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"os"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Texture struct {
	handle  uint32
	target  uint32
	texUnit uint32
	VAO     uint32
	Width   int
	Height  int
}

var errUnsupportedStride = errors.New("unsupported stride, only 32-bit colors supported")

var errTextureNotBound = errors.New("texture not bound")

func NewTextureFromFile(file string) (*Texture, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}
	return NewTexture(img)
}

func NewTexture(img image.Image) (*Texture, error) {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	if rgba.Stride != rgba.Rect.Size().X*4 { // TODO-cs: why?
		return nil, errUnsupportedStride
	}

	var handle uint32
	gl.GenTextures(1, &handle)

	target := uint32(gl.TEXTURE_2D)
	internalFmt := int32(gl.RGBA)
	format := uint32(gl.RGBA)
	width := int32(rgba.Rect.Size().X)
	height := int32(rgba.Rect.Size().Y)
	pixType := uint32(gl.UNSIGNED_BYTE)
	dataPtr := gl.Ptr(rgba.Pix)

	texture := Texture{
		handle: handle,
		target: target,
	}

	texture.Bind(gl.TEXTURE0)
	defer texture.UnBind()

	gl.TexParameteri(texture.target, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(texture.target, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(texture.target, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(texture.target, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(target, 0, internalFmt, width, height, 0, format, pixType, dataPtr)

	// gl.GenerateMipmap(texture.handle)

	return &texture, nil
}

func (t *Texture) Clear(color color.Color) {
	rgba := image.NewRGBA(image.Rect(0, 0, t.Width, t.Height))
	draw.Draw(rgba, rgba.Bounds(), &image.Uniform{color}, image.Point{}, draw.Src)
	t.UpdateFromImage(rgba)
}

func (t *Texture) UpdateFromImage(img image.Image) error {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return errUnsupportedStride
	}

	t.Bind(gl.TEXTURE0)
	defer t.UnBind()

	width := int32(rgba.Rect.Size().X)
	height := int32(rgba.Rect.Size().Y)
	dataPtr := gl.Ptr(rgba.Pix)

	gl.TexSubImage2D(t.target, 0, 0, 0, width, height, gl.RGBA, gl.UNSIGNED_BYTE, dataPtr)

	return nil
}

func (t *Texture) Handle() uint32 {
	return t.handle
}

func (t *Texture) SetWrap(wrapR, wrapS int32) {
	gl.TexParameteri(t.target, gl.TEXTURE_WRAP_R, wrapR)
	gl.TexParameteri(t.target, gl.TEXTURE_WRAP_S, wrapS)
}

func (tex *Texture) Bind(texUnit uint32) {
	gl.ActiveTexture(texUnit)
	gl.BindTexture(tex.target, tex.handle)
	tex.texUnit = texUnit
}

func (tex *Texture) UnBind() {
	tex.texUnit = 0
	gl.BindTexture(tex.target, 0)
}

func (tex *Texture) SetUniform(uniformLoc int32) error {
	if tex.texUnit == 0 {
		return errTextureNotBound
	}
	gl.Uniform1i(uniformLoc, int32(tex.texUnit-gl.TEXTURE0))
	return nil
}
