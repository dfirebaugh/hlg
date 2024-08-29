package primitives

import (
	"image/color"
	"unsafe"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/graphics"
)

type Segment struct {
	x1, y1, x2, y2 int
	thickness      float32
	color          color.Color
	segmentQuad    graphics.ShaderRenderable
}

func NewSegment(x1, y1, x2, y2 int, thickness float32, color color.Color) *Segment {
	s := &Segment{
		x1:        x1,
		y1:        y1,
		x2:        x2,
		y2:        y2,
		thickness: thickness,
		color:     color,
	}
	s.init()
	return s
}

func (s *Segment) init() {
	ww, wh := hlg.GetWindowSize()
	fullScreenQuad := makeFullScreenQuad(float32(ww), float32(wh))

	s.segmentQuad = hlg.CreateRenderable(LineShader, fullScreenQuad, s.getUniforms(), s.getDataMap())
	if s.segmentQuad == nil {
		panic("Failed to create segment renderable")
	}
}

func (s *Segment) SetStart(x, y int) {
	s.x1 = x
	s.y1 = y
	s.UpdateUniforms()
}

func (s *Segment) SetEnd(x, y int) {
	s.x2 = x
	s.y2 = y
	s.UpdateUniforms()
}

func (s *Segment) SetThickness(thickness float32) {
	s.thickness = thickness
	s.UpdateUniforms()
}

func (s *Segment) SetColor(color color.Color) {
	s.color = color
	s.UpdateUniforms()
}

func (s *Segment) getUniforms() map[string]hlg.Uniform {
	start := Position{float32(s.x1), float32(s.y1)}
	end := Position{float32(s.x2), float32(s.y2)}
	colorSlice := colorToFloatSlice(s.color)
	thickness := s.thickness

	return map[string]hlg.Uniform{
		"line_start": {
			Binding: 0,
			Size:    uint64(unsafe.Sizeof(start)),
		},
		"line_end": {
			Binding: 1,
			Size:    uint64(unsafe.Sizeof(end)),
		},
		"line_color": {
			Binding: 2,
			Size:    uint64(unsafe.Sizeof(colorSlice)),
		},
		"line_thickness": {
			Binding: 3,
			Size:    uint64(unsafe.Sizeof(thickness)),
		},
	}
}

func (s *Segment) getDataMap() map[string][]byte {
	start := Position{float32(s.x1), float32(s.y1)}
	end := Position{float32(s.x2), float32(s.y2)}
	colorSlice := colorToFloatSlice(s.color)
	thickness := s.thickness

	return map[string][]byte{
		"line_start":     unsafe.Slice((*byte)(unsafe.Pointer(&start)), int(unsafe.Sizeof(start))),
		"line_end":       unsafe.Slice((*byte)(unsafe.Pointer(&end)), int(unsafe.Sizeof(end))),
		"line_color":     unsafe.Slice((*byte)(unsafe.Pointer(&colorSlice[0])), int(unsafe.Sizeof(colorSlice))),
		"line_thickness": unsafe.Slice((*byte)(unsafe.Pointer(&thickness)), int(unsafe.Sizeof(thickness))),
	}
}

func (s *Segment) UpdateUniforms() {
	if s.segmentQuad != nil {
		s.segmentQuad.UpdateUniforms(s.getDataMap())
	}
}

func (s *Segment) Render() {
	s.UpdateUniforms()
	s.segmentQuad.Render()
}
