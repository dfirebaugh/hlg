package main

import (
	"fmt"
	"image/color"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/gui"
)

var (
	bgDark   = color.RGBA{30, 30, 35, 255}
	textMute = color.RGBA{140, 140, 150, 255}
)

func main() {
	screenWidth := 800
	screenHeight := 600
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetScreenSize(screenWidth, screenHeight)
	hlg.SetTitle("GUI Demo")
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

	// State for three panels (imgui style - caller owns state)
	type WidgetState struct {
		toggleOn     bool
		checkboxOn   bool
		sliderValue1 float32
		sliderValue2 float32
		inputText    string
		inputState   gui.TextInputState
	}

	widgetStates := []WidgetState{
		{sliderValue1: 50, sliderValue2: 25},
		{sliderValue1: 50, sliderValue2: 25},
		{sliderValue1: 50, sliderValue2: 25},
	}

	// Panel states for draggable panels
	panelStates := []gui.PanelState{
		{X: 20, Y: 20},
		{X: 290, Y: 20},
		{X: 560, Y: 20},
	}

	hlg.Run(func() {
		inputCtx.Update()
	}, func() {
		hlg.Clear(bgDark)

		ctx.Begin()

		// Render three draggable panels
		for i := 0; i < 3; i++ {
			pw, ph := 250, 320
			title := fmt.Sprintf("Panel %d", i+1)

			if ctx.Panel(title, &panelStates[i], pw, ph) {
				px, py := panelStates[i].X, panelStates[i].Y

				// Button
				if ctx.Button(fmt.Sprintf("Click Me##%d", i), px+10, py+40, 100, 30) {
					fmt.Printf("%s: Button clicked\n", title)
				}

				// Toggle
				ctx.Toggle(fmt.Sprintf("toggle_%d", i), &widgetStates[i].toggleOn, px+130, py+45, 50, 20)
				hlg.Text("Toggle", px+188, py+48, 12, textMute)

				// Checkbox
				if ctx.Checkbox(fmt.Sprintf("Enable feature##%d", i), &widgetStates[i].checkboxOn, px+10, py+85, 20) {
					fmt.Printf("%s: Checkbox is now %v\n", title, widgetStates[i].checkboxOn)
				}

				// Input box
				hlg.Text("Input:", px+10, py+118, 12, textMute)
				changed, submitted := ctx.InputText(fmt.Sprintf("input_%d", i), &widgetStates[i].inputText, &widgetStates[i].inputState, px+10, py+135, 220, 28)
				if changed {
					fmt.Printf("%s: Input changed: %s\n", title, widgetStates[i].inputText)
				}
				if submitted {
					fmt.Printf("%s: Submitted: %s\n", title, widgetStates[i].inputText)
				}

				// Slider 1 with value display
				hlg.Text("Slider 1:", px+10, py+180, 12, textMute)
				if ctx.Slider(fmt.Sprintf("slider1_%d", i), &widgetStates[i].sliderValue1, 0, 100, px+10, py+200, 180, 12) {
					fmt.Printf("%s: Slider 1 value: %.0f\n", title, widgetStates[i].sliderValue1)
				}
				hlg.Text(fmt.Sprintf("%.0f", widgetStates[i].sliderValue1), px+200, py+198, 12, textMute)

				// Slider 2
				hlg.Text("Slider 2:", px+10, py+230, 12, textMute)
				ctx.Slider(fmt.Sprintf("slider2_%d", i), &widgetStates[i].sliderValue2, 0, 100, px+10, py+250, 180, 12)
				hlg.Text(fmt.Sprintf("%.0f", widgetStates[i].sliderValue2), px+200, py+248, 12, textMute)

				// Info text
				hlg.Text("Drag title bar to move", px+10, py+290, 11, textMute)
			}
		}

		ctx.End()
	})

	font.Dispose()
}
