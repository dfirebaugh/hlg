//go:build !js

package renderer

import (
	"image"
	"image/color"
	"log"
	"math"
	"unsafe"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/context"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/pipelines"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/shader"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

// wgpuRenderable is an internal interface for types that can render via wgpu render pass.
// This is used internally by the webgpu backend; the public graphics.Renderable interface
// does not expose this method.
type wgpuRenderable interface {
	RenderPass(pass *wgpu.RenderPassEncoder)
}

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

	// Callback invoked before adding a renderable to flush pending batched primitives
	onBeforeAddToQueue func()
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

	// Using storage buffer approach - no vertex layout needed
	// Primitives are stored in storage buffer and vertices are constructed in shader
	rq.PrimitiveBuffer = pipelines.NewPrimitiveBuffer(rq.RenderContext, nil)

	return rq
}

func (rq *RenderQueue) SetPriority(priority int) {
	rq.Priority = priority
}

func (rq *RenderQueue) SetShouldClear(shouldClear bool) {
	rq.shouldClear = shouldClear
}

func (rq *RenderQueue) RenderClear() {
	// Clear slices while preserving underlying capacity
	rq.currentFrame = rq.currentFrame[:0]
	rq.queue = rq.queue[:0]
}

func (rq *RenderQueue) AddToRenderQueue(r graphics.Renderable) {
	// Flush pending batched primitives before adding a new renderable
	// This preserves draw order between batched primitives and Shapes
	if rq.onBeforeAddToQueue != nil {
		rq.onBeforeAddToQueue()
	}
	rq.queue = append(rq.queue, r)
}

// SetOnBeforeAddToQueue sets a callback that's invoked before adding a renderable to the queue.
// This is used to flush pending batched primitives to preserve draw order.
func (rq *RenderQueue) SetOnBeforeAddToQueue(callback func()) {
	rq.onBeforeAddToQueue = callback
}

func (rq *RenderQueue) PrepareFrame() {
	// Reuse currentFrame capacity if possible to avoid allocation
	if cap(rq.currentFrame) >= len(rq.queue) {
		rq.currentFrame = rq.currentFrame[:len(rq.queue)]
	} else {
		rq.currentFrame = make([]graphics.Renderable, len(rq.queue))
	}
	copy(rq.currentFrame, rq.queue)
}

func (rq *RenderQueue) RenderFrame(pass *wgpu.RenderPassEncoder) {
	rq.PrimitiveBuffer.RenderPass(pass)
	for _, renderable := range rq.currentFrame {
		if renderable != nil {
			// Use type assertion to call the internal RenderPass method
			if wr, ok := renderable.(wgpuRenderable); ok {
				wr.RenderPass(pass)
			}
		}
	}
}

// Present renders this queue's contents immediately.
// Note: WebGPU backend requires render pass management - this is a no-op.
// Use the standard rendering flow for WebGPU.
func (rq *RenderQueue) Present() {
	log.Println("RenderQueue.Present() is not fully supported in WebGPU backend")
}

func (rq *RenderQueue) CreateTextureFromImage(img image.Image) (graphics.Texture, error) {
	tex := NewTexture(rq.RenderContext, img, rq)
	handle := uintptr(unsafe.Pointer(tex))
	tex.SetHandle(textureHandle(handle))
	rq.Textures[textureHandle(handle)] = tex
	return tex, nil
}

func (rq *RenderQueue) UpdateTextureFromImage(texture graphics.Texture, img image.Image) {
	_ = texture.UpdateImage(img)
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

	sw, sh := rq.GetSurfaceSize()
	vertices := graphics.MakeSolidTriangle(x1, y1, x2, y2, x3, y3, c, sw, sh)
	screenPos := [][2]float32{
		{float32(x1), float32(y1)},
		{float32(x2), float32(y2)},
		{float32(x3), float32(y3)},
	}

	return pipelines.NewPrimitiveShape(rq.PrimitiveBuffer, vertices, screenPos)
}

// AddRectangle creates a new Rectangle renderable and adds it to the RenderQueue.
// It returns a reference to the created Rectangle.
func (rq *RenderQueue) AddRectangle(x, y, width, height int, c color.Color) graphics.Shape {
	sw, sh := rq.GetSurfaceSize()
	vertices := graphics.MakeSolidRectangle(x, y, width, height, c, sw, sh)
	// Rectangle has 6 vertices (2 triangles), but we store 4 corner positions for transforms
	screenPos := [][2]float32{
		{float32(x), float32(y)},                  // top-left
		{float32(x), float32(y + height)},         // bottom-left
		{float32(x + width), float32(y)},          // top-right
		{float32(x), float32(y + height)},         // bottom-left (repeated for 2nd triangle)
		{float32(x + width), float32(y + height)}, // bottom-right
		{float32(x + width), float32(y)},          // top-right (repeated for 2nd triangle)
	}

	return pipelines.NewPrimitiveShape(rq.PrimitiveBuffer, vertices, screenPos)
}

// AddRoundedRectangle creates a new rounded rectangle renderable and adds it to the RenderQueue.
// It returns a reference to the created shape.
func (rq *RenderQueue) AddRoundedRectangle(x, y, width, height, radius int, c color.Color) graphics.Shape {
	sw, sh := rq.GetSurfaceSize()
	vertices := graphics.MakeRoundedRectangle(x, y, width, height, radius, c, sw, sh)
	// Rectangle has 6 vertices (2 triangles), store all for transforms
	screenPos := [][2]float32{
		{float32(x), float32(y + height)},         // bottom-left
		{float32(x + width), float32(y + height)}, // bottom-right
		{float32(x), float32(y)},                  // top-left
		{float32(x), float32(y)},                  // top-left (repeated)
		{float32(x + width), float32(y + height)}, // bottom-right (repeated)
		{float32(x + width), float32(y)},          // top-right
	}

	return pipelines.NewPrimitiveShape(rq.PrimitiveBuffer, vertices, screenPos)
}

// AddCircle creates a new Circle renderable and adds it to the RenderQueue.
// It returns a reference to the created Circle.
func (rq *RenderQueue) AddCircle(cx, cy int, radius float32, c color.Color, segments int) graphics.Shape {
	// note: we could probably more efficiently draw circles with a custom shader -- but this is a good start
	circle := rq.AddPolygon(cx, cy, radius*2, c, segments)
	return circle
}

func (rq *RenderQueue) AddPolygonFromVertices(cx, cy int, width float32, vertices []graphics.Vertex) graphics.Shape {
	sw, sh := rq.GetSurfaceSize()
	swf, shf := float32(sw), float32(sh)

	// Convert graphics.Vertex to graphics.PrimitiveVertex with OpCodeSolid
	guiVertices := make([]graphics.PrimitiveVertex, len(vertices))
	screenPos := make([][2]float32, len(vertices))

	// Barycentric coordinates for triangle vertices (for edge anti-aliasing)
	// Vertex 0: (1,0) -> bary(1,0,0), Vertex 1: (0,1) -> bary(0,1,0), Vertex 2: (0,0) -> bary(0,0,1)
	baryCoords := [][2]float32{{1, 0}, {0, 1}, {0, 0}}

	for i, v := range vertices {
		screenPos[i] = [2]float32{v.Position[0], v.Position[1]}
		// Cycle through barycentric coords for each triangle (3 vertices per triangle)
		bary := baryCoords[i%3]
		guiVertices[i] = graphics.PrimitiveVertex{
			Position:      screenToNDC(v.Position[0], v.Position[1], swf, shf),
			LocalPosition: bary,
			OpCode:        graphics.OpCodeSolid,
			Radius:        0,
			Color:         v.Color,
			TexCoords:     [2]float32{0, 0},
		}
	}

	return pipelines.NewPrimitiveShape(rq.PrimitiveBuffer, guiVertices, screenPos)
}

// AddPolygon creates a new Polygon renderable and adds it to the RenderQueue.
// It returns a reference to the created Polygon.
func (rq *RenderQueue) AddPolygon(cx, cy int, width float32, c color.Color, sides int) graphics.Shape {
	sw, sh := rq.GetSurfaceSize()
	vertices := graphics.MakeSolidPolygon(cx, cy, width, sides, c, sw, sh)

	// Generate screen positions (sides*3 vertices: each triangle is center, current, next)
	screenPos := make([][2]float32, 0, sides*3)
	angleStep := 2 * math.Pi / float64(sides)
	radius := float64(width / 2)
	fcx, fcy := float64(cx), float64(cy)

	x := fcx + radius*math.Cos(0)
	y := fcy + radius*math.Sin(0)

	for i := 0; i < sides; i++ {
		nextAngle := float64(i+1) * angleStep
		nextX := fcx + radius*math.Cos(nextAngle)
		nextY := fcy + radius*math.Sin(nextAngle)

		// center, current, next
		screenPos = append(screenPos,
			[2]float32{float32(fcx), float32(fcy)},
			[2]float32{float32(x), float32(y)},
			[2]float32{float32(nextX), float32(nextY)},
		)

		x, y = nextX, nextY
	}

	return pipelines.NewPrimitiveShape(rq.PrimitiveBuffer, vertices, screenPos)
}

// AddLine creates a new Line renderable and adds it to the RenderQueue.
// It returns a reference to the created Line.
func (rq *RenderQueue) AddLine(x1, y1, x2, y2 int, width float32, c color.Color) graphics.Shape {
	sw, sh := rq.GetSurfaceSize()
	vertices := graphics.MakeSolidLine(x1, y1, x2, y2, width, c, sw, sh)
	if vertices == nil {
		// Zero-length line
		return pipelines.NewPrimitiveShape(rq.PrimitiveBuffer, nil, nil)
	}

	// Calculate corner positions for transforms
	dx := float32(x2 - x1)
	dy := float32(y2 - y1)
	length := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	sin := dy / length
	cos := dx / length
	halfWidth := width / 2

	// Line has 6 vertices (2 triangles)
	screenPos := [][2]float32{
		{float32(x1) - sin*halfWidth, float32(y1) + cos*halfWidth},
		{float32(x2) - sin*halfWidth, float32(y2) + cos*halfWidth},
		{float32(x2) + sin*halfWidth, float32(y2) - cos*halfWidth},
		{float32(x1) - sin*halfWidth, float32(y1) + cos*halfWidth},
		{float32(x2) + sin*halfWidth, float32(y2) - cos*halfWidth},
		{float32(x1) + sin*halfWidth, float32(y1) - cos*halfWidth},
	}

	return pipelines.NewPrimitiveShape(rq.PrimitiveBuffer, vertices, screenPos)
}

func (rq *RenderQueue) DrawPrimitiveBuffer(vertices []graphics.PrimitiveVertex) {
	if len(vertices) == 0 {
		return
	}
	rq.PrimitiveBuffer.UpdateVertexBuffer(vertices)
}

// DrawPrimitiveBufferWithClipRects draws vertices with per-vertex clip rects.
// Note: WebGPU clip rect support requires render pass integration.
func (rq *RenderQueue) DrawPrimitiveBufferWithClipRects(vertices []graphics.PrimitiveVertex, clipRects []*[4]int) {
	if len(vertices) == 0 {
		return
	}
	// TODO: Implement clip rect support for WebGPU via SetScissorRect
	rq.PrimitiveBuffer.UpdateVertexBuffer(vertices)
}

// DrawPrimitives uploads primitives directly to the storage buffer.
// This is the new efficient approach using 64 bytes per primitive instead of 312 bytes.
func (rq *RenderQueue) DrawPrimitives(primitives []graphics.Primitive) {
	if len(primitives) == 0 {
		return
	}
	rq.PrimitiveBuffer.UpdatePrimitives(primitives)
}

// FlushPrimitiveBuffer is a no-op for WebGPU since immediate rendering
// is not supported in the same way as OpenGL. Clip rect support in WebGPU
// requires render pass integration.
func (rq *RenderQueue) FlushPrimitiveBuffer() {
	// WebGPU uses render passes, so immediate flush isn't directly supported.
	// The clip rect would need to be applied via SetScissorRect on the render pass.
}

func (rq *RenderQueue) SetMSDFAtlas(atlasImg image.Image, pxRange float64) {
	rq.PrimitiveBuffer.SetMSDFAtlas(atlasImg, pxRange)
}

func (rq *RenderQueue) SetMSDFMode(mode int) {
	rq.PrimitiveBuffer.SetMSDFMode(mode)
}

func (rq *RenderQueue) EnableSnapMSDFToPixels(enable bool) {
	rq.PrimitiveBuffer.EnableSnapMSDFToPixels(enable)
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

// screenToNDC converts screen coordinates to normalized device coordinates.
func screenToNDC(x, y, screenWidth, screenHeight float32) [3]float32 {
	ndcX := (x/screenWidth)*2 - 1
	ndcY := 1 - (y/screenHeight)*2
	return [3]float32{ndcX, ndcY, 0}
}

// CreateMSDFAtlas creates a new MSDF atlas from an image
func (rq *RenderQueue) CreateMSDFAtlas(atlasImg image.Image, distanceRange float64) (graphics.MSDFAtlas, error) {
	return pipelines.NewMSDFAtlas(rq.RenderContext, atlasImg, distanceRange)
}

// GetCurrentClipRect returns the current clip rect.
// Note: WebGPU clip rect support requires render pass integration.
func (rq *RenderQueue) GetCurrentClipRect() *[4]int {
	// TODO: Implement proper clip rect tracking for WebGPU
	return nil
}
