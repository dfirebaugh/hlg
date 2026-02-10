package gui

import (
	"image/color"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/pkg/input"
)

// Button renders a button and returns true if clicked.
func (c *Context) Button(label string, x, y, w, h int) bool {
	id := c.GetID(label)
	c.registerFocusable(id)

	mx, my := c.input.MousePosition()
	blocked := c.isInputBlocked(mx, my)
	hovered := !blocked && pointInRect(mx, my, x, y, w, h)
	focused := c.isFocused(id)

	if hovered {
		c.setHot(id)
	}

	clicked := false
	if hovered && c.input.IsButtonJustPressed(input.MouseButtonLeft) {
		c.setActive(id)
		c.setFocused(id)
	}

	if c.isActive(id) && c.input.IsButtonJustReleased(input.MouseButtonLeft) {
		if hovered {
			clicked = true
		}
		c.clearActive()
	}

	if focused {
		if c.input.IsKeyJustPressed(input.KeyEnter) || c.input.IsKeyJustPressed(input.KeySpace) {
			clicked = true
		}
	}

	bgColor := color.RGBA{70, 70, 80, 255}
	if c.isActive(id) {
		bgColor = color.RGBA{50, 50, 60, 255}
	} else if hovered || focused {
		bgColor = color.RGBA{90, 90, 100, 255}
	}

	hlg.RoundedRect(x, y, w, h, 4, bgColor)

	if focused {
		hlg.RoundedRectOutline(x+1, y+1, w-2, h-2, 3, 1, color.RGBA{100, 150, 255, 255}, color.RGBA{100, 150, 255, 255})
	}

	textColor := color.RGBA{220, 220, 220, 255}
	textX := x + (w-int(hlg.MeasureText(label, 14)))/2
	textY := y + (h-14)/2 + 2
	hlg.Text(label, textX, textY, 14, textColor)

	return clicked
}

// Checkbox renders a checkbox and returns true if the state changed.
func (c *Context) Checkbox(label string, checked *bool, x, y, size int) bool {
	id := c.GetID(label)
	c.registerFocusable(id)

	mx, my := c.input.MousePosition()
	blocked := c.isInputBlocked(mx, my)
	labelWidth := len(label) * 8
	hitW := size + 8 + labelWidth
	hovered := !blocked && pointInRect(mx, my, x, y, hitW, size)
	focused := c.isFocused(id)

	if hovered {
		c.setHot(id)
	}

	changed := false

	if hovered && c.input.IsButtonJustPressed(input.MouseButtonLeft) {
		c.setActive(id)
		c.setFocused(id)
	}

	if c.isActive(id) && c.input.IsButtonJustReleased(input.MouseButtonLeft) {
		if hovered {
			*checked = !*checked
			changed = true
		}
		c.clearActive()
	}

	if focused {
		if c.input.IsKeyJustPressed(input.KeyEnter) || c.input.IsKeyJustPressed(input.KeySpace) {
			*checked = !*checked
			changed = true
		}
	}

	boxColor := color.RGBA{60, 60, 70, 255}
	if hovered || focused {
		boxColor = color.RGBA{80, 80, 90, 255}
	}

	hlg.RoundedRect(x, y, size, size, 3, boxColor)

	if focused {
		hlg.RoundedRectOutline(x+2, y+2, size-4, size-4, 1, 2, color.RGBA{100, 150, 255, 255}, color.RGBA{100, 150, 255, 255})
	}

	if *checked {
		checkColor := color.RGBA{100, 200, 100, 255}
		margin := size / 4
		hlg.RoundedRect(x+margin, y+margin, size-margin*2, size-margin*2, 2, checkColor)
	}

	textColor := color.RGBA{220, 220, 220, 255}
	hlg.Text(label, x+size+8, y+(size-14)/2, 14, textColor)

	return changed
}

// Toggle renders a toggle switch and returns true if the state changed.
func (c *Context) Toggle(label string, on *bool, x, y, w, h int) bool {
	id := c.GetID(label + "_toggle")
	c.registerFocusable(id)

	mx, my := c.input.MousePosition()
	blocked := c.isInputBlocked(mx, my)
	hovered := !blocked && pointInRect(mx, my, x, y, w, h)
	focused := c.isFocused(id)

	if hovered {
		c.setHot(id)
	}

	changed := false

	if hovered && c.input.IsButtonJustPressed(input.MouseButtonLeft) {
		c.setActive(id)
		c.setFocused(id)
	}

	if c.isActive(id) && c.input.IsButtonJustReleased(input.MouseButtonLeft) {
		if hovered {
			*on = !*on
			changed = true
		}
		c.clearActive()
	}

	if focused {
		if c.input.IsKeyJustPressed(input.KeyEnter) || c.input.IsKeyJustPressed(input.KeySpace) {
			*on = !*on
			changed = true
		}
	}

	trackColor := color.RGBA{60, 60, 70, 255}
	if *on {
		trackColor = color.RGBA{60, 140, 80, 255}
	}

	radius := h / 2
	hlg.RoundedRect(x, y, w, h, radius, trackColor)

	if focused {
		hlg.RoundedRectOutline(x+2, y+2, w-4, h-4, radius-2, 2, color.RGBA{100, 150, 255, 255}, color.RGBA{100, 150, 255, 255})
	}

	knobSize := h - 4
	knobX := x + 2
	if *on {
		knobX = x + w - knobSize - 2
	}

	knobColor := color.RGBA{220, 220, 220, 255}
	if hovered {
		knobColor = color.RGBA{255, 255, 255, 255}
	}

	hlg.RoundedRect(knobX, y+2, knobSize, knobSize, knobSize/2, knobColor)

	return changed
}

// Slider renders a horizontal slider and returns true if the value changed.
// The label is used for widget identification (not displayed).
func (c *Context) Slider(label string, value *float32, min, max float32, x, y, w, h int) bool {
	id := c.GetID(label + "_slider")
	c.registerFocusable(id)

	mx, my := c.input.MousePosition()
	blocked := c.isInputBlocked(mx, my)
	hovered := !blocked && pointInRect(mx, my, x, y-2, w, h+8)
	focused := c.isFocused(id)

	if hovered {
		c.setHot(id)
	}

	changed := false
	oldValue := *value

	if hovered && c.input.IsButtonJustPressed(input.MouseButtonLeft) {
		c.setActive(id)
		c.setFocused(id)
	}

	if c.isActive(id) {
		if c.input.IsButtonPressed(input.MouseButtonLeft) {
			ratio := float32(mx-x) / float32(w)
			if ratio < 0 {
				ratio = 0
			}
			if ratio > 1 {
				ratio = 1
			}
			*value = min + ratio*(max-min)
			if *value != oldValue {
				changed = true
			}
		} else {
			c.clearActive()
		}
	}

	if focused {
		step := (max - min) / 20
		if c.input.IsKeyJustPressed(input.KeyLeft) {
			*value -= step
			if *value < min {
				*value = min
			}
			changed = true
		}
		if c.input.IsKeyJustPressed(input.KeyRight) {
			*value += step
			if *value > max {
				*value = max
			}
			changed = true
		}
	}

	trackColor := color.RGBA{50, 50, 60, 255}
	hlg.RoundedRect(x, y+2, w, h, h/2, trackColor)

	if focused {
		hlg.RoundedRectOutline(x+2, y+4, w-4, h-4, h/2-2, 2, color.RGBA{100, 150, 255, 255}, color.RGBA{100, 150, 255, 255})
	}

	ratio := (*value - min) / (max - min)
	filledWidth := int(float32(w) * ratio)
	if filledWidth > 0 {
		fillColor := color.RGBA{80, 140, 200, 255}
		hlg.RoundedRect(x, y+2, filledWidth, h, h/2, fillColor)
	}

	handleRadius := h + 4
	handleX := x + int(float32(w)*ratio) - handleRadius/2
	handleY := y

	handleColor := color.RGBA{200, 200, 210, 255}
	if c.isActive(id) || hovered {
		handleColor = color.RGBA{240, 240, 250, 255}
	}

	hlg.RoundedRect(handleX, handleY, handleRadius, handleRadius, handleRadius/2, handleColor)

	return changed
}

// Panel renders a draggable panel with a title bar.
// The label is used for both the title and widget identification.
// Returns true if the panel is expanded (not collapsed).
// The state parameter holds position and collapse state.
func (c *Context) Panel(label string, state *PanelState, w, h int) bool {
	id := c.GetID(label + "_panel")
	titleBarHeight := 28

	actualHeight := h
	if state.Collapsed {
		actualHeight = titleBarHeight
	}

	c.registerPanelBounds(id, state.X, state.Y, w, actualHeight)

	mx, my := c.input.MousePosition()

	blocked := c.isBlockedByLaterPanel(id, mx, my)

	collapseButtonX := state.X + w - 30
	collapseButtonW := 28
	collapseHovered := !blocked && pointInRect(mx, my, collapseButtonX, state.Y, collapseButtonW, titleBarHeight)

	titleHovered := !blocked && pointInRect(mx, my, state.X, state.Y, w-collapseButtonW, titleBarHeight)

	if titleHovered || collapseHovered {
		c.setHot(id)
	}

	if collapseHovered && c.input.IsButtonJustPressed(input.MouseButtonLeft) {
		state.Collapsed = !state.Collapsed
	}

	if titleHovered && c.input.IsButtonJustPressed(input.MouseButtonLeft) {
		c.setActive(id)
		state.dragging = true
		state.dragOffX = mx - state.X
		state.dragOffY = my - state.Y
	}

	if state.dragging {
		if c.input.IsButtonPressed(input.MouseButtonLeft) {
			state.X = mx - state.dragOffX
			state.Y = my - state.dragOffY

			sw, sh := hlg.GetScreenSize()
			if state.X < 0 {
				state.X = 0
			}
			if state.Y < 0 {
				state.Y = 0
			}
			if state.X+w > sw {
				state.X = sw - w
			}
			if state.Y+titleBarHeight > sh {
				state.Y = sh - titleBarHeight
			}
		} else {
			state.dragging = false
			c.clearActive()
		}
	}

	bgColor := color.RGBA{45, 45, 50, 240}
	titleBgColor := color.RGBA{55, 55, 65, 255}
	if c.isActive(id) || titleHovered {
		titleBgColor = color.RGBA{65, 65, 75, 255}
	}
	titleTextColor := color.RGBA{220, 220, 225, 255}

	hlg.RoundedRect(state.X, state.Y, w, actualHeight, 8, bgColor)

	hlg.RoundedRect(state.X, state.Y, w, titleBarHeight, 8, titleBgColor)
	if !state.Collapsed {
		hlg.FilledRect(state.X, state.Y+titleBarHeight-8, w, 8, titleBgColor)
	}

	hlg.Text(label, state.X+10, state.Y+10, 14, titleTextColor)

	indicatorX := state.X + w - 20
	indicatorY := state.Y + 10
	indicatorColor := color.RGBA{150, 150, 160, 255}
	if collapseHovered {
		indicatorColor = color.RGBA{200, 200, 210, 255}
	}
	if state.Collapsed {
		hlg.FilledTriangle(indicatorX, indicatorY, indicatorX, indicatorY+8, indicatorX+6, indicatorY+4, indicatorColor)
		c.ClearCurrentPanel()
	} else {
		hlg.FilledTriangle(indicatorX, indicatorY, indicatorX+8, indicatorY, indicatorX+4, indicatorY+6, indicatorColor)
		c.SetCurrentPanel(id)
	}

	return !state.Collapsed
}
