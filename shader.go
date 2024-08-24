package hlg

// CompileShader takes wgsl shader code as an argument and returns a hash
// that can be used internally to reference a shader
func CompileShader(shaderCode string) int {
	ensureSetupCompletion()
	return int(hlg.graphicsBackend.CompileShader(shaderCode))
}
