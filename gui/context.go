package gui

import (
	"hash/fnv"
	"image/color"
	"time"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/pkg/input"
)

// ID is a unique identifier for a widget.
type ID uint64

// InputContext defines the input interface needed by gui.
type InputContext interface {
	MousePosition() (int, int)
	IsButtonPressed(button input.MouseButton) bool
	IsButtonJustPressed(button input.MouseButton) bool
	IsButtonJustReleased(button input.MouseButton) bool
	IsKeyPressed(key input.Key) bool
	IsKeyJustPressed(key input.Key) bool
	GetCharInput() []rune
	Update()
}

// panelBound stores a panel's bounds for input blocking.
type panelBound struct {
	id         ID
	x, y, w, h int
}

// Context holds the state for an immediate mode GUI frame.
type Context struct {
	input      InputContext
	hotID      ID // Widget being hovered
	activeID   ID // Widget being interacted with
	focusedID  ID // Widget receiving keyboard input
	idStack    []ID
	frameTime  time.Time
	focusables []ID // List of focusable widgets in render order
	tabPressed bool
	shiftHeld  bool

	// Panel input blocking (uses previous frame's bounds)
	panelBounds     []panelBound // Current frame's panel bounds
	prevPanelBounds []panelBound // Previous frame's panel bounds (used for blocking)
	currentPanelID  ID           // The panel currently being rendered into (for widget blocking)
}

// NewContext creates a new gui context.
func NewContext(input InputContext) *Context {
	return &Context{
		input:   input,
		idStack: make([]ID, 0, 8),
	}
}

// Begin starts a new immediate mode frame.
func (c *Context) Begin() {
	c.frameTime = time.Now()
	c.hotID = 0
	c.focusables = c.focusables[:0]

	// Swap panel bounds for input blocking (use previous frame's data)
	c.prevPanelBounds = c.panelBounds
	c.panelBounds = c.panelBounds[:0]
	c.currentPanelID = 0

	// Handle tab navigation
	c.tabPressed = c.input.IsKeyJustPressed(input.KeyTab)
	c.shiftHeld = c.input.IsKeyPressed(input.KeyLeftShift) || c.input.IsKeyPressed(input.KeyRightShift)

	// Start batched drawing
	hlg.BeginDraw()
}

// End finishes the immediate mode frame.
func (c *Context) End() {
	// Handle tab navigation after all widgets have registered
	if c.tabPressed && len(c.focusables) > 0 {
		c.handleTabNavigation()
	}

	// If we clicked somewhere and nothing is hot, clear focus
	if c.input.IsButtonJustPressed(input.MouseButtonLeft) && c.hotID == 0 {
		c.focusedID = 0
	}

	// Submit batched drawing
	hlg.EndDraw()
}

// handleTabNavigation moves focus between widgets.
func (c *Context) handleTabNavigation() {
	if len(c.focusables) == 0 {
		return
	}

	currentIndex := -1
	for i, id := range c.focusables {
		if id == c.focusedID {
			currentIndex = i
			break
		}
	}

	if c.shiftHeld {
		// Move backwards
		if currentIndex <= 0 {
			c.focusedID = c.focusables[len(c.focusables)-1]
		} else {
			c.focusedID = c.focusables[currentIndex-1]
		}
	} else {
		// Move forwards
		if currentIndex < 0 || currentIndex >= len(c.focusables)-1 {
			c.focusedID = c.focusables[0]
		} else {
			c.focusedID = c.focusables[currentIndex+1]
		}
	}
}

// registerFocusable adds a widget to the tab order.
func (c *Context) registerFocusable(id ID) {
	c.focusables = append(c.focusables, id)
}

// PushID pushes an identifier onto the ID stack.
func (c *Context) PushID(id string) {
	h := fnv.New64a()
	if len(c.idStack) > 0 {
		parent := c.idStack[len(c.idStack)-1]
		h.Write([]byte{byte(parent), byte(parent >> 8), byte(parent >> 16), byte(parent >> 24)})
	}
	h.Write([]byte(id))
	c.idStack = append(c.idStack, ID(h.Sum64()))
}

// PopID removes the top identifier from the ID stack.
func (c *Context) PopID() {
	if len(c.idStack) > 0 {
		c.idStack = c.idStack[:len(c.idStack)-1]
	}
}

// GetID computes a unique ID for a widget label.
func (c *Context) GetID(label string) ID {
	h := fnv.New64a()
	if len(c.idStack) > 0 {
		parent := c.idStack[len(c.idStack)-1]
		h.Write([]byte{byte(parent), byte(parent >> 8), byte(parent >> 16), byte(parent >> 24)})
	}
	h.Write([]byte(label))
	return ID(h.Sum64())
}

// isActive returns true if the given widget is the active widget.
func (c *Context) isActive(id ID) bool {
	return c.activeID == id
}

// isFocused returns true if the given widget has keyboard focus.
func (c *Context) isFocused(id ID) bool {
	return c.focusedID == id
}

// setHot marks a widget as hot (hovered).
func (c *Context) setHot(id ID) {
	c.hotID = id
}

// setActive marks a widget as active (being interacted with).
func (c *Context) setActive(id ID) {
	c.activeID = id
}

// setFocused gives keyboard focus to a widget.
func (c *Context) setFocused(id ID) {
	c.focusedID = id
}

// clearActive clears the active widget.
func (c *Context) clearActive() {
	c.activeID = 0
}

// pointInRect returns true if the point (px, py) is inside the rectangle.
func pointInRect(px, py, x, y, w, h int) bool {
	return px >= x && px < x+w && py >= y && py < y+h
}

// Label renders static text at the given position.
func (c *Context) Label(text string, x, y int) {
	hlg.Text(text, x, y, 14, color.RGBA{220, 220, 220, 255})
}

// LabelWithSize renders text at the given position with a specific font size.
func (c *Context) LabelWithSize(text string, x, y int, size float32) {
	hlg.Text(text, x, y, size, color.RGBA{220, 220, 220, 255})
}

// LabelWithColor renders text at the given position with a specific color.
func (c *Context) LabelWithColor(text string, x, y int, col color.Color) {
	hlg.Text(text, x, y, 14, col)
}

// registerPanelBounds records a panel's bounds for input blocking.
func (c *Context) registerPanelBounds(id ID, x, y, w, h int) {
	c.panelBounds = append(c.panelBounds, panelBound{id: id, x: x, y: y, w: w, h: h})
}

// isBlockedByLaterPanel checks if a point would be blocked by a panel
// rendered after the given panel ID (higher z-order).
func (c *Context) isBlockedByLaterPanel(id ID, px, py int) bool {
	foundSelf := false
	for _, pb := range c.prevPanelBounds {
		if pb.id == id {
			foundSelf = true
			continue
		}
		if foundSelf && pointInRect(px, py, pb.x, pb.y, pb.w, pb.h) {
			return true
		}
	}
	return false
}

// isInputBlocked checks if input at the given point should be blocked.
// This is used by widgets to check if they're covered by a later panel.
func (c *Context) isInputBlocked(px, py int) bool {
	if c.currentPanelID == 0 {
		return false // Not inside a panel
	}
	return c.isBlockedByLaterPanel(c.currentPanelID, px, py)
}

// SetCurrentPanel sets the current panel context for widget input blocking.
func (c *Context) SetCurrentPanel(id ID) {
	c.currentPanelID = id
}

// ClearCurrentPanel clears the current panel context.
func (c *Context) ClearCurrentPanel() {
	c.currentPanelID = 0
}
