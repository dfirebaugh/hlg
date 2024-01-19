package component

import (
	"fmt"

	"github.com/dfirebaugh/hlg/pkg/math/geom"
)

type Coordinate struct {
	X float32
	Y float32
}

func (c Coordinate) String() string {
	return fmt.Sprintf("%d, %d", int(c.X), int(c.Y))
}

func (c *Coordinate) SetCoordinate(newCoord Coordinate) {
	c.X = newCoord.X
	c.Y = newCoord.Y
}

func (c Coordinate) GetDistance(other Coordinate) float32 {
	a := geom.MakeVector(c.X, c.Y)
	b := geom.MakeVector(other.X, other.Y)
	return a.GetDistance(b)
}

func (c Coordinate) GetDirection(other Coordinate) float32 {
	a := geom.MakeVector(c.X, c.Y)
	b := geom.MakeVector(other.X, other.Y)
	return a.GetDirection(b)
}

func (c Coordinate) Add(other Coordinate) Coordinate {
	return Coordinate{
		X: c.X + other.X,
		Y: c.Y + other.Y,
	}
}

func (c Coordinate) Subtract(other Coordinate) Coordinate {
	return Coordinate{
		X: c.X - other.X,
		Y: c.Y - other.Y,
	}
}

func (c Coordinate) TranslateXY(offset Coordinate, pixelSize float32) (float32, float32) {
	x := (c.X - offset.X) / pixelSize
	y := (c.Y - offset.Y) / pixelSize
	return x, y
}
