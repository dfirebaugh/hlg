package main

import (
	"fmt"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/gui"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

// AppState holds all the state for our immediate mode GUI demo.
// In immediate mode, the caller owns all state.
type AppState struct {
	// Simple widget state (just values)
	Volume       float32
	Brightness   float32
	Fullscreen   bool
	VSync        bool
	ShowFPS      bool
	MusicEnabled bool

	// Counter for button demo
	Counter int

	// Text input state (needs the TextInputState struct)
	PlayerName string
	NameInput  gui.TextInputState
}

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetScreenSize(screenWidth, screenHeight)
	hlg.SetTitle("Immediate Mode GUI Demo")
	hlg.EnableFPS()

	font, err := hlg.LoadDefaultFont()
	if err != nil {
		fmt.Printf("Failed to load font: %v\n", err)
		return
	}
	font.SetAsActiveAtlas()
	hlg.SetDefaultFont(font)

	// Initialize our application state
	state := &AppState{
		Volume:       0.5,
		Brightness:   0.8,
		Fullscreen:   false,
		VSync:        true,
		ShowFPS:      true,
		MusicEnabled: true,
		PlayerName:   "",
	}

	// Create our input context wrapper
	inputCtx := gui.NewDefaultInputContext()

	// Create the gui context once (it's reused each frame)
	ctx := gui.NewContext(inputCtx)

	hlg.Run(func() {
		// Update input context
		inputCtx.Update()
	}, func() {
		hlg.Clear(colornames.Darkslategray)

		// Begin immediate mode frame
		ctx.Begin()

		// Section: Title
		ctx.LabelWithSize("Immediate Mode GUI Demo", 20, 20, 18)

		// Section: Buttons
		ctx.Label("Buttons:", 20, 60)

		if ctx.Button("Click Me!", 20, 85, 100, 35) {
			state.Counter++
			fmt.Printf("Button clicked! Counter: %d\n", state.Counter)
		}

		if ctx.Button("Reset", 130, 85, 80, 35) {
			state.Counter = 0
			fmt.Println("Counter reset!")
		}

		ctx.Label(fmt.Sprintf("Counter: %d", state.Counter), 220, 95)

		// Section: Checkboxes
		ctx.Label("Checkboxes:", 20, 140)

		if ctx.Checkbox("Show FPS", &state.ShowFPS, 20, 165, 20) {
			fmt.Printf("Show FPS: %v\n", state.ShowFPS)
		}

		if ctx.Checkbox("Music Enabled", &state.MusicEnabled, 20, 195, 20) {
			fmt.Printf("Music Enabled: %v\n", state.MusicEnabled)
		}

		// Section: Toggles
		ctx.Label("Toggles:", 250, 140)

		if ctx.Toggle("Fullscreen", &state.Fullscreen, 250, 165, 44, 24) {
			fmt.Printf("Fullscreen: %v\n", state.Fullscreen)
		}

		if ctx.Toggle("VSync", &state.VSync, 250, 200, 44, 24) {
			fmt.Printf("VSync: %v\n", state.VSync)
		}

		// Section: Sliders
		ctx.Label("Sliders:", 20, 245)

		ctx.Label("Volume:", 20, 275)
		if ctx.Slider("volume", &state.Volume, 0, 1, 90, 270, 200, 12) {
			fmt.Printf("Volume: %.2f\n", state.Volume)
		}
		ctx.Label(fmt.Sprintf("%.0f%%", state.Volume*100), 300, 275)

		ctx.Label("Brightness:", 20, 305)
		if ctx.Slider("brightness", &state.Brightness, 0, 1, 90, 300, 200, 12) {
			fmt.Printf("Brightness: %.2f\n", state.Brightness)
		}
		ctx.Label(fmt.Sprintf("%.0f%%", state.Brightness*100), 300, 305)

		// Section: Text Input
		ctx.Label("Text Input:", 20, 350)

		ctx.Label("Name:", 20, 380)
		changed, submitted := ctx.InputText("name", &state.PlayerName, &state.NameInput, 80, 375, 200, 30)
		if changed {
			fmt.Printf("Name changed: %s\n", state.PlayerName)
		}
		if submitted {
			fmt.Printf("Name submitted: %s\n", state.PlayerName)
		}

		// Instructions
		ctx.Label("Press Tab to navigate between widgets", 20, 430)
		ctx.Label("Press Space/Enter to activate focused widget", 20, 450)

		// End immediate mode frame (submits all drawing)
		ctx.End()
	})

	font.Dispose()
}
