package geom

import (
	"math"
)

type Circle struct {
	X float32
	Y float32
	R float32
}

func MakeCircle(x, y, r float32) Circle {
	return Circle{
		X: x,
		Y: y,
		R: r,
	}
}

func (c Circle) Diameter() float32 {
	return 2 * c.R
}

func (c Circle) HasOverlap(other Circle) bool {
	return math.Pow(float64(c.X-other.X), 2)+math.Pow(float64(c.Y-other.Y), 2) <=
		math.Pow(2*float64(c.R+other.R), 2)
}

func (c Circle) ContainsPoint(p Point) bool {
	return math.Abs(float64((c.X-p.X)*(c.X-p.X)+(c.Y-p.Y)*(c.Y-p.Y))) < float64((c.Diameter())*(c.Diameter()))
}
