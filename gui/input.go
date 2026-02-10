package gui

import (
	"image/color"
	"time"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/pkg/input"
)

const (
	blinkInterval = 500 * time.Millisecond
)

// InputText renders a text input field and returns (changed, submitted).
// The text is stored in the provided pointer.
// The state parameter holds cursor and selection information.
// The label is used for widget identification (not displayed).
func (c *Context) InputText(label string, text *string, state *TextInputState, x, y, w, h int) (changed, submitted bool) {
	id := c.GetID(label + "_input")
	c.registerFocusable(id)

	mx, my := c.input.MousePosition()
	blocked := c.isInputBlocked(mx, my)
	hovered := !blocked && pointInRect(mx, my, x, y, w, h)
	focused := c.isFocused(id)

	if hovered {
		c.setHot(id)
	}

	if hovered && c.input.IsButtonJustPressed(input.MouseButtonLeft) {
		c.setFocused(id)
		state.CursorPos = c.calculateCursorPosition(*text, mx-x-4)
		state.ClearSelection()
	}

	if focused {
		changed, submitted = c.handleTextInput(text, state)
	}

	state.focused = focused

	bgColor := color.RGBA{40, 40, 50, 255}
	if focused {
		bgColor = color.RGBA{50, 50, 60, 255}
	}
	hlg.RoundedRect(x, y, w, h, 4, bgColor)

	if focused {
		hlg.RoundedRectOutline(x+2, y+2, w-4, h-4, 2, 2, color.RGBA{100, 150, 255, 255}, color.RGBA{100, 150, 255, 255})
	}

	hlg.PushClipRect(x+2, y+2, w-4, h-4)

	textColor := color.RGBA{220, 220, 220, 255}
	textY := y + (h-14)/2
	hlg.Text(*text, x+4, textY, 14, textColor)

	if focused {
		elapsed := time.Since(c.frameTime)
		showCursor := (elapsed/blinkInterval)%2 == 0

		if showCursor || c.isActive(id) {
			textBeforeCursor := (*text)[:state.CursorPos]
			cursorX := x + 4 + int(hlg.MeasureText(textBeforeCursor, 14))
			cursorColor := color.RGBA{200, 200, 200, 255}
			hlg.FilledRect(cursorX, y+4, 2, h-8, cursorColor)
		}
	}

	hlg.PopClipRect()

	return changed, submitted
}

// calculateCursorPosition determines cursor position from x offset.
func (c *Context) calculateCursorPosition(text string, xOffset int) int {
	if xOffset <= 0 {
		return 0
	}

	for i := 1; i <= len(text); i++ {
		width := hlg.MeasureText(text[:i], 14)
		if int(width) >= xOffset {
			prevWidth := float32(0)
			if i > 1 {
				prevWidth = hlg.MeasureText(text[:i-1], 14)
			}
			if float32(xOffset)-prevWidth < width-float32(xOffset) {
				return i - 1
			}
			return i
		}
	}
	return len(text)
}

// handleTextInput processes keyboard input for text editing.
func (c *Context) handleTextInput(text *string, state *TextInputState) (changed, submitted bool) {
	chars := c.input.GetCharInput()
	for _, ch := range chars {
		if ch >= 32 && ch < 127 { // Printable ASCII
			*text = (*text)[:state.CursorPos] + string(ch) + (*text)[state.CursorPos:]
			state.CursorPos++
			changed = true
		}
	}

	if c.input.IsKeyJustPressed(input.KeyBackspace) {
		if state.CursorPos > 0 {
			*text = (*text)[:state.CursorPos-1] + (*text)[state.CursorPos:]
			state.CursorPos--
			changed = true
		}
	}

	if c.input.IsKeyJustPressed(input.KeyDelete) {
		if state.CursorPos < len(*text) {
			*text = (*text)[:state.CursorPos] + (*text)[state.CursorPos+1:]
			changed = true
		}
	}

	if c.input.IsKeyJustPressed(input.KeyLeft) {
		if state.CursorPos > 0 {
			state.CursorPos--
		}
	}

	if c.input.IsKeyJustPressed(input.KeyRight) {
		if state.CursorPos < len(*text) {
			state.CursorPos++
		}
	}

	if c.input.IsKeyJustPressed(input.KeyHome) {
		state.CursorPos = 0
	}

	if c.input.IsKeyJustPressed(input.KeyEnd) {
		state.CursorPos = len(*text)
	}

	if c.input.IsKeyJustPressed(input.KeyEnter) {
		submitted = true
	}

	return changed, submitted
}
