package main

import (
	"bytes"
	"image"
	"math"
	"strings"
	"time"

	"github.com/dfirebaugh/ggez"
	"github.com/dfirebaugh/ggez/assets"
	"github.com/dfirebaugh/ggez/pkg/load"
	"github.com/dfirebaugh/ggez/pkg/math/geom"
	"golang.org/x/image/colornames"
)

var (
	teapot                         *geom.Model
	t                              *ggez.Texture
	mouseX, mouseY                 int
	rotationSpeedX, rotationSpeedY float32
	zoomSpeed                      float32
	lastUpdateTime                 time.Time
)

const (
	rotationDecayRate = 0.10
	minRotationSpeed  = 0.05
	maxRotationSpeed  = 1.0
	zoomRate          = 0.005
)

func update() {
	ggez.Clear(colornames.Black)
	ggez.DrawModel(teapot, t)

	if ggez.IsKeyJustPressed(ggez.KeyZ) {
		ggez.ToggleWireFrame()
	}

	if !ggez.IsButtonPressed(ggez.MouseButton1) {
		elapsed := time.Since(lastUpdateTime).Seconds()
		rotationSpeedX *= float32(math.Pow(rotationDecayRate, elapsed))
		rotationSpeedY *= float32(math.Pow(rotationDecayRate, elapsed))
	}

	if ggez.IsButtonPressed(ggez.MouseButton1) {
		prevX, prevY := ggez.GetCursorPosition()
		rotationSpeedX = float32(mouseX - prevX)
		rotationSpeedY = float32(mouseY - prevY)
		mouseX, mouseY = prevX, prevY
	}

	if zoomSpeed != 0 {
		teapot.ScaleFactor += zoomSpeed
		zoomSpeed = 0
	}

	teapot.Rotate(rotationSpeedX, geom.Vector3D{0, 1, 0})
	teapot.Rotate(rotationSpeedY, geom.Vector3D{1, 0, 0})
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
	rotationSpeedX, rotationSpeedY = 0, 0

	ggez.SetScrollCallback(func(_, y float64) {
		zoomSpeed += float32(y) * zoomRate
	})

	ggez.Update(update)
}
