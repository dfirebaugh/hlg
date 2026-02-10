//go:build !js

package hlg

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"os"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/dfirebaugh/msdf/msdf"
)

//go:embed assets/fonts/Noto/NotoSansNerdFont-Regular.ttf
var defaultFontData []byte

// SDFType specifies the type of signed distance field to generate
type SDFType int

const (
	// TypeSDF generates a single-channel signed distance field
	TypeSDF SDFType = iota
	// TypePSDF generates a pseudo signed distance field
	TypePSDF
	// TypeMSDF generates a multi-channel signed distance field (3 channels)
	TypeMSDF
	// TypeMTSDF generates a multi-channel true signed distance field (4 channels)
	TypeMTSDF
)

// FontConfig contains configuration for font atlas generation
type FontConfig struct {
	// Size is the glyph size in pixels for the atlas
	Size int
	// PixelRange is the distance range in pixels for the SDF
	PixelRange float64
	// Charset specifies which characters to include in the atlas
	// If empty, ASCII characters are used by default
	Charset string
	// Type specifies the SDF type (SDF, MSDF, or MTSDF)
	// Defaults to MTSDF
	Type SDFType
}

// DefaultFontConfig returns the default font configuration
func DefaultFontConfig() FontConfig {
	return FontConfig{
		Size:       64,  // Matches banana-c default (64pt is sufficient for MSDF)
		PixelRange: 8.0, // Higher range for better AA (banana-c default, minimum should be 4.0)
		Charset:    "",  // Will use ASCII
		Type:       TypeMTSDF,
	}
}

// Font represents a loaded font with an MSDF atlas for text rendering
type Font struct {
	atlas    graphics.MSDFAtlas
	config   FontConfig
	fontData []byte // Store font data for re-generation

	// Atlas metadata for text layout
	atlasWidth  int
	atlasHeight int
	emSize      float64

	// Store the atlas image for debugging
	atlasImage image.Image

	// Cached space glyph for missing character fallback
	spaceGlyph *graphics.GlyphInfo
}

// LoadFont loads a font from a file path and generates an MSDF atlas
func LoadFont(fontPath string) (*Font, error) {
	return LoadFontWithConfig(fontPath, DefaultFontConfig())
}

// LoadFontWithConfig loads a font with custom configuration
func LoadFontWithConfig(fontPath string, config FontConfig) (*Font, error) {
	ensureSetupCompletion()

	// Read font file
	data, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read font file: %w", err)
	}

	return createFontFromData(data, config)
}

// LoadFontFromBytes loads a font from raw bytes
func LoadFontFromBytes(data []byte) (*Font, error) {
	return LoadFontFromBytesWithConfig(data, DefaultFontConfig())
}

// LoadFontFromBytesWithConfig loads a font from raw bytes with custom configuration
func LoadFontFromBytesWithConfig(data []byte, config FontConfig) (*Font, error) {
	ensureSetupCompletion()
	return createFontFromData(data, config)
}

// LoadDefaultFont loads the embedded default font (Noto Sans)
func LoadDefaultFont() (*Font, error) {
	return LoadFontFromBytes(defaultFontData)
}

// LoadDefaultFontWithConfig loads the embedded default font with custom configuration
func LoadDefaultFontWithConfig(config FontConfig) (*Font, error) {
	return LoadFontFromBytesWithConfig(defaultFontData, config)
}

// AtlasMetadata represents the JSON format for atlas metadata
type AtlasMetadata struct {
	Atlas struct {
		Type          string `json:"type"`
		DistanceRange int    `json:"distanceRange"`
		Size          int    `json:"size"`
		Width         int    `json:"width"`
		Height        int    `json:"height"`
		YOrigin       string `json:"yOrigin"`
	} `json:"atlas"`
	Metrics struct {
		EmSize     float64 `json:"emSize"`
		LineHeight float64 `json:"lineHeight"`
		Ascender   float64 `json:"ascender"`
		Descender  float64 `json:"descender"`
	} `json:"metrics"`
	Glyphs []struct {
		Unicode     int     `json:"unicode"`
		Advance     float64 `json:"advance"`
		PlaneBounds *struct {
			Left   float64 `json:"left"`
			Bottom float64 `json:"bottom"`
			Right  float64 `json:"right"`
			Top    float64 `json:"top"`
		} `json:"planeBounds,omitempty"`
		AtlasBounds *struct {
			Left   float64 `json:"left"`
			Bottom float64 `json:"bottom"`
			Right  float64 `json:"right"`
			Top    float64 `json:"top"`
		} `json:"atlasBounds,omitempty"`
	} `json:"glyphs"`
}

// LoadFontFromAtlas loads a font from pregenerated atlas PNG and JSON metadata files.
// This allows using a pregenerated atlas instead of generating one at runtime.
func LoadFontFromAtlas(atlasPNGPath, atlasJSONPath string) (*Font, error) {
	ensureSetupCompletion()

	// Load the atlas image
	imgFile, err := os.Open(atlasPNGPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open atlas image: %w", err)
	}
	defer imgFile.Close()

	atlasImg, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode atlas image: %w", err)
	}

	// Load the JSON metadata
	jsonFile, err := os.Open(atlasJSONPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open atlas JSON: %w", err)
	}
	defer jsonFile.Close()

	var meta AtlasMetadata
	if err := json.NewDecoder(jsonFile).Decode(&meta); err != nil {
		return nil, fmt.Errorf("failed to decode atlas JSON: %w", err)
	}

	// Create GPU atlas through FontManager interface
	fontManager := getFontManager()
	atlas, err := fontManager.CreateMSDFAtlas(atlasImg, float64(meta.Atlas.DistanceRange))
	if err != nil {
		return nil, fmt.Errorf("failed to create MSDF atlas: %w", err)
	}

	// Set font metrics from JSON
	atlas.SetMetrics(graphics.FontMetrics{
		EmSize:     meta.Metrics.EmSize,
		LineHeight: meta.Metrics.LineHeight,
		Ascender:   meta.Metrics.Ascender,
		Descender:  meta.Metrics.Descender,
	})

	// Add glyph information from JSON
	atlasWidth := meta.Atlas.Width
	atlasHeight := meta.Atlas.Height

	for _, g := range meta.Glyphs {
		var info *graphics.GlyphInfo

		if g.PlaneBounds != nil && g.AtlasBounds != nil {
			// Calculate UV coordinates from atlas bounds
			// The JSON stores atlas bounds with Y-origin at bottom, need to convert to UV coords
			s0 := g.AtlasBounds.Left / float64(atlasWidth)
			t0 := 1.0 - g.AtlasBounds.Top/float64(atlasHeight)
			s1 := g.AtlasBounds.Right / float64(atlasWidth)
			t1 := 1.0 - g.AtlasBounds.Bottom/float64(atlasHeight)

			info = &graphics.GlyphInfo{
				Unicode: g.Unicode,
				Quad: graphics.GlyphQuad{
					S0:      s0,
					T0:      t0,
					S1:      s1,
					T1:      t1,
					PL:      g.PlaneBounds.Left,
					PB:      g.PlaneBounds.Bottom,
					PR:      g.PlaneBounds.Right,
					PT:      g.PlaneBounds.Top,
					Advance: g.Advance,
				},
			}
		} else {
			// Space or empty glyph - just has advance
			info = &graphics.GlyphInfo{
				Unicode: g.Unicode,
				Quad: graphics.GlyphQuad{
					Advance: g.Advance,
				},
			}
		}

		atlas.AddGlyph(rune(g.Unicode), info)
	}

	// Determine pixel range from atlas type
	pixelRange := float64(meta.Atlas.DistanceRange)

	f := &Font{
		atlas:       atlas,
		fontData:    nil, // No font data when loading from pregenerated atlas
		config:      FontConfig{Size: meta.Atlas.Size, PixelRange: pixelRange},
		atlasWidth:  atlasWidth,
		atlasHeight: atlasHeight,
		emSize:      meta.Metrics.EmSize,
		atlasImage:  atlasImg,
	}

	// Cache the space glyph for missing character fallback
	f.spaceGlyph = f.atlas.GetGlyph(' ')

	return f, nil
}

// getFontManager returns the FontManager from the graphics backend
func getFontManager() graphics.FontManager {
	return hlg.graphicsBackend.(graphics.FontManager)
}

func createFontFromData(fontData []byte, config FontConfig) (*Font, error) {
	// Initialize FreeType
	ft, err := msdf.NewFreetype()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize FreeType: %w", err)
	}
	defer ft.Close()

	// Load font from memory
	msdfFont, err := ft.LoadFontMemory(fontData)
	if err != nil {
		return nil, fmt.Errorf("failed to load font: %w", err)
	}
	defer msdfFont.Close()

	// Create font geometry
	geom, err := msdf.NewFontGeometry()
	if err != nil {
		return nil, fmt.Errorf("failed to create font geometry: %w", err)
	}
	defer geom.Close()

	// Load glyphs
	var loaded int
	if config.Charset == "" {
		// Load ASCII by default
		loaded, err = geom.LoadASCII(msdfFont, 1.0)
	} else {
		loaded, err = geom.LoadCharset(msdfFont, 1.0, config.Charset)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load glyphs: %w", err)
	}
	if loaded == 0 {
		return nil, fmt.Errorf("no glyphs loaded")
	}

	// Configure atlas generation
	// Settings based on banana-c library for improved text quality
	atlasConfig := msdf.AtlasConfig{
		ImageType:            toMsdfImageType(config.Type),
		ImageFormat:          msdf.ImageFormatPNG,
		YDirection:           msdf.YDownward,
		EmSize:               float64(config.Size),
		PxRange:              config.PixelRange,
		AngleThreshold:       3.0,
		MiterLimit:           2.0, // Higher miter limit for sharper corners (banana-c uses 2.0)
		ThreadCount:          0,   // Auto
		DimensionsConstraint: msdf.DimensionsMultipleOfFourSquare,
	}

	// Generate atlas
	result, err := geom.GenerateAtlas(atlasConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate atlas: %w", err)
	}

	// Save to temporary files and load via LoadFontFromAtlas
	// This ensures correct coordinate handling that matches the wrapper's export format
	tmpPNG := os.TempDir() + "/hlg_atlas_" + fmt.Sprintf("%d", os.Getpid()) + ".png"
	tmpJSON := os.TempDir() + "/hlg_atlas_" + fmt.Sprintf("%d", os.Getpid()) + ".json"

	if err := msdf.SaveImage(result, tmpPNG, msdf.ImageFormatPNG); err != nil {
		return nil, fmt.Errorf("failed to save temp atlas image: %w", err)
	}
	defer os.Remove(tmpPNG)

	if err := geom.ExportJSON(result, tmpJSON, false); err != nil {
		os.Remove(tmpPNG)
		return nil, fmt.Errorf("failed to save temp atlas JSON: %w", err)
	}
	defer os.Remove(tmpJSON)

	// Load the atlas using the standard loader
	font, err := LoadFontFromAtlas(tmpPNG, tmpJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to load generated atlas: %w", err)
	}

	// Store the font data for potential re-generation
	font.fontData = fontData
	font.config = config

	return font, nil
}

func toMsdfImageType(t SDFType) msdf.ImageType {
	switch t {
	case TypeSDF:
		return msdf.ImageTypeSDF
	case TypePSDF:
		return msdf.ImageTypePSDF
	case TypeMSDF:
		return msdf.ImageTypeMSDF
	case TypeMTSDF:
		return msdf.ImageTypeMTSDF
	default:
		return msdf.ImageTypeMTSDF
	}
}

// GetAtlasImage returns the atlas as an image (useful for debugging)
func (f *Font) GetAtlasImage() image.Image {
	return f.atlasImage
}

// Dispose releases the font resources
func (f *Font) Dispose() {
	if f.atlas != nil {
		f.atlas.Dispose()
	}
}

// GetMetrics returns the font metrics
func (f *Font) GetMetrics() graphics.FontMetrics {
	return f.atlas.GetMetrics()
}

// DebugPrintMetrics prints font metrics and glyph information for debugging
func (f *Font) DebugPrintMetrics() {
	metrics := f.atlas.GetMetrics()
	fmt.Printf("=== Font Metrics ===\n")
	fmt.Printf("EmSize: %.2f\n", metrics.EmSize)
	fmt.Printf("LineHeight: %.2f\n", metrics.LineHeight)
	fmt.Printf("Ascender: %.2f\n", metrics.Ascender)
	fmt.Printf("Descender: %.2f\n", metrics.Descender)

	// Print a few sample glyphs
	for _, ch := range "HAe" {
		g := f.atlas.GetGlyph(ch)
		if g != nil {
			fmt.Printf("\nGlyph '%c' (U+%04X):\n", ch, ch)
			fmt.Printf("  Advance: %.3f\n", g.Quad.Advance)
			fmt.Printf("  Plane bounds (em): L=%.4f B=%.4f R=%.4f T=%.4f\n",
				g.Quad.PL, g.Quad.PB, g.Quad.PR, g.Quad.PT)
			fmt.Printf("  Tex coords: S0=%.4f T0=%.4f S1=%.4f T1=%.4f\n",
				g.Quad.S0, g.Quad.T0, g.Quad.S1, g.Quad.T1)
		}
	}
	fmt.Println("")
}

// SaveAtlasDebug saves the atlas to a PNG file for debugging purposes
func (f *Font) SaveAtlasDebug(filename string) error {
	if f.fontData == nil {
		return fmt.Errorf("no font data available for re-generation")
	}

	// Re-generate the atlas for saving
	ft, err := msdf.NewFreetype()
	if err != nil {
		return fmt.Errorf("failed to initialize FreeType: %w", err)
	}
	defer ft.Close()

	msdfFont, err := ft.LoadFontMemory(f.fontData)
	if err != nil {
		return fmt.Errorf("failed to load font: %w", err)
	}
	defer msdfFont.Close()

	geom, err := msdf.NewFontGeometry()
	if err != nil {
		return fmt.Errorf("failed to create font geometry: %w", err)
	}
	defer geom.Close()

	if f.config.Charset == "" {
		_, err = geom.LoadASCII(msdfFont, 1.0)
	} else {
		_, err = geom.LoadCharset(msdfFont, 1.0, f.config.Charset)
	}
	if err != nil {
		return fmt.Errorf("failed to load glyphs: %w", err)
	}

	atlasConfig := msdf.AtlasConfig{
		ImageType:            toMsdfImageType(f.config.Type),
		ImageFormat:          msdf.ImageFormatPNG,
		YDirection:           msdf.YDownward,
		EmSize:               float64(f.config.Size),
		PxRange:              f.config.PixelRange,
		AngleThreshold:       3.0,
		MiterLimit:           2.0, // Higher miter limit for sharper corners
		ThreadCount:          0,
		DimensionsConstraint: msdf.DimensionsMultipleOfFourSquare,
	}

	result, err := geom.GenerateAtlas(atlasConfig)
	if err != nil {
		return fmt.Errorf("failed to generate atlas: %w", err)
	}

	return msdf.SaveImage(result, filename, msdf.ImageFormatPNG)
}

// DrawText draws text to a GlyphDrawer (such as gui.DrawContext), allowing it to be layered with other primitives.
// This method requires that SetMSDFAtlas has been called with this font's atlas first.
func (f *Font) DrawText(dc graphics.GlyphDrawer, text string, x, y, fontSize float32, c color.Color) {
	metrics := f.atlas.GetMetrics()
	advanceScale := float32(float64(fontSize) / metrics.EmSize)

	cursorX := x
	// Position Y is treated as the top of the text area (cap height), so we add
	// an approximation of cap height to get the baseline position.
	// Cap height is typically ~70-75% of ascender for most fonts.
	capHeightScale := float32(metrics.Ascender / metrics.EmSize * 0.68)
	cursorY := y + capHeightScale*fontSize

	for _, ch := range text {
		if ch == '\n' {
			cursorX = x
			cursorY += fontSize * 1.2
			continue
		}

		glyph := f.atlas.GetGlyph(ch)
		if glyph == nil {
			// Use cached space glyph advance for missing characters
			if f.spaceGlyph != nil {
				cursorX += float32(f.spaceGlyph.Quad.Advance) * advanceScale
			}
			continue
		}

		q := glyph.Quad

		// Calculate screen positions for this glyph
		// Pixel-perfect rounding: only round top-left corner, add unrounded dimensions
		// This preserves exact glyph size while snapping to pixel grid (banana-c approach)
		glyphX := float32(math.Floor(float64(cursorX+float32(q.PL)*fontSize) + 0.5))
		glyphY := float32(math.Floor(float64(cursorY-float32(q.PT)*fontSize) + 0.5))
		glyphW := float32(q.PR-q.PL) * fontSize
		glyphH := float32(q.PT-q.PB) * fontSize

		x0 := glyphX
		y0 := glyphY
		x1 := glyphX + glyphW
		y1 := glyphY + glyphH

		// UV coordinates from atlas
		u0 := float32(q.S0)
		v0 := float32(q.T0)
		u1 := float32(q.S1)
		v1 := float32(q.T1)

		// Draw the glyph quad
		dc.DrawGlyph(x0, y0, x1, y1, u0, v0, u1, v1, c)

		// Advance cursor for next glyph
		cursorX += float32(q.Advance) * advanceScale
	}
}

// RenderText renders text and returns vertices directly without needing a DrawContext.
// screenWidth and screenHeight are needed for NDC conversion.
// DEPRECATED: Use RenderTextPrimitives for ~5x memory reduction.
func (f *Font) RenderText(text string, x, y, fontSize float32, c color.Color, screenWidth, screenHeight int) []graphics.PrimitiveVertex {
	metrics := f.atlas.GetMetrics()
	advanceScale := float32(float64(fontSize) / metrics.EmSize)

	cursorX := x
	// Use cap height approximation for more intuitive Y positioning
	capHeightScale := float32(metrics.Ascender / metrics.EmSize * 0.68)
	cursorY := y + capHeightScale*fontSize

	// Pre-allocate vertices (6 vertices per glyph for 2 triangles)
	vertices := make([]graphics.PrimitiveVertex, 0, len(text)*6)

	sw := float32(screenWidth)
	sh := float32(screenHeight)

	// Convert color once
	r, g, b, a := c.RGBA()
	colorVec := [4]float32{
		float32(r) / 0xffff,
		float32(g) / 0xffff,
		float32(b) / 0xffff,
		float32(a) / 0xffff,
	}

	for _, ch := range text {
		if ch == '\n' {
			cursorX = x
			cursorY += fontSize * 1.2
			continue
		}

		glyph := f.atlas.GetGlyph(ch)
		if glyph == nil {
			if f.spaceGlyph != nil {
				cursorX += float32(f.spaceGlyph.Quad.Advance) * advanceScale
			}
			continue
		}

		q := glyph.Quad

		// Pixel-perfect rounding: only round top-left corner, add unrounded dimensions
		glyphX := float32(math.Floor(float64(cursorX+float32(q.PL)*fontSize) + 0.5))
		glyphY := float32(math.Floor(float64(cursorY-float32(q.PT)*fontSize) + 0.5))
		glyphW := float32(q.PR-q.PL) * fontSize
		glyphH := float32(q.PT-q.PB) * fontSize

		x0 := glyphX
		y0 := glyphY
		x1 := glyphX + glyphW
		y1 := glyphY + glyphH

		// UV coordinates
		u0 := float32(q.S0)
		v0 := float32(q.T0)
		u1 := float32(q.S1)
		v1 := float32(q.T1)

		// Convert to NDC
		ndcX0 := (x0/sw)*2.0 - 1.0
		ndcY0 := 1.0 - (y0/sh)*2.0
		ndcX1 := (x1/sw)*2.0 - 1.0
		ndcY1 := 1.0 - (y1/sh)*2.0

		// Create 6 vertices for 2 triangles
		vertices = append(vertices,
			graphics.PrimitiveVertex{Position: [3]float32{ndcX0, ndcY0, 0.0}, LocalPosition: [2]float32{0, 0}, OpCode: graphics.OpCodeMSDF, Radius: 0.0, Color: colorVec, TexCoords: [2]float32{u0, v0}},
			graphics.PrimitiveVertex{Position: [3]float32{ndcX0, ndcY1, 0.0}, LocalPosition: [2]float32{0, 0}, OpCode: graphics.OpCodeMSDF, Radius: 0.0, Color: colorVec, TexCoords: [2]float32{u0, v1}},
			graphics.PrimitiveVertex{Position: [3]float32{ndcX1, ndcY0, 0.0}, LocalPosition: [2]float32{0, 0}, OpCode: graphics.OpCodeMSDF, Radius: 0.0, Color: colorVec, TexCoords: [2]float32{u1, v0}},
			graphics.PrimitiveVertex{Position: [3]float32{ndcX0, ndcY1, 0.0}, LocalPosition: [2]float32{0, 0}, OpCode: graphics.OpCodeMSDF, Radius: 0.0, Color: colorVec, TexCoords: [2]float32{u0, v1}},
			graphics.PrimitiveVertex{Position: [3]float32{ndcX1, ndcY1, 0.0}, LocalPosition: [2]float32{0, 0}, OpCode: graphics.OpCodeMSDF, Radius: 0.0, Color: colorVec, TexCoords: [2]float32{u1, v1}},
			graphics.PrimitiveVertex{Position: [3]float32{ndcX1, ndcY0, 0.0}, LocalPosition: [2]float32{0, 0}, OpCode: graphics.OpCodeMSDF, Radius: 0.0, Color: colorVec, TexCoords: [2]float32{u1, v0}},
		)

		cursorX += float32(q.Advance) * advanceScale
	}

	return vertices
}

// RenderTextPrimitives renders text and returns Primitives directly (one per glyph).
// This is the efficient storage buffer approach using 64 bytes per glyph instead of 312 bytes.
func (f *Font) RenderTextPrimitives(text string, x, y, fontSize float32, c color.Color) []graphics.Primitive {
	metrics := f.atlas.GetMetrics()
	advanceScale := float32(float64(fontSize) / metrics.EmSize)

	cursorX := x
	// Use cap height approximation for more intuitive Y positioning
	capHeightScale := float32(metrics.Ascender / metrics.EmSize * 0.68)
	cursorY := y + capHeightScale*fontSize

	// Pre-allocate primitives (1 per glyph)
	primitives := make([]graphics.Primitive, 0, len(text))

	// Convert color once
	r, g, b, a := c.RGBA()
	colorVec := [4]float32{
		float32(r) / 0xffff,
		float32(g) / 0xffff,
		float32(b) / 0xffff,
		float32(a) / 0xffff,
	}

	for _, ch := range text {
		if ch == '\n' {
			cursorX = x
			cursorY += fontSize * 1.2
			continue
		}

		glyph := f.atlas.GetGlyph(ch)
		if glyph == nil {
			if f.spaceGlyph != nil {
				cursorX += float32(f.spaceGlyph.Quad.Advance) * advanceScale
			}
			continue
		}

		q := glyph.Quad

		// Pixel-perfect rounding: only round top-left corner, add unrounded dimensions
		// This preserves exact glyph size while snapping to pixel grid (banana-c approach)
		// Pixel-perfect rounding: only round top-left corner, add unrounded dimensions
		glyphX := float32(math.Floor(float64(cursorX+float32(q.PL)*fontSize) + 0.5))
		glyphY := float32(math.Floor(float64(cursorY-float32(q.PT)*fontSize) + 0.5))
		glyphW := float32(q.PR-q.PL) * fontSize
		glyphH := float32(q.PT-q.PB) * fontSize

		// UV coordinates
		u0 := float32(q.S0)
		v0 := float32(q.T0)
		u1 := float32(q.S1)
		v1 := float32(q.T1)

		// Skip glyphs with invalid sizes
		if glyphW <= 0 || glyphH <= 0 {
			cursorX += float32(q.Advance) * advanceScale
			continue
		}

		primitives = append(primitives, graphics.Primitive{
			X:      glyphX,
			Y:      glyphY,
			W:      glyphW,
			H:      glyphH,
			Color:  colorVec,
			Radius: 0,
			OpCode: graphics.OpCodeMSDF,
			Extra:  [4]float32{u0, v0, u1 - u0, v1 - v0}, // UV: base + size
		})

		cursorX += float32(q.Advance) * advanceScale
	}

	return primitives
}

// SetAsActiveAtlas sets this font's atlas as the active MSDF atlas for the primitive buffer.
// This must be called before using DrawText.
func (f *Font) SetAsActiveAtlas() {
	SetMSDFAtlas(f.atlasImage, f.config.PixelRange)
}

// MeasureText returns the width of the text in pixels at the given font size.
// For multiline text, returns the width of the longest line.
func (f *Font) MeasureText(text string, fontSize float32) float32 {
	metrics := f.atlas.GetMetrics()
	advanceScale := float32(float64(fontSize) / metrics.EmSize)

	var maxWidth float32
	var currentWidth float32
	for _, ch := range text {
		if ch == '\n' {
			if currentWidth > maxWidth {
				maxWidth = currentWidth
			}
			currentWidth = 0
			continue
		}

		glyph := f.atlas.GetGlyph(ch)
		if glyph == nil {
			// Use cached space glyph advance for missing characters
			if f.spaceGlyph != nil {
				currentWidth += float32(f.spaceGlyph.Quad.Advance) * advanceScale
			}
			continue
		}

		currentWidth += float32(glyph.Quad.Advance) * advanceScale
	}

	// Check final line
	if currentWidth > maxWidth {
		maxWidth = currentWidth
	}

	return maxWidth
}

// SaveAtlasJSON saves the atlas metadata to a JSON file for debugging/interop
func (f *Font) SaveAtlasJSON(filename string) error {
	if f.fontData == nil {
		return fmt.Errorf("no font data available for re-generation")
	}

	// Re-generate the atlas for saving
	ft, err := msdf.NewFreetype()
	if err != nil {
		return fmt.Errorf("failed to initialize FreeType: %w", err)
	}
	defer ft.Close()

	msdfFont, err := ft.LoadFontMemory(f.fontData)
	if err != nil {
		return fmt.Errorf("failed to load font: %w", err)
	}
	defer msdfFont.Close()

	geom, err := msdf.NewFontGeometry()
	if err != nil {
		return fmt.Errorf("failed to create font geometry: %w", err)
	}
	defer geom.Close()

	if f.config.Charset == "" {
		_, err = geom.LoadASCII(msdfFont, 1.0)
	} else {
		_, err = geom.LoadCharset(msdfFont, 1.0, f.config.Charset)
	}
	if err != nil {
		return fmt.Errorf("failed to load glyphs: %w", err)
	}

	atlasConfig := msdf.AtlasConfig{
		ImageType:            toMsdfImageType(f.config.Type),
		ImageFormat:          msdf.ImageFormatPNG,
		YDirection:           msdf.YDownward,
		EmSize:               float64(f.config.Size),
		PxRange:              f.config.PixelRange,
		AngleThreshold:       3.0,
		MiterLimit:           2.0, // Higher miter limit for sharper corners
		ThreadCount:          0,
		DimensionsConstraint: msdf.DimensionsMultipleOfFourSquare,
	}

	result, err := geom.GenerateAtlas(atlasConfig)
	if err != nil {
		return fmt.Errorf("failed to generate atlas: %w", err)
	}

	return geom.ExportJSON(result, filename, true)
}
