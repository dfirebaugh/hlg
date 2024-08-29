package primitives

import (
	"image/color"
	"math"
	"unsafe"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/graphics"
)

type Circle struct {
	x, y, radius     int
	color            color.Color
	outlineColor     color.Color
	outlineThickness float32
	fillQuad         graphics.ShaderRenderable
	drawQuad         graphics.ShaderRenderable
	dirtyFlags       circleDirtyFlags
}

type circleDirtyFlags struct {
	position     bool
	radius       bool
	color        bool
	outlineColor bool
}

func NewCircle(x, y, radius int, color color.Color, outlineThickness float32) *Circle {
	c := &Circle{
		x:                x,
		y:                y,
		radius:           radius,
		color:            color,
		outlineThickness: outlineThickness,
		dirtyFlags: circleDirtyFlags{
			position:     true,
			radius:       true,
			color:        true,
			outlineColor: true,
		},
	}
	c.init()
	return c
}

func (c *Circle) init() {
	ww, wh := hlg.GetWindowSize()
	fullScreenQuad := makeFullScreenQuad(float32(ww), float32(wh))

	c.fillQuad = hlg.CreateRenderable(CircleFillShader, fullScreenQuad, c.getUniforms(false), c.getDataMap(false))
	if c.fillQuad == nil {
		panic("Failed to create fill renderable for circle")
	}

	c.drawQuad = hlg.CreateRenderable(CircleOutlineShader, fullScreenQuad, c.getUniforms(true), c.getDataMap(true))
	if c.drawQuad == nil {
		panic("Failed to create draw renderable for circle")
	}
}

func (c *Circle) SetOutlineThickness(thickness float32) {
	if c.outlineThickness != thickness {
		c.outlineThickness = thickness
		c.dirtyFlags.outlineColor = true
	}
}

func (c *Circle) SetX(x int) {
	if c.x != x {
		c.x = x
		c.dirtyFlags.position = true
	}
}

func (c *Circle) SetY(y int) {
	if c.y != y {
		c.y = y
		c.dirtyFlags.position = true
	}
}

func (c *Circle) SetRadius(radius int) {
	if c.radius != radius {
		c.radius = radius
		c.dirtyFlags.radius = true
	}
}

func (c *Circle) SetColor(color color.Color) {
	if c.color != color {
		c.color = color
		c.dirtyFlags.color = true
	}
}

func (c *Circle) SetOutlineColor(outlineColor color.Color) {
	if c.outlineColor != outlineColor {
		c.outlineColor = outlineColor
		c.dirtyFlags.outlineColor = true
	}
}

func (c *Circle) GetX() int {
	return c.x
}

func (c *Circle) GetY() int {
	return c.y
}

func (c *Circle) GetRadius() int {
	return c.radius
}

func (c *Circle) IsPointWithin(px, py int) bool {
	dx := float64(px - c.x)
	dy := float64(py - c.y)
	distance := math.Sqrt(dx*dx + dy*dy)
	return distance <= float64(c.radius)
}

func (c *Circle) getUniforms(isOutline bool) map[string]hlg.Uniform {
	pos := Position{float32(c.x), float32(c.y)}
	colorSlice := colorToFloatSlice(c.color)
	if isOutline {
		colorSlice = colorToFloatSlice(c.outlineColor)
	}
	radius := float32(c.radius)
	thickness := c.outlineThickness

	return map[string]hlg.Uniform{
		"circle_pos": {
			Binding: 0,
			Size:    uint64(unsafe.Sizeof(pos)),
		},
		"circle_radius": {
			Binding: 1,
			Size:    uint64(unsafe.Sizeof(radius)),
		},
		"circle_color": {
			Binding: 2,
			Size:    uint64(unsafe.Sizeof(colorSlice)),
		},
		"outline_thickness": {
			Binding: 3,
			Size:    uint64(unsafe.Sizeof(thickness)),
		},
	}
}

func (c *Circle) getDataMap(isOutline bool) map[string][]byte {
	pos := Position{float32(c.x), float32(c.y)}
	colorSlice := colorToFloatSlice(c.color)
	if isOutline {
		colorSlice = colorToFloatSlice(c.outlineColor)
	}
	radius := float32(c.radius)
	thickness := c.outlineThickness

	return map[string][]byte{
		"circle_pos":        unsafe.Slice((*byte)(unsafe.Pointer(&pos)), int(unsafe.Sizeof(pos))),
		"circle_radius":     unsafe.Slice((*byte)(unsafe.Pointer(&radius)), int(unsafe.Sizeof(radius))),
		"circle_color":      unsafe.Slice((*byte)(unsafe.Pointer(&colorSlice[0])), int(unsafe.Sizeof(colorSlice))),
		"outline_thickness": unsafe.Slice((*byte)(unsafe.Pointer(&thickness)), int(unsafe.Sizeof(thickness))),
	}
}

func (c *Circle) UpdateUniforms() {
	if c.dirtyFlags.position || c.dirtyFlags.radius || c.dirtyFlags.color || c.dirtyFlags.outlineColor {
		if c.fillQuad != nil && (c.dirtyFlags.position || c.dirtyFlags.radius || c.dirtyFlags.color) {
			c.fillQuad.UpdateUniforms(c.getDataMap(false))
		}

		if c.drawQuad != nil && (c.dirtyFlags.position || c.dirtyFlags.radius || c.dirtyFlags.outlineColor) {
			c.drawQuad.UpdateUniforms(c.getDataMap(true))
		}

		c.dirtyFlags = circleDirtyFlags{}
	}
}

func (c *Circle) Update() {}
func (c *Circle) Render() {
	c.Fill()
	c.Draw()
}

func (c *Circle) Fill() {
	c.UpdateUniforms()
	c.fillQuad.Render()
}

func (c *Circle) Draw() {
	c.UpdateUniforms()
	c.drawQuad.Render()
}
