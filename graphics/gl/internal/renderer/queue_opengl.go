//go:build !js

package renderer

import (
	"image"
	"image/color"
	"math"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/glapi"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/pipelines"
	"github.com/dfirebaugh/hlg/graphics/gl/internal/shader"
)

type textureHandle int

// RenderQueue manages rendering operations for OpenGL
type RenderQueue struct {
	ctx           *glapi.Context
	surface       *Surface
	ShaderManager *shader.ShaderManager

	primitiveBuffer *pipelines.PrimitiveBuffer
	Textures        map[textureHandle]*Texture

	renderQueue        []graphics.Renderable
	clipRectStack      [][4]int
	nextTextureHandle  textureHandle
	isDisposed         bool
	presentedThisFrame bool // tracks if Present() was called this frame

	onBeforeAddToQueue func()
}

// NewRenderQueue creates a new render queue
func NewRenderQueue(ctx *glapi.Context, surface *Surface) *RenderQueue {
	rq := &RenderQueue{
		ctx:               ctx,
		surface:           surface,
		Textures:          make(map[textureHandle]*Texture),
		renderQueue:       make([]graphics.Renderable, 0),
		clipRectStack:     make([][4]int, 0),
		nextTextureHandle: 1,
	}

	rq.ShaderManager = shader.NewShaderManager(ctx)
	shader.CompileShaders(rq.ShaderManager)

	// Pass rq (which implements pipelines.Surface and pipelines.RenderQueue) so that
	// PrimitiveShape.Render() can access GetCurrentClipRect() and AddToRenderQueue()
	rq.primitiveBuffer = pipelines.NewPrimitiveBuffer(ctx, rq, rq.ShaderManager)

	return rq
}

// GetGL returns the GL context
func (rq *RenderQueue) GetGL() *glapi.Context {
	return rq.ctx
}

// GetSurfaceSize returns the surface size
func (rq *RenderQueue) GetSurfaceSize() (int, int) {
	return rq.surface.GetSurfaceSize()
}

// AddToRenderQueue adds a renderable to the queue
func (rq *RenderQueue) AddToRenderQueue(r graphics.Renderable) {
	// Flush pending batched primitives before adding a new renderable
	// This preserves draw order between batched primitives and Shapes
	if rq.onBeforeAddToQueue != nil {
		rq.onBeforeAddToQueue()
		// After flushing, render the shape immediately instead of queueing
		// This preserves draw order when using BeginDraw/EndDraw
		if gr, ok := r.(glRenderable); ok {
			gr.GLRender()
			return
		}
	}
	rq.renderQueue = append(rq.renderQueue, r)
}

// SetOnBeforeAddToQueue sets a callback to be called before adding to queue
func (rq *RenderQueue) SetOnBeforeAddToQueue(fn func()) {
	rq.onBeforeAddToQueue = fn
}

// PrepareFrame clears the render queue for a new frame
func (rq *RenderQueue) PrepareFrame() {
	rq.renderQueue = rq.renderQueue[:0]
	rq.presentedThisFrame = false
}

// RenderFrame renders all queued items
func (rq *RenderQueue) RenderFrame() {
	// Skip if already presented this frame (via Present())
	if rq.presentedThisFrame {
		return
	}

	// First render primitive buffer
	rq.primitiveBuffer.GLRender()

	// Then render other items
	for _, r := range rq.renderQueue {
		if glr, ok := r.(glRenderable); ok {
			glr.GLRender()
		}
	}
}

// glRenderable is the interface for items that can render via OpenGL
type glRenderable interface {
	GLRender()
}

// CreateTexture creates a new texture from an image
func (rq *RenderQueue) CreateTexture(img image.Image) graphics.Texture {
	tex := NewTexture(rq, img)
	handle := rq.nextTextureHandle
	rq.nextTextureHandle++
	tex.SetHandle(handle)
	rq.Textures[handle] = tex
	return tex
}

// CreateMSDFAtlas creates a new MSDF font atlas
func (rq *RenderQueue) CreateMSDFAtlas() graphics.MSDFAtlas {
	atlas, _ := pipelines.NewMSDFAtlas(nil, 4.0)
	return atlas
}

// GetPrimitiveBuffer returns the primitive buffer
func (rq *RenderQueue) GetPrimitiveBuffer() *pipelines.PrimitiveBuffer {
	return rq.primitiveBuffer
}

// PushClipRect pushes a clip rectangle
func (rq *RenderQueue) PushClipRect(x, y, width, height int) {
	rq.clipRectStack = append(rq.clipRectStack, [4]int{x, y, width, height})
}

// PopClipRect pops a clip rectangle
func (rq *RenderQueue) PopClipRect() {
	if len(rq.clipRectStack) > 0 {
		rq.clipRectStack = rq.clipRectStack[:len(rq.clipRectStack)-1]
	}
}

// GetCurrentClipRect returns the current clip rectangle
func (rq *RenderQueue) GetCurrentClipRect() *[4]int {
	if len(rq.clipRectStack) == 0 {
		return nil
	}
	rect := rq.clipRectStack[len(rq.clipRectStack)-1]
	return &rect
}

// CompileShader compiles a shader
func (rq *RenderQueue) CompileShader(code string) graphics.ShaderHandle {
	return rq.ShaderManager.CompileShader(code)
}

// Dispose cleans up resources
func (rq *RenderQueue) Dispose() {
	if rq.isDisposed {
		return
	}
	rq.isDisposed = true

	if rq.primitiveBuffer != nil {
		rq.primitiveBuffer.Dispose()
	}

	for _, tex := range rq.Textures {
		tex.Dispose()
	}
	rq.Textures = nil

	if rq.ShaderManager != nil {
		rq.ShaderManager.ReleaseShaders()
	}
}

// SetPriority sets the render priority (lower values render first)
func (rq *RenderQueue) SetPriority(priority int) {
	// Priority is not used in GL backend currently
}

// Present renders this queue's contents immediately
func (rq *RenderQueue) Present() {
	// First render primitive buffer
	rq.primitiveBuffer.GLRender()

	// Then render other items
	for _, r := range rq.renderQueue {
		if glr, ok := r.(glRenderable); ok {
			glr.GLRender()
		}
	}

	// Mark as presented so RenderFrame() skips this queue in the main render pass
	rq.presentedThisFrame = true
}

// AddTriangle creates a new triangle shape
func (rq *RenderQueue) AddTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) graphics.Shape {
	sw, sh := rq.GetSurfaceSize()
	vertices := graphics.MakeSolidTriangle(x1, y1, x2, y2, x3, y3, c, sw, sh)
	screenPos := [][2]float32{
		{float32(x1), float32(y1)},
		{float32(x2), float32(y2)},
		{float32(x3), float32(y3)},
	}
	return pipelines.NewPrimitiveShape(rq, vertices, screenPos)
}

// AddRectangle creates a new rectangle shape
func (rq *RenderQueue) AddRectangle(x, y, width, height int, c color.Color) graphics.Shape {
	sw, sh := rq.GetSurfaceSize()
	vertices := graphics.MakeSolidRectangle(x, y, width, height, c, sw, sh)
	screenPos := [][2]float32{
		{float32(x), float32(y)},
		{float32(x), float32(y + height)},
		{float32(x + width), float32(y)},
		{float32(x), float32(y + height)},
		{float32(x + width), float32(y + height)},
		{float32(x + width), float32(y)},
	}
	return pipelines.NewPrimitiveShape(rq, vertices, screenPos)
}

// AddRoundedRectangle creates a new rounded rectangle shape
func (rq *RenderQueue) AddRoundedRectangle(x, y, width, height, radius int, c color.Color) graphics.Shape {
	sw, sh := rq.GetSurfaceSize()
	vertices := graphics.MakeRoundedRectangle(x, y, width, height, radius, c, sw, sh)
	screenPos := [][2]float32{
		{float32(x), float32(y + height)},
		{float32(x + width), float32(y + height)},
		{float32(x), float32(y)},
		{float32(x), float32(y)},
		{float32(x + width), float32(y + height)},
		{float32(x + width), float32(y)},
	}
	return pipelines.NewPrimitiveShape(rq, vertices, screenPos)
}

// AddCircle creates a new circle shape
func (rq *RenderQueue) AddCircle(cx, cy int, radius float32, c color.Color, segments int) graphics.Shape {
	circle := rq.AddPolygon(cx, cy, radius*2, c, segments)
	return circle
}

// AddPolygon creates a new polygon shape
func (rq *RenderQueue) AddPolygon(cx, cy int, width float32, c color.Color, sides int) graphics.Shape {
	sw, sh := rq.GetSurfaceSize()
	vertices := graphics.MakeSolidPolygon(cx, cy, width, sides, c, sw, sh)

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

		screenPos = append(screenPos,
			[2]float32{float32(fcx), float32(fcy)},
			[2]float32{float32(x), float32(y)},
			[2]float32{float32(nextX), float32(nextY)},
		)

		x, y = nextX, nextY
	}

	return pipelines.NewPrimitiveShape(rq, vertices, screenPos)
}

// AddPolygonFromVertices creates a new polygon from vertices
func (rq *RenderQueue) AddPolygonFromVertices(cx, cy int, width float32, vertices []graphics.Vertex) graphics.Shape {
	sw, sh := rq.GetSurfaceSize()
	swf, shf := float32(sw), float32(sh)

	guiVertices := make([]graphics.PrimitiveVertex, len(vertices))
	screenPos := make([][2]float32, len(vertices))

	baryCoords := [][2]float32{{1, 0}, {0, 1}, {0, 0}}

	for i, v := range vertices {
		screenPos[i] = [2]float32{v.Position[0], v.Position[1]}
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

	return pipelines.NewPrimitiveShape(rq, guiVertices, screenPos)
}

// AddLine creates a new line shape
func (rq *RenderQueue) AddLine(x1, y1, x2, y2 int, width float32, c color.Color) graphics.Shape {
	sw, sh := rq.GetSurfaceSize()
	vertices := graphics.MakeSolidLine(x1, y1, x2, y2, width, c, sw, sh)
	if vertices == nil {
		return pipelines.NewPrimitiveShape(rq, nil, nil)
	}

	dx := float32(x2 - x1)
	dy := float32(y2 - y1)
	length := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	sin := dy / length
	cos := dx / length
	halfWidth := width / 2

	screenPos := [][2]float32{
		{float32(x1) - sin*halfWidth, float32(y1) + cos*halfWidth},
		{float32(x2) - sin*halfWidth, float32(y2) + cos*halfWidth},
		{float32(x2) + sin*halfWidth, float32(y2) - cos*halfWidth},
		{float32(x1) - sin*halfWidth, float32(y1) + cos*halfWidth},
		{float32(x2) + sin*halfWidth, float32(y2) - cos*halfWidth},
		{float32(x1) + sin*halfWidth, float32(y1) - cos*halfWidth},
	}

	return pipelines.NewPrimitiveShape(rq, vertices, screenPos)
}

// AddDynamicRenderable creates a new dynamic renderable
func (rq *RenderQueue) AddDynamicRenderable(vertexData []byte, layout graphics.VertexBufferLayout, shaderHandle int, uniforms map[string]graphics.Uniform, dataMap map[string][]byte) graphics.ShaderRenderable {
	program := rq.ShaderManager.GetProgram(graphics.ShaderHandle(shaderHandle))
	return pipelines.NewShaderRenderable(rq.ctx, rq, vertexData, layout, program, uniforms, dataMap)
}

// DrawPrimitiveBuffer draws vertices to the primitive buffer
func (rq *RenderQueue) DrawPrimitiveBuffer(vertices []graphics.PrimitiveVertex) {
	if len(vertices) == 0 {
		return
	}
	rq.primitiveBuffer.UpdateVertexBuffer(vertices)
}

// DrawPrimitiveBufferWithClipRects draws vertices with per-vertex clip rects
func (rq *RenderQueue) DrawPrimitiveBufferWithClipRects(vertices []graphics.PrimitiveVertex, clipRects []*[4]int) {
	if len(vertices) == 0 {
		return
	}
	rq.primitiveBuffer.UpdateVertexBufferWithClipRects(vertices, clipRects)
}

// DrawPrimitives draws primitives to the buffer
func (rq *RenderQueue) DrawPrimitives(primitives []graphics.Primitive) {
	if len(primitives) == 0 {
		return
	}
	rq.primitiveBuffer.UpdatePrimitives(primitives)
}

// FlushPrimitiveBuffer forces immediate render of pending primitives
func (rq *RenderQueue) FlushPrimitiveBuffer() {
	rq.primitiveBuffer.FlushImmediate()
}

// SetMSDFAtlas sets the MSDF atlas texture
func (rq *RenderQueue) SetMSDFAtlas(atlasImg image.Image, pxRange float64) {
	rq.primitiveBuffer.SetMSDFAtlas(atlasImg, pxRange)
}

// SetMSDFMode sets the MSDF rendering mode
func (rq *RenderQueue) SetMSDFMode(mode int) {
	rq.primitiveBuffer.SetMSDFMode(mode)
}

// EnableSnapMSDFToPixels enables pixel snapping for MSDF
func (rq *RenderQueue) EnableSnapMSDFToPixels(enable bool) {
	rq.primitiveBuffer.EnableSnapMSDFToPixels(enable)
}

// CreateTextureFromImage creates a texture from an image
func (rq *RenderQueue) CreateTextureFromImage(img image.Image) (graphics.Texture, error) {
	return rq.CreateTexture(img), nil
}

// DisposeTexture disposes a texture by handle
func (rq *RenderQueue) DisposeTexture(h uintptr) {
	if tex, ok := rq.Textures[textureHandle(h)]; ok {
		tex.Dispose()
		delete(rq.Textures, textureHandle(h))
	}
}

// screenToNDC converts screen coordinates to NDC
func screenToNDC(x, y, screenWidth, screenHeight float32) [3]float32 {
	ndcX := (x/screenWidth)*2 - 1
	ndcY := 1 - (y/screenHeight)*2
	return [3]float32{ndcX, ndcY, 0}
}
