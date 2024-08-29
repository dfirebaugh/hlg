package components

import (
	"time"

	ui "github.com/dfirebaugh/hlg/pkg/grugui"
	"github.com/dfirebaugh/hlg/pkg/input"
)

type InputBox struct {
	X, Y, Width, Height int
	Text                string
	IsFocused           bool
	cursorPosition      int
	lastBlinkTime       time.Time
	cursorVisible       bool
}

func (ib *InputBox) Update(ctx ui.Context) {
	mouseX, mouseY := ctx.GetCursorPosition()

	if ctx.IsButtonJustPressed(input.MouseButtonLeft) &&
		mouseX > ib.X && mouseX < ib.X+ib.Width &&
		mouseY > ib.Y && mouseY < ib.Y+ib.Height {
		ib.IsFocused = true
	} else if ctx.IsButtonJustPressed(input.MouseButtonLeft) {
		ib.IsFocused = false
	}

	if ib.IsFocused {
		ib.handleKeyInput(ctx)
	}

	if time.Since(ib.lastBlinkTime) >= 500*time.Millisecond {
		ib.cursorVisible = !ib.cursorVisible
		ib.lastBlinkTime = time.Now()
	}
}

func (ib *InputBox) Render(ctx ui.Context) {
	theme := ctx.Theme()

	ctx.FillRect(ib.X, ib.Y, ib.Width, ib.Height, theme.BackgroundColor)
	ctx.DrawRect(ib.X, ib.Y, ib.Width, ib.Height, theme.PrimaryColor)

	textColor := theme.TextColor
	textY := ib.Y + (ib.Height/2 - 5)
	ctx.DrawText(ib.X+5, textY, ib.Text, textColor)

	if ib.IsFocused && ib.cursorVisible {
		cursorX := ib.X + 5 + ctx.TextWidth(ib.Text[:ib.cursorPosition])
		ctx.DrawLine(cursorX, ib.Y+5, cursorX, ib.Y+ib.Height-5, textColor)
	}
}

func (ib *InputBox) handleKeyInput(ctx ui.Context) {
	if ctx.IsKeyJustPressed(input.KeyBackspace) {
		if len(ib.Text) > 0 && ib.cursorPosition > 0 {
			ib.Text = ib.Text[:ib.cursorPosition-1] + ib.Text[ib.cursorPosition:]
			ib.cursorPosition--
		}
	}

	if ctx.IsKeyJustPressed(input.KeyLeft) {
		if ib.cursorPosition > 0 {
			ib.cursorPosition--
		}
	}

	if ctx.IsKeyJustPressed(input.KeyRight) {
		if ib.cursorPosition < len(ib.Text) {
			ib.cursorPosition++
		}
	}

	if ctx.IsKeyJustPressed(input.KeyEnter) {
		ib.IsFocused = false
	}

	for key := input.KeyA; key <= input.KeyZ; key++ {
		if ctx.IsKeyJustPressed(key) {
			letter := string('a' + key - input.KeyA)
			ib.Text = ib.Text[:ib.cursorPosition] + letter + ib.Text[ib.cursorPosition:]
			ib.cursorPosition++
		}
	}

	for key := input.Key0; key <= input.Key9; key++ {
		if ctx.IsKeyJustPressed(key) {
			number := string('0' + key - input.Key0)
			ib.Text = ib.Text[:ib.cursorPosition] + number + ib.Text[ib.cursorPosition:]
			ib.cursorPosition++
		}
	}

	if ctx.IsKeyJustPressed(input.KeySpace) {
		ib.Text = ib.Text[:ib.cursorPosition] + " " + ib.Text[ib.cursorPosition:]
		ib.cursorPosition++
	}
}
