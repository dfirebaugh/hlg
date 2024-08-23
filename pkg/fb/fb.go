package fb

import (
	"image"
	"image/color"
)

// Displayer defines the methods that any display-compatible framebuffer must implement.
type Displayer interface {
	SetPixel(x, y int16, c color.RGBA)
	Display() error
	Size() (int16, int16)
}

// ImageFB encapsulates an in-memory image framebuffer.
type ImageFB struct {
	img *image.RGBA
}

// New creates and returns a new instance of ImageFB with the specified width and height.
func New(width, height int) *ImageFB {
	return &ImageFB{
		img: image.NewRGBA(image.Rect(0, 0, width, height)),
	}
}

func (fb *ImageFB) Width() int {
	return fb.img.Rect.Dx()
}

func (fb *ImageFB) Height() int {
	return fb.img.Rect.Dy()
}

// SetPixel sets the pixel color at the specified x and y coordinates within the framebuffer.
// It does nothing if the coordinates are out of bounds.
func (i *ImageFB) SetPixel(x, y int16, c color.RGBA) {
	bounds := i.img.Bounds()
	if !image.Pt(int(x), int(y)).In(bounds) {
		return
	}
	i.img.SetRGBA(int(x), int(y), c)
}

// Display implements the Display method for the displayer interface.
// It currently does nothing and returns nil.
func (i *ImageFB) Display() error {
	// TODO: Implement display logic, if necessary.
	return nil
}

// Size returns the width and height dimensions of the framebuffer.
func (i *ImageFB) Size() (int16, int16) {
	bounds := i.img.Bounds()
	return int16(bounds.Dx()), int16(bounds.Dy())
}

// ToImage returns the internal image.RGBA image.
func (i *ImageFB) ToImage() *image.RGBA {
	return i.img
}
