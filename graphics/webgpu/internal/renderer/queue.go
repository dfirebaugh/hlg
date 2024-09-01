package renderer

import (
	"image"
	"image/color"
	"log"
	"unsafe"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/context"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/pipelines"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/primitives"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/shader"
	"github.com/dfirebaugh/hlg/gui"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type RenderQueue struct {
	surface context.Surface
	*wgpu.Device
	*wgpu.SwapChainDescriptor
	context.RenderContext
	*shader.ShaderManager

	*pipelines.PrimitiveBuffer

	Textures     map[textureHandle]*Texture
	queue        []graphics.Renderable
	currentFrame []graphics.Renderable

	Priority    int
	shouldClear bool
}

func NewRenderQueue(surface context.Surface, d *wgpu.Device, scd *wgpu.SwapChainDescriptor) *RenderQueue {
	rq := &RenderQueue{
		surface:             surface,
		Device:              d,
		SwapChainDescriptor: scd,
		Textures:            make(map[textureHandle]*Texture),
		currentFrame:        []graphics.Renderable{},
		queue:               []graphics.Renderable{},
	}

	shaderManager := shader.NewShaderManager(d)
	rq.RenderContext = context.NewRenderContext(surface, d, scd, rq, shaderManager)
	rq.ShaderManager = shaderManager
  
	vertexLayout := graphics.VertexBufferLayout{
		ArrayStride: 4*3 + 4*2 + 4 + 4 + 4*4,
		Attributes: []graphics.VertexAttributeLayout{
			{
				ShaderLocation: 0,
				Offset:         0,
				Format:         "float32x3", // Position (Clip space)
			},
			{
				ShaderLocation: 1,
				Offset:         3 * 4,
				Format:         "float32x2", // Local Position
			},
			{
				ShaderLocation: 2,
				Offset:         (3 + 2) * 4,
				Format:         "float32", // OpCode
			},
			{
				ShaderLocation: 3,
				Offset:         (3 + 2 + 1) * 4,
				Format:         "float32", // Radius
			},
			{
				ShaderLocation: 4,
				Offset:         (3 + 2 + 1 + 1) * 4,
				Format:         "float32x4", // Color
			},
		},
	}
	rq.PrimitiveBuffer = pipelines.NewPrimitiveBuffer(rq.RenderContext, nil, vertexLayout)

	return rq
}

func (rq *RenderQueue) SetPriority(priority int) {
	rq.Priority = priority
}

func (rq *RenderQueue) SetShouldClear(shouldClear bool) {
	rq.shouldClear = shouldClear
}

func (rq *RenderQueue) RenderClear() {
	rq.currentFrame = nil
	rq.queue = nil
}

func (rq *RenderQueue) AddToRenderQueue(r graphics.Renderable) {
	rq.queue = append(rq.queue, r)
}

func (rq *RenderQueue) PrepareFrame() {
	rq.currentFrame = make([]graphics.Renderable, len(rq.queue))
	copy(rq.currentFrame, rq.queue)
}

func (rq *RenderQueue) RenderFrame(pass *wgpu.RenderPassEncoder) {
	rq.PrimitiveBuffer.RenderPass(pass)
	for _, renderable := range rq.currentFrame {
		if renderable != nil {
			renderable.RenderPass(pass)
		}
	}
}

func (rq *RenderQueue) CreateTextureFromImage(img image.Image) (graphics.Texture, error) {
	tex := NewTexture(rq.RenderContext, img, rq)
	handle := uintptr(unsafe.Pointer(tex))
	tex.SetHandle(textureHandle(handle))
	rq.Textures[textureHandle(handle)] = tex
	return tex, nil
}

func (rq *RenderQueue) UpdateTextureFromImage(texture graphics.Texture, img image.Image) {
	texture.UpdateImage(img)
}

func (rq *RenderQueue) DisposeTexture(h uintptr) {
	rq.Textures[textureHandle(h)].gpuTexture.Destroy()
	delete(rq.Textures, textureHandle(h))
}

// AddTriangle creates a new Triangle renderable and adds it to the RenderQueue.
// It returns a reference to the created Triangle.
func (rq *RenderQueue) AddTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) graphics.Shape {
	if rq == nil {
		log.Fatal("RenderQueue is nil")
	}
	if rq.RenderContext == nil {
		log.Fatal("RenderContext in RenderQueue is nil")
	}

	vertices := primitives.MakeTriangle(x1, y1, x2, y2, x3, y3, c)
	triangle := pipelines.NewPolygon(rq.RenderContext, vertices)

	return triangle
}

// AddRectangle creates a new Rectangle renderable and adds it to the RenderQueue.
// It returns a reference to the created Rectangle.
func (rq *RenderQueue) AddRectangle(x, y, width, height int, c color.Color) graphics.Shape {
	rectangleVertices := primitives.MakeRectangle(x, y, width, height, c)
	rectangle := pipelines.NewPolygon(rq.RenderContext, rectangleVertices)
	return rectangle
}

// AddCircle creates a new Circle renderable and adds it to the RenderQueue.
// It returns a reference to the created Circle.
func (rq *RenderQueue) AddCircle(cx, cy int, radius float32, c color.Color, segments int) graphics.Shape {
	// note: we could probably more efficiently draw circles with a custom shader -- but this is a good start
	circle := rq.AddPolygon(cx, cy, radius*2, c, segments)
	return circle
}

func (rq *RenderQueue) AddPolygonFromVertices(cx, cy int, width float32, vertices []graphics.Vertex) graphics.Shape {
	v := primitives.MakePolygonFromVertices(cx, cy, width, vertices)
	polygon := pipelines.NewPolygon(rq.RenderContext, v)
	return polygon
}

// AddPolygon creates a new Polygon renderable and adds it to the RenderQueue.
// It returns a reference to the created Polygon.
func (rq *RenderQueue) AddPolygon(cx, cy int, width float32, c color.Color, sides int) graphics.Shape {
	vertices := primitives.MakePolygon(cx, cy, width, c, sides)
	polygon := pipelines.NewPolygon(rq.RenderContext, vertices)
	return polygon
}

// AddLine creates a new Line renderable and adds it to the RenderQueue.
// It returns a reference to the created Line.
func (rq *RenderQueue) AddLine(x1, y1, x2, y2 int, width float32, c color.Color) graphics.Shape {
	lineVertices := primitives.MakeLine(x1, y1, x2, y2, width, c)
	line := pipelines.NewPolygon(rq.RenderContext, lineVertices)
	return line
}

func (rq *RenderQueue) DrawPrimitiveBuffer(vertices []gui.Vertex) {
	if len(vertices) == 0 {
		return
	}
	rq.PrimitiveBuffer.UpdateVertexBuffer(vertices)
}

func (rq *RenderQueue) AddDynamicRenderable(vertexData []byte, layout graphics.VertexBufferLayout, shaderHandle int, uniforms map[string]graphics.Uniform, dataMap map[string][]byte) graphics.ShaderRenderable {
	if rq == nil {
		log.Println("RenderQueue is nil, cannot add to queue")
		return nil
	}

	u := convertUniforms(rq.Device, uniforms, dataMap)
	r := pipelines.NewRenderable(rq.RenderContext, vertexData, layout, shaderHandle, u)

	return r
}

func convertVertex(gv graphics.Vertex) primitives.Vertex {
	return primitives.Vertex{
		Position:  gv.Position,
		Color:     gv.Color,
		TexCoords: [2]float32{},
	}
}

func convertVertices(gvs []graphics.Vertex) []primitives.Vertex {
	pvs := make([]primitives.Vertex, len(gvs))
	for i, gv := range gvs {
		pvs[i] = convertVertex(gv)
	}
	return pvs
}

func convertUniform(device *wgpu.Device, gu graphics.Uniform, data []byte) pipelines.Uniform {
	buffer, err := device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Uniform Buffer",
		Contents: data,
		Usage:    wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		panic(err)
	}

	return pipelines.Uniform{
		Binding: gu.Binding,
		Buffer:  buffer,
		Size:    gu.Size,
	}
}

func convertUniforms(device *wgpu.Device, gus map[string]graphics.Uniform, dataMap map[string][]byte) map[string]pipelines.Uniform {
	rus := make(map[string]pipelines.Uniform)
	for name, gu := range gus {
		data, exists := dataMap[name]
		if !exists {
			panic("Uniform data not provided for " + name)
		}
		rus[name] = convertUniform(device, gu, data)
	}
	return rus
}
