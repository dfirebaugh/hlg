package gui

import "image/color"

type Theme struct {
	BackgroundColor color.Color
	PrimaryColor    color.Color
	SecondaryColor  color.Color
	TextColor       color.Color
	HandleColor     color.Color
}

func (d *Draw) GetTheme() *Theme {
	return d.Theme
}

func (d *Draw) SetTheme(t *Theme) {
	d.Theme = t
}

func DefaultTheme() *Theme {
	return &Theme{
		BackgroundColor: color.RGBA{30, 30, 30, 255},
		PrimaryColor:    color.RGBA{60, 120, 180, 255},
		SecondaryColor:  color.RGBA{90, 90, 90, 255},
		TextColor:       color.RGBA{255, 255, 255, 255},
		HandleColor:     color.RGBA{200, 200, 200, 255},
	}
}
