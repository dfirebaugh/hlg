package matrix

import "math"

type Matrix [16]float32

func (m Matrix) Multiply(other Matrix) Matrix {
	var result Matrix

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

func (m Matrix) Scale(sx, sy float32) Matrix {
	m[0] *= sx
	m[5] *= sy

	return m
}

func (m Matrix) Rotate(angle float32) Matrix {
	c := float32(math.Cos(float64(angle)))
	s := float32(math.Sin(float64(angle)))

	rotationMatrix := Matrix([16]float32{
		c, -s, 0, 0,
		s, c, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	})

	return m.Multiply(rotationMatrix)
}

func (m Matrix) Translate(dx, dy float32) Matrix {
	translationMatrix := CreateTranslationMatrix(dx, dy, 0)
	return m.Multiply(translationMatrix)
}

func (m Matrix) IsZero() bool {
	for _, value := range m {
		if value != 0 {
			return false
		}
	}
	return true
}

type Vector4 [4]float32

func (m Matrix) MultiplyVec(v Vector4) Vector4 {
	return Vector4{
		m[0]*v[0] + m[4]*v[1] + m[8]*v[2] + m[12]*v[3],
		m[1]*v[0] + m[5]*v[1] + m[9]*v[2] + m[13]*v[3],
		m[2]*v[0] + m[6]*v[1] + m[10]*v[2] + m[14]*v[3],
		m[3]*v[0] + m[7]*v[1] + m[11]*v[2] + m[15]*v[3],
	}
}

func MatrixIdentity() Matrix {
	matrix := Matrix{}
	for i := 0; i < 16; i++ {
		if i%5 == 0 {
			matrix[i] = 1.0
		} else {
			matrix[i] = 0.0
		}
	}
	return matrix
}

func MatrixRotationX(angle float32) Matrix {
	cosAngle := float32(math.Cos(float64(angle)))
	sinAngle := float32(math.Sin(float64(angle)))

	matrix := MatrixIdentity()
	matrix[5] = cosAngle
	matrix[6] = -sinAngle
	matrix[9] = sinAngle
	matrix[10] = cosAngle

	return matrix
}

func MatrixRotationY(angle float32) Matrix {
	cosAngle := float32(math.Cos(float64(angle)))
	sinAngle := float32(math.Sin(float64(angle)))

	matrix := MatrixIdentity()
	matrix[0] = cosAngle
	matrix[2] = sinAngle
	matrix[8] = -sinAngle
	matrix[10] = cosAngle

	return matrix
}

func MatrixRotationZ(angle float32) Matrix {
	cosAngle := float32(math.Cos(float64(angle)))
	sinAngle := float32(math.Sin(float64(angle)))

	matrix := MatrixIdentity()
	matrix[0] = cosAngle
	matrix[1] = -sinAngle
	matrix[4] = sinAngle
	matrix[5] = cosAngle

	return matrix
}

func CreateRotationMatrix(x, y, z float32) Matrix {
	rotationMatrixX := MatrixRotationX(x)
	rotationMatrixY := MatrixRotationY(y)
	rotationMatrixZ := MatrixRotationZ(z)

	finalRotationMatrix := rotationMatrixX.Multiply(rotationMatrixY).Multiply(rotationMatrixZ)

	return finalRotationMatrix
}

func CreateTranslationMatrix(dx, dy, dz float32) Matrix {
	return Matrix{
		1, 0, 0, dx,
		0, 1, 0, dy,
		0, 0, 1, dz,
		0, 0, 0, 1,
	}
}
