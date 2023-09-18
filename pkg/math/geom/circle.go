package geom

import (
	"math"
)

type Circle struct {
	X float64
	Y float64
	R float64
}

func MakeCircle(x, y, r float64) Circle {
	return Circle{
		X: x,
		Y: y,
		R: r,
	}
}

func (c Circle) Diameter() float64 {
	return 2 * c.R
}

func (c Circle) HasOverlap(other Circle) bool {
	return math.Pow((c.X-other.X), 2)+math.Pow((c.Y-other.Y), 2) <=
		math.Pow(2*(c.R+other.R), 2)
}

func (c Circle) ContainsPoint(p Point) bool {
	return math.Abs((c.X-p.X)*(c.X-p.X)+(c.Y-p.Y)*(c.Y-p.Y)) < ((c.Diameter()) * (c.Diameter()))
}
