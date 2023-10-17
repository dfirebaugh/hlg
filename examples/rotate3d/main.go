package main

import (
	"bytes"
	"image"
	"strings"
	"time"

	"github.com/dfirebaugh/ggez"
	"github.com/dfirebaugh/ggez/assets"
	"github.com/dfirebaugh/ggez/pkg/load"
	"github.com/dfirebaugh/ggez/pkg/math/geom"
	"golang.org/x/image/colornames"
)

var (
	teapot         *geom.Model
	t              *ggez.Texture
	mouseX, mouseY int
	rotation       geom.Vector3D
	lastUpdateTime time.Time
	rotationSpeed  geom.Vector3D
)

const (
	rotationDecayRate = 0.10
	minRotationSpeed  = 0.05
	maxRotationSpeed  = 1.0
	mouseSensitivity  = 0.005
)

func update() {
	ggez.Clear(colornames.Black)
	ggez.DrawModel(teapot, t)

	if ggez.IsKeyJustPressed(ggez.KeyZ) {
		ggez.ToggleWireFrame()
	}

	newMouseX, newMouseY := ggez.GetCursorPosition()

	if ggez.IsButtonPressed(ggez.MouseButton1) {
		deltaX := float32(newMouseX-mouseX) * mouseSensitivity
		deltaY := float32(newMouseY-mouseY) * mouseSensitivity

		rotation[0] -= deltaY
		rotation[1] -= deltaX

		teapot.SetRotation(rotation)
	} else {
		rotationSpeed.Scaled(rotationDecayRate)

		deltaX := float32(newMouseX-mouseX) * mouseSensitivity
		deltaY := float32(newMouseY-mouseY) * mouseSensitivity

		rotation[0] += deltaY * rotationSpeed[0]
		rotation[1] -= deltaX * rotationSpeed[1]

		teapot.SetRotation(rotation)
	}

	mouseX, mouseY = newMouseX, newMouseY
	lastUpdateTime = time.Now()
}

func main() {
	var err error

	ggez.SetScreenSize(240*2, 160*2)

	reader := strings.NewReader(assets.TeaPot)
	teapot, err = load.LoadOBJModelFromReader("teapot", reader)
	if err != nil {
		panic(err)
	}
	imgReader := bytes.NewReader(assets.DefaultTextureImage)
	img, _, err := image.Decode(imgReader)
	if err != nil {
		panic(err)
	}

	t, err = ggez.CreateTextureFromImage(img)
	if err != nil {
		panic(err)
	}

	teapot.ScaleFactor = 0.01
	teapot.Position[1] = -0.5

	rotationSpeed = geom.Vector3D{0, 0, 0}

	ggez.Update(update)
}
