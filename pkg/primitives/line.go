package primitives

import (
	"image/color"
)

type Line struct {
	points    []Position
	segments  []*Segment
	thickness float32
	color     color.Color
}

func NewLine(points []Position, thickness float32, color color.Color) *Line {
	l := &Line{
		points:    points,
		thickness: thickness,
		color:     color,
	}
	l.init()
	return l
}

func (l *Line) init() {
	for i := 0; i < len(l.points)-1; i++ {
		segment := NewSegment(
			int(l.points[i].X),
			int(l.points[i].Y),
			int(l.points[i+1].X),
			int(l.points[i+1].Y),
			l.thickness,
			l.color,
		)
		l.segments = append(l.segments, segment)
	}
}

func (l *Line) SetPoints(points []Position) {
	l.points = points
	l.init()
}

func (l *Line) SetThickness(thickness float32) {
	l.thickness = thickness
	for _, segment := range l.segments {
		segment.SetThickness(thickness)
	}
}

func (l *Line) SetColor(color color.Color) {
	l.color = color
	for _, segment := range l.segments {
		segment.SetColor(color)
	}
}

func (l *Line) Render() {
	for _, segment := range l.segments {
		segment.Render()
	}
}
