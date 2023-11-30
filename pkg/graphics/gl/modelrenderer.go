package gl

import (
	"github.com/dfirebaugh/ggez/pkg/graphics"
	"github.com/dfirebaugh/ggez/pkg/math/geom"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type ModelRenderer struct {
	cam                *Camera
	lightIntensity     float32
	program            graphics.ShaderProgram
	nonTexturedProgram graphics.ShaderProgram
}

func NewModelRenderer(cam *Camera, program graphics.ShaderProgram, nonTexturedProgram graphics.ShaderProgram) ModelRenderer {
	modelMatrix = mgl32.Ident4()

	return ModelRenderer{
		cam:                cam,
		lightIntensity:     4,
		program:            program,
		nonTexturedProgram: nonTexturedProgram,
	}
}

var (
	modelMatrix mgl32.Mat4
)

func (mr ModelRenderer) createVertexArrayObject(vertices []float32, indices []uint32) uint32 {
	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, uintptr(0))
	gl.BindVertexArray(0)

	return vao
}

func (mr ModelRenderer) RenderNonTexturedModel(m graphics.Model) {
	mr.nonTexturedProgram.Use()

	modelMatrix = mgl32.Ident4()
	modelMatrix = modelMatrix.Mul4(mgl32.Scale3D(m.GetScaleFactor(), m.GetScaleFactor(), m.GetScaleFactor()))
	modelMatrix = modelMatrix.Mul4(mgl32.Translate3D(m.GetPosition()[0], m.GetPosition()[1], m.GetPosition()[2]))

	rotationAngle := float32(glfw.GetTime())
	modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3D(rotationAngle, mgl32.Vec3{0, 1, 0}))

	modelLocation := mr.nonTexturedProgram.GetUniformLocation("model")
	gl.UniformMatrix4fv(modelLocation, 1, false, &modelMatrix[0])

	for _, mesh := range m.GetMeshes() {
		vao := mr.createVertexArrayObject(mesh.Vertices, mesh.Indices)

		gl.BindVertexArray(vao)
		gl.DrawElementsWithOffset(gl.TRIANGLES, int32(len(mesh.Indices)), gl.UNSIGNED_INT, uintptr(0))
		gl.BindVertexArray(0)
	}
}

func (mr ModelRenderer) RenderModel(m graphics.Model, t graphics.Texture) {
	if t == nil {
		mr.RenderNonTexturedModel(m)
		return
	}

	mr.program.Use()

	texture := textures[t.Handle()]
	texture.Bind(gl.TEXTURE0)
	defer texture.UnBind()

	texture.SetUniform(mr.program.GetUniformLocation("texture_diffuse1"))
	lightPosLoc := mr.program.GetUniformLocation("lightPos")
	lightColorLoc := mr.program.GetUniformLocation("lightColor")

	lightX, lightY, lightZ := float32(100), float32(100), float32(100)

	gl.Uniform3f(lightPosLoc, lightX, lightY, lightZ)
	gl.Uniform3f(lightColorLoc, 1*mr.lightIntensity, 1*mr.lightIntensity, 1*mr.lightIntensity)

	model, _ := TranslateGeomModelToGLModel(m.(*geom.Model))
	model.Draw(mr.program, mr.cam)
}
