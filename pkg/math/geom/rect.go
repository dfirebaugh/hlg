package geom

import (
	"math"
)

// Rect a float32 slice with 4 elements []float32{x, y, width, height}
type Rect [4]float32

func MakeRect(x, y, width, height float32) Rect {
	return Rect{x, y, width, height}
}

func (r Rect) IsAxisAlignedCollision(other Rect) bool {
	ax := r[0]
	ay := r[1]
	aw := r[2]
	ah := r[3]

	bx := other[0]
	by := other[1]
	bw := other[2]
	bh := other[3]

	return ax < bx+bw &&
		ax+aw > bx &&
		ay < by+bh &&
		ah+ay > by
}

func (r Rect) GetOverlap(other Rect) (float32, float32) {
	return (r[0] + r[2]) - other[0], (r[1] + r[3]) - other[1]
}

// Dimensions returns the total number of dimensions
func (r Rect) Dimensions() int {
	return 4
}

func (r Rect) GetCenter() (float32, float32) {
	return (r[0] + r[0] + r[2]) / 2, (r[1] + r[1] + r[3]) / 2
}

// Dimension returns the value of the i-th dimension
func (r Rect) Dimension(i int) float32 {
	return r[i]
}

func (r Rect) ContainsPoint(p Point) bool {
	return p.X >= r[0] && p.Y >= r[1] && p.X < r[0]+r[2] && p.Y < r[1]+r[3]
}

type Normal struct{}
type Collision struct {
	Point    Vector
	FarPoint Vector
	Normal   Vector
	TimeNear Vector
	TimeFar  Vector
	HitNear  float64
	HitFar   float64
}

// HasRayIntersection returns true if an intersection exists
// the collision argument will contain information about the collision
func (r Rect) HasRayIntersection(ray Ray, collision *Collision) bool {
	invdir := MakeVector(1.0, 1.0).Divide(ray.Direction)

	targetPoint := MakeVector(r[0], r[1])
	targetSize := MakeVector(r[2], r[3])

	// Calculate intersections with rectangle bounding axes
	collision.TimeNear = targetPoint.Subtract(ray.Origin).Multiply(invdir)
	collision.TimeFar = targetPoint.Add(targetSize).Subtract(ray.Origin).Multiply(invdir)

	if math.IsNaN(float64(collision.TimeNear[0])) ||
		math.IsNaN(float64(collision.TimeNear[1])) ||
		math.IsNaN(float64(collision.TimeFar[0])) ||
		math.IsNaN(float64(collision.TimeFar[1])) {
		return false
	}

	// sort distances
	if collision.TimeNear[0] > collision.TimeFar[0] {
		collision.TimeNear[0], collision.TimeFar[0] = collision.TimeFar[0], collision.TimeNear[0]
	}
	if collision.TimeNear[1] > collision.TimeFar[1] {
		collision.TimeNear[1], collision.TimeFar[1] = collision.TimeFar[1], collision.TimeNear[1]
	}

	// Early Rejection
	if collision.TimeNear[0] > collision.TimeFar[1] || collision.TimeNear[1] > collision.TimeFar[0] {
		return false
	}

	// closest 'time' will be the first contact
	collision.HitNear = math.Max(float64(collision.TimeNear[0]), float64(collision.TimeNear[1]))

	// furthest 'time' is contact on opposite side of target
	collision.HitFar = math.Min(float64(collision.TimeFar[0]), float64(collision.TimeFar[1]))

	if collision.HitFar < 0 {
		return false
	}
	if collision.HitNear > 1 {
		return false
	}

	collision.Point = ray.Origin.Add(MakeVector(float32(collision.HitNear), float32(collision.HitNear)).Multiply(ray.Direction))

	if collision.TimeNear[0] > collision.TimeNear[1] {
		if invdir[0] < 0 {
			collision.Normal = MakeVector(1, 0)
		} else {
			collision.Normal = MakeVector(-1, 0)
		}
	} else if collision.TimeNear[0] < collision.TimeNear[1] {
		if invdir[1] < 0 {
			collision.Normal = MakeVector(0, 1)
		} else {
			collision.Normal = MakeVector(0, -1)
		}
	}

	return true
}
