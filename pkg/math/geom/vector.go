package geom

import (
	"fmt"
	"math"
)

// Vector a float32 slice with 2 elements [2]float32{x, y}
type Vector [2]float32

func MakeVector(x, y float32) Vector {
	return Vector([2]float32{x, y})
}

func (v Vector) String() string {
	return fmt.Sprintf("%f, %f", v[0], v[1])
}

func (v Vector) ToPoint() Point {
	return MakePoint(v[0], v[1])
}

func (v Vector) Offset(o Vector) Vector {
	return Vector([2]float32{v[0] - o[0], v[1] - o[1]})
}

func (v Vector) GetDistance(b Vector) float32 {
	return float32(math.Sqrt(math.Pow(float64(v[0]-b[0]), 2) + math.Pow(float64(v[1]-b[1]), 2)))
}

func (v Vector) GetDirection(b Vector) float32 {
	return float32(math.Atan2(float64(b[1]-v[1]), float64(b[0]-v[0])))
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

type Vector3D [3]float32

func MakeVector3D(x, y, z float32) Vector3D {
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

func (v Vector3D) DistanceTo(other Vector3D) float32 {
	return float32(math.Sqrt(math.Pow(float64(v[0]-other[0]), 2) + math.Pow(float64(v[1]-other[1]), 2) + math.Pow(float64(v[2]-other[2]), 2)))
}

// Scaled scales the vector by a scalar value.
func (v Vector3D) Scaled(scalar float32) Vector3D {
	return Vector3D{v[0] * float32(scalar), v[1] * float32(scalar), v[2] * float32(scalar)}
}

func (v Vector3D) normalize() Vector3D {
	mag := float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])))
	return Vector3D{v[0] / float32(mag), v[1] / float32(mag), v[2] / float32(mag)}
}
