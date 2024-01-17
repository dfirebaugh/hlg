package main

import (
	"image/color"

	"github.com/dfirebaugh/ggez"
	"github.com/dfirebaugh/ggez/pkg/draw"
	"github.com/dfirebaugh/ggez/pkg/fb"
	"github.com/dfirebaugh/ggez/pkg/math/geom"
)

type Screen struct {
	geom.Rect
	*fb.ImageFB
	*ggez.Texture
}

func NewScreen(width, height int) Screen {
	s := Screen{Rect: geom.MakeRect(0, 0, float32(width), float32(height))}
	s.ImageFB = fb.New(width, height)
	texture, _ := ggez.CreateTextureFromImage(s.ToImage())
	s.Texture = texture
	return s
}

func (s Screen) Clear(c color.Color) {
	draw.Rect(s.Rect).Fill(s, c)
}

func (s Screen) Render() {
	s.UpdateImage(s.ToImage())
	s.Texture.Render()
}
