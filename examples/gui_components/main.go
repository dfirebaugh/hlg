package main

import (
	"fmt"
	"image/color"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/gui"
)

const (
	screenWidth  = 900
	screenHeight = 550
)

var (
	bgDark      = color.RGBA{24, 24, 27, 255}
	bgCard      = color.RGBA{39, 39, 42, 255}
	textPrimary = color.RGBA{244, 244, 245, 255}
	textMuted   = color.RGBA{161, 161, 170, 255}
	accent      = color.RGBA{99, 102, 241, 255}

	statusText = "Ready"
)

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetScreenSize(screenWidth, screenHeight)
	hlg.SetTitle("GUI Components Demo")
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

	// State for widgets (imgui style - caller owns state)
	sliderValues := []float32{50, 5}
	toggleStates := []bool{true, false, false}
	checkboxStates := []bool{true, false, false}

	// InputText state
	var inputText string
	var inputState gui.TextInputState

	// Panel widget state
	var panelToggle bool
	var panelSlider float32 = 75

	hlg.Run(func() {
		inputCtx.Update()
	}, func() {
		hlg.Clear(bgDark)

		ctx.Begin()

		// Section: Labels
		ctx.LabelWithSize("Labels & Text", 24, 24, 16)
		ctx.Label("This is a simple label with default styling", 24, 50)
		ctx.LabelWithColor("This label has custom color", 24, 70, textMuted)

		// Section: Buttons
		ctx.LabelWithSize("Buttons", 24, 110, 16)
		buttonLabels := []string{"Save", "Cancel", "Delete"}
		for i, lbl := range buttonLabels {
			if ctx.Button(lbl, 24+i*110, 140, 100, 36) {
				statusText = fmt.Sprintf("Button '%s' clicked", lbl)
			}
		}

		// Section: Sliders
		ctx.LabelWithSize("Sliders", 24, 200, 16)
		hlg.Text("Volume:", 24, 230, 12, textPrimary)
		ctx.Slider("volume", &sliderValues[0], 0, 100, 24, 250, 200, 8)
		hlg.Text(fmt.Sprintf("%.0f", sliderValues[0]), 230, 250, 12, textMuted)

		hlg.Text("Level:", 24, 280, 12, textPrimary)
		ctx.Slider("level", &sliderValues[1], 0, 10, 24, 300, 200, 8)
		hlg.Text(fmt.Sprintf("%.0f", sliderValues[1]), 230, 300, 12, textMuted)

		// Section: Toggles
		ctx.LabelWithSize("Toggles", 280, 200, 16)
		toggleLabels := []string{"Notifications", "Dark Mode", "Disabled"}
		for i, lbl := range toggleLabels {
			ctx.Toggle(lbl, &toggleStates[i], 280, 230+i*35, 44, 24)
			hlg.Text(lbl, 332, 235+i*35, 12, textPrimary)
		}

		// Section: Checkboxes
		ctx.LabelWithSize("Checkboxes", 480, 200, 16)
		checkLabels := []string{"Option A", "Option B", "Option C"}
		for i, lbl := range checkLabels {
			if ctx.Checkbox(lbl, &checkboxStates[i], 480, 230+i*30, 20) {
				statusText = fmt.Sprintf("Checkbox '%s' changed to %v", lbl, checkboxStates[i])
			}
		}

		// Section: Input
		ctx.LabelWithSize("Input Fields", 24, 350, 16)
		hlg.Text("Name:", 24, 385, 12, textPrimary)
		changed, submitted := ctx.InputText("name", &inputText, &inputState, 80, 380, 250, 30)
		if changed {
			statusText = fmt.Sprintf("Input changed: %s", inputText)
		}
		if submitted {
			statusText = fmt.Sprintf("Input submitted: %s", inputText)
		}

		// Section: Panel (static background)
		ctx.LabelWithSize("Panel", 400, 350, 16)
		hlg.RoundedRect(400, 380, 280, 120, 8, bgCard)
		hlg.Text("Settings Panel", 412, 392, 14, textPrimary)

		ctx.Toggle("panel_feature", &panelToggle, 412, 420, 40, 22)
		hlg.Text("Enable Feature", 460, 424, 12, textPrimary)

		ctx.Slider("panel_slider", &panelSlider, 0, 100, 412, 460, 200, 8)

		// Status bar
		hlg.FilledRect(0, screenHeight-32, screenWidth, 32, bgCard)
		hlg.Text(statusText, 24, screenHeight-22, 12, textMuted)

		ctx.End()
	})

	font.Dispose()
}
