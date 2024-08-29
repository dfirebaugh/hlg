package components

import (
	ui "github.com/dfirebaugh/hlg/pkg/grugui"
	"github.com/dfirebaugh/hlg/pkg/input"
)

type SurfaceHandle struct {
	Width, Height    int
	Label            string
	isDragging       bool
	isCollapsed      bool
	offsetX, offsetY int
	originalWidth    int
	originalHeight   int
}

func (sh *SurfaceHandle) Update(ctx ui.Context) {
	surfaceX := ctx.GetX()
	surfaceY := ctx.GetY()

	mouseX, mouseY := ctx.GetGlobalCursorPosition()

	if ctx.IsButtonJustPressed(input.MouseButtonLeft) {
		buttonX := sh.Width - 20
		buttonY := 5

		if mouseX > surfaceX+buttonX && mouseX < surfaceX+buttonX+15 &&
			mouseY > surfaceY+buttonY && mouseY < surfaceY+buttonY+15 {
			sh.toggleCollapse(ctx)
			return
		}

		if mouseX > surfaceX && mouseX < surfaceX+sh.Width &&
			mouseY > surfaceY && mouseY < surfaceY+sh.Height {
			sh.isDragging = true
			sh.offsetX = mouseX - surfaceX
			sh.offsetY = mouseY - surfaceY
		}
	}

	if sh.isDragging {
		if ctx.IsButtonPressed(input.MouseButtonLeft) {
			newX := mouseX - sh.offsetX
			newY := mouseY - sh.offsetY
			ctx.Move(newX, newY)
		} else {
			sh.isDragging = false
		}
	}
}

func (sh *SurfaceHandle) Render(ctx ui.Context) {
	theme := ctx.Theme()

	ctx.FillRect(0, 0, sh.Width, sh.Height, theme.SecondaryColor)

	buttonX := sh.Width - 20
	buttonY := 5
	buttonColor := theme.PrimaryColor
	if sh.isCollapsed {
		buttonColor = theme.SecondaryColor
	}

	ctx.FillRect(buttonX, buttonY, 15, 15, buttonColor)
	ctx.DrawRect(buttonX, buttonY, 15, 15, theme.PrimaryColor)

	labelX := 10
	labelY := sh.Height/2 - ctx.TextHeight(sh.Label)/2
	ctx.DrawText(labelX, labelY, sh.Label, theme.TextColor)
}

func (sh *SurfaceHandle) toggleCollapse(ctx ui.Context) {
	if sh.isCollapsed {
		ctx.Clip(0, 0, sh.originalWidth, sh.originalHeight)
		sh.isCollapsed = false
	} else {
		sh.originalWidth = ctx.Width()
		sh.originalHeight = ctx.Height()
		ctx.Clip(0, 0, sh.Width, sh.Height)
		sh.isCollapsed = true
	}
}
