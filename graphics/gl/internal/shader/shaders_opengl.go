//go:build !js

package shader

import (
	_ "embed"

	"github.com/dfirebaugh/hlg/graphics"
)

var (
	//go:embed texture_opengl.vert
	textureVertexShaderCode string

	//go:embed texture_opengl.frag
	textureFragmentShaderCode string

	//go:embed primitive_buffer_opengl.vert
	primitiveBufferVertexShaderCode string

	//go:embed primitive_buffer_opengl.frag
	primitiveBufferFragmentShaderCode string

	TextureShader         graphics.ShaderHandle
	PrimitiveBufferShader graphics.ShaderHandle
)

func CompileShaders(sm *ShaderManager) {
	TextureShader = sm.CompileShaderFromSource(textureVertexShaderCode, textureFragmentShaderCode)
	PrimitiveBufferShader = sm.CompileShaderFromSource(primitiveBufferVertexShaderCode, primitiveBufferFragmentShaderCode)
}
