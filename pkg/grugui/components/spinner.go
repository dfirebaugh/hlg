package components

import (
	"image/color"
	"math"
	"time"

	ui "github.com/dfirebaugh/hlg/pkg/grugui"
)

type Spinner struct {
	X, Y, Radius int
	Color        color.Color
	Speed        float64
	startTime    time.Time
	angle        float64
}

func NewSpinner(x, y, radius int, color color.Color) *Spinner {
	return &Spinner{
		X:         x,
		Y:         y,
		Radius:    radius,
		Color:     color,
		Speed:     150,
		startTime: time.Now(),
	}
}

func (s *Spinner) Update(ctx ui.Context) {
	elapsed := time.Since(s.startTime).Seconds()
	s.startTime = time.Now()

	s.angle += elapsed * s.Speed

	if s.angle > 2*math.Pi {
		s.angle -= 2 * math.Pi
	}
}

func (s *Spinner) Render(ctx ui.Context) {
	theme := ctx.Theme()
	spinnerColor := s.Color
	if spinnerColor == nil {
		spinnerColor = theme.TextColor
	}

	ctx.FillRect(s.X-s.Radius, s.Y-s.Radius, s.Radius*2, s.Radius*2, theme.BackgroundColor)

	for i := 0; i < 12; i++ {
		theta := s.angle + float64(i)*math.Pi/6
		x := s.X + int(float64(s.Radius)*math.Cos(theta))
		y := s.Y + int(float64(s.Radius)*math.Sin(theta))

		ctx.FillCircle(x, y, s.Radius/5, spinnerColor)
	}
}

func applyAlpha(c color.Color, alpha float64) color.Color {
	r, g, b, a := c.RGBA()
	a = uint32(float64(a) * alpha)
	return color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
}
