package geom

import (
	"fmt"
	"math"
)

// Vector a float64 slice with 2 elements [2]float64{x, y}
type Vector [2]float64

func MakeVector(x, y float64) Vector {
	return Vector([2]float64{x, y})
}

func (v Vector) String() string {
	return fmt.Sprintf("%f, %f", v[0], v[1])
}

func (v Vector) ToPoint() Point {
	return MakePoint(v[0], v[1])
}

func (v Vector) Offset(o Vector) Vector {
	return Vector([2]float64{v[0] - o[0], v[1] - o[1]})
}

func (v Vector) GetDistance(b Vector) float64 {
	return math.Sqrt(math.Pow(v[0]-b[0], 2) + math.Pow(v[1]-b[1], 2))
}

func (v Vector) GetDirection(b Vector) float64 {
	return math.Atan2(b[1]-v[1], b[0]-v[0])
}

func (v Vector) Subtract(other Vector) Vector {
	return MakeVector(v[0]-other[0], v[1]-other[1])
}

func (v Vector) Add(other Vector) Vector {
	return MakeVector(v[0]+other[0], v[1]+other[1])
}

func (v Vector) Multiply(other Vector) Vector {
	return MakeVector(v[0]*other[0], v[1]*other[1])
}

func (v Vector) Divide(other Vector) Vector {
	return MakeVector(v[0]/other[0], v[1]/other[1])
}

type Vector3D [3]float64

func MakeVector3D(x, y, z float64) Vector3D {
	return Vector3D{x, y, z}
}

func (v Vector3D) ToPoint3D() Point3D {
	return MakePoint3D(v[0], v[1], v[2])
}

func (v Vector3D) Subtract(other Vector3D) Vector3D {
	return MakeVector3D(v[0]-other[0], v[1]-other[1], v[2]-other[2])
}

func (v Vector3D) Add(other Vector3D) Vector3D {
	return MakeVector3D(v[0]+other[0], v[1]+other[1], v[2]+other[2])
}

func (v Vector3D) DistanceTo(other Vector3D) float64 {
	return math.Sqrt(math.Pow(v[0]-other[0], 2) + math.Pow(v[1]-other[1], 2) + math.Pow(v[2]-other[2], 2))
}
