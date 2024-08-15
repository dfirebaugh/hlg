package draw

import (
	"image/color"

	"github.com/dfirebaugh/hlg/pkg/math/geom"
)

type Cube [8]geom.Point3D

func (c Cube) Draw(d displayer, clr color.Color) {
	// Orthographically project 3D points to 2D
	projected := func(p geom.Point3D) geom.Point {
		return geom.MakePoint(p.X, p.Y)
	}

	viewDirection := geom.Point3D{0, 0, -1}

	sides := [][]geom.Point3D{
		{c[0], c[1], c[2], c[3]},
		{c[0], c[3], c[7], c[4]},
		{c[1], c[2], c[6], c[5]},
	}

	for _, side := range sides {
		normal := computeNormal(side[0], side[1], side[2])

		if dotProduct(normal, viewDirection) < 0 {
			polygon := geom.MakePolygon(projected(side[0]), projected(side[1]), projected(side[2]), projected(side[3]))
			Polygon(polygon).Draw(d, clr)
		}
	}
}

func (c Cube) Fill(d displayer, clr color.Color) {
	projected := func(p geom.Point3D) geom.Point {
		return geom.MakePoint(p.X, p.Y)
	}

	viewDirection := geom.Point3D{0, 0, -1}

	sides := [][]geom.Point3D{
		{c[0], c[1], c[2], c[3]},
		{c[0], c[3], c[7], c[4]},
		{c[1], c[2], c[6], c[5]},
	}

	for _, side := range sides {
		normal := computeNormal(side[0], side[1], side[2])

		if dotProduct(normal, viewDirection) < 0 {
			polygon := geom.MakePolygon(projected(side[0]), projected(side[1]), projected(side[2]), projected(side[3]))
			Polygon(polygon).Fill(d, clr)
		}
	}
}

func computeNormal(p1, p2, p3 geom.Point3D) geom.Point3D {
	v1 := p2.Subtract(p1)
	v2 := p3.Subtract(p1)
	return v1.Cross(v2).Normalize()
}

func dotProduct(p1, p2 geom.Point3D) float32 {
	return p1.X*p2.X + p1.Y*p2.Y + p1.Z*p2.Z
}
