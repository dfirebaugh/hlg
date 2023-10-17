package geom

import "math"

type Matrix4 [16]float32

func (m Matrix4) Multiply(other Matrix4) Matrix4 {
	var result Matrix4

	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			var sum float32
			for i := 0; i < 4; i++ {
				sum += m[row*4+i] * other[i*4+col]
			}
			result[row*4+col] = sum
		}
	}

	return result
}

func Matrix4RotationX(angle float32) Matrix4 {
	cosAngle := float32(math.Cos(float64(angle)))
	sinAngle := float32(math.Sin(float64(angle)))

	matrix := Matrix4Identity()
	matrix[5] = cosAngle
	matrix[6] = -sinAngle
	matrix[9] = sinAngle
	matrix[10] = cosAngle

	return matrix
}

func Matrix4RotationY(angle float32) Matrix4 {
	cosAngle := float32(math.Cos(float64(angle)))
	sinAngle := float32(math.Sin(float64(angle)))

	matrix := Matrix4Identity()
	matrix[0] = cosAngle
	matrix[2] = sinAngle
	matrix[8] = -sinAngle
	matrix[10] = cosAngle

	return matrix
}

func Matrix4RotationZ(angle float32) Matrix4 {
	cosAngle := float32(math.Cos(float64(angle)))
	sinAngle := float32(math.Sin(float64(angle)))

	matrix := Matrix4Identity()
	matrix[0] = cosAngle
	matrix[1] = -sinAngle
	matrix[4] = sinAngle
	matrix[5] = cosAngle

	return matrix
}

func Matrix4Identity() Matrix4 {
	matrix := Matrix4{}
	for i := 0; i < 16; i++ {
		if i%5 == 0 {
			matrix[i] = 1.0
		} else {
			matrix[i] = 0.0
		}
	}
	return matrix
}

func createRotationMatrix(x, y, z float32) Matrix4 {
	rotationMatrixX := Matrix4RotationX(x)
	rotationMatrixY := Matrix4RotationY(y)
	rotationMatrixZ := Matrix4RotationZ(z)

	finalRotationMatrix := rotationMatrixX.Multiply(rotationMatrixY).Multiply(rotationMatrixZ)

	return finalRotationMatrix
}

func createTranslationMatrix(dx, dy, dz float32) Matrix4 {
	return Matrix4{
		1, 0, 0, dx,
		0, 1, 0, dy,
		0, 0, 1, dz,
		0, 0, 0, 1,
	}
}
