package draw

import "image/color"

type displayer interface {
	SetPixel(x, y int16, c color.RGBA)
	Display() error
	Size() (int16, int16)
}
