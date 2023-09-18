package geom

type Point struct {
	X float64
	Y float64
}

func MakePoint(x, y float64) Point {
	return Point{X: x, Y: y}
}

func (p Point) ToVector() Vector {
	return MakeVector(p.X, p.Y)
}
