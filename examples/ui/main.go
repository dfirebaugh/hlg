package main

import (
	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/pkg/ui"
	"github.com/dfirebaugh/hlg/pkg/ui/components"
)

func main() {
	screenWidth := 800
	screenHeight := 600
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetTitle("imgui test")
	hlg.EnableFPS()

	a := ui.NewApp(screenWidth, screenHeight)

	generateSurface(a, screenWidth/2, screenHeight/2)
	generateSurface(a, 200, 200)
	generateSurface(a, 0, 0)

	hlg.RunApp(a)
}

func generateSurface(app *ui.App, x, y int) *ui.Surface {
	s := app.CreateSurface(260, 300)
	s.Move(x, y)

	s.Add(&components.SurfaceHandle{
		Width:  260,
		Height: 30,
		Label:  "hello",
	})

	s.Add(&components.Button{
		X:      10,
		Y:      40,
		Width:  100,
		Height: 30,
		Text:   "Click Me",
	})

	s.Add(&components.Toggle{
		X:      120,
		Y:      45,
		Width:  50,
		Height: 20,
		IsOn:   false,
	})

	s.Add(&components.Checkbox{
		X:         10,
		Y:         95,
		Size:      20,
		IsChecked: false,
	})

	s.Add(&components.InputBox{
		X:         40,
		Y:         90,
		Width:     200,
		Height:    30,
		Text:      "",
		IsFocused: false,
	})

	s.Add(&components.ValueSlider{
		X:        10,
		Y:        130,
		Width:    240,
		Height:   20,
		Value:    0.5,
		MinValue: 0.0,
		MaxValue: 100.0,
	})

	s.Add(&components.Slider{
		X:      10,
		Y:      160,
		Width:  240,
		Height: 20,
		Value:  0.5,
	})

	s.Add(&components.Spinner{
		X:      50,
		Y:      210,
		Radius: 20,
		Color:  nil,
		Speed:  1.0,
	})

	return s
}
