package ggez

import (
	"image"
)

type Texture struct {
	ptr        uintptr
	X, Y, W, H int
	img        image.Image
	Angle      float64
	Center     struct {
		X int
		Y int
	}
	Flip   FlipType
	ScaleX float64
	ScaleY float64
}

// FlipType is a type to represent flip operations on textures
type FlipType int

const (
	FLIP_NONE FlipType = iota
	FLIP_HORIZONTAL
	FLIP_VERTICAL
	FLIP_BOTH
)

func (t *Texture) Image() image.Image {
	return t.img
}

func CreateTexture(x, y, w, h int) (*Texture, error) {
	ensureSetupCompletion()
	img := image.NewRGBA(image.Rect(x, y, w, h))
	return CreateTextureFromImage(img)
}

// CreateTextureFromImage creates an SDL texture from an image.Image.
func CreateTextureFromImage(img image.Image) (*Texture, error) {
	ensureSetupCompletion()
	var err error

	// Convert img to RGBA, this ensures the texture always works with an RGBA image.
	rgba := image.NewRGBA(img.Bounds())
	for y := 0; y < rgba.Bounds().Dy(); y++ {
		for x := 0; x < rgba.Bounds().Dx(); x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	t := Texture{}
	t.img = rgba
	t.W = rgba.Bounds().Max.X
	t.H = rgba.Bounds().Max.Y
	t.ScaleX = 1.0
	t.ScaleY = 1.0
	t.ptr, err = graphicsBackend.CreateTextureFromImage(rgba)
	return &t, err
}

func (t *Texture) SetScale(scaleX, scaleY float64) {
	t.ScaleX = scaleX
	t.ScaleY = scaleY
}

func (t *Texture) SetRotation(angle float64) {
	t.Angle = angle
}

func (t *Texture) SetFlip(flip FlipType) {
	t.Flip = flip
}

// Render renders a texture to the screen considering the scale factor.
func (t Texture) Render() {
	ensureSetupCompletion()
	renderW := int(float64(t.W) * t.ScaleX)
	renderH := int(float64(t.H) * t.ScaleY)
	graphicsBackend.RenderTexture(t.ptr, t.X, t.Y, renderW, renderH, t.Angle, t.Center.X, t.Center.Y, int(t.Flip))
}

// Destroy removes the texture from the renderer
func (t Texture) Destroy() {
	ensureSetupCompletion()
	graphicsBackend.DestroyTexture(t.ptr)
}
