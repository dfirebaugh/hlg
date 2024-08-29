package primitives

import (
	"image/color"
	"unsafe"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/graphics"
)

type Rect struct {
	x, y, width, height int
	cornerRadius        float32
	color               color.Color
	outlineColor        color.Color
	outlineWidth        float32
	fillQuad            graphics.ShaderRenderable
	drawQuad            graphics.ShaderRenderable
	dirtyFlags          dirtyFlags
}

type dirtyFlags struct {
	position     bool
	size         bool
	color        bool
	outlineColor bool
	cornerRadius bool
	outlineWidth bool
}

func NewRect(x, y, width, height int, cornerRadius float32, color color.Color) *Rect {
	r := &Rect{
		x:            x,
		y:            y,
		width:        width,
		height:       height,
		cornerRadius: cornerRadius,
		color:        color,
		outlineWidth: 1.0,
		dirtyFlags: dirtyFlags{
			position:     true,
			size:         true,
			color:        true,
			outlineColor: true,
			cornerRadius: true,
			outlineWidth: true,
		},
	}
	r.init()
	return r
}

func (r *Rect) init() {
	ww, wh := hlg.GetWindowSize()
	fullScreenQuad := makeFullScreenQuad(float32(ww), float32(wh))

	r.fillQuad = hlg.CreateRenderable(RectFillShader, fullScreenQuad, r.getUniforms(false), r.getDataMap(false))
	if r.fillQuad == nil {
		panic("Failed to create fill renderable for rect")
	}

	r.drawQuad = hlg.CreateRenderable(RectOutlineShader, fullScreenQuad, r.getUniforms(true), r.getDataMap(true))
	if r.drawQuad == nil {
		panic("Failed to create draw renderable for rect")
	}
}

func (r *Rect) SetX(x int) {
	if r.x != x {
		r.x = x
		r.dirtyFlags.position = true
	}
}

func (r *Rect) SetY(y int) {
	if r.y != y {
		r.y = y
		r.dirtyFlags.position = true
	}
}

func (r *Rect) SetWidth(width int) {
	if r.width != width {
		r.width = width
		r.dirtyFlags.size = true
	}
}

func (r *Rect) SetHeight(height int) {
	if r.height != height {
		r.height = height
		r.dirtyFlags.size = true
	}
}

func (r *Rect) SetCornerRadius(cornerRadius float32) {
	if r.cornerRadius != cornerRadius {
		r.cornerRadius = cornerRadius
		r.dirtyFlags.cornerRadius = true
	}
}

func (r *Rect) SetColor(color color.Color) {
	if r.color != color {
		r.color = color
		r.dirtyFlags.color = true
	}
}

func (r *Rect) SetOutlineColor(outlineColor color.Color) {
	if r.outlineColor != outlineColor {
		r.outlineColor = outlineColor
		r.dirtyFlags.outlineColor = true
	}
}

func (r *Rect) SetOutlineWidth(outlineWidth float32) {
	if r.outlineWidth != outlineWidth {
		r.outlineWidth = outlineWidth
		r.dirtyFlags.outlineWidth = true
	}
}

func (r *Rect) GetX() int {
	return r.x
}

func (r *Rect) GetY() int {
	return r.y
}

func (r *Rect) GetWidth() int {
	return r.width
}

func (r *Rect) GetHeight() int {
	return r.height
}

func (r *Rect) IsPointWithin(px, py int) bool {
	left := r.x
	right := r.x + r.width
	top := r.y
	bottom := r.y + r.height

	return px >= left && px <= right && py >= top && py <= bottom
}

func (r *Rect) getUniforms(isOutline bool) map[string]hlg.Uniform {
	pos := Position{
		X: float32(r.x),
		Y: float32(r.y),
	}
	size := Size{
		Width:  float32(r.width),
		Height: float32(r.height),
	}
	colorSlice := colorToFloatSlice(r.color)
	if isOutline {
		colorSlice = colorToFloatSlice(r.outlineColor)
	}
	radius := r.cornerRadius
	outlineWidth := r.outlineWidth

	return map[string]hlg.Uniform{
		"rect_pos": {
			Binding: 0,
			Size:    uint64(unsafe.Sizeof(pos)),
		},
		"rect_size": {
			Binding: 1,
			Size:    uint64(unsafe.Sizeof(size)),
		},
		"rect_color": {
			Binding: 2,
			Size:    uint64(unsafe.Sizeof(colorSlice)),
		},
		"corner_radius": {
			Binding: 3,
			Size:    uint64(unsafe.Sizeof(radius)),
		},
		"outline_width": {
			Binding: 4,
			Size:    uint64(unsafe.Sizeof(outlineWidth)),
		},
	}
}

// getDataMap generates the data map for rendering the rectangle.
func (r *Rect) getDataMap(isOutline bool) map[string][]byte {
	pos := Position{float32(r.x), float32(r.y)}
	size := Size{float32(r.width), float32(r.height)}
	colorSlice := colorToFloatSlice(r.color)
	if isOutline {
		colorSlice = colorToFloatSlice(r.outlineColor)
	}
	radius := r.cornerRadius
	outlineWidth := r.outlineWidth

	return map[string][]byte{
		"rect_pos":      unsafe.Slice((*byte)(unsafe.Pointer(&pos)), int(unsafe.Sizeof(pos))),
		"rect_size":     unsafe.Slice((*byte)(unsafe.Pointer(&size)), int(unsafe.Sizeof(size))),
		"rect_color":    unsafe.Slice((*byte)(unsafe.Pointer(&colorSlice[0])), int(unsafe.Sizeof(colorSlice))),
		"corner_radius": unsafe.Slice((*byte)(unsafe.Pointer(&radius)), int(unsafe.Sizeof(radius))),
		"outline_width": unsafe.Slice((*byte)(unsafe.Pointer(&outlineWidth)), int(unsafe.Sizeof(outlineWidth))),
	}
}

// UpdateUniforms updates the uniforms for the rectangle renderables if they are dirty.
func (r *Rect) UpdateUniforms() {
	if r.dirtyFlags.position || r.dirtyFlags.size || r.dirtyFlags.color || r.dirtyFlags.cornerRadius || r.dirtyFlags.outlineWidth || r.dirtyFlags.outlineColor {
		if r.fillQuad != nil && (r.dirtyFlags.position || r.dirtyFlags.size || r.dirtyFlags.color || r.dirtyFlags.cornerRadius) {
			r.fillQuad.UpdateUniforms(r.getDataMap(false))
		}

		if r.drawQuad != nil && (r.dirtyFlags.position || r.dirtyFlags.size || r.dirtyFlags.outlineColor || r.dirtyFlags.cornerRadius || r.dirtyFlags.outlineWidth) {
			r.drawQuad.UpdateUniforms(r.getDataMap(true))
		}

		// Reset dirty flags after updating
		r.dirtyFlags = dirtyFlags{}
	}
}

func (r *Rect) Render() {
	r.Fill()
	r.Draw()
}

func (r *Rect) Fill() {
	r.UpdateUniforms()
	r.fillQuad.Render()
}

func (r *Rect) Draw() {
	r.UpdateUniforms()
	r.drawQuad.Render()
}
