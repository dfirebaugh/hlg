// Package glapi provides a unified GL abstraction layer for OpenGL and WebGL
package glapi

// GL constants shared between OpenGL and WebGL
const (
	// Buffer targets
	ARRAY_BUFFER         uint32 = 0x8892
	ELEMENT_ARRAY_BUFFER uint32 = 0x8893

	// Buffer usage
	STATIC_DRAW  uint32 = 0x88E4
	DYNAMIC_DRAW uint32 = 0x88E8

	// Data types
	BYTE           uint32 = 0x1400
	UNSIGNED_BYTE  uint32 = 0x1401
	SHORT          uint32 = 0x1402
	UNSIGNED_SHORT uint32 = 0x1403
	INT            uint32 = 0x1404
	UNSIGNED_INT   uint32 = 0x1405
	FLOAT          uint32 = 0x1406

	// Primitive types
	POINTS         uint32 = 0x0000
	LINES          uint32 = 0x0001
	LINE_LOOP      uint32 = 0x0002
	LINE_STRIP     uint32 = 0x0003
	TRIANGLES      uint32 = 0x0004
	TRIANGLE_STRIP uint32 = 0x0005
	TRIANGLE_FAN   uint32 = 0x0006

	// Shader types
	FRAGMENT_SHADER uint32 = 0x8B30
	VERTEX_SHADER   uint32 = 0x8B31

	// Shader parameters
	COMPILE_STATUS  uint32 = 0x8B81
	LINK_STATUS     uint32 = 0x8B82
	INFO_LOG_LENGTH uint32 = 0x8B84

	// Texture targets
	TEXTURE_2D uint32 = 0x0DE1

	// Texture parameters
	TEXTURE_MAG_FILTER uint32 = 0x2800
	TEXTURE_MIN_FILTER uint32 = 0x2801
	TEXTURE_WRAP_S     uint32 = 0x2802
	TEXTURE_WRAP_T     uint32 = 0x2803

	// Texture filter modes
	NEAREST uint32 = 0x2600
	LINEAR  uint32 = 0x2601

	// Texture wrap modes
	REPEAT          uint32 = 0x2901
	CLAMP_TO_EDGE   uint32 = 0x812F
	MIRRORED_REPEAT uint32 = 0x8370

	// Texture formats
	ALPHA           uint32 = 0x1906
	RGB             uint32 = 0x1907
	RGBA            uint32 = 0x1908
	LUMINANCE       uint32 = 0x1909
	LUMINANCE_ALPHA uint32 = 0x190A
	SRGB_ALPHA      uint32 = 0x8C42
	SRGB8_ALPHA8    uint32 = 0x8C43

	// Texture units
	TEXTURE0  uint32 = 0x84C0
	TEXTURE1  uint32 = 0x84C1
	TEXTURE2  uint32 = 0x84C2
	TEXTURE3  uint32 = 0x84C3
	TEXTURE4  uint32 = 0x84C4
	TEXTURE5  uint32 = 0x84C5
	TEXTURE6  uint32 = 0x84C6
	TEXTURE7  uint32 = 0x84C7
	TEXTURE8  uint32 = 0x84C8
	TEXTURE9  uint32 = 0x84C9
	TEXTURE10 uint32 = 0x84CA

	// Enable/Disable
	BLEND            uint32 = 0x0BE2
	CULL_FACE        uint32 = 0x0B44
	DEPTH_TEST       uint32 = 0x0B71
	SCISSOR_TEST     uint32 = 0x0C11
	FRAMEBUFFER_SRGB uint32 = 0x8DB9

	// Blend functions
	ZERO                     uint32 = 0
	ONE                      uint32 = 1
	SRC_COLOR                uint32 = 0x0300
	ONE_MINUS_SRC_COLOR      uint32 = 0x0301
	SRC_ALPHA                uint32 = 0x0302
	ONE_MINUS_SRC_ALPHA      uint32 = 0x0303
	DST_ALPHA                uint32 = 0x0304
	ONE_MINUS_DST_ALPHA      uint32 = 0x0305
	DST_COLOR                uint32 = 0x0306
	ONE_MINUS_DST_COLOR      uint32 = 0x0307
	SRC_ALPHA_SATURATE       uint32 = 0x0308
	CONSTANT_COLOR           uint32 = 0x8001
	ONE_MINUS_CONSTANT_COLOR uint32 = 0x8002
	CONSTANT_ALPHA           uint32 = 0x8003
	ONE_MINUS_CONSTANT_ALPHA uint32 = 0x8004

	// Clear buffer bits
	DEPTH_BUFFER_BIT   uint32 = 0x00000100
	STENCIL_BUFFER_BIT uint32 = 0x00000400
	COLOR_BUFFER_BIT   uint32 = 0x00004000

	// Boolean values
	FALSE uint32 = 0
	TRUE  uint32 = 1

	// Pixel storage parameters (WebGL-specific but included for compatibility)
	UNPACK_FLIP_Y_WEBGL                uint32 = 0x9240
	UNPACK_PREMULTIPLY_ALPHA_WEBGL     uint32 = 0x9241
	UNPACK_COLORSPACE_CONVERSION_WEBGL uint32 = 0x9243

	// Get parameters
	VIEWPORT uint32 = 0x0BA2
)
