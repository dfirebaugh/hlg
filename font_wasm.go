//go:build js && wasm

package hlg

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io"

	"github.com/dfirebaugh/hlg/graphics"
)

//go:embed assets/fonts/Noto/noto_atlas.png
var defaultAtlasPNG []byte

//go:embed assets/fonts/Noto/noto_atlas.json
var defaultAtlasJSON []byte

// SDFType specifies the type of signed distance field to generate
type SDFType int

const (
	TypeSDF SDFType = iota
	TypePSDF
	TypeMSDF
	TypeMTSDF
)

// FontConfig contains configuration for font atlas generation
type FontConfig struct {
	Size       int
	PixelRange float64
	Charset    string
	Type       SDFType
}

// DefaultFontConfig returns the default font configuration
func DefaultFontConfig() FontConfig {
	return FontConfig{
		Size:       128,
		PixelRange: 3.0,
		Charset:    "",
		Type:       TypeMTSDF,
	}
}

// Font represents a loaded font with an MSDF atlas for text rendering
type Font struct {
	atlas    graphics.MSDFAtlas
	config   FontConfig
	fontData []byte

	atlasWidth  int
	atlasHeight int
	emSize      float64
	atlasImage  image.Image
	spaceGlyph  *graphics.GlyphInfo
}

// LoadFont is not supported in WASM builds - use LoadFontFromAtlasBytes instead
func LoadFont(fontPath string) (*Font, error) {
	return nil, errors.New("LoadFont is not supported in WASM builds - use LoadFontFromAtlasBytes with embedded assets")
}

// LoadFontWithConfig is not supported in WASM builds
func LoadFontWithConfig(fontPath string, config FontConfig) (*Font, error) {
	return nil, errors.New("LoadFontWithConfig is not supported in WASM builds - use LoadFontFromAtlasBytes")
}

// LoadFontFromBytes is not supported in WASM builds - runtime font generation requires native libraries
func LoadFontFromBytes(data []byte) (*Font, error) {
	return nil, errors.New("LoadFontFromBytes is not supported in WASM builds - use LoadFontFromAtlasBytes with pregenerated atlas")
}

// LoadFontFromBytesWithConfig is not supported in WASM builds
func LoadFontFromBytesWithConfig(data []byte, config FontConfig) (*Font, error) {
	return nil, errors.New("LoadFontFromBytesWithConfig is not supported in WASM builds - use LoadFontFromAtlasBytes")
}

// LoadDefaultFont loads the default font using the embedded pregenerated atlas
func LoadDefaultFont() (*Font, error) {
	return LoadFontFromAtlasBytes(defaultAtlasPNG, defaultAtlasJSON)
}

// LoadDefaultFontWithConfig is not supported in WASM builds
func LoadDefaultFontWithConfig(config FontConfig) (*Font, error) {
	return nil, errors.New("LoadDefaultFontWithConfig is not supported in WASM builds")
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

// LoadFontFromAtlas is not supported in WASM - use LoadFontFromAtlasBytes instead
func LoadFontFromAtlas(atlasPNGPath, atlasJSONPath string) (*Font, error) {
	return nil, errors.New("LoadFontFromAtlas is not supported in WASM builds - use LoadFontFromAtlasBytes")
}

// LoadFontFromAtlasBytes loads a font from pregenerated atlas PNG and JSON metadata bytes.
// This is the primary way to load fonts in WASM builds.
func LoadFontFromAtlasBytes(atlasPNG, atlasJSON []byte) (*Font, error) {
	ensureSetupCompletion()

	atlasImg, _, err := image.Decode(bytes.NewReader(atlasPNG))
	if err != nil {
		return nil, fmt.Errorf("failed to decode atlas image: %w", err)
	}

	var meta AtlasMetadata
	if err := json.NewDecoder(bytes.NewReader(atlasJSON)).Decode(&meta); err != nil {
		return nil, fmt.Errorf("failed to decode atlas JSON: %w", err)
	}

	fontManager := getFontManager()
	atlas, err := fontManager.CreateMSDFAtlas(atlasImg, float64(meta.Atlas.DistanceRange))
	if err != nil {
		return nil, fmt.Errorf("failed to create MSDF atlas: %w", err)
	}

	atlas.SetMetrics(graphics.FontMetrics{
		EmSize:     meta.Metrics.EmSize,
		LineHeight: meta.Metrics.LineHeight,
		Ascender:   meta.Metrics.Ascender,
		Descender:  meta.Metrics.Descender,
	})

	atlasWidth := meta.Atlas.Width
	atlasHeight := meta.Atlas.Height

	for _, g := range meta.Glyphs {
		var info *graphics.GlyphInfo

		if g.PlaneBounds != nil && g.AtlasBounds != nil {
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
			info = &graphics.GlyphInfo{
				Unicode: g.Unicode,
				Quad: graphics.GlyphQuad{
					Advance: g.Advance,
				},
			}
		}

		atlas.AddGlyph(rune(g.Unicode), info)
	}

	pixelRange := float64(meta.Atlas.DistanceRange)

	f := &Font{
		atlas:       atlas,
		fontData:    nil,
		config:      FontConfig{Size: meta.Atlas.Size, PixelRange: pixelRange},
		atlasWidth:  atlasWidth,
		atlasHeight: atlasHeight,
		emSize:      meta.Metrics.EmSize,
		atlasImage:  atlasImg,
	}

	f.spaceGlyph = f.atlas.GetGlyph(' ')

	return f, nil
}

// LoadFontFromAtlasReaders loads a font from io.Reader sources
func LoadFontFromAtlasReaders(atlasPNGReader, atlasJSONReader io.Reader) (*Font, error) {
	pngData, err := io.ReadAll(atlasPNGReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read atlas PNG: %w", err)
	}

	jsonData, err := io.ReadAll(atlasJSONReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read atlas JSON: %w", err)
	}

	return LoadFontFromAtlasBytes(pngData, jsonData)
}

func getFontManager() graphics.FontManager {
	return hlg.graphicsBackend.(graphics.FontManager)
}

func (f *Font) GetAtlasImage() image.Image {
	return f.atlasImage
}

func (f *Font) Dispose() {
	if f.atlas != nil {
		f.atlas.Dispose()
	}
}

func (f *Font) GetMetrics() graphics.FontMetrics {
	return f.atlas.GetMetrics()
}

func (f *Font) DebugPrintMetrics() {
	metrics := f.atlas.GetMetrics()
	fmt.Printf("=== Font Metrics ===\n")
	fmt.Printf("EmSize: %.2f\n", metrics.EmSize)
	fmt.Printf("LineHeight: %.2f\n", metrics.LineHeight)
	fmt.Printf("Ascender: %.2f\n", metrics.Ascender)
	fmt.Printf("Descender: %.2f\n", metrics.Descender)

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

func (f *Font) SaveAtlasDebug(filename string) error {
	return errors.New("SaveAtlasDebug is not supported in WASM builds")
}

func (f *Font) DrawText(dc graphics.GlyphDrawer, text string, x, y, fontSize float32, c color.Color) {
	metrics := f.atlas.GetMetrics()
	advanceScale := float32(float64(fontSize) / metrics.EmSize)

	cursorX := x
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
			if f.spaceGlyph != nil {
				cursorX += float32(f.spaceGlyph.Quad.Advance) * advanceScale
			}
			continue
		}

		q := glyph.Quad

		x0 := cursorX + float32(q.PL)*fontSize
		y0 := cursorY - float32(q.PT)*fontSize
		x1 := cursorX + float32(q.PR)*fontSize
		y1 := cursorY - float32(q.PB)*fontSize

		u0 := float32(q.S0)
		v0 := float32(q.T0)
		u1 := float32(q.S1)
		v1 := float32(q.T1)

		dc.DrawGlyph(x0, y0, x1, y1, u0, v0, u1, v1, c)

		cursorX += float32(q.Advance) * advanceScale
	}
}

func (f *Font) RenderText(text string, x, y, fontSize float32, c color.Color, screenWidth, screenHeight int) []graphics.PrimitiveVertex {
	metrics := f.atlas.GetMetrics()
	advanceScale := float32(float64(fontSize) / metrics.EmSize)

	cursorX := x
	capHeightScale := float32(metrics.Ascender / metrics.EmSize * 0.68)
	cursorY := y + capHeightScale*fontSize

	vertices := make([]graphics.PrimitiveVertex, 0, len(text)*6)

	sw := float32(screenWidth)
	sh := float32(screenHeight)

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

		x0 := cursorX + float32(q.PL)*fontSize
		y0 := cursorY - float32(q.PT)*fontSize
		x1 := cursorX + float32(q.PR)*fontSize
		y1 := cursorY - float32(q.PB)*fontSize

		u0 := float32(q.S0)
		v0 := float32(q.T0)
		u1 := float32(q.S1)
		v1 := float32(q.T1)

		ndcX0 := (x0/sw)*2.0 - 1.0
		ndcY0 := 1.0 - (y0/sh)*2.0
		ndcX1 := (x1/sw)*2.0 - 1.0
		ndcY1 := 1.0 - (y1/sh)*2.0

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

func (f *Font) RenderTextPrimitives(text string, x, y, fontSize float32, c color.Color) []graphics.Primitive {
	metrics := f.atlas.GetMetrics()
	advanceScale := float32(float64(fontSize) / metrics.EmSize)

	cursorX := x
	capHeightScale := float32(metrics.Ascender / metrics.EmSize * 0.68)
	cursorY := y + capHeightScale*fontSize

	primitives := make([]graphics.Primitive, 0, len(text))

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

		x0 := cursorX + float32(q.PL)*fontSize
		y0 := cursorY - float32(q.PT)*fontSize
		x1 := cursorX + float32(q.PR)*fontSize
		y1 := cursorY - float32(q.PB)*fontSize

		u0 := float32(q.S0)
		v0 := float32(q.T0)
		u1 := float32(q.S1)
		v1 := float32(q.T1)

		w := x1 - x0
		h := y1 - y0
		if w <= 0 || h <= 0 {
			cursorX += float32(q.Advance) * advanceScale
			continue
		}

		primitives = append(primitives, graphics.Primitive{
			X:      x0,
			Y:      y0,
			W:      w,
			H:      h,
			Color:  colorVec,
			Radius: 0,
			OpCode: graphics.OpCodeMSDF,
			Extra:  [4]float32{u0, v0, u1 - u0, v1 - v0},
		})

		cursorX += float32(q.Advance) * advanceScale
	}

	return primitives
}

func (f *Font) SetAsActiveAtlas() {
	SetMSDFAtlas(f.atlasImage, f.config.PixelRange)
}

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
			if f.spaceGlyph != nil {
				currentWidth += float32(f.spaceGlyph.Quad.Advance) * advanceScale
			}
			continue
		}

		currentWidth += float32(glyph.Quad.Advance) * advanceScale
	}

	if currentWidth > maxWidth {
		maxWidth = currentWidth
	}

	return maxWidth
}

func (f *Font) SaveAtlasJSON(filename string) error {
	return errors.New("SaveAtlasJSON is not supported in WASM builds")
}
