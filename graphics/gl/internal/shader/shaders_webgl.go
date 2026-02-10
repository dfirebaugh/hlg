//go:build js && wasm

package shader

import (
	_ "embed"

	"github.com/dfirebaugh/hlg/graphics"
)

var (
	//go:embed texture_webgl.vert
	textureVertexShaderCode string

	//go:embed texture_webgl.frag
	textureFragmentShaderCode string

	//go:embed primitive_buffer_webgl.vert
	primitiveBufferVertexShaderCode string

	//go:embed primitive_buffer_webgl.frag
	primitiveBufferFragmentShaderCode string

	TextureShader         graphics.ShaderHandle
	PrimitiveBufferShader graphics.ShaderHandle
)

// CompileShaders compiles all built-in shaders
func CompileShaders(sm *ShaderManager) {
	TextureShader = sm.CompileShaderFromSource(textureVertexShaderCode, textureFragmentShaderCode)
	PrimitiveBufferShader = sm.CompileShaderFromSource(primitiveBufferVertexShaderCode, primitiveBufferFragmentShaderCode)
}
