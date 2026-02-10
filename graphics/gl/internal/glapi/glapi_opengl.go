//go:build !js

package glapi

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// Context wraps OpenGL state for desktop platforms
type Context struct {
	// OpenGL doesn't need a context wrapper, but we use this for API consistency
}

// NewContext creates a new OpenGL context wrapper
// Note: OpenGL must be initialized via gl.Init() before calling this
func NewContext() *Context {
	return &Context{}
}

// Buffer operations

func (c *Context) CreateBuffer() Buffer {
	var buf uint32
	gl.GenBuffers(1, &buf)
	return Buffer(buf)
}

func (c *Context) DeleteBuffer(buffer Buffer) {
	buf := uint32(buffer)
	gl.DeleteBuffers(1, &buf)
}

func (c *Context) BindBuffer(target uint32, buffer Buffer) {
	gl.BindBuffer(target, uint32(buffer))
}

func (c *Context) BufferData(target uint32, data []byte, usage uint32) {
	if len(data) == 0 {
		gl.BufferData(target, 0, nil, usage)
		return
	}
	gl.BufferData(target, len(data), gl.Ptr(data), usage)
}

func (c *Context) BufferDataSize(target uint32, size int, usage uint32) {
	gl.BufferData(target, size, nil, usage)
}

func (c *Context) BufferSubData(target uint32, offset int, data []byte) {
	if len(data) == 0 {
		return
	}
	gl.BufferSubData(target, offset, len(data), gl.Ptr(data))
}

// Shader operations

func (c *Context) CreateShader(shaderType uint32) Shader {
	return Shader(gl.CreateShader(shaderType))
}

func (c *Context) DeleteShader(shader Shader) {
	gl.DeleteShader(uint32(shader))
}

func (c *Context) ShaderSource(shader Shader, source string) {
	csource, free := gl.Strs(source + "\x00")
	gl.ShaderSource(uint32(shader), 1, csource, nil)
	free()
}

func (c *Context) CompileShader(shader Shader) {
	gl.CompileShader(uint32(shader))
}

func (c *Context) GetShaderParameter(shader Shader, pname uint32) bool {
	var status int32
	gl.GetShaderiv(uint32(shader), pname, &status)
	return status != gl.FALSE
}

func (c *Context) GetShaderInfoLog(shader Shader) string {
	var logLength int32
	gl.GetShaderiv(uint32(shader), gl.INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		return ""
	}
	log := make([]byte, logLength)
	gl.GetShaderInfoLog(uint32(shader), logLength, nil, &log[0])
	return string(log[:logLength-1]) // Exclude null terminator
}

// Program operations

func (c *Context) CreateProgram() Program {
	return Program(gl.CreateProgram())
}

func (c *Context) DeleteProgram(program Program) {
	gl.DeleteProgram(uint32(program))
}

func (c *Context) AttachShader(program Program, shader Shader) {
	gl.AttachShader(uint32(program), uint32(shader))
}

func (c *Context) LinkProgram(program Program) {
	gl.LinkProgram(uint32(program))
}

func (c *Context) GetProgramParameter(program Program, pname uint32) bool {
	var status int32
	gl.GetProgramiv(uint32(program), pname, &status)
	return status != gl.FALSE
}

func (c *Context) GetProgramInfoLog(program Program) string {
	var logLength int32
	gl.GetProgramiv(uint32(program), gl.INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		return ""
	}
	log := make([]byte, logLength)
	gl.GetProgramInfoLog(uint32(program), logLength, nil, &log[0])
	return string(log[:logLength-1])
}

func (c *Context) UseProgram(program Program) {
	gl.UseProgram(uint32(program))
}

// Uniform operations

func (c *Context) GetUniformLocation(program Program, name string) UniformLocation {
	return UniformLocation(gl.GetUniformLocation(uint32(program), gl.Str(name+"\x00")))
}

func (c *Context) Uniform1i(location UniformLocation, v int) {
	gl.Uniform1i(int32(location), int32(v))
}

func (c *Context) Uniform1f(location UniformLocation, v float32) {
	gl.Uniform1f(int32(location), v)
}

func (c *Context) Uniform2f(location UniformLocation, v0, v1 float32) {
	gl.Uniform2f(int32(location), v0, v1)
}

func (c *Context) Uniform3f(location UniformLocation, v0, v1, v2 float32) {
	gl.Uniform3f(int32(location), v0, v1, v2)
}

func (c *Context) Uniform4f(location UniformLocation, v0, v1, v2, v3 float32) {
	gl.Uniform4f(int32(location), v0, v1, v2, v3)
}

func (c *Context) Uniform2fv(location UniformLocation, data []float32) {
	gl.Uniform2fv(int32(location), 1, &data[0])
}

func (c *Context) Uniform4fv(location UniformLocation, data []float32) {
	gl.Uniform4fv(int32(location), 1, &data[0])
}

func (c *Context) UniformMatrix4fv(location UniformLocation, transpose bool, data []float32) {
	gl.UniformMatrix4fv(int32(location), 1, transpose, &data[0])
}

// Vertex Array operations

func (c *Context) CreateVertexArray() VertexArray {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	return VertexArray(vao)
}

func (c *Context) DeleteVertexArray(vao VertexArray) {
	v := uint32(vao)
	gl.DeleteVertexArrays(1, &v)
}

func (c *Context) BindVertexArray(vao VertexArray) {
	gl.BindVertexArray(uint32(vao))
}

func (c *Context) UnbindVertexArray() {
	gl.BindVertexArray(0)
}

// Vertex attribute operations

func (c *Context) VertexAttribPointer(index uint32, size int, dataType uint32, normalized bool, stride int, offset int) {
	gl.VertexAttribPointerWithOffset(index, int32(size), dataType, normalized, int32(stride), uintptr(offset))
}

func (c *Context) EnableVertexAttribArray(index uint32) {
	gl.EnableVertexAttribArray(index)
}

func (c *Context) DisableVertexAttribArray(index uint32) {
	gl.DisableVertexAttribArray(index)
}

// Texture operations

func (c *Context) CreateTexture() Texture {
	var tex uint32
	gl.GenTextures(1, &tex)
	return Texture(tex)
}

func (c *Context) DeleteTexture(texture Texture) {
	tex := uint32(texture)
	gl.DeleteTextures(1, &tex)
}

func (c *Context) BindTexture(target uint32, texture Texture) {
	gl.BindTexture(target, uint32(texture))
}

func (c *Context) ActiveTexture(unit uint32) {
	gl.ActiveTexture(unit)
}

func (c *Context) TexParameteri(target, pname, param uint32) {
	gl.TexParameteri(target, pname, int32(param))
}

func (c *Context) TexImage2D(target uint32, level int, internalFormat uint32, width, height int, border int, format, dataType uint32, data []byte) {
	if len(data) == 0 {
		gl.TexImage2D(target, int32(level), int32(internalFormat), int32(width), int32(height), int32(border), format, dataType, nil)
		return
	}
	gl.TexImage2D(target, int32(level), int32(internalFormat), int32(width), int32(height), int32(border), format, dataType, gl.Ptr(data))
}

func (c *Context) TexSubImage2D(target uint32, level, xoffset, yoffset, width, height int, format, dataType uint32, data []byte) {
	gl.TexSubImage2D(target, int32(level), int32(xoffset), int32(yoffset), int32(width), int32(height), format, dataType, gl.Ptr(data))
}

// Drawing operations

func (c *Context) DrawArrays(mode uint32, first, count int) {
	gl.DrawArrays(mode, int32(first), int32(count))
}

func (c *Context) DrawElements(mode uint32, count int, dataType uint32, offset int) {
	gl.DrawElementsWithOffset(mode, int32(count), dataType, uintptr(offset))
}

// State operations

func (c *Context) Enable(cap uint32) {
	gl.Enable(cap)
}

func (c *Context) Disable(cap uint32) {
	gl.Disable(cap)
}

func (c *Context) BlendFunc(sfactor, dfactor uint32) {
	gl.BlendFunc(sfactor, dfactor)
}

func (c *Context) Viewport(x, y, width, height int) {
	gl.Viewport(int32(x), int32(y), int32(width), int32(height))
}

func (c *Context) Clear(mask uint32) {
	gl.Clear(mask)
}

func (c *Context) ClearColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

func (c *Context) Scissor(x, y, width, height int) {
	gl.Scissor(int32(x), int32(y), int32(width), int32(height))
}

func (c *Context) GetViewport() [4]int {
	var viewport [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &viewport[0])
	return [4]int{int(viewport[0]), int(viewport[1]), int(viewport[2]), int(viewport[3])}
}

// Ptr returns a GL-compatible pointer to data
func Ptr(data interface{}) unsafe.Pointer {
	return gl.Ptr(data)
}
