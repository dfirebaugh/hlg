//go:build !js

package shader

import (
	_ "embed"

	"github.com/dfirebaugh/hlg/graphics"
)

var (
	//go:embed texture.wgsl
	textureShaderCode string

	//go:embed primitive_buffer.wgsl
	primitiveBufferShaderCode string

	//go:embed solid_shape.wgsl
	solidShapeShaderCode string

	TextureShader         graphics.ShaderHandle
	PrimitiveBufferShader graphics.ShaderHandle
	SolidShapeShader      graphics.ShaderHandle
)

func CompileShaders(sm *ShaderManager) {
	TextureShader = sm.CompileShader(textureShaderCode)
	PrimitiveBufferShader = sm.CompileShader(primitiveBufferShaderCode)
	SolidShapeShader = sm.CompileShader(solidShapeShaderCode)
}
