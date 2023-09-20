package renderer

import (
	"image/color"
	"unsafe"
)

func rgbaFromColor(c color.Color) (r, g, b, a uint8) {
	R, G, B, A := c.RGBA()
	return uint8(R >> 8), uint8(G >> 8), uint8(B >> 8), uint8(A >> 8)
}

// Convert a C-style string (null-terminated) to a Go string
func goStringFromCString(str *byte) string {
	if str == nil {
		return ""
	}
	var buffer []byte
	for i := str; *i != 0; i = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(i)) + 1)) {
		buffer = append(buffer, *i)
	}
	return string(buffer)
}
