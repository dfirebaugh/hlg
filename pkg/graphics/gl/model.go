package gl

import (
	"github.com/dfirebaugh/ggez/pkg/graphics"
	"github.com/dfirebaugh/ggez/pkg/math/geom"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Mesh struct {
	Vertices []float32
	Indices  []uint32
	VAO      uint32
}

type Model struct {
	Meshes   []*Mesh
	scale    float32
	position mgl32.Vec3
	rotation mgl32.Mat4
}

func TranslateGeomModelToGLModel(m graphics.Model) (*Model, error) {
	glModel := &Model{
		Meshes:   make([]*Mesh, len(m.GetMeshes())),
		scale:    m.GetScaleFactor(),
		position: Vector3ToMgl32(m.GetPosition()),
		rotation: Matrix4ToMgl32(m.GetRotation()),
	}

	for i, gMesh := range m.GetMeshes() {
		glMesh := &Mesh{
			Vertices: gMesh.Vertices,
			Indices:  gMesh.Indices,
		}

		glMesh.VAO = glMesh.createVAO()
		glModel.Meshes[i] = glMesh
	}

	return glModel, nil
}

func NewMesh(vertices []float32, indices []uint32) *Mesh {
	m := &Mesh{
		Vertices: vertices,
		Indices:  indices,
	}
	m.VAO = m.createVAO()
	return m
}

func (m *Mesh) createVAO() uint32 {
	var VAO, VBO, EBO uint32

	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.GenBuffers(1, &EBO)

	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.Vertices)*4, gl.Ptr(m.Vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.Indices)*4, gl.Ptr(m.Indices), gl.STATIC_DRAW)

	stride := int32(5 * 4)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, stride, uintptr(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, stride, uintptr(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	return VAO
}

func (m *Model) Translate(v mgl32.Vec3) {
	m.position = m.position.Add(v)
}

func (m *Model) SetPosition(v mgl32.Vec3) {
	m.position = v
}

func (m *Model) Scale(n float32) {
	m.scale = n
}

func (m *Model) Rotate(angle float32, axis mgl32.Vec3) {
	m.rotation = mgl32.HomogRotate3D(mgl32.DegToRad(angle), axis).Mul4(m.rotation)
}

func (m *Model) Draw(program graphics.ShaderProgram) {
	program.Use()

	model := mgl32.Translate3D(m.position.X(), m.position.Y(), m.position.Z()).
		Mul4(m.rotation).
		Mul4(mgl32.Scale3D(m.scale, m.scale, m.scale))
	view := mgl32.LookAtV(cameraPos, cameraPos.Add(cameraFront), cameraUp)
	projection := mgl32.Perspective(mgl32.DegToRad(zoom), float32(windowWidth/windowHeight), 0.1, 100.0)

	modelLoc := program.GetUniformLocation("model")
	viewLoc := program.GetUniformLocation("view")
	projLoc := program.GetUniformLocation("projection")

	gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])
	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
	gl.UniformMatrix4fv(projLoc, 1, false, &projection[0])

	for _, mesh := range m.Meshes {
		gl.BindVertexArray(mesh.VAO)
		gl.DrawElements(gl.TRIANGLES, int32(len(mesh.Indices)), gl.UNSIGNED_INT, nil)
		gl.BindVertexArray(0)
	}
}

func Vector3ToMgl32(v geom.Vector3) mgl32.Vec3 {
	return mgl32.Vec3{v.X, v.Y, v.Z}
}

func Matrix4ToMgl32(m geom.Matrix4) mgl32.Mat4 {
	return mgl32.Mat4{
		m[0], m[1], m[2], m[3],
		m[4], m[5], m[6], m[7],
		m[8], m[9], m[10], m[11],
		m[12], m[13], m[14], m[15],
	}
}
