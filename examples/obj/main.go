package main

import (
	"image"
	"os"
	"strings"

	"github.com/dfirebaugh/ggez"
	"github.com/dfirebaugh/ggez/assets"
	"github.com/dfirebaugh/ggez/pkg/load"
	"github.com/dfirebaugh/ggez/pkg/math/geom"
	"golang.org/x/image/colornames"
)

var (
	teapot *geom.Model
	t      *ggez.Texture
)

func update() {
	ggez.Clear(colornames.Black)
	teapot.Rotate(.3, geom.Vector3{0, 1, 0})
	ggez.DrawModel(teapot, t)

	if ggez.IsKeyJustPressed(ggez.KeyZ) {
		ggez.ToggleWireFrame()
	}
}

func main() {
	var err error

	ggez.SetScreenSize(240*2, 160*2)

	reader := strings.NewReader(assets.TeaPot)
	teapot, err = load.LoadOBJModelFromReader("teapot", reader)
	if err != nil {
		panic(err)
	}

	imgFile, err := os.Open("assets/models/the-utah-teapot/source/default.png")
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	t, err = ggez.CreateTextureFromImage(img)
	if err != nil {
		panic(err)
	}

	teapot.ScaleFactor = 0.01
	teapot.Position.Y = -.5

	ggez.Update(update)
}
