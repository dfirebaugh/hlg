package geom

type Polygon []Point

// MakePolygon creates a new polygon from a slice of points.
func MakePolygon(points ...Point) Polygon {
	return Polygon(points)
}

// AddPoint appends a point to the end of the polygon.
func (p *Polygon) AddPoint(pt Point) {
	*p = append(*p, pt)
}

// NumPoints returns the number of points in the polygon.
func (p Polygon) NumPoints() int {
	return len(p)
}

// At returns the i-th point of the polygon.
func (p Polygon) At(i int) Point {
	return p[i]
}
