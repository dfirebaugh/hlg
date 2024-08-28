package shader

import (
	_ "embed"

	"github.com/dfirebaugh/hlg/graphics"
)

var (
	//go:embed shapes.wgsl
	shapeShaderCode string

	//go:embed texture.wgsl
	textureShaderCode string

	ShapeShader   graphics.ShaderHandle
	TextureShader graphics.ShaderHandle
)

func CompileShaders(sm *ShaderManager) {
	ShapeShader = sm.CompileShader(shapeShaderCode)
	TextureShader = sm.CompileShader(textureShaderCode)
}
