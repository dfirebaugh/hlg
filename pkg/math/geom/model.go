package geom

import "math"

type Texture struct {
	ID uint32
}

type Mesh struct {
	Vertices []float32
	Indices  []uint32
}

type Model struct {
	Meshes      []*Mesh
	Matrix      Matrix4
	ScaleFactor float32
	Position    Vector3D
	Rotation    Matrix4
}

func NewModel(mesh *Mesh) *Model {
	return &Model{
		Meshes:      []*Mesh{mesh},
		ScaleFactor: 1.0,
		Position:    Vector3D{0, 0, 0},
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

func (m *Model) Translate(v Vector3D) {
	m.Position[0] += v[0]
	m.Position[1] += v[1]
	m.Position[2] += v[2]
}

func (m *Model) ApplyMatrix(matrix Matrix4) {
	currentMatrix := m.Matrix

	m.Matrix = currentMatrix.Multiply(matrix)
}

func (m *Model) SetPosition(v Vector3D) {
	m.Position = v
}

func (m *Model) Scale(factor float32) {
	m.ScaleFactor *= factor
}

func (m *Model) SetRotation(v Vector3D) {
	translationMatrixToOrigin := createTranslationMatrix(-m.Position[0], -m.Position[1], -m.Position[2])
	m.ApplyMatrix(translationMatrixToOrigin)

	// Apply rotation
	rotationMatrix := createRotationMatrix(v[0], v[1], v[2])
	m.Rotation = rotationMatrix

	translationMatrixBack := createTranslationMatrix(m.Position[0], m.Position[1], m.Position[2])
	m.ApplyMatrix(translationMatrixBack)
}

func (m *Model) Rotate(angle float32, axis Vector3D) {
	translationMatrixToOrigin := createTranslationMatrix(-m.Position[0], -m.Position[1], -m.Position[2])
	m.ApplyMatrix(translationMatrixToOrigin)

	radAngle := float32(angle) * math.Pi / 180.0
	s := float32(math.Sin(float64(radAngle)))
	c := float32(math.Cos(float64(radAngle)))
	axis = axis.normalize()
	ux, uy, uz := axis[0], axis[1], axis[2]
	r := Matrix4{
		c + ux*ux*(1-c), ux*uy*(1-c) - uz*s, ux*uz*(1-c) + uy*s, 0,
		uy*ux*(1-c) + uz*s, c + uy*uy*(1-c), uy*uz*(1-c) - ux*s, 0,
		uz*ux*(1-c) - uy*s, uz*uy*(1-c) + ux*s, c + uz*uz*(1-c), 0,
		0, 0, 0, 1,
	}
	m.Rotation = m.Rotation.Multiply(r)

	translationMatrixBack := createTranslationMatrix(m.Position[0], m.Position[1], m.Position[2])
	m.ApplyMatrix(translationMatrixBack)
}

func (m *Model) GetMeshes() []*Mesh {
	return m.Meshes
}

func (m *Model) GetScaleFactor() float32 {
	return m.ScaleFactor
}

func (m *Model) GetPosition() Vector3D {
	return m.Position
}

func (m *Model) GetRotation() Matrix4 {
	return m.Rotation
}
