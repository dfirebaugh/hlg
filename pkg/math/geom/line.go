package geom

import (
	"errors"
)

type Segment struct {
	V0 Vector
	V1 Vector
}

type Line struct {
	slope   float32
	yint    float32
	Segment Segment
}

func MakeLine(a, b Point) Line {
	slope := (b.Y - a.Y) / (b.X - a.X)
	yint := a.Y - slope*a.X
	return Line{
		slope: slope,
		yint:  yint,
		Segment: Segment{
			V0: MakeVector(a.X, a.Y),
			V1: MakeVector(b.X, b.Y),
		},
	}
}

func (l Line) EvalX(x float32) float32 {
	return l.slope*x + l.yint
}

func (l Line) IsParrallel(l1, l2 Line) bool {
	return l1.slope == l2.slope
}

func (l Line) Intersection(l2 Line) (Point, error) {
	if l.slope == l2.slope {
		return Point{}, errors.New("the lines do not intersect")
	}
	x := (l2.yint - l.yint) / (l.slope - l2.slope)
	y := l.EvalX(x)

	return Point{x, y}, nil
}
