package hlg

import (
	"image"
	"image/color"

	"github.com/dfirebaugh/hlg/graphics"
)

// defaultFont is an optional default font that can be set for PrintAt
var defaultFont *Font

// framePrimitives holds batched primitives for the current frame (new efficient format)
var framePrimitives []graphics.Primitive

// frameVertices holds raw vertices for shapes that can't use the Primitive format (e.g., triangles)
var frameVertices []graphics.PrimitiveVertex

// frameClipRects holds clip rects for frameVertices (one per vertex, parallel to frameVertices)
var frameClipRects []*[4]int

// frameScreenWidth and frameScreenHeight cache the screen dimensions for the current frame
var frameScreenWidth, frameScreenHeight int

// clipRectStack tracks the current clip rect stack at the hlg level
// Each primitive captures the current clip rect when created
var clipRectStack [][4]int

// getCurrentClipRect returns the current clip rect, or nil if none is active
func getCurrentClipRect() *[4]int {
	if len(clipRectStack) == 0 {
		return nil
	}
	// Return a copy to avoid aliasing
	rect := clipRectStack[len(clipRectStack)-1]
	return &rect
}

// pushClipRectToStack adds a clip rect to the local stack (called from hlg.go)
func pushClipRectToStack(x, y, width, height int) {
	clipRectStack = append(clipRectStack, [4]int{x, y, width, height})
}

// popClipRectFromStack removes a clip rect from the local stack (called from hlg.go)
func popClipRectFromStack() {
	if len(clipRectStack) > 0 {
		clipRectStack = clipRectStack[:len(clipRectStack)-1]
	}
}

// BeginDraw starts a new frame for batched primitive drawing.
// Call this at the start of your render function, then use Text(), RoundedRect(), etc.
// Call EndDraw() to submit all batched primitives.
func BeginDraw() {
	ensureSetupCompletion()
	frameScreenWidth, frameScreenHeight = GetScreenSize()
	// Reset slices while preserving capacity
	if framePrimitives != nil {
		framePrimitives = framePrimitives[:0]
	} else {
		framePrimitives = make([]graphics.Primitive, 0, 1024)
	}
	if frameVertices != nil {
		frameVertices = frameVertices[:0]
	} else {
		frameVertices = make([]graphics.PrimitiveVertex, 0, 256)
	}
	if frameClipRects != nil {
		frameClipRects = frameClipRects[:0]
	} else {
		frameClipRects = make([]*[4]int, 0, 256)
	}
}

// EndDraw submits all batched primitives and vertices from the current frame.
func EndDraw() {
	flushBatch()
}

// flushBatch submits the current batch of primitives without ending the draw session.
// This is called automatically when a Shape is rendered to preserve draw order.
func flushBatch() {
	if len(framePrimitives) > 0 || len(frameVertices) > 0 {
		sw, sh := float32(frameScreenWidth), float32(frameScreenHeight)
		vertices := graphics.ConvertPrimitivesToVertices(framePrimitives, sw, sh)
		// Extract clip rects from primitives and merge with raw vertex clip rects
		clipRects := graphics.ExtractClipRectsFromPrimitives(framePrimitives)
		clipRects = append(clipRects, frameClipRects...)

		vertices = append(vertices, frameVertices...)
		SubmitDrawBufferWithClipRects(vertices, clipRects)

		// Render immediately to preserve draw order with Shapes
		hlg.graphicsBackend.FlushPrimitiveBuffer()

		// Reset slices while preserving capacity
		framePrimitives = framePrimitives[:0]
		frameVertices = frameVertices[:0]
		frameClipRects = frameClipRects[:0]
	}
}

// Text draws text at the specified position with the given font size and color.
// Must be called between BeginDraw() and EndDraw().
func Text(s string, x, y int, fontSize float32, c color.Color) {
	if defaultFont == nil {
		font, err := LoadDefaultFont()
		if err != nil {
			return
		}
		SetDefaultFont(font)
		font.SetAsActiveAtlas()
	}
	primitives := defaultFont.RenderTextPrimitives(s, float32(x), float32(y), fontSize, c)
	if len(primitives) == 0 {
		return
	}
	// Attach current clip rect to all primitives
	clipRect := getCurrentClipRect()
	for i := range primitives {
		primitives[i].ClipRect = clipRect
	}
	if framePrimitives == nil {
		// Fallback to immediate mode if not in batched mode
		SubmitPrimitives(primitives)
		return
	}
	framePrimitives = append(framePrimitives, primitives...)
}

// RoundedRect draws a filled rounded rectangle.
// Must be called between BeginDraw() and EndDraw().
func RoundedRect(x, y, width, height, cornerRadius int, c color.Color) {
	primitive := graphics.MakeRoundedRectPrimitive(x, y, width, height, cornerRadius, c)
	primitive.ClipRect = getCurrentClipRect()
	if framePrimitives == nil {
		SubmitPrimitives([]graphics.Primitive{primitive})
		return
	}
	framePrimitives = append(framePrimitives, primitive)
}

// RoundedRectOutline draws a rounded rectangle with an outline.
// Must be called between BeginDraw() and EndDraw().
func RoundedRectOutline(x, y, width, height, cornerRadius, outlineWidth int, fillColor, outlineColor color.Color) {
	clipRect := getCurrentClipRect()

	// Draw outer rectangle (outline)
	outerX := x - outlineWidth
	outerY := y - outlineWidth
	outerWidth := width + 2*outlineWidth
	outerHeight := height + 2*outlineWidth
	outerRadius := cornerRadius + outlineWidth

	outerPrimitive := graphics.MakeRoundedRectPrimitive(outerX, outerY, outerWidth, outerHeight, outerRadius, outlineColor)
	outerPrimitive.ClipRect = clipRect

	// Draw inner rectangle (fill)
	innerX := x + outlineWidth
	innerY := y + outlineWidth
	innerWidth := width - 2*outlineWidth
	innerHeight := height - 2*outlineWidth
	innerRadius := cornerRadius - outlineWidth
	if innerRadius < 0 {
		innerRadius = 0
	}

	innerPrimitive := graphics.MakeRoundedRectPrimitive(innerX, innerY, innerWidth, innerHeight, innerRadius, fillColor)
	innerPrimitive.ClipRect = clipRect

	if framePrimitives == nil {
		SubmitPrimitives([]graphics.Primitive{outerPrimitive, innerPrimitive})
		return
	}
	framePrimitives = append(framePrimitives, outerPrimitive, innerPrimitive)
}

// FilledCircle draws a filled circle.
// Must be called between BeginDraw() and EndDraw().
func FilledCircle(x, y, radius int, c color.Color) {
	primitive := graphics.MakeCirclePrimitive(x, y, radius, c)
	primitive.ClipRect = getCurrentClipRect()
	if framePrimitives == nil {
		SubmitPrimitives([]graphics.Primitive{primitive})
		return
	}
	framePrimitives = append(framePrimitives, primitive)
}

// FilledRect draws a filled rectangle (no rounded corners).
// Must be called between BeginDraw() and EndDraw().
func FilledRect(x, y, width, height int, c color.Color) {
	primitive := graphics.MakeRectPrimitive(x, y, width, height, c)
	primitive.ClipRect = getCurrentClipRect()
	if framePrimitives == nil {
		SubmitPrimitives([]graphics.Primitive{primitive})
		return
	}
	framePrimitives = append(framePrimitives, primitive)
}

// FilledTriangle draws a filled triangle.
// Must be called between BeginDraw() and EndDraw().
func FilledTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color) {
	// Triangles use the vertex-based approach since the Primitive format
	// doesn't support arbitrary triangle vertices
	sw, sh := frameScreenWidth, frameScreenHeight
	if sw == 0 || sh == 0 {
		sw, sh = GetScreenSize()
	}
	vertices := graphics.MakeSolidTriangle(x1, y1, x2, y2, x3, y3, c, sw, sh)
	// Track clip rect for these vertices
	clipRect := getCurrentClipRect()
	for range vertices {
		frameClipRects = append(frameClipRects, clipRect)
	}
	if frameVertices == nil {
		// Fallback if not in BeginDraw/EndDraw block
		ensureSetupCompletion()
		SubmitDrawBuffer(vertices)
		return
	}
	frameVertices = append(frameVertices, vertices...)
}

// Segment draws a line segment between two points.
// Must be called between BeginDraw() and EndDraw().
func Segment(x1, y1, x2, y2, thickness int, c color.Color) {
	primitive := graphics.MakeLinePrimitive(x1, y1, x2, y2, float32(thickness), c)
	// Check for zero-length line (MakeLinePrimitive returns empty Primitive)
	if primitive.W == 0 && primitive.H == 0 {
		return
	}
	primitive.ClipRect = getCurrentClipRect()
	if framePrimitives == nil {
		SubmitPrimitives([]graphics.Primitive{primitive})
		return
	}
	framePrimitives = append(framePrimitives, primitive)
}

type Shape interface {
	graphics.Shape
}

// SetDefaultFont sets the default font used by PrintAt.
// The font should be loaded using LoadFont or LoadFontFromBytes.
func SetDefaultFont(font *Font) {
	defaultFont = font
}

// GetDefaultFont returns the currently set default font, or nil if not set.
func GetDefaultFont() *Font {
	return defaultFont
}

// MeasureText returns the width of the text in pixels at the given font size.
// If no default font is set, the embedded default font will be loaded automatically.
func MeasureText(text string, fontSize float32) float32 {
	if defaultFont == nil {
		font, err := LoadDefaultFont()
		if err != nil {
			return 0
		}
		SetDefaultFont(font)
	}
	return defaultFont.MeasureText(text, fontSize)
}

// PolygonFromVertices creates a polygon shape using a specified array of vertices.
// The vertices should be defined with their positions and colors. The function converts
// the input vertices from the local Vertex type to the graphics.Vertex type required by
// the graphics backend. The created polygon is then added to the render queue.
//
// Parameters:
//   - x, y: The x and y coordinates of the polygon's position.
//   - width: The width of the polygon (not used directly in this function, included for interface consistency).
//   - vertices: A slice of Vertex that defines the positions and colors of the polygon's vertices.
//
// Returns:
//   - A graphics.Shape that represents the created polygon.
func PolygonFromVertices(x, y int, width float32, vertices []Vertex) Shape {
	ensureSetupCompletion()

	graphicsVertices := make([]graphics.Vertex, len(vertices))
	for i, v := range vertices {
		graphicsVertices[i] = graphics.Vertex{
			Position: v.Position,
			Color:    toRGBA(v.Color),
		}
	}

	return hlg.graphicsBackend.AddPolygonFromVertices(x, y, width, graphicsVertices)
}

// Polygon creates a polygon shape with a specified number of sides, position, width, and color.
// x, y define the center of the polygon.
// width defines the diameter of the circumcircle of the polygon.
// sides specify the number of sides (vertices) of the polygon.
// c specifies the color of the polygon.
func Polygon(x, y int, width float32, sides int, c color.Color) Shape {
	ensureSetupCompletion()
	return hlg.graphicsBackend.AddPolygon(x, y, width, c, sides)
}

// Triangle creates a triangle shape with specified vertices and color.
// x1, y1, x2, y2, x3, y3 define the coordinates of the three vertices of the triangle.
// c specifies the color of the triangle.
func Triangle(x1, y1, x2, y2, x3, y3 int, c color.Color) Shape {
	ensureSetupCompletion()
	return hlg.graphicsBackend.AddTriangle(x1, y1, x2, y2, x3, y3, c)
}

// Rectangle creates a rectangle shape with specified position, dimensions, and color.
// x, y define the top-left corner of the rectangle.
// width, height define the dimensions of the rectangle.
// c specifies the color of the rectangle.
func Rectangle(x, y, width, height int, c color.Color) Shape {
	ensureSetupCompletion()
	return hlg.graphicsBackend.AddRectangle(x, y, width, height, c)
}

// RoundedRectangleShape creates a rounded rectangle shape with SDF rendering.
// x, y define the top-left corner of the rectangle.
// width, height define the dimensions of the rectangle.
// radius defines the corner radius.
// c specifies the color of the rectangle.
// This is a first-class shape that can be transformed (move, rotate, scale).
func RoundedRectangleShape(x, y, width, height, radius int, c color.Color) Shape {
	ensureSetupCompletion()
	return hlg.graphicsBackend.AddRoundedRectangle(x, y, width, height, radius, c)
}

// Circle creates a circle shape with specified center, radius, and color.
// x, y define the center of the circle.
// radius defines the radius of the circle.
// c specifies the color of the circle.
func Circle(x, y int, radius float32, c color.Color) Shape {
	ensureSetupCompletion()
	return hlg.graphicsBackend.AddCircle(x, y, radius, c, 32)
}

// Line creates a line with specified start and end points, width, and color.
// x1, y1 define the start point of the line.
// x2, y2 define the end point of the line.
// width defines the thickness of the line.
// c specifies the color of the line.
func Line(x1, y1, x2, y2 int, width float32, c color.Color) Shape {
	ensureSetupCompletion()
	return hlg.graphicsBackend.AddLine(x1, y1, x2, y2, width, c)
}

// PrintAt renders text at a specified position with a specified color.
// s is the string to be rendered.
// x, y define the position where the text will be rendered.
// c specifies the color of the text.
//
// If no default font is set, the embedded default font will be loaded automatically.
func PrintAt(s string, x int, y int, c color.Color) {
	PrintAtWithSize(s, x, y, 16, c)
}

// PrintAtWithSize renders text at a specified position with a specified color and font size.
// s is the string to be rendered.
// x, y define the position where the text will be rendered.
// fontSize defines the size of the text in pixels.
// c specifies the color of the text.
//
// If no default font is set, the embedded default font will be loaded automatically.
func PrintAtWithSize(s string, x int, y int, fontSize float32, c color.Color) {
	if defaultFont == nil {
		font, err := LoadDefaultFont()
		if err != nil {
			return
		}
		SetDefaultFont(font)
		font.SetAsActiveAtlas()
	}
	primitives := defaultFont.RenderTextPrimitives(s, float32(x), float32(y), fontSize, c)
	if len(primitives) > 0 {
		SubmitPrimitives(primitives)
	}
}

// SubmitDrawBuffer submits vertices using the legacy PrimitiveVertex format.
// DEPRECATED: Use SubmitPrimitives for ~5x memory reduction.
func SubmitDrawBuffer(vertices []graphics.PrimitiveVertex) {
	hlg.graphicsBackend.DrawPrimitiveBuffer(vertices)
}

// SubmitDrawBufferWithClipRects submits vertices with per-vertex clip rects.
func SubmitDrawBufferWithClipRects(vertices []graphics.PrimitiveVertex, clipRects []*[4]int) {
	hlg.graphicsBackend.DrawPrimitiveBufferWithClipRects(vertices, clipRects)
}

// SubmitPrimitives submits primitives using the new efficient storage buffer format.
func SubmitPrimitives(primitives []graphics.Primitive) {
	hlg.graphicsBackend.DrawPrimitives(primitives)
}

// SetMSDFAtlas sets the MSDF font atlas for the primitive buffer.
// This must be called before using DrawGlyph in the GUI draw context.
func SetMSDFAtlas(atlasImg image.Image, pxRange float64) {
	ensureSetupCompletion()
	hlg.graphicsBackend.SetMSDFAtlas(atlasImg, pxRange)
}

// SetMSDFMode sets the MSDF rendering mode.
// Mode 0: median(RGB) - MSDF reconstruction for sharp corners (default)
// Mode 1: alpha channel only (true SDF fallback)
// Mode 2: visualize RGB channels directly (for debugging atlas)
func SetMSDFMode(mode int) {
	ensureSetupCompletion()
	hlg.graphicsBackend.SetMSDFMode(mode)
}

// EnableSnapMSDFToPixels enables integer pixel snapping for MSDF primitives.
// This is intended for debugging subpixel placement issues that can make text look fuzzy.
func EnableSnapMSDFToPixels(enable bool) {
	ensureSetupCompletion()
	hlg.graphicsBackend.EnableSnapMSDFToPixels(enable)
}
