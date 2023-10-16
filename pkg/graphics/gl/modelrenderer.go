package gl

import (
	"github.com/dfirebaugh/ggez/pkg/graphics"
	"github.com/dfirebaugh/ggez/pkg/math/geom"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type ModelRenderer struct {
	lightIntensity float32
	program        graphics.ShaderProgram
}

func NewModelRenderer(program graphics.ShaderProgram) ModelRenderer {
	return ModelRenderer{
		lightIntensity: 4,
		program:        program,
	}
}

func (mr ModelRenderer) RenderModel(m graphics.Model, t graphics.Texture) {
	mr.program.Use()

	texture := textures[t.Handle()]
	texture.Bind(gl.TEXTURE0)
	defer texture.UnBind()

	texture.SetUniform(mr.program.GetUniformLocation("texture_diffuse1"))

	model, _ := TranslateGeomModelToGLModel(m.(*geom.Model))
	model.Draw(mr.program)

	lightPosLoc := mr.program.GetUniformLocation("lightPos")
	lightColorLoc := mr.program.GetUniformLocation("lightColor")

	lightX, lightY, lightZ := float32(100), float32(100), float32(100)

	gl.Uniform3f(lightPosLoc, lightX, lightY, lightZ)
	gl.Uniform3f(lightColorLoc, 1*mr.lightIntensity, 1*mr.lightIntensity, 1*mr.lightIntensity)
}
