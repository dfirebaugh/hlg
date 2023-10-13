package geom

import "math"

type Vector3 struct {
	X, Y, Z float32
}

type Matrix4 [16]float32

type Texture struct {
	ID uint32
}

type Mesh struct {
	Vertices []float32
	Indices  []uint32
}

type Model struct {
	Meshes      []*Mesh
	ScaleFactor float32
	Position    Vector3
	Rotation    Matrix4
}

func NewModel(mesh *Mesh) *Model {
	return &Model{
		Meshes:      []*Mesh{mesh},
		ScaleFactor: 1.0,
		Position:    Vector3{0, 0, 0},
		Rotation: Matrix4{
			1, 0, 0, 0,
			0, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		},
	}
}

func NewMesh(vertices []float32, indices []uint32) *Mesh {
	return &Mesh{
		Vertices: vertices,
		Indices:  indices,
	}
}

func (m *Model) Translate(v Vector3) {
	m.Position.X += v.X
	m.Position.Y += v.Y
	m.Position.Z += v.Z
}

func (m *Model) SetPosition(v Vector3) {
	m.Position = v
}

func (m *Model) Scale(factor float32) {
	m.ScaleFactor *= factor
}

func (m *Model) Rotate(angle float32, axis Vector3) {
	radAngle := float64(angle) * math.Pi / 180.0

	s := float32(math.Sin(radAngle))
	c := float32(math.Cos(radAngle))

	axis = normalize(axis)
	ux, uy, uz := axis.X, axis.Y, axis.Z

	r := Matrix4{
		c + ux*ux*(1-c), ux*uy*(1-c) - uz*s, ux*uz*(1-c) + uy*s, 0,
		uy*ux*(1-c) + uz*s, c + uy*uy*(1-c), uy*uz*(1-c) - ux*s, 0,
		uz*ux*(1-c) - uy*s, uz*uy*(1-c) + ux*s, c + uz*uz*(1-c), 0,
		0, 0, 0, 1,
	}

	m.Rotation = multiplyMatrices(m.Rotation, r)
}

func (m *Model) GetMeshes() []*Mesh {
	return m.Meshes
}

func (m *Model) GetScaleFactor() float32 {
	return m.ScaleFactor
}

func (m *Model) GetPosition() Vector3 {
	return m.Position
}

func (m *Model) GetRotation() Matrix4 {
	return m.Rotation
}

func multiplyMatrices(a, b Matrix4) Matrix4 {
	var c Matrix4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			c[i*4+j] = a[i*4+0]*b[0*4+j] +
				a[i*4+1]*b[1*4+j] +
				a[i*4+2]*b[2*4+j] +
				a[i*4+3]*b[3*4+j]
		}
	}
	return c
}

func normalize(v Vector3) Vector3 {
	mag := float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
	return Vector3{v.X / mag, v.Y / mag, v.Z / mag}
}
