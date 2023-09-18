package component

import "math"

type Position struct {
	X, Y float64
}

type Size struct {
	Width, Height float64
}

type Rotation float64

func (r Rotation) Cos() float64 {
	return math.Cos(float64(r))
}

func (r Rotation) Sin() float64 {
	return math.Sin(float64(r))
}

func (r Rotation) Rad() float64 {
	return float64(r) * math.Pi
}
