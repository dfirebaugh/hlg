package vector

import "math"

type Matrix3x3 [3][3]float64

func (m *Matrix3x3) MultiplyVector(v Vector3D) Vector3D {
	return Vector3D{
		X: m[0][0]*v.X + m[0][1]*v.Y + m[0][2]*v.Z,
		Y: m[1][0]*v.X + m[1][1]*v.Y + m[1][2]*v.Z,
		Z: m[2][0]*v.X + m[2][1]*v.Y + m[2][2]*v.Z,
	}
}

func RotationMatrixX(angle float64) Matrix3x3 {
	return Matrix3x3{
		{1, 0, 0},
		{0, math.Cos(angle), -math.Sin(angle)},
		{0, math.Sin(angle), math.Cos(angle)},
	}
}

func RotationMatrixY(angle float64) Matrix3x3 {
	return Matrix3x3{
		{math.Cos(angle), 0, math.Sin(angle)},
		{0, 1, 0},
		{-math.Sin(angle), 0, math.Cos(angle)},
	}
}

func RotationMatrixZ(angle float64) Matrix3x3 {
	return Matrix3x3{
		{math.Cos(angle), -math.Sin(angle), 0},
		{math.Sin(angle), math.Cos(angle), 0},
		{0, 0, 1},
	}
}
