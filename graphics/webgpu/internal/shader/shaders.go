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

	//go:embed primitive_buffer.wgsl
	primitiveBufferShaderCode string

	ShapeShader           graphics.ShaderHandle
	TextureShader         graphics.ShaderHandle
	PrimitiveBufferShader graphics.ShaderHandle
)

func CompileShaders(sm *ShaderManager) {
	ShapeShader = sm.CompileShader(shapeShaderCode)
	TextureShader = sm.CompileShader(textureShaderCode)
	PrimitiveBufferShader = sm.CompileShader(primitiveBufferShaderCode)
}
