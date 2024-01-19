package hlg

import "image"

type Sprite struct {
	*Texture
	frameSize    image.Point
	sheetSize    image.Point
	currentFrame image.Point
}

func NewSprite(img image.Image, frameSize, sheetSize image.Point) *Sprite {
	texture, err := CreateTextureFromImage(img)
	if err != nil {
		panic(err)
	}
	s := &Sprite{
		Texture:   texture,
		frameSize: frameSize,
		sheetSize: sheetSize,
	}
	s.Clip(0, 0, float32(s.frameSize.X), float32(s.frameSize.Y))
	return s
}

func (s *Sprite) NextFrame() {
	s.currentFrame.X++
	if s.currentFrame.X >= s.sheetSize.X {
		s.currentFrame.X = 0
		s.currentFrame.Y++
		if s.currentFrame.Y >= s.sheetSize.Y {
			s.currentFrame.Y = 0
		}
	}

	x0 := s.currentFrame.X * s.frameSize.X
	y0 := s.currentFrame.Y * s.frameSize.Y
	x1 := x0 + s.frameSize.X
	y1 := y0 + s.frameSize.Y
	s.Clip(float32(x0), float32(y0), float32(x1), float32(y1))

	s.Render()
}
