package pipelines

import (
	"math"

	"github.com/dfirebaugh/hlg/graphics"
)

// verticesToBytes converts PrimitiveVertex slice to bytes
func verticesToBytes(vertices []graphics.PrimitiveVertex) []byte {
	bytes := make([]byte, len(vertices)*60)
	for i, v := range vertices {
		offset := i * 60
		// Position
		writeFloat32(bytes[offset:], v.Position[0])
		writeFloat32(bytes[offset+4:], v.Position[1])
		writeFloat32(bytes[offset+8:], v.Position[2])
		// LocalPosition
		writeFloat32(bytes[offset+12:], v.LocalPosition[0])
		writeFloat32(bytes[offset+16:], v.LocalPosition[1])
		// OpCode
		writeFloat32(bytes[offset+20:], v.OpCode)
		// Radius
		writeFloat32(bytes[offset+24:], v.Radius)
		// Color
		writeFloat32(bytes[offset+28:], v.Color[0])
		writeFloat32(bytes[offset+32:], v.Color[1])
		writeFloat32(bytes[offset+36:], v.Color[2])
		writeFloat32(bytes[offset+40:], v.Color[3])
		// TexCoords
		writeFloat32(bytes[offset+44:], v.TexCoords[0])
		writeFloat32(bytes[offset+48:], v.TexCoords[1])
		// HalfSize
		writeFloat32(bytes[offset+52:], v.HalfSize[0])
		writeFloat32(bytes[offset+56:], v.HalfSize[1])
	}
	return bytes
}

func writeFloat32(b []byte, f float32) {
	u := math.Float32bits(f)
	b[0] = byte(u)
	b[1] = byte(u >> 8)
	b[2] = byte(u >> 16)
	b[3] = byte(u >> 24)
}

// Convert Primitive to PrimitiveVertex format
func convertPrimitivesToVertices(primitives []graphics.Primitive, screenW, screenH float32) []graphics.PrimitiveVertex {
	vertices := make([]graphics.PrimitiveVertex, 0, len(primitives)*6)

	for _, prim := range primitives {
		// Convert screen coords to NDC
		x0 := (prim.X/screenW)*2 - 1
		y0 := 1 - (prim.Y/screenH)*2
		x1 := ((prim.X+prim.W)/screenW)*2 - 1
		y1 := 1 - ((prim.Y+prim.H)/screenH)*2

		halfW := prim.W / 2
		halfH := prim.H / 2

		var texCoords [4][2]float32
		switch prim.OpCode {
		case graphics.OpCodeMSDF:
			// MSDF: extra contains UV base (xy) and UV size (zw)
			u0, v0 := prim.Extra[0], prim.Extra[1]
			us, vs := prim.Extra[2], prim.Extra[3]
			texCoords = [4][2]float32{
				{u0, v0 + vs},      // bottom-left
				{u0 + us, v0 + vs}, // bottom-right
				{u0, v0},           // top-left
				{u0 + us, v0},      // top-right
			}
		case graphics.OpCodeLine:
			// Line: extra contains half-direction vector (cos*halfLen, sin*halfLen)
			dirX, dirY := prim.Extra[0], prim.Extra[1]
			texCoords = [4][2]float32{
				{dirX, dirY},
				{dirX, dirY},
				{dirX, dirY},
				{dirX, dirY},
			}
		default:
			// Shapes: pass half_size via tex_coords
			texCoords = [4][2]float32{
				{halfW, halfH},
				{halfW, halfH},
				{halfW, halfH},
				{halfW, halfH},
			}
		}

		// 6 vertices per primitive (2 triangles)
		cornerPositions := [6][3]float32{
			{x0, y1, 0}, // 0: bottom-left
			{x1, y1, 0}, // 1: bottom-right
			{x0, y0, 0}, // 2: top-left
			{x0, y0, 0}, // 3: top-left
			{x1, y1, 0}, // 4: bottom-right
			{x1, y0, 0}, // 5: top-right
		}

		localPositions := [6][2]float32{
			{-1, -1}, // bottom-left
			{1, -1},  // bottom-right
			{-1, 1},  // top-left
			{-1, 1},  // top-left
			{1, -1},  // bottom-right
			{1, 1},   // top-right
		}

		texCoordIndices := [6]int{0, 1, 2, 2, 1, 3}
		halfSize := [2]float32{halfW, halfH}

		for i := range 6 {
			vertices = append(vertices, graphics.PrimitiveVertex{
				Position:      cornerPositions[i],
				LocalPosition: localPositions[i],
				OpCode:        prim.OpCode,
				Radius:        prim.Radius,
				Color:         prim.Color,
				TexCoords:     texCoords[texCoordIndices[i]],
				HalfSize:      halfSize,
			})
		}
	}

	return vertices
}
