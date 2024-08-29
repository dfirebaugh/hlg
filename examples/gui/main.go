package main

import (
	"fmt"
	"image/color"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/gui"
	"github.com/dfirebaugh/hlg/gui/components"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 840
	screenHeight = 700
)

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetScreenSize(screenWidth, screenHeight)
	hlg.SetTitle("UI Components Example")
	hlg.EnableFPS()

	buttons := []*components.Button{}
	for i := 0; i < 16; i++ {
		x := 32 + (i%4)*200
		y := 20 + (i/4)*60
		index := i
		button := components.NewButton(x, y, 120, 40, 5, fmt.Sprintf("Btn %d", i+1), func() {
			fmt.Printf("Button %d clicked\n", index+1)
		})
		buttons = append(buttons, button)
	}

	sliders := []*components.Slider{}
	for i := 0; i < 10; i++ {
		x := 32 + (i%5)*160
		y := 280 + (i/5)*46
		slider := components.NewSlider(x, y, 100, 10, 0, 100, func(value float32) {
			fmt.Printf("Slider %d value: %f\n", i+1, value)
		})
		sliders = append(sliders, slider)
	}

	toggles := []*components.Toggle{}
	for i := 0; i < 10; i++ {
		x := 32 + (i%5)*160
		y := 360 + (i/5)*40
		toggle := components.NewToggle(x, y, 40, 20, false, func(isOn bool) {
			fmt.Printf("Toggle %d state: %t\n", i+1, isOn)
		})
		toggles = append(toggles, toggle)
	}

	rectColors := []color.Color{
		colornames.Mediumvioletred,
		colornames.Darkorange,
		colornames.Gold,
		colornames.Mediumseagreen,
	}

	hlg.Run(func() {
		for _, button := range buttons {
			button.Update()
		}
		for _, slider := range sliders {
			slider.Update()
		}
		for _, toggle := range toggles {
			toggle.Update()
		}
	}, func() {
		hlg.Clear(colornames.Gainsboro)
		w, h := hlg.GetWindowSize()
		drawContext := gui.NewDrawContext(w, h)

		for _, button := range buttons {
			button.Render(drawContext)
		}
		for _, slider := range sliders {
			slider.Render(drawContext)
		}
		for _, toggle := range toggles {
			toggle.Render(drawContext)
		}

		for i := 0; i < 8; i++ {
			x := 32 + (i%4)*200
			y := 460 + (i/4)*60
			color := rectColors[i%4]
			drawContext.DrawRectangle(x, y, 120, 60, &gui.DrawOptions{
				Style: gui.Style{
					FillColor:    color,
					OutlineColor: colornames.White,
					OutlineSize:  2,
					CornerRadius: 6,
				},
			})
		}

		hlg.SubmitDrawBuffer(drawContext.Encode())
	})
}
