package draw

import (
	"image/color"
	"sort"

	"github.com/dfirebaugh/hlg/pkg/math/geom"
)

type Polygon geom.Polygon

func (p Polygon) Draw(d displayer, clr color.Color) {
	numPoints := len(p)
	for i := 0; i < numPoints; i++ {
		start := p[i]
		end := p[(i+1)%numPoints]

		line := Line(geom.MakeLine(start, end))
		line.Draw(d, clr)
	}
}

func (p Polygon) Fill(d displayer, clr color.Color) {
	minY, maxY := p[0].Y, p[0].Y
	for _, pt := range p {
		if pt.Y < minY {
			minY = pt.Y
		}
		if pt.Y > maxY {
			maxY = pt.Y
		}
	}

	for y := minY; y <= maxY; y++ {
		var intersections []float64
		for i := 0; i < len(p); i++ {
			next := (i + 1) % len(p)
			if intersects(p[i], p[next], y) {
				x := float64(intersectionX(p[i], p[next], y))
				intersections = append(intersections, x)
			}
		}

		// Sort the intersections.
		sort.Float64s(intersections)

		// Draw horizontal lines between pairs of intersections.
		for i := 0; i < len(intersections)-1; i += 2 {
			start := geom.MakePoint(float32(intersections[i]), y)
			end := geom.MakePoint(float32(intersections[i+1]), y)
			line := Line(geom.MakeLine(start, end))
			line.Draw(d, clr)
		}
	}
}

// Determines if a line segment intersects with a horizontal line at y.
func intersects(p1, p2 geom.Point, y float32) bool {
	return (p1.Y <= y && p2.Y > y) || (p1.Y > y && p2.Y <= y)
}

// Gets the x-coordinate of intersection of a line segment with a horizontal line at y.
func intersectionX(p1, p2 geom.Point, y float32) float32 {
	if p1.Y == p2.Y {
		return p1.X
	}
	return p1.X + (y-p1.Y)*(p2.X-p1.X)/(p2.Y-p1.Y)
}
