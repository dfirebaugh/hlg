package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/assets"
	"github.com/dfirebaugh/hlg/gui"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetScreenSize(screenWidth, screenHeight)
	hlg.SetTitle("Texture Controls Example")
	hlg.EnableFPS()
	hlg.SetVSync(true)

	// Load font for labels
	font, err := hlg.LoadDefaultFont()
	if err != nil {
		fmt.Println("Warning: Could not load font:", err)
	} else {
		font.SetAsActiveAtlas()
		hlg.SetDefaultFont(font)
	}

	// Create gui context
	inputCtx := gui.NewDefaultInputContext()
	ctx := gui.NewContext(inputCtx)

	// Create a separate render queue for the texture (renders behind UI when Draw() is called first)
	textureRenderQueue := hlg.CreateRenderQueue()

	// Load texture from embedded data
	img, _, err := image.Decode(bytes.NewReader(assets.RTS_Crate))
	if err != nil {
		panic(err)
	}
	texture, err := hlg.CreateTextureFromImage(img)
	if err != nil {
		panic(err)
	}

	// Get original texture dimensions
	bounds := img.Bounds()
	origWidth := float32(bounds.Dx())
	origHeight := float32(bounds.Dy())

	// Texture control values (imgui style - caller owns state)
	var (
		scaleX float32 = 1.0
		scaleY float32 = 1.0
		posX   float32 = 150
		posY   float32 = 150
		flipH  bool    = false
		flipV  bool    = false
	)

	// Track previous values to detect changes
	var (
		prevScaleX, prevScaleY float32 = scaleX, scaleY
		prevPosX, prevPosY     float32 = posX, posY
		prevFlipH, prevFlipV   bool    = flipH, flipV
	)

	// Control panel position
	panelX := 560

	// Apply initial texture transforms
	texture.Resize(origWidth*scaleX, origHeight*scaleY)
	texture.Move(posX, posY)
	texture.SetFlipHorizontal(flipH)
	texture.SetFlipVertical(flipV)

	hlg.Run(func() {
		inputCtx.Update()

		// Only apply transforms when values change
		scaleChanged := scaleX != prevScaleX || scaleY != prevScaleY
		posChanged := posX != prevPosX || posY != prevPosY
		flipChanged := flipH != prevFlipH || flipV != prevFlipV

		if scaleChanged {
			texture.Resize(origWidth*scaleX, origHeight*scaleY)
			prevScaleX, prevScaleY = scaleX, scaleY
		}

		if posChanged || scaleChanged {
			texture.Move(posX, posY)
			prevPosX, prevPosY = posX, posY
		}

		if flipChanged {
			texture.SetFlipHorizontal(flipH)
			texture.SetFlipVertical(flipV)
			prevFlipH, prevFlipV = flipH, flipV
		}
	}, func() {
		hlg.Clear(color.RGBA{40, 44, 52, 255})

		// Add texture to queue and present it first (renders behind UI)
		texture.RenderToQueue(textureRenderQueue)
		textureRenderQueue.Present()

		ctx.Begin()

		// Panel background
		hlg.RoundedRect(550, 10, 240, 510, 10, color.RGBA{60, 63, 70, 255})

		// Panel title
		hlg.Text("Texture Controls", panelX, 20, 18, colornames.White)

		// Scale X slider
		hlg.Text(fmt.Sprintf("Scale X: %.2f", scaleX), panelX, 50, 14, colornames.White)
		ctx.Slider("scaleX", &scaleX, 0.1, 3.0, panelX, 70, 200, 12)

		// Scale Y slider
		hlg.Text(fmt.Sprintf("Scale Y: %.2f", scaleY), panelX, 110, 14, colornames.White)
		ctx.Slider("scaleY", &scaleY, 0.1, 3.0, panelX, 130, 200, 12)

		// Position X slider
		hlg.Text(fmt.Sprintf("Position X: %.0f", posX), panelX, 170, 14, colornames.White)
		ctx.Slider("posX", &posX, 0, 400, panelX, 190, 200, 12)

		// Position Y slider
		hlg.Text(fmt.Sprintf("Position Y: %.0f", posY), panelX, 230, 14, colornames.White)
		ctx.Slider("posY", &posY, 0, 350, panelX, 250, 200, 12)

		// Flip Horizontal toggle
		ctx.Toggle("flipH", &flipH, panelX, 310, 50, 24)
		hlg.Text("Flip Horizontal", panelX+60, 314, 14, colornames.White)

		// Flip Vertical toggle
		ctx.Toggle("flipV", &flipV, panelX, 360, 50, 24)
		hlg.Text("Flip Vertical", panelX+60, 364, 14, colornames.White)

		// Reset button
		if ctx.Button("Reset", panelX, 430, 120, 36) {
			scaleX, scaleY = 1.0, 1.0
			posX, posY = 150, 150
			flipH, flipV = false, false
		}

		// Info text
		hlg.Text("Drag sliders to modify", panelX, 480, 12, colornames.Gray)

		ctx.End()
	})
}
