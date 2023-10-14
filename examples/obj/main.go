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
	teapot         *geom.Model
	t              *ggez.Texture
	mouseX         int
	rotationSpeed  float32
	lastUpdateTime time.Time
)

const (
	rotationDecayRate = 0.10
	minRotationSpeed  = 0.05
	maxRotationSpeed  = 1.0
)

func update() {
	ggez.Clear(colornames.Black)
	ggez.DrawModel(teapot, t)

	if ggez.IsKeyJustPressed(ggez.KeyZ) {
		ggez.ToggleWireFrame()
	}

	if !ggez.IsButtonPressed(ggez.MouseButton1) {
		// Gradually decrease the rotation speed
		elapsed := time.Since(lastUpdateTime).Seconds()
		rotationSpeed *= float32(math.Pow(rotationDecayRate, elapsed))
		if math.Abs(float64(rotationSpeed)) < minRotationSpeed {
			rotationSpeed = minRotationSpeed
		}
	}

	if ggez.IsButtonPressed(ggez.MouseButton1) {
		prevX, _ := ggez.GetCursorPosition()
		if prevX < mouseX {
			rotationSpeed = float32(mouseX - prevX)
		}
		if prevX > mouseX {
			rotationSpeed = -float32(prevX - mouseX)
		}
		mouseX = prevX
	}

	teapot.Rotate(rotationSpeed, geom.Vector3{0, 1, 0})
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
	teapot.Position.Y = -0.5

	rotationSpeed = 0

	ggez.Update(update)
}
