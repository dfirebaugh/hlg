package webgpu

import (
	"image"
	"image/color"
	"math"
	"unsafe"

	"github.com/dfirebaugh/ggez/graphics"
	"github.com/dfirebaugh/ggez/graphics/webgpu/internal/shapes"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type textureHandle uintptr

type RenderQueue struct {
	*wgpu.Device
	*wgpu.SwapChainDescriptor
	Textures     map[textureHandle]*Texture
	renderQueue  []graphics.Renderable
	nextFrame    []graphics.Renderable
	currentFrame []graphics.Renderable
}

func NewRenderQueue(d *wgpu.Device, scd *wgpu.SwapChainDescriptor) *RenderQueue {
	return &RenderQueue{
		Device:              d,
		SwapChainDescriptor: scd,
		Textures:            make(map[textureHandle]*Texture),
		nextFrame:           []graphics.Renderable{},
		currentFrame:        []graphics.Renderable{},
	}
}

func (rq *RenderQueue) RenderClear() {
	for _, r := range rq.renderQueue {
		r.Hide()
	}
}
func (rq *RenderQueue) AddToRenderQueue(r graphics.Renderable) {
	rq.renderQueue = append(rq.renderQueue, r)
}
func (rq *RenderQueue) Pop() (graphics.Renderable, bool) {
	if len(rq.renderQueue) == 0 {
		return nil, false
	}

	renderable := rq.renderQueue[0]
	rq.renderQueue = rq.renderQueue[1:]

	return renderable, true
}

func (rq *RenderQueue) PrepareFrame() {
	if len(rq.nextFrame) > 64 {
		rq.nextFrame = rq.nextFrame[:64]
	}
	for {
		if renderable, ok := rq.Pop(); ok {
			rq.nextFrame = append(rq.nextFrame, renderable)
			continue
		}
		break
	}
	rq.currentFrame = rq.nextFrame
}

func (rq *RenderQueue) RenderFrame(pass *wgpu.RenderPassEncoder) {
	for _, renderable := range rq.currentFrame {
		if renderable == nil {
			continue
		}
		renderable.RenderPass(pass)
	}
}

func (rq *RenderQueue) CreateTextureFromImage(img image.Image) (graphics.Texture, error) {
	tex := NewTexture(rq.Device, rq.SwapChainDescriptor, img, rq)
	handle := textureHandle(uintptr(unsafe.Pointer(tex)))
	tex.SetHandle(handle)
	rq.Textures[handle] = tex
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
	r, g, b, a := c.RGBA()
	triangle := shapes.NewPolygon(rq.Device, rq.SwapChainDescriptor, rq, []shapes.Vertex{
		{
			Position: [3]float32{float32(x1), float32(y1), 0},
			Color:    [4]float32{float32(r) / 0xffff, float32(g) / 0xffff, float32(b) / 0xffff, float32(a) / 0xffff},
		},
		{
			Position: [3]float32{float32(x2), float32(y2), 0},
			Color:    [4]float32{float32(r) / 0xffff, float32(g) / 0xffff, float32(b) / 0xffff, float32(a) / 0xffff},
		},
		{
			Position: [3]float32{float32(x3), float32(y3), 0},
			Color:    [4]float32{float32(r) / 0xffff, float32(g) / 0xffff, float32(b) / 0xffff, float32(a) / 0xffff},
		},
	}, shapes.RenderFilled)

	rq.AddToRenderQueue(triangle)

	return triangle
}

// AddRectangle creates a new Rectangle renderable and adds it to the RenderQueue.
// It returns a reference to the created Rectangle.
func (rq *RenderQueue) AddRectangle(x, y, width, height int, c color.Color) graphics.Shape {
	r, g, b, a := c.RGBA()
	colorArray := [4]float32{float32(r) / 0xffff, float32(g) / 0xffff, float32(b) / 0xffff, float32(a) / 0xffff}

	topLeft := shapes.Vertex{
		Position: [3]float32{float32(x), float32(y), 0},
		Color:    colorArray,
	}
	topRight := shapes.Vertex{
		Position: [3]float32{float32(x + width), float32(y), 0},
		Color:    colorArray,
	}
	bottomLeft := shapes.Vertex{
		Position: [3]float32{float32(x), float32(y + height), 0},
		Color:    colorArray,
	}
	bottomRight := shapes.Vertex{
		Position: [3]float32{float32(x + width), float32(y + height), 0},
		Color:    colorArray,
	}

	rectangleVertices := []shapes.Vertex{topLeft, bottomLeft, topRight, bottomLeft, bottomRight, topRight}

	rectangle := shapes.NewPolygon(rq.Device, rq.SwapChainDescriptor, rq, rectangleVertices, shapes.RenderFilled)

	rq.AddToRenderQueue(rectangle)
	return rectangle
}

// AddCircle creates a new Circle renderable and adds it to the RenderQueue.
// It returns a reference to the created Circle.
// note: we could probably more efficiently draw circles with a custom shader -- but this is a good start
func (rq *RenderQueue) AddCircle(cx, cy int, radius float32, c color.Color, segments int) graphics.Shape {
	r, g, b, a := c.RGBA()
	colorArray := [4]float32{float32(r) / 0xffff, float32(g) / 0xffff, float32(b) / 0xffff, float32(a) / 0xffff}
	var vertices []shapes.Vertex
	center := shapes.Vertex{
		Position: [3]float32{float32(cx), float32(cy), 0},
		Color:    colorArray,
	}

	for i := 0; i <= segments; i++ {
		angle := float32(i) * 2 * float32(math.Pi) / float32(segments)
		x := float32(cx) + radius*float32(math.Cos(float64(angle)))
		y := float32(cy) + radius*float32(math.Sin(float64(angle)))

		vertex := shapes.Vertex{
			Position: [3]float32{x, y, 0},
			Color:    colorArray,
		}

		vertices = append(vertices, center, vertex)

		if i < segments {
			nextAngle := float32(i+1) * 2 * float32(math.Pi) / float32(segments)
			nextX := float32(cx) + radius*float32(math.Cos(float64(nextAngle)))
			nextY := float32(cy) + radius*float32(math.Sin(float64(nextAngle)))

			nextVertex := shapes.Vertex{
				Position: [3]float32{nextX, nextY, 0},
				Color:    colorArray,
			}

			vertices = append(vertices, nextVertex)
		}
	}

	circle := shapes.NewPolygon(rq.Device, rq.SwapChainDescriptor, rq, vertices, shapes.RenderFilled)
	rq.AddToRenderQueue(circle)
	return circle
}

// AddLine creates a new Line renderable and adds it to the RenderQueue.
// It returns a reference to the created Line.
func (rq *RenderQueue) AddLine(x1, y1, x2, y2 int, width float32, c color.Color) graphics.Shape {
	r, g, b, a := c.RGBA()
	colorArray := [4]float32{float32(r) / 0xffff, float32(g) / 0xffff, float32(b) / 0xffff, float32(a) / 0xffff}

	dx := float32(x2 - x1)
	dy := float32(y2 - y1)
	len := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	sin := dy / len
	cos := dx / len

	// Calculate the four corners of the line (as a very thin rectangle)
	halfWidth := width / 2
	vertices := []shapes.Vertex{
		{Position: [3]float32{float32(x1) - sin*halfWidth, float32(y1) + cos*halfWidth, 0}, Color: colorArray},
		{Position: [3]float32{float32(x2) - sin*halfWidth, float32(y2) + cos*halfWidth, 0}, Color: colorArray},
		{Position: [3]float32{float32(x2) + sin*halfWidth, float32(y2) - cos*halfWidth, 0}, Color: colorArray},
		{Position: [3]float32{float32(x1) + sin*halfWidth, float32(y1) - cos*halfWidth, 0}, Color: colorArray},
	}

	// Creating two triangles to form the line
	lineVertices := []shapes.Vertex{vertices[0], vertices[1], vertices[2], vertices[0], vertices[2], vertices[3]}

	line := shapes.NewPolygon(rq.Device, rq.SwapChainDescriptor, rq, lineVertices, shapes.RenderFilled)
	rq.AddToRenderQueue(line)
	return line
}
