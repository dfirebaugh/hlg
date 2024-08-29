package components

import (
	"image/color"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/gui"
	"github.com/dfirebaugh/hlg/pkg/input"
)

type Button struct {
	Node
	x, y, width, height int
	cornerRadius        int
	onClick             func()
	isHovered           bool
}

func NewButton(x, y, width, height int, cornerRadius int, label string, onClick func()) *Button {
	b := &Button{
		x:            x,
		y:            y,
		width:        width,
		height:       height,
		cornerRadius: cornerRadius,
		onClick:      onClick,
	}

	return b
}

func (b *Button) Render(ctx gui.DrawContext) {
	ctx.DrawRectangle(b.x, b.y, b.width, b.height, &gui.DrawOptions{
		Style: gui.Style{
			FillColor:    b.getFillColor(ctx),
			OutlineColor: ctx.GetTheme().SecondaryColor,
			OutlineSize:  2,
			CornerRadius: b.cornerRadius,
		},
	})

	b.Node.Render(ctx)
}

func (b *Button) Update() {
	b.Node.Update()
	mouseX, mouseY := hlg.GetCursorPosition()
	b.isHovered = b.isPointWithin(mouseX, mouseY)

	if hlg.IsButtonPressed(input.MouseButtonLeft) && b.isHovered {
		b.onClick()
	}
}

func (b *Button) getFillColor(ctx gui.DrawContext) color.Color {
	if b.isHovered {
		return ctx.GetTheme().BackgroundColor
	}
	return ctx.GetTheme().PrimaryColor
}

func (b *Button) isPointWithin(x, y int) bool {
	return x >= b.x && x <= b.x+b.width && y >= b.y && y <= b.y+b.height
}
