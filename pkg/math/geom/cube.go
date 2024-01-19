package geom

type Cube [8]Point3D

func MakeCube(points [8]Point3D) Cube {
	return Cube(points)
}

func (c *Cube) Centroid() Point3D {
	var sumX, sumY, sumZ float32 = 0.0, 0.0, 0.0
	for _, p := range c {
		sumX += p.X
		sumY += p.Y
		sumZ += p.Z
	}
	return MakePoint3D(sumX/8, sumY/8, sumZ/8)
}

func (c *Cube) Translate(vector Vector3D) {
	for i := range c {
		c[i].X += vector[0]
		c[i].Y += vector[1]
		c[i].Z += vector[2]
	}
}

func (c *Cube) Scale(factor float32) {
	centroid := c.Centroid()
	for i := range c {
		dir := c[i].ToVector3D().Subtract(centroid.ToVector3D())
		c[i] = Point3D{
			X: centroid.X + dir[0]*factor,
			Y: centroid.Y + dir[1]*factor,
			Z: centroid.Z + dir[2]*factor,
		}
	}
}

func (c Cube) GetVertices() []float32 {
	var vertices []float32
	for _, point := range c {
		vertices = append(vertices, float32(point.X), float32(point.Y), float32(point.Z))
	}
	return vertices
}
