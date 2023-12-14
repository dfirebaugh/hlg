package ggez

import (
	"image"

	"github.com/dfirebaugh/ggez/graphics"
)

type Texture struct {
	graphics.Texture
}

func CreateTexture(x, y, w, h int) (*Texture, error) {
	ensureSetupCompletion()
	img := image.NewRGBA(image.Rect(x, y, w, h))
	return CreateTextureFromImage(img)
}

func CreateTextureFromImage(img image.Image) (*Texture, error) {
	ensureSetupCompletion()
	var err error
	t := Texture{}
	t.Texture, err = ggez.graphicsBackend.CreateTextureFromImage(img)
	return &t, err
}

func (t *Texture) UpdateTextureFromImage(img image.Image) {
	ggez.graphicsBackend.UpdateTextureFromImage(t, img)
}

// Destroy removes the texture from the renderer
func (t Texture) Destroy() {
	ensureSetupCompletion()
	ggez.graphicsBackend.DisposeTexture(t.Handle())
}
