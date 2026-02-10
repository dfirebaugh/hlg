// Package pipelines provides rendering pipelines
package pipelines

import (
	"image"

	"github.com/dfirebaugh/hlg/graphics"
)

// MSDFAtlas represents an MSDF font atlas
type MSDFAtlas struct {
	glyphs        map[rune]*graphics.GlyphInfo
	metrics       graphics.FontMetrics
	distanceRange float64
	isDisposed    bool
}

func NewMSDFAtlas(atlasImg image.Image, distanceRange float64) (*MSDFAtlas, error) {
	return &MSDFAtlas{
		glyphs:        make(map[rune]*graphics.GlyphInfo),
		distanceRange: distanceRange,
	}, nil
}

func (a *MSDFAtlas) AddGlyph(r rune, info *graphics.GlyphInfo) {
	a.glyphs[r] = info
}

func (a *MSDFAtlas) GetGlyph(r rune) *graphics.GlyphInfo {
	return a.glyphs[r]
}

func (a *MSDFAtlas) SetMetrics(metrics graphics.FontMetrics) {
	a.metrics = metrics
}

func (a *MSDFAtlas) GetMetrics() graphics.FontMetrics {
	return a.metrics
}

func (a *MSDFAtlas) Dispose() {
	a.isDisposed = true
	a.glyphs = nil
}

func (a *MSDFAtlas) IsDisposed() bool {
	return a.isDisposed
}

// Ensure MSDFAtlas implements graphics.MSDFAtlas
var _ graphics.MSDFAtlas = (*MSDFAtlas)(nil)
