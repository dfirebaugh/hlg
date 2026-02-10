//go:build !js

package pipelines

import (
	"encoding/binary"
	"log"
	"math"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/glapi"
)

// ShaderRenderable implements graphics.ShaderRenderable and glRenderable for OpenGL.
type ShaderRenderable struct {
	ctx          *glapi.Context
	vao          glapi.VertexArray
	vbo          glapi.Buffer
	program      glapi.Program
	vertexData   []byte
	arrayStride  uint64
	uniforms     map[string]uniformInfo
	shouldRender bool
	isDisposed   bool
	rq           RenderQueue
}

type uniformInfo struct {
	location glapi.UniformLocation
	size     uint64
	data     []byte
}

// NewShaderRenderable creates a new ShaderRenderable for the GL backend.
func NewShaderRenderable(ctx *glapi.Context, rq RenderQueue, vertexData []byte, layout graphics.VertexBufferLayout, program glapi.Program, uniforms map[string]graphics.Uniform, dataMap map[string][]byte) *ShaderRenderable {
	r := &ShaderRenderable{
		ctx:         ctx,
		program:     program,
		vertexData:  vertexData,
		arrayStride: layout.ArrayStride,
		uniforms:    make(map[string]uniformInfo),
		rq:          rq,
	}

	// Create VAO and VBO
	r.vao = ctx.CreateVertexArray()
	ctx.BindVertexArray(r.vao)

	r.vbo = ctx.CreateBuffer()
	ctx.BindBuffer(glapi.ARRAY_BUFFER, r.vbo)
	ctx.BufferData(glapi.ARRAY_BUFFER, vertexData, glapi.STATIC_DRAW)

	// Set up vertex attributes from layout
	for _, attr := range layout.Attributes {
		size, glType := translateFormat(attr.Format)
		ctx.VertexAttribPointer(attr.ShaderLocation, size, glType, false, int(layout.ArrayStride), int(attr.Offset))
		ctx.EnableVertexAttribArray(attr.ShaderLocation)
	}

	ctx.UnbindVertexArray()

	// Resolve uniform locations and store initial data
	for name, u := range uniforms {
		loc := ctx.GetUniformLocation(program, name)
		info := uniformInfo{
			location: loc,
			size:     u.Size,
		}
		if data, ok := dataMap[name]; ok {
			info.data = make([]byte, len(data))
			copy(info.data, data)
		}
		r.uniforms[name] = info
	}

	return r
}

// translateFormat converts a vertex format string to GL size and type.
func translateFormat(format string) (size int, glType uint32) {
	switch format {
	case "float32":
		return 1, glapi.FLOAT
	case "float32x2":
		return 2, glapi.FLOAT
	case "float32x3":
		return 3, glapi.FLOAT
	case "float32x4":
		return 4, glapi.FLOAT
	default:
		log.Fatalf("Unknown vertex format: %s", format)
		return 4, glapi.FLOAT
	}
}

// GLRender performs the actual OpenGL draw call.
func (r *ShaderRenderable) GLRender() {
	if !r.shouldRender || r.isDisposed {
		return
	}
	if r.vao == 0 {
		return
	}

	r.ctx.UseProgram(r.program)

	// Upload uniforms
	for _, u := range r.uniforms {
		if u.data == nil || u.location == glapi.InvalidUniformLocation {
			continue
		}
		uploadUniform(r.ctx, u.location, u.size, u.data)
	}

	r.ctx.BindVertexArray(r.vao)

	vertexCount := len(r.vertexData) / int(r.arrayStride)
	r.ctx.DrawArrays(glapi.TRIANGLES, 0, vertexCount)

	r.ctx.UnbindVertexArray()
}

// uploadUniform uploads uniform data based on its size.
func uploadUniform(ctx *glapi.Context, loc glapi.UniformLocation, size uint64, data []byte) {
	switch size {
	case 4:
		ctx.Uniform1f(loc, bytesToFloat32(data))
	case 8:
		ctx.Uniform2f(loc, bytesToFloat32(data[0:4]), bytesToFloat32(data[4:8]))
	case 12:
		ctx.Uniform3f(loc, bytesToFloat32(data[0:4]), bytesToFloat32(data[4:8]), bytesToFloat32(data[8:12]))
	case 16:
		ctx.Uniform4f(loc, bytesToFloat32(data[0:4]), bytesToFloat32(data[4:8]), bytesToFloat32(data[8:12]), bytesToFloat32(data[12:16]))
	case 64:
		floats := make([]float32, 16)
		for i := 0; i < 16; i++ {
			floats[i] = bytesToFloat32(data[i*4 : (i+1)*4])
		}
		ctx.UniformMatrix4fv(loc, false, floats)
	default:
		log.Printf("Unsupported uniform size: %d", size)
	}
}

func bytesToFloat32(b []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(b))
}

// Render marks this renderable for drawing and adds it to the render queue.
func (r *ShaderRenderable) Render() {
	if r.isDisposed {
		return
	}
	r.shouldRender = true
	r.rq.AddToRenderQueue(r)
}

// UpdateUniform updates a single uniform's cached data.
func (r *ShaderRenderable) UpdateUniform(name string, data []byte) {
	if u, ok := r.uniforms[name]; ok {
		u.data = make([]byte, len(data))
		copy(u.data, data)
		r.uniforms[name] = u
	} else {
		log.Printf("Uniform %s does not exist", name)
	}
}

// UpdateUniforms updates multiple uniforms' cached data.
func (r *ShaderRenderable) UpdateUniforms(dataMap map[string][]byte) {
	for name, data := range dataMap {
		r.UpdateUniform(name, data)
	}
}

// Dispose releases GL resources.
func (r *ShaderRenderable) Dispose() {
	if r.isDisposed {
		return
	}
	r.isDisposed = true

	if r.vao != 0 {
		r.ctx.DeleteVertexArray(r.vao)
		r.vao = 0
	}
	if r.vbo != 0 {
		r.ctx.DeleteBuffer(r.vbo)
		r.vbo = 0
	}
}

// IsDisposed returns whether this renderable has been disposed.
func (r *ShaderRenderable) IsDisposed() bool {
	return r.isDisposed
}

var _ graphics.ShaderRenderable = (*ShaderRenderable)(nil)
