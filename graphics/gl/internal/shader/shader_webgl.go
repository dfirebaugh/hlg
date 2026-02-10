//go:build js && wasm

package shader

import (
	"fmt"
	"strings"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/glapi"
)

// ShaderManager manages WebGL shader programs
type ShaderManager struct {
	ctx        *glapi.Context
	programs   map[graphics.ShaderHandle]glapi.Program
	nextHandle graphics.ShaderHandle
}

// NewShaderManager creates a new shader manager
func NewShaderManager(ctx *glapi.Context) *ShaderManager {
	return &ShaderManager{
		ctx:        ctx,
		programs:   make(map[graphics.ShaderHandle]glapi.Program),
		nextHandle: 1,
	}
}

// CompileShader compiles a combined GLSL shader (vertex + fragment with markers)
func (sm *ShaderManager) CompileShader(shaderCode string) graphics.ShaderHandle {
	vertexSource, fragmentSource, err := parseShaderCode(shaderCode)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse shader: %v", err))
	}

	program, err := sm.compileProgram(vertexSource, fragmentSource)
	if err != nil {
		panic(fmt.Sprintf("Failed to compile shader: %v", err))
	}

	handle := sm.nextHandle
	sm.nextHandle++
	sm.programs[handle] = program

	return handle
}

// parseShaderCode parses combined shader code with #vertex and #fragment markers
func parseShaderCode(shaderCode string) (vertexSource, fragmentSource string, err error) {
	vertexMarker := "#vertex"
	fragmentMarker := "#fragment"

	vertexIdx := strings.Index(shaderCode, vertexMarker)
	fragmentIdx := strings.Index(shaderCode, fragmentMarker)

	if vertexIdx == -1 || fragmentIdx == -1 {
		return "", "", fmt.Errorf("shader code must contain both #vertex and #fragment markers")
	}

	if vertexIdx < fragmentIdx {
		vertexSource = strings.TrimSpace(shaderCode[vertexIdx+len(vertexMarker) : fragmentIdx])
		fragmentSource = strings.TrimSpace(shaderCode[fragmentIdx+len(fragmentMarker):])
	} else {
		fragmentSource = strings.TrimSpace(shaderCode[fragmentIdx+len(fragmentMarker) : vertexIdx])
		vertexSource = strings.TrimSpace(shaderCode[vertexIdx+len(vertexMarker):])
	}

	// Convert OpenGL desktop GLSL to WebGL 2 GLSL ES
	vertexSource = convertToGLSLES(vertexSource)
	fragmentSource = convertToGLSLES(fragmentSource)

	// Fragment shaders need precision declaration
	if !strings.Contains(fragmentSource, "precision ") {
		fragmentSource = addPrecisionToFragment(fragmentSource)
	}

	return vertexSource, fragmentSource, nil
}

// convertToGLSLES converts OpenGL desktop GLSL to WebGL 2 GLSL ES 3.00
func convertToGLSLES(source string) string {
	// Replace OpenGL version with WebGL 2 version
	source = strings.Replace(source, "#version 410 core", "#version 300 es", 1)
	source = strings.Replace(source, "#version 330 core", "#version 300 es", 1)
	source = strings.Replace(source, "#version 400 core", "#version 300 es", 1)
	source = strings.Replace(source, "#version 450 core", "#version 300 es", 1)

	return source
}

// addPrecisionToFragment adds precision declaration after #version for fragment shaders
func addPrecisionToFragment(source string) string {
	versionEnd := strings.Index(source, "\n")
	if versionEnd == -1 {
		return source
	}
	// Insert precision declaration after version line
	return source[:versionEnd+1] + "precision mediump float;\n" + source[versionEnd+1:]
}

// CompileShaderFromSource compiles separate vertex and fragment shader sources
func (sm *ShaderManager) CompileShaderFromSource(vertexSource, fragmentSource string) graphics.ShaderHandle {
	program, err := sm.compileProgram(vertexSource, fragmentSource)
	if err != nil {
		panic(fmt.Sprintf("Failed to compile shader: %v", err))
	}

	handle := sm.nextHandle
	sm.nextHandle++
	sm.programs[handle] = program

	return handle
}

// GetProgram returns the WebGL program for a shader handle
func (sm *ShaderManager) GetProgram(handle graphics.ShaderHandle) glapi.Program {
	return sm.programs[handle]
}

// ReleaseShaders deletes all compiled shaders
func (sm *ShaderManager) ReleaseShaders() {
	for _, program := range sm.programs {
		sm.ctx.DeleteProgram(program)
	}
	sm.programs = make(map[graphics.ShaderHandle]glapi.Program)
}

// compileProgram compiles and links vertex and fragment shaders into a program
func (sm *ShaderManager) compileProgram(vertexSource, fragmentSource string) (glapi.Program, error) {
	vertexShader, err := sm.compileShader(vertexSource, glapi.VERTEX_SHADER)
	if err != nil {
		return glapi.InvalidProgram, fmt.Errorf("vertex shader: %w", err)
	}

	fragmentShader, err := sm.compileShader(fragmentSource, glapi.FRAGMENT_SHADER)
	if err != nil {
		sm.ctx.DeleteShader(vertexShader)
		return glapi.InvalidProgram, fmt.Errorf("fragment shader: %w", err)
	}

	program := sm.ctx.CreateProgram()
	sm.ctx.AttachShader(program, vertexShader)
	sm.ctx.AttachShader(program, fragmentShader)
	sm.ctx.LinkProgram(program)

	// Clean up shaders after linking
	sm.ctx.DeleteShader(vertexShader)
	sm.ctx.DeleteShader(fragmentShader)

	if !sm.ctx.GetProgramParameter(program, glapi.LINK_STATUS) {
		log := sm.ctx.GetProgramInfoLog(program)
		sm.ctx.DeleteProgram(program)
		return glapi.InvalidProgram, fmt.Errorf("link error: %s", log)
	}

	return program, nil
}

// compileShader compiles a single shader
func (sm *ShaderManager) compileShader(source string, shaderType uint32) (glapi.Shader, error) {
	shader := sm.ctx.CreateShader(shaderType)
	sm.ctx.ShaderSource(shader, source)
	sm.ctx.CompileShader(shader)

	if !sm.ctx.GetShaderParameter(shader, glapi.COMPILE_STATUS) {
		log := sm.ctx.GetShaderInfoLog(shader)
		sm.ctx.DeleteShader(shader)
		return glapi.InvalidShader, fmt.Errorf("compile error: %s", log)
	}

	return shader, nil
}
