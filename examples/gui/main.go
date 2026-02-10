package main

import (
	"fmt"
	"image/color"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/gui"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 840
	screenHeight = 700
)

func main() {
	// hlg.SetBackend(hlg.BackendWebGPU)
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetScreenSize(screenWidth, screenHeight)
	hlg.SetTitle("UI Components Example")
	hlg.EnableFPS()

	font, err := hlg.LoadDefaultFont()
	if err != nil {
		fmt.Printf("Failed to load font: %v\n", err)
		return
	}
	font.SetAsActiveAtlas()
	hlg.SetDefaultFont(font)

	// Create gui context
	inputCtx := gui.NewDefaultInputContext()
	ctx := gui.NewContext(inputCtx)

	// State for sliders and toggles (imgui style - caller owns state)
	sliderValues := make([]float32, 10)
	toggleStates := make([]bool, 10)

	rectColors := []color.Color{
		colornames.Mediumvioletred,
		colornames.Darkorange,
		colornames.Gold,
		colornames.Mediumseagreen,
	}

	hlg.Run(func() {
		inputCtx.Update()
	}, func() {
		hlg.Clear(colornames.Gainsboro)

		ctx.Begin()

		// Buttons
		for i := range 16 {
			x := 32 + (i%4)*200
			y := 20 + (i/4)*60
			if ctx.Button(fmt.Sprintf("Btn %d", i+1), x, y, 120, 40) {
				fmt.Printf("Button %d clicked\n", i+1)
			}
		}

		// Sliders
		for i := range 10 {
			x := 32 + (i%5)*160
			y := 280 + (i/5)*46
			if ctx.Slider(fmt.Sprintf("slider_%d", i), &sliderValues[i], 0, 100, x, y, 100, 10) {
				fmt.Printf("Slider %d value: %f\n", i+1, sliderValues[i])
			}
		}

		// Toggles
		for i := range 10 {
			x := 32 + (i%5)*160
			y := 360 + (i/5)*40
			if ctx.Toggle(fmt.Sprintf("toggle_%d", i), &toggleStates[i], x, y, 40, 20) {
				fmt.Printf("Toggle %d state: %t\n", i+1, toggleStates[i])
			}
		}

		// Colored rectangles with rounded corners and outline
		for i := range 8 {
			x := 32 + (i%4)*200
			y := 460 + (i/4)*60
			c := rectColors[i%4]
			hlg.RoundedRectOutline(x, y, 120, 60, 15, 2, c, colornames.White)
		}

		ctx.End()
	})

	font.Dispose()
}
