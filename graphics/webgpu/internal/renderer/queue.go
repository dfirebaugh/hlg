package renderer

import (
	"image"
	"image/color"
	"log"
	"unsafe"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/context"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/primitives"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/renderable"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/shapes"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type RenderQueue struct {
	surface context.Surface
	*wgpu.Device
	*wgpu.SwapChainDescriptor
	context.RenderContext

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

	rq.RenderContext = context.NewRenderContext(surface, d, scd, rq)

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
	triangle := shapes.NewPolygon(rq.RenderContext, vertices)

	return triangle
}

// AddRectangle creates a new Rectangle renderable and adds it to the RenderQueue.
// It returns a reference to the created Rectangle.
func (rq *RenderQueue) AddRectangle(x, y, width, height int, c color.Color) graphics.Shape {
	rectangleVertices := primitives.MakeRectangle(x, y, width, height, c)
	rectangle := shapes.NewPolygon(rq.RenderContext, rectangleVertices)
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
	polygon := shapes.NewPolygon(rq.RenderContext, v)
	return polygon
}

// AddPolygon creates a new Polygon renderable and adds it to the RenderQueue.
// It returns a reference to the created Polygon.
func (rq *RenderQueue) AddPolygon(cx, cy int, width float32, c color.Color, sides int) graphics.Shape {
	vertices := primitives.MakePolygon(cx, cy, width, c, sides)
	polygon := shapes.NewPolygon(rq.RenderContext, vertices)
	return polygon
}

// AddLine creates a new Line renderable and adds it to the RenderQueue.
// It returns a reference to the created Line.
func (rq *RenderQueue) AddLine(x1, y1, x2, y2 int, width float32, c color.Color) graphics.Shape {
	lineVertices := primitives.MakeLine(x1, y1, x2, y2, width, c)
	line := shapes.NewPolygon(rq.RenderContext, lineVertices)
	return line
}

func (rq *RenderQueue) AddDynamicRenderable(vertices []graphics.Vertex, shaderCode string, uniforms map[string]graphics.Uniform, dataMap map[string][]byte) graphics.ShaderRenderable {
	if rq == nil {
		log.Println("RenderQueue is nil, cannot add to queue")
		return nil
	}
	primitivesVertices := convertVertices(vertices)
	renderableUniforms := convertUniforms(rq.Device, uniforms, dataMap)
	r := renderable.NewRenderable(rq.RenderContext, primitivesVertices, shaderCode, renderableUniforms)

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

func convertUniform(device *wgpu.Device, gu graphics.Uniform, data []byte) renderable.Uniform {
	buffer, err := device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Uniform Buffer",
		Contents: data,
		Usage:    wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		panic(err)
	}

	return renderable.Uniform{
		Binding: gu.Binding,
		Buffer:  buffer,
		Size:    gu.Size,
	}
}

func convertUniforms(device *wgpu.Device, gus map[string]graphics.Uniform, dataMap map[string][]byte) map[string]renderable.Uniform {
	rus := make(map[string]renderable.Uniform)
	for name, gu := range gus {
		data, exists := dataMap[name]
		if !exists {
			panic("Uniform data not provided for " + name)
		}
		rus[name] = convertUniform(device, gu, data)
	}
	return rus
}
