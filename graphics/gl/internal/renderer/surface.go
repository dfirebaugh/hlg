// Package renderer provides the core rendering components
package renderer

import "sync"

// Surface represents the rendering surface (shared between platforms)
type Surface struct {
	mu            sync.RWMutex
	logicalWidth  int
	logicalHeight int
	surfaceLocked bool // true if surface size was explicitly set
}

// NewSurface creates a new surface with the given dimensions
func NewSurface(width, height int) *Surface {
	return &Surface{
		logicalWidth:  width,
		logicalHeight: height,
		surfaceLocked: false,
	}
}

// GetSurfaceSize returns the logical surface size (coordinate system dimensions)
func (s *Surface) GetSurfaceSize() (int, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.logicalWidth, s.logicalHeight
}

// SetSurfaceSize sets the logical surface size and locks it
// Once locked, Resize() will not change the size
func (s *Surface) SetSurfaceSize(width, height int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logicalWidth = width
	s.logicalHeight = height
	s.surfaceLocked = true
}

// Resize updates the surface size if it's not locked
func (s *Surface) Resize(width, height int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.surfaceLocked {
		return
	}
	s.logicalWidth = width
	s.logicalHeight = height
}

// IsLocked returns whether the surface size is locked
func (s *Surface) IsLocked() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.surfaceLocked
}
