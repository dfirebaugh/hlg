package geom

import "math"

type Point struct {
	X float32
	Y float32
}

func MakePoint(x, y float32) Point {
	return Point{X: x, Y: y}
}

func (p Point) ToVector() Vector {
	return MakeVector(p.X, p.Y)
}

type Point3D struct {
	X, Y, Z float32
}

func MakePoint3D(x, y, z float32) Point3D {
	return Point3D{X: x, Y: y, Z: z}
}

func (p Point3D) Add(other Point3D) Point3D {
	return Point3D{
		X: p.X + other.X,
		Y: p.Y + other.Y,
		Z: p.Z + other.Z,
	}
}

func (p Point3D) Subtract(other Point3D) Point3D {
	return Point3D{
		X: p.X - other.X,
		Y: p.Y - other.Y,
		Z: p.Z - other.Z,
	}
}

func (p Point3D) Cross(other Point3D) Point3D {
	return Point3D{
		X: p.Y*other.Z - p.Z*other.Y,
		Y: p.Z*other.X - p.X*other.Z,
		Z: p.X*other.Y - p.Y*other.X,
	}
}

func (p Point3D) Magnitude() float32 {
	return float32(math.Sqrt(float64(p.X*p.X + p.Y*p.Y + p.Z*p.Z)))
}

func (p Point3D) Normalize() Point3D {
	mag := p.Magnitude()
	if mag == 0 {
		return Point3D{0, 0, 0}
	}
	return Point3D{
		X: p.X / mag,
		Y: p.Y / mag,
		Z: p.Z / mag,
	}
}

func (p Point3D) ToVector3D() Vector3D {
	return MakeVector3D(p.X, p.Y, p.Z)
}
