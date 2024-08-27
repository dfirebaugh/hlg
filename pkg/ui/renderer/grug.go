package renderer

import (
	_ "embed"
	"image"
	"image/color"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/pkg/fb"
	"github.com/dfirebaugh/hlg/pkg/math/geom"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	img_draw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type grugRenderer struct {
	*fb.ImageFB
	*hlg.Texture
	dirty       bool
	x, y        int
	visible     bool
	isCollapsed bool
}

func NewGrugRenderer(width, height int) *grugRenderer {
	r := &grugRenderer{
		ImageFB: fb.New(width, height),
		visible: true,
	}
	r.Texture, _ = hlg.CreateTexture(0, 0, width, height)
	r.Texture.Move(float32(width/2), float32(height/2))

	return r
}

func (r *grugRenderer) getBounds() image.Rectangle {
	width, height := r.Size()
	return image.Rect(0, 0, int(width), int(height))
}

func (r *grugRenderer) Clear(c color.Color) {
	transparentColor := color.RGBA{0, 0, 0, 0}
	if c != nil {
		transparentColor = c.(color.RGBA)
	}

	gc := draw2dimg.NewGraphicContext(r.ToImage())
	gc.SetFillColor(transparentColor)

	radius := 20.0

	width, height := r.Size()
	draw2dkit.RoundedRectangle(gc, 0, 0, float64(width), float64(height), radius, radius)
	gc.Fill()

	r.dirty = true
}

func (r *grugRenderer) Render() {
	// r.Clear()
	if !r.visible {
		return
	}
	r.Texture.UpdateImage(r.ImageFB.ToImage())
	r.Texture.Render()
}

func (r *grugRenderer) SetVisibility(visible bool) {
	r.visible = visible
}

func (r *grugRenderer) DrawPolygon(points []geom.Point, c color.Color) {
	gc := draw2dimg.NewGraphicContext(r.ToImage())
	gc.SetStrokeColor(c)
	gc.SetFillColor(c)

	gc.MoveTo(float64(points[0].X), float64(points[0].Y))

	for _, point := range points[1:] {
		gc.LineTo(float64(point.X), float64(point.Y))
	}

	gc.Close()
	gc.FillStroke()

	r.dirty = true
}

func (r *grugRenderer) FillPolygon(points []geom.Point, c color.Color) {
	gc := draw2dimg.NewGraphicContext(r.ToImage())
	gc.SetFillColor(c)

	gc.MoveTo(float64(points[0].X), float64(points[0].Y))

	for _, point := range points[1:] {
		gc.LineTo(float64(point.X), float64(point.Y))
	}

	gc.Close()
	gc.Fill()

	r.dirty = true
}

func (r *grugRenderer) DrawRoundedRectangle(x, y, width, height int, c color.Color) {
	radius := float64(height) / 2.0

	gc := draw2dimg.NewGraphicContext(r.ToImage())
	gc.SetStrokeColor(c)
	gc.SetFillColor(c)

	draw2dkit.RoundedRectangle(gc, float64(x), float64(y), float64(x+width), float64(y+height), radius, radius)

	gc.Stroke()

	r.dirty = true
}

func (r *grugRenderer) FillRoundedRectangle(x, y, width, height int, c color.Color) {
	radius := float64(height) / 2.0

	gc := draw2dimg.NewGraphicContext(r.ToImage())
	gc.SetFillColor(c)

	draw2dkit.RoundedRectangle(gc, float64(x), float64(y), float64(x+width), float64(y+height), radius, radius)

	gc.Fill()

	r.dirty = true
}

func (r *grugRenderer) DrawRect(x, y, width, height int, c color.Color) {
	gc := draw2dimg.NewGraphicContext(r.ToImage())
	gc.SetStrokeColor(c)
	draw2dkit.Rectangle(gc, float64(x), float64(y), float64(x+width), float64(y+height))
	gc.Stroke()
	r.dirty = true
}

func (r *grugRenderer) FillRect(x, y, width, height int, c color.Color) {
	gc := draw2dimg.NewGraphicContext(r.ToImage())
	gc.SetFillColor(c)
	draw2dkit.Rectangle(gc, float64(x), float64(y), float64(x+width), float64(y+height))
	gc.Fill()
	r.dirty = true
}

func (r *grugRenderer) DrawCircle(x, y, radius int, c color.Color) {
	gc := draw2dimg.NewGraphicContext(r.ToImage())
	gc.SetFillColor(c)
	gc.SetStrokeColor(c)
	draw2dkit.Circle(gc, float64(x), float64(y), float64(radius))
	gc.FillStroke()
	r.dirty = true
}

func (r *grugRenderer) FillCircle(x, y, radius int, c color.Color) {
	gc := draw2dimg.NewGraphicContext(r.ToImage())
	gc.SetFillColor(c)
	draw2dkit.Circle(gc, float64(x), float64(y), float64(radius))
	gc.Fill()
	r.dirty = true
}

func (r *grugRenderer) DrawTriangle(points [3]geom.Point, c color.Color) {
	gc := draw2dimg.NewGraphicContext(r.ToImage())
	gc.SetFillColor(c)
	gc.SetStrokeColor(c)
	gc.MoveTo(float64(points[0].X), float64(points[0].Y))
	gc.LineTo(float64(points[1].X), float64(points[1].Y))
	gc.LineTo(float64(points[2].X), float64(points[2].Y))
	gc.Close()
	gc.FillStroke()
	r.dirty = true
}

func (r *grugRenderer) FillTriangle(points [3]geom.Point, c color.Color) {
	gc := draw2dimg.NewGraphicContext(r.ToImage())
	gc.SetFillColor(c)
	gc.MoveTo(float64(points[0].X), float64(points[0].Y))
	gc.LineTo(float64(points[1].X), float64(points[1].Y))
	gc.LineTo(float64(points[2].X), float64(points[2].Y))
	gc.Close()
	gc.Fill()
	r.dirty = true
}

func (r *grugRenderer) DrawLine(x1, y1, x2, y2 int, c color.Color) {
	gc := draw2dimg.NewGraphicContext(r.ToImage())
	gc.SetStrokeColor(c)
	gc.MoveTo(float64(x1), float64(y1))
	gc.LineTo(float64(x2), float64(y2))
	gc.Stroke()
	r.dirty = true
}

func (r *grugRenderer) DrawText(x, y int, text string, c color.Color) {
	tempImg := image.NewRGBA(r.getBounds())
	col := image.NewUniform(c)

	point := fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y + basicfont.Face7x13.Ascent),
	}

	d := &font.Drawer{
		Dst:  tempImg,
		Src:  col,
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)

	img_draw.Draw(r.ToImage(), r.getBounds(), tempImg, image.Point{}, img_draw.Over)

	r.dirty = true
}

func (r *grugRenderer) TextWidth(text string) int {
	face := basicfont.Face7x13

	d := &font.Drawer{
		Face: face,
	}

	width := d.MeasureString(text)

	return width.Round()
}

func (r *grugRenderer) TextHeight(text string) int {
	face := basicfont.Face7x13

	height := face.Metrics().Ascent + face.Metrics().Descent

	return height.Round()
}

func (r *grugRenderer) Clip(x, y, width, height int) {
	originalWidth := r.Width()
	originalHeight := r.Height()

	r.Texture.Clip(float32(x), float32(y), float32(width), float32(height))

	if width < originalWidth || height < originalHeight {
		r.Move(r.GetX(), r.GetY())
		r.isCollapsed = true
	} else {
		r.isCollapsed = false
		r.Move(r.GetX(), r.GetY())
	}
}

func (r *grugRenderer) Resize(width int, height int) {
	r.Texture.Resize(float32(width/2), float32(height/2))
}

func (r *grugRenderer) Move(x int, y int) {
	if r.isCollapsed {
		return
	}
	r.x = x
	r.y = y
	r.Texture.Move(float32(x), float32(y))
}

func (r *grugRenderer) GetX() int {
	return r.x
}

func (r *grugRenderer) GetY() int {
	return r.y
}
