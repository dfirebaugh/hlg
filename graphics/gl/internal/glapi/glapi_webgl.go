//go:build js && wasm

package glapi

import (
	"errors"
	"math"
	"syscall/js"
)

// Context wraps a WebGL2RenderingContext
type Context struct {
	gl     js.Value
	canvas js.Value
}

// NewContextFromCanvas creates a new WebGL 2.0 context from a canvas element ID
func NewContextFromCanvas(canvasID string) (*Context, error) {
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", canvasID)
	if canvas.IsNull() || canvas.IsUndefined() {
		return nil, errors.New("canvas element not found: " + canvasID)
	}

	// Try to get WebGL 2.0 context
	gl := canvas.Call("getContext", "webgl2")
	if gl.IsNull() || gl.IsUndefined() {
		return nil, errors.New("WebGL 2.0 not supported")
	}

	return &Context{
		gl:     gl,
		canvas: canvas,
	}, nil
}

// GetCanvas returns the canvas element
func (c *Context) GetCanvas() js.Value {
	return c.canvas
}

// GetGL returns the raw WebGL context
func (c *Context) GetGL() js.Value {
	return c.gl
}

// GetCanvasSize returns the current canvas dimensions
func (c *Context) GetCanvasSize() (int, int) {
	width := c.canvas.Get("width").Int()
	height := c.canvas.Get("height").Int()
	return width, height
}

// SetCanvasSize sets the canvas drawing buffer size
func (c *Context) SetCanvasSize(width, height int) {
	c.canvas.Set("width", width)
	c.canvas.Set("height", height)
}

// Buffer operations

func (c *Context) CreateBuffer() Buffer {
	val := c.gl.Call("createBuffer")
	return Buffer(jsValueToUint32(val))
}

func (c *Context) DeleteBuffer(buffer Buffer) {
	c.gl.Call("deleteBuffer", uint32ToJsValue(uint32(buffer)))
}

func (c *Context) BindBuffer(target uint32, buffer Buffer) {
	c.gl.Call("bindBuffer", target, uint32ToJsValue(uint32(buffer)))
}

func (c *Context) BufferData(target uint32, data []byte, usage uint32) {
	if len(data) == 0 {
		c.gl.Call("bufferData", target, 0, usage)
		return
	}
	jsArray := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(jsArray, data)
	c.gl.Call("bufferData", target, jsArray, usage)
}

func (c *Context) BufferDataSize(target uint32, size int, usage uint32) {
	c.gl.Call("bufferData", target, size, usage)
}

func (c *Context) BufferSubData(target uint32, offset int, data []byte) {
	if len(data) == 0 {
		return
	}
	jsArray := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(jsArray, data)
	c.gl.Call("bufferSubData", target, offset, jsArray)
}

// BufferDataFloat32 uploads float32 data to a buffer
func (c *Context) BufferDataFloat32(target uint32, data []float32, usage uint32) {
	if len(data) == 0 {
		c.gl.Call("bufferData", target, 0, usage)
		return
	}
	jsArray := js.Global().Get("Float32Array").New(len(data))
	js.CopyBytesToJS(
		js.Global().Get("Uint8Array").New(jsArray.Get("buffer")),
		float32SliceToBytes(data),
	)
	c.gl.Call("bufferData", target, jsArray, usage)
}

// BufferSubDataFloat32 updates a portion of a buffer with float32 data
func (c *Context) BufferSubDataFloat32(target uint32, offset int, data []float32) {
	if len(data) == 0 {
		return
	}
	jsArray := js.Global().Get("Float32Array").New(len(data))
	js.CopyBytesToJS(
		js.Global().Get("Uint8Array").New(jsArray.Get("buffer")),
		float32SliceToBytes(data),
	)
	c.gl.Call("bufferSubData", target, offset, jsArray)
}

// Shader operations

func (c *Context) CreateShader(shaderType uint32) Shader {
	val := c.gl.Call("createShader", shaderType)
	return Shader(jsValueToUint32(val))
}

func (c *Context) DeleteShader(shader Shader) {
	c.gl.Call("deleteShader", uint32ToJsValue(uint32(shader)))
}

func (c *Context) ShaderSource(shader Shader, source string) {
	c.gl.Call("shaderSource", uint32ToJsValue(uint32(shader)), source)
}

func (c *Context) CompileShader(shader Shader) {
	c.gl.Call("compileShader", uint32ToJsValue(uint32(shader)))
}

func (c *Context) GetShaderParameter(shader Shader, pname uint32) bool {
	return c.gl.Call("getShaderParameter", uint32ToJsValue(uint32(shader)), pname).Bool()
}

func (c *Context) GetShaderInfoLog(shader Shader) string {
	return c.gl.Call("getShaderInfoLog", uint32ToJsValue(uint32(shader))).String()
}

// Program operations

func (c *Context) CreateProgram() Program {
	val := c.gl.Call("createProgram")
	return Program(jsValueToUint32(val))
}

func (c *Context) DeleteProgram(program Program) {
	c.gl.Call("deleteProgram", uint32ToJsValue(uint32(program)))
}

func (c *Context) AttachShader(program Program, shader Shader) {
	c.gl.Call("attachShader", uint32ToJsValue(uint32(program)), uint32ToJsValue(uint32(shader)))
}

func (c *Context) LinkProgram(program Program) {
	c.gl.Call("linkProgram", uint32ToJsValue(uint32(program)))
}

func (c *Context) GetProgramParameter(program Program, pname uint32) bool {
	return c.gl.Call("getProgramParameter", uint32ToJsValue(uint32(program)), pname).Bool()
}

func (c *Context) GetProgramInfoLog(program Program) string {
	return c.gl.Call("getProgramInfoLog", uint32ToJsValue(uint32(program))).String()
}

func (c *Context) UseProgram(program Program) {
	c.gl.Call("useProgram", uint32ToJsValue(uint32(program)))
}

// Uniform operations

func (c *Context) GetUniformLocation(program Program, name string) UniformLocation {
	val := c.gl.Call("getUniformLocation", uint32ToJsValue(uint32(program)), name)
	if val.IsNull() || val.IsUndefined() {
		return InvalidUniformLocation
	}
	return UniformLocation(jsValueToUint32(val))
}

func (c *Context) Uniform1i(location UniformLocation, v int) {
	c.gl.Call("uniform1i", uint32ToJsValue(uint32(location)), v)
}

func (c *Context) Uniform1f(location UniformLocation, v float32) {
	c.gl.Call("uniform1f", uint32ToJsValue(uint32(location)), v)
}

func (c *Context) Uniform2f(location UniformLocation, v0, v1 float32) {
	c.gl.Call("uniform2f", uint32ToJsValue(uint32(location)), v0, v1)
}

func (c *Context) Uniform3f(location UniformLocation, v0, v1, v2 float32) {
	c.gl.Call("uniform3f", uint32ToJsValue(uint32(location)), v0, v1, v2)
}

func (c *Context) Uniform4f(location UniformLocation, v0, v1, v2, v3 float32) {
	c.gl.Call("uniform4f", uint32ToJsValue(uint32(location)), v0, v1, v2, v3)
}

func (c *Context) Uniform2fv(location UniformLocation, data []float32) {
	jsArray := js.Global().Get("Float32Array").New(len(data))
	js.CopyBytesToJS(
		js.Global().Get("Uint8Array").New(jsArray.Get("buffer")),
		float32SliceToBytes(data),
	)
	c.gl.Call("uniform2fv", uint32ToJsValue(uint32(location)), jsArray)
}

func (c *Context) Uniform4fv(location UniformLocation, data []float32) {
	jsArray := js.Global().Get("Float32Array").New(len(data))
	js.CopyBytesToJS(
		js.Global().Get("Uint8Array").New(jsArray.Get("buffer")),
		float32SliceToBytes(data),
	)
	c.gl.Call("uniform4fv", uint32ToJsValue(uint32(location)), jsArray)
}

func (c *Context) UniformMatrix4fv(location UniformLocation, transpose bool, data []float32) {
	jsArray := js.Global().Get("Float32Array").New(len(data))
	js.CopyBytesToJS(
		js.Global().Get("Uint8Array").New(jsArray.Get("buffer")),
		float32SliceToBytes(data),
	)
	c.gl.Call("uniformMatrix4fv", uint32ToJsValue(uint32(location)), transpose, jsArray)
}

// Vertex Array operations

func (c *Context) CreateVertexArray() VertexArray {
	val := c.gl.Call("createVertexArray")
	return VertexArray(jsValueToUint32(val))
}

func (c *Context) DeleteVertexArray(vao VertexArray) {
	c.gl.Call("deleteVertexArray", uint32ToJsValue(uint32(vao)))
}

func (c *Context) BindVertexArray(vao VertexArray) {
	c.gl.Call("bindVertexArray", uint32ToJsValue(uint32(vao)))
}

func (c *Context) UnbindVertexArray() {
	c.gl.Call("bindVertexArray", nil)
}

// Vertex attribute operations

func (c *Context) VertexAttribPointer(index uint32, size int, dataType uint32, normalized bool, stride int, offset int) {
	c.gl.Call("vertexAttribPointer", index, size, dataType, normalized, stride, offset)
}

func (c *Context) EnableVertexAttribArray(index uint32) {
	c.gl.Call("enableVertexAttribArray", index)
}

func (c *Context) DisableVertexAttribArray(index uint32) {
	c.gl.Call("disableVertexAttribArray", index)
}

// Texture operations

func (c *Context) CreateTexture() Texture {
	val := c.gl.Call("createTexture")
	return Texture(jsValueToUint32(val))
}

func (c *Context) DeleteTexture(texture Texture) {
	c.gl.Call("deleteTexture", uint32ToJsValue(uint32(texture)))
}

func (c *Context) BindTexture(target uint32, texture Texture) {
	c.gl.Call("bindTexture", target, uint32ToJsValue(uint32(texture)))
}

func (c *Context) ActiveTexture(unit uint32) {
	c.gl.Call("activeTexture", unit)
}

func (c *Context) TexParameteri(target, pname, param uint32) {
	c.gl.Call("texParameteri", target, pname, param)
}

func (c *Context) PixelStorei(pname uint32, param int) {
	c.gl.Call("pixelStorei", pname, param)
}

func (c *Context) TexImage2D(target uint32, level int, internalFormat uint32, width, height int, border int, format, dataType uint32, data []byte) {
	if len(data) == 0 {
		c.gl.Call("texImage2D", target, level, internalFormat, width, height, border, format, dataType, nil)
		return
	}
	jsArray := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(jsArray, data)
	c.gl.Call("texImage2D", target, level, internalFormat, width, height, border, format, dataType, jsArray)
}

func (c *Context) TexSubImage2D(target uint32, level, xoffset, yoffset, width, height int, format, dataType uint32, data []byte) {
	jsArray := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(jsArray, data)
	c.gl.Call("texSubImage2D", target, level, xoffset, yoffset, width, height, format, dataType, jsArray)
}

// Drawing operations

func (c *Context) DrawArrays(mode uint32, first, count int) {
	c.gl.Call("drawArrays", mode, first, count)
}

func (c *Context) DrawElements(mode uint32, count int, dataType uint32, offset int) {
	c.gl.Call("drawElements", mode, count, dataType, offset)
}

// State operations

func (c *Context) Enable(cap uint32) {
	c.gl.Call("enable", cap)
}

func (c *Context) Disable(cap uint32) {
	c.gl.Call("disable", cap)
}

func (c *Context) BlendFunc(sfactor, dfactor uint32) {
	c.gl.Call("blendFunc", sfactor, dfactor)
}

func (c *Context) Viewport(x, y, width, height int) {
	c.gl.Call("viewport", x, y, width, height)
}

func (c *Context) Clear(mask uint32) {
	c.gl.Call("clear", mask)
}

func (c *Context) ClearColor(r, g, b, a float32) {
	c.gl.Call("clearColor", r, g, b, a)
}

func (c *Context) Scissor(x, y, width, height int) {
	c.gl.Call("scissor", x, y, width, height)
}

func (c *Context) GetViewport() [4]int {
	// WebGL doesn't have getIntegerv, use canvas size as viewport
	w, h := c.GetCanvasSize()
	return [4]int{0, 0, w, h}
}

// Helper functions

// jsValueToUint32 converts a js.Value to uint32 for storing GL objects
// We use a counter-based approach since we can't safely convert js.Value to uint32
var (
	jsValueCounter uint32 = 1
	jsValueMap            = make(map[uint32]js.Value)
)

func jsValueToUint32(val js.Value) uint32 {
	if val.IsNull() || val.IsUndefined() {
		return 0
	}
	id := jsValueCounter
	jsValueCounter++
	jsValueMap[id] = val
	return id
}

func uint32ToJsValue(id uint32) js.Value {
	if id == 0 {
		return js.Null()
	}
	val, ok := jsValueMap[id]
	if !ok {
		return js.Null()
	}
	return val
}

// Helper to convert float32 slice to bytes
func float32SliceToBytes(data []float32) []byte {
	bytes := make([]byte, len(data)*4)
	for i, f := range data {
		u := math.Float32bits(f)
		bytes[i*4+0] = byte(u)
		bytes[i*4+1] = byte(u >> 8)
		bytes[i*4+2] = byte(u >> 16)
		bytes[i*4+3] = byte(u >> 24)
	}
	return bytes
}
