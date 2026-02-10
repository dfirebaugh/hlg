//go:build !js

package pipelines

import (
	"image"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/hlg/graphics/webgpu/internal/context"
)

// MSDFAtlas represents an MSDF font atlas for text rendering.
// The atlas stores glyph metadata and metrics; actual rendering is done
// via the primitive buffer using DrawText.
type MSDFAtlas struct {
	context.RenderContext

	// Atlas dimensions
	width, height int

	// Distance range used when generating the atlas
	distanceRange float64

	// Font metrics
	metrics graphics.FontMetrics

	// Glyph lookup
	glyphs map[rune]*graphics.GlyphInfo

	isDisposed bool
}

// NewMSDFAtlas creates a new MSDF atlas from glyph data.
// The actual texture is managed by the primitive buffer pipeline via SetMSDFAtlas.
func NewMSDFAtlas(ctx context.RenderContext, atlasImg image.Image, distanceRange float64) (*MSDFAtlas, error) {
	r := atlasImg.Bounds()
	width := r.Dx()
	height := r.Dy()

	atlas := &MSDFAtlas{
		RenderContext: ctx,
		width:         width,
		height:        height,
		distanceRange: distanceRange,
		glyphs:        make(map[rune]*graphics.GlyphInfo),
	}

	return atlas, nil
}

// AddGlyph adds glyph information to the atlas
func (a *MSDFAtlas) AddGlyph(r rune, info *graphics.GlyphInfo) {
	a.glyphs[r] = info
}

// GetGlyph returns the glyph info for a rune, or nil if not found
func (a *MSDFAtlas) GetGlyph(r rune) *graphics.GlyphInfo {
	return a.glyphs[r]
}

// SetMetrics sets the font metrics
func (a *MSDFAtlas) SetMetrics(metrics graphics.FontMetrics) {
	a.metrics = metrics
}

// GetMetrics returns the font metrics
func (a *MSDFAtlas) GetMetrics() graphics.FontMetrics {
	return a.metrics
}

// Dispose marks the atlas as disposed
func (a *MSDFAtlas) Dispose() {
	a.isDisposed = true
}

// IsDisposed returns true if the atlas has been disposed
func (a *MSDFAtlas) IsDisposed() bool {
	return a.isDisposed
}
