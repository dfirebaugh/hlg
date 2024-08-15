package draw

import "image/color"

type Drawable interface {
	Draw(d displayer, clr color.Color)
}

type Fillableable interface {
	Fill(d displayer, clr color.Color)
}

func Draw[T Drawable](shape T, d displayer, clr color.Color) {
	shape.Draw(d, clr)
}

func Fill[T Fillableable](shape T, d displayer, clr color.Color) {
	shape.Fill(d, clr)
}
