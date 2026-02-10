package main

import (
	"fmt"
	"image/color"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/gui"
)

// Clean dark theme
var (
	bgDark      = color.RGBA{24, 24, 27, 255}    // zinc-900
	bgCard      = color.RGBA{39, 39, 42, 255}    // zinc-800
	textPrimary = color.RGBA{244, 244, 245, 255} // zinc-100
	textMuted   = color.RGBA{161, 161, 170, 255} // zinc-400
	accent      = color.RGBA{99, 102, 241, 255}  // indigo-500
	accentHover = color.RGBA{129, 140, 248, 255} // indigo-400
	success     = color.RGBA{34, 197, 94, 255}   // green-500
	warning     = color.RGBA{251, 191, 36, 255}  // amber-400
)

const (
	screenWidth  = 860
	screenHeight = 540
)

var statusText = "Ready"

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetScreenSize(screenWidth, screenHeight)
	hlg.SetTitle("UI Components")
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
	buttonLabels := []string{"New", "Open", "Save", "Export", "Settings", "Help"}
	sliderLabels := []string{"Volume", "Brightness", "Speed", "Scale"}
	sliderValues := []float32{75, 50, 60, 100}
	toggleLabels := []string{"Notifications", "Auto-save", "Dark mode", "Animations"}
	toggleStates := []bool{true, false, true, false}

	// Card colors
	cardColors := []color.Color{accent, success, warning, accentHover}

	hlg.Run(func() {
		inputCtx.Update()
	}, func() {
		hlg.Clear(bgDark)

		ctx.Begin()

		// Section: Buttons
		ctx.LabelWithColor("Actions", 24, 24, textMuted)
		for i := 0; i < 6; i++ {
			x := 24 + (i%3)*140
			y := 50 + (i/3)*42
			if ctx.Button(buttonLabels[i], x, y, 124, 34) {
				statusText = fmt.Sprintf("%s clicked", buttonLabels[i])
			}
		}

		// Section: Sliders
		ctx.LabelWithColor("Adjustments", 24, 150, textMuted)
		for i := 0; i < 4; i++ {
			x := 24 + (i%2)*210
			y := 195 + (i/2)*44
			// Label
			ctx.LabelWithColor(sliderLabels[i], x, y-16, textPrimary)
			// Slider
			ctx.Slider(fmt.Sprintf("slider_%d", i), &sliderValues[i], 0, 100, x, y, 160, 8)
			// Value
			ctx.LabelWithColor(fmt.Sprintf("%.0f", sliderValues[i]), x+168, y-4, textMuted)
		}

		// Section: Toggles
		ctx.LabelWithColor("Preferences", 24, 296, textMuted)
		for i := 0; i < 4; i++ {
			x := 24 + (i%2)*210
			y := 340 + (i/2)*36
			if ctx.Toggle(fmt.Sprintf("toggle_%d", i), &toggleStates[i], x, y, 40, 22) {
				state := "off"
				if toggleStates[i] {
					state = "on"
				}
				statusText = fmt.Sprintf("%s %s", toggleLabels[i], state)
			}
			// Label
			ctx.LabelWithColor(toggleLabels[i], x+48, y+4, textPrimary)
		}

		// Section: Cards
		ctx.LabelWithColor("Status", 24, 420, textMuted)
		cardLabels := []string{"Active", "Online", "Pending", "New"}
		for i := 0; i < 4; i++ {
			x := 24 + (i%4)*105
			y := 446
			hlg.RoundedRect(x, y, 96, 56, 6, cardColors[i])
			hlg.Text(cardLabels[i], x+12, y+20, 14, bgDark)
		}

		// Status bar
		hlg.FilledRect(0, screenHeight-32, screenWidth, 32, bgCard)
		hlg.Text(statusText, 24, screenHeight-22, 12, textMuted)

		ctx.End()
	})

	font.Dispose()
}
