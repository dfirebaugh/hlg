package graphics

import (
	"image/color"
	"math"
)

// colorToFloat32 converts a color.Color to a [4]float32 RGBA array.
func colorToFloat32(c color.Color) [4]float32 {
	r, g, b, a := c.RGBA()
	return [4]float32{
		float32(r) / 0xffff,
		float32(g) / 0xffff,
		float32(b) / 0xffff,
		float32(a) / 0xffff,
	}
}

// screenToNDC transforms screen space coordinates to NDC.
func screenToNDC(x, y, screenWidth, screenHeight float32) [3]float32 {
	normalizedX := x / screenWidth
	normalizedY := y / screenHeight
	ndcX := normalizedX*2 - 1
	ndcY := 1 - normalizedY*2
	return [3]float32{ndcX, ndcY, 0}
}

// MakeSolidTriangle creates vertices for a solid triangle using OpCodeSolid.
// Coordinates are in screen space and will be converted to NDC.
// LocalPosition stores barycentric coordinates for edge anti-aliasing.
// Returns 6 vertices (2 identical triangles) to match the primitive buffer format.
func MakeSolidTriangle(x1, y1, x2, y2, x3, y3 int, c color.Color, screenW, screenH int) []PrimitiveVertex {
	col := colorToFloat32(c)
	sw, sh := float32(screenW), float32(screenH)

	// Barycentric coordinates: v1(1,0)->bary(1,0,0), v2(0,1)->bary(0,1,0), v3(0,0)->bary(0,0,1)
	v1 := PrimitiveVertex{
		Position:      screenToNDC(float32(x1), float32(y1), sw, sh),
		LocalPosition: [2]float32{1, 0},
		OpCode:        OpCodeSolid,
		Radius:        0,
		Color:         col,
		TexCoords:     [2]float32{0, 0},
	}
	v2 := PrimitiveVertex{
		Position:      screenToNDC(float32(x2), float32(y2), sw, sh),
		LocalPosition: [2]float32{0, 1},
		OpCode:        OpCodeSolid,
		Radius:        0,
		Color:         col,
		TexCoords:     [2]float32{0, 0},
	}
	v3 := PrimitiveVertex{
		Position:      screenToNDC(float32(x3), float32(y3), sw, sh),
		LocalPosition: [2]float32{0, 0},
		OpCode:        OpCodeSolid,
		Radius:        0,
		Color:         col,
		TexCoords:     [2]float32{0, 0},
	}

	// Return 6 vertices (triangle duplicated) to match primitive buffer format
	return []PrimitiveVertex{v1, v2, v3, v1, v2, v3}
}

// MakeSolidRectangle creates vertices for a solid rectangle using OpCodeSolid.
// Coordinates are in screen space and will be converted to NDC.
// LocalPosition stores barycentric coordinates for edge anti-aliasing.
func MakeSolidRectangle(x, y, w, h int, c color.Color, screenW, screenH int) []PrimitiveVertex {
	col := colorToFloat32(c)
	sw, sh := float32(screenW), float32(screenH)

	topLeft := screenToNDC(float32(x), float32(y), sw, sh)
	topRight := screenToNDC(float32(x+w), float32(y), sw, sh)
	bottomLeft := screenToNDC(float32(x), float32(y+h), sw, sh)
	bottomRight := screenToNDC(float32(x+w), float32(y+h), sw, sh)

	makeVertex := func(pos [3]float32, bary [2]float32) PrimitiveVertex {
		return PrimitiveVertex{
			Position:      pos,
			LocalPosition: bary,
			OpCode:        OpCodeSolid,
			Radius:        0,
			Color:         col,
			TexCoords:     [2]float32{0, 0},
		}
	}

	// Two triangles with barycentric coordinates:
	// Triangle 1: topLeft(1,0), bottomLeft(0,1), topRight(0,0)
	// Triangle 2: bottomLeft(1,0), bottomRight(0,1), topRight(0,0)
	return []PrimitiveVertex{
		makeVertex(topLeft, [2]float32{1, 0}),
		makeVertex(bottomLeft, [2]float32{0, 1}),
		makeVertex(topRight, [2]float32{0, 0}),
		makeVertex(bottomLeft, [2]float32{1, 0}),
		makeVertex(bottomRight, [2]float32{0, 1}),
		makeVertex(topRight, [2]float32{0, 0}),
	}
}

// MakeSolidPolygon creates vertices for a solid regular polygon using OpCodeSolid.
// cx, cy is the center, width is the diameter, sides is the number of sides.
// LocalPosition stores barycentric coordinates for edge anti-aliasing.
func MakeSolidPolygon(cx, cy int, width float32, sides int, c color.Color, screenW, screenH int) []PrimitiveVertex {
	col := colorToFloat32(c)
	sw, sh := float32(screenW), float32(screenH)

	vertices := make([]PrimitiveVertex, 0, sides*3)

	centerNDC := screenToNDC(float32(cx), float32(cy), sw, sh)
	// Barycentric (1, 0) for center vertex -> bary(1, 0, 0)
	centerVertex := PrimitiveVertex{
		Position:      centerNDC,
		LocalPosition: [2]float32{1, 0},
		OpCode:        OpCodeSolid,
		Radius:        0,
		Color:         col,
		TexCoords:     [2]float32{0, 0},
	}

	angleStep := 2 * math.Pi / float64(sides)
	radius := float64(width / 2)
	fcx, fcy := float64(cx), float64(cy)

	x := fcx + radius*math.Cos(0)
	y := fcy + radius*math.Sin(0)

	for i := 0; i < sides; i++ {
		currentNDC := screenToNDC(float32(x), float32(y), sw, sh)
		// Barycentric (0, 1) for current vertex -> bary(0, 1, 0)
		currentVertex := PrimitiveVertex{
			Position:      currentNDC,
			LocalPosition: [2]float32{0, 1},
			OpCode:        OpCodeSolid,
			Radius:        0,
			Color:         col,
			TexCoords:     [2]float32{0, 0},
		}

		nextAngle := float64(i+1) * angleStep
		nextX := fcx + radius*math.Cos(nextAngle)
		nextY := fcy + radius*math.Sin(nextAngle)

		nextNDC := screenToNDC(float32(nextX), float32(nextY), sw, sh)
		// Barycentric (0, 0) for next vertex -> bary(0, 0, 1)
		nextVertex := PrimitiveVertex{
			Position:      nextNDC,
			LocalPosition: [2]float32{0, 0},
			OpCode:        OpCodeSolid,
			Radius:        0,
			Color:         col,
			TexCoords:     [2]float32{0, 0},
		}

		vertices = append(vertices, centerVertex, currentVertex, nextVertex)
		x, y = nextX, nextY
	}

	return vertices
}

// MakeCircle creates vertices for a circle using OpCodeCircle (SDF rendering).
// x, y is the center of the circle, radius is the radius.
func MakeCircle(x, y, radius int, c color.Color, screenW, screenH int) []PrimitiveVertex {
	col := colorToFloat32(c)
	sw, sh := float32(screenW), float32(screenH)

	centerX := float32(x)
	centerY := float32(y)
	r := float32(radius)

	ndcCenterX := (centerX/sw)*2.0 - 1.0
	ndcCenterY := 1.0 - (centerY/sh)*2.0

	ndcRadiusX := r / sw * 2.0
	ndcRadiusY := r / sh * 2.0

	left := ndcCenterX - ndcRadiusX
	right := ndcCenterX + ndcRadiusX
	bottom := ndcCenterY - ndcRadiusY
	top := ndcCenterY + ndcRadiusY

	return []PrimitiveVertex{
		{Position: [3]float32{left, bottom, 0.0}, LocalPosition: [2]float32{-1.0, -1.0}, OpCode: OpCodeCircle, Radius: r, Color: col, TexCoords: [2]float32{r, r}},
		{Position: [3]float32{right, bottom, 0.0}, LocalPosition: [2]float32{1.0, -1.0}, OpCode: OpCodeCircle, Radius: r, Color: col, TexCoords: [2]float32{r, r}},
		{Position: [3]float32{left, top, 0.0}, LocalPosition: [2]float32{-1.0, 1.0}, OpCode: OpCodeCircle, Radius: r, Color: col, TexCoords: [2]float32{r, r}},

		{Position: [3]float32{left, top, 0.0}, LocalPosition: [2]float32{-1.0, 1.0}, OpCode: OpCodeCircle, Radius: r, Color: col, TexCoords: [2]float32{r, r}},
		{Position: [3]float32{right, bottom, 0.0}, LocalPosition: [2]float32{1.0, -1.0}, OpCode: OpCodeCircle, Radius: r, Color: col, TexCoords: [2]float32{r, r}},
		{Position: [3]float32{right, top, 0.0}, LocalPosition: [2]float32{1.0, 1.0}, OpCode: OpCodeCircle, Radius: r, Color: col, TexCoords: [2]float32{r, r}},
	}
}

// MakeRoundedRectangle creates vertices for a rounded rectangle using OpCodeRoundedRect.
// This uses SDF rendering for smooth corners.
// x, y is the top-left corner, w, h are width and height, radius is the corner radius.
func MakeRoundedRectangle(x, y, w, h, radius int, c color.Color, screenW, screenH int) []PrimitiveVertex {
	col := colorToFloat32(c)
	sw, sh := float32(screenW), float32(screenH)

	centerX := float32(x + w/2)
	centerY := float32(y + h/2)

	ndcCenterX := (centerX/sw)*2.0 - 1.0
	ndcCenterY := 1.0 - (centerY/sh)*2.0
	ndcWidth := (float32(w) / sw) * 2.0
	ndcHeight := (float32(h) / sh) * 2.0

	left := ndcCenterX - ndcWidth/2
	right := ndcCenterX + ndcWidth/2
	bottom := ndcCenterY - ndcHeight/2
	top := ndcCenterY + ndcHeight/2

	halfWidth := float32(w) / 2.0
	halfHeight := float32(h) / 2.0
	cornerRadius := float32(radius)

	return []PrimitiveVertex{
		{Position: [3]float32{left, bottom, 0.0}, LocalPosition: [2]float32{-1.0, -1.0}, OpCode: OpCodeRoundedRect, Radius: cornerRadius, Color: col, TexCoords: [2]float32{halfWidth, halfHeight}},
		{Position: [3]float32{right, bottom, 0.0}, LocalPosition: [2]float32{1.0, -1.0}, OpCode: OpCodeRoundedRect, Radius: cornerRadius, Color: col, TexCoords: [2]float32{halfWidth, halfHeight}},
		{Position: [3]float32{left, top, 0.0}, LocalPosition: [2]float32{-1.0, 1.0}, OpCode: OpCodeRoundedRect, Radius: cornerRadius, Color: col, TexCoords: [2]float32{halfWidth, halfHeight}},

		{Position: [3]float32{left, top, 0.0}, LocalPosition: [2]float32{-1.0, 1.0}, OpCode: OpCodeRoundedRect, Radius: cornerRadius, Color: col, TexCoords: [2]float32{halfWidth, halfHeight}},
		{Position: [3]float32{right, bottom, 0.0}, LocalPosition: [2]float32{1.0, -1.0}, OpCode: OpCodeRoundedRect, Radius: cornerRadius, Color: col, TexCoords: [2]float32{halfWidth, halfHeight}},
		{Position: [3]float32{right, top, 0.0}, LocalPosition: [2]float32{1.0, 1.0}, OpCode: OpCodeRoundedRect, Radius: cornerRadius, Color: col, TexCoords: [2]float32{halfWidth, halfHeight}},
	}
}

// MakeSolidLine creates vertices for a solid line using OpCodeLine.
// The line is rendered using an SDF capsule shape.
func MakeSolidLine(x1, y1, x2, y2 int, width float32, c color.Color, screenW, screenH int) []PrimitiveVertex {
	col := colorToFloat32(c)
	sw, sh := float32(screenW), float32(screenH)

	dx := float32(x2 - x1)
	dy := float32(y2 - y1)
	length := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	if length == 0 {
		return nil
	}

	sin := dy / length
	cos := dx / length

	halfWidth := width / 2

	// Compute the 4 corners of the thick line (for bounding box)
	c0x := float32(x1) - sin*halfWidth
	c0y := float32(y1) + cos*halfWidth
	c1x := float32(x2) - sin*halfWidth
	c1y := float32(y2) + cos*halfWidth
	c2x := float32(x2) + sin*halfWidth
	c2y := float32(y2) - cos*halfWidth
	c3x := float32(x1) + sin*halfWidth
	c3y := float32(y1) - cos*halfWidth

	p0 := screenToNDC(c0x, c0y, sw, sh)
	p1 := screenToNDC(c1x, c1y, sw, sh)
	p2 := screenToNDC(c2x, c2y, sw, sh)
	p3 := screenToNDC(c3x, c3y, sw, sh)

	// Store line endpoints in TexCoords for the shader
	// TexCoords[0] = normalized x1, TexCoords[1] = normalized y1 (will be expanded in Extra)
	// We use LocalPosition to store the normalized line direction
	// Store direction and half-length for the shader
	halfLength := length / 2

	makeVertex := func(pos [3]float32, localX, localY float32) PrimitiveVertex {
		return PrimitiveVertex{
			Position:      pos,
			LocalPosition: [2]float32{localX, localY},
			OpCode:        OpCodeLine,
			Radius:        halfWidth,
			Color:         col,
			// Store line direction (cos, sin) and half_length in TexCoords
			// TexCoords[0] = cos, TexCoords[1] = sin
			// We'll encode half_length by storing it scaled - see conversion
			TexCoords: [2]float32{cos * halfLength, sin * halfLength},
		}
	}

	// LocalPosition corners for the quad: BL, BR, TL, TR pattern
	return []PrimitiveVertex{
		makeVertex(p0, -1, -1), // BL
		makeVertex(p1, 1, -1),  // BR
		makeVertex(p2, 1, 1),   // TR (note: this is actually BR of line end)
		makeVertex(p0, -1, -1), // BL
		makeVertex(p2, 1, 1),   // TR
		makeVertex(p3, -1, 1),  // TL
	}
}

// === New Primitive constructors (storage buffer format) ===
// These return a single Primitive struct instead of 6 PrimitiveVertex.
// Screen-space coordinates - the shader handles NDC conversion.

// MakeRectPrimitive creates a Primitive for a solid rectangle.
// x, y is the top-left corner, w, h are width and height in screen pixels.
func MakeRectPrimitive(x, y, w, h int, c color.Color) Primitive {
	col := colorToFloat32(c)
	halfW := float32(w) / 2.0
	halfH := float32(h) / 2.0
	return Primitive{
		X:      float32(x),
		Y:      float32(y),
		W:      float32(w),
		H:      float32(h),
		Color:  col,
		Radius: 0,
		OpCode: OpCodeSolid,
		Extra:  [4]float32{halfW, halfH, 0, 0},
	}
}

// MakeCirclePrimitive creates a Primitive for a circle using SDF rendering.
// x, y is the center, radius is the radius in screen pixels.
func MakeCirclePrimitive(x, y, radius int, c color.Color) Primitive {
	col := colorToFloat32(c)
	r := float32(radius)
	return Primitive{
		X:      float32(x) - r,
		Y:      float32(y) - r,
		W:      r * 2,
		H:      r * 2,
		Color:  col,
		Radius: r,
		OpCode: OpCodeCircle,
		Extra:  [4]float32{r, r, 0, 0},
	}
}

// MakeRoundedRectPrimitive creates a Primitive for a rounded rectangle.
// x, y is the top-left corner, w, h are dimensions, radius is corner radius.
func MakeRoundedRectPrimitive(x, y, w, h, radius int, c color.Color) Primitive {
	col := colorToFloat32(c)
	halfW := float32(w) / 2.0
	halfH := float32(h) / 2.0
	return Primitive{
		X:      float32(x),
		Y:      float32(y),
		W:      float32(w),
		H:      float32(h),
		Color:  col,
		Radius: float32(radius),
		OpCode: OpCodeRoundedRect,
		Extra:  [4]float32{halfW, halfH, 0, 0},
	}
}

// MakeLinePrimitive creates a Primitive for a line segment.
// x1, y1, x2, y2 are endpoints, width is line thickness.
func MakeLinePrimitive(x1, y1, x2, y2 int, width float32, c color.Color) Primitive {
	col := colorToFloat32(c)

	dx := float32(x2 - x1)
	dy := float32(y2 - y1)
	length := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	if length == 0 {
		return Primitive{}
	}

	sin := dy / length
	cos := dx / length
	halfWidth := width / 2
	halfLength := length / 2

	// Compute bounding box
	c0x := float32(x1) - sin*halfWidth
	c0y := float32(y1) + cos*halfWidth
	c1x := float32(x2) - sin*halfWidth
	c1y := float32(y2) + cos*halfWidth
	c2x := float32(x2) + sin*halfWidth
	c2y := float32(y2) - cos*halfWidth
	c3x := float32(x1) + sin*halfWidth
	c3y := float32(y1) - cos*halfWidth

	minX := min(min(c0x, c1x), min(c2x, c3x))
	maxX := max(max(c0x, c1x), max(c2x, c3x))
	minY := min(min(c0y, c1y), min(c2y, c3y))
	maxY := max(max(c0y, c1y), max(c2y, c3y))

	return Primitive{
		X:      minX,
		Y:      minY,
		W:      maxX - minX,
		H:      maxY - minY,
		Color:  col,
		Radius: halfWidth,
		OpCode: OpCodeLine,
		Extra:  [4]float32{cos * halfLength, sin * halfLength, 0, 0},
	}
}

// ConvertPrimitivesToVertices converts Primitive structs to PrimitiveVertex format.
// This is used by EndDraw to merge primitives with raw vertices.
func ConvertPrimitivesToVertices(primitives []Primitive, screenW, screenH float32) []PrimitiveVertex {
	if len(primitives) == 0 {
		return nil
	}

	vertices := make([]PrimitiveVertex, 0, len(primitives)*6)

	for _, prim := range primitives {
		// Convert screen coords to NDC
		x0 := (prim.X/screenW)*2 - 1
		y0 := 1 - (prim.Y/screenH)*2
		x1 := ((prim.X+prim.W)/screenW)*2 - 1
		y1 := 1 - ((prim.Y+prim.H)/screenH)*2

		halfW := prim.W / 2
		halfH := prim.H / 2

		var texCoords [4][2]float32
		if prim.OpCode == OpCodeMSDF {
			u0, v0 := prim.Extra[0], prim.Extra[1]
			us, vs := prim.Extra[2], prim.Extra[3]
			texCoords = [4][2]float32{
				{u0, v0 + vs},
				{u0 + us, v0 + vs},
				{u0, v0},
				{u0 + us, v0},
			}
		} else if prim.OpCode == OpCodeLine {
			dirX, dirY := prim.Extra[0], prim.Extra[1]
			texCoords = [4][2]float32{
				{dirX, dirY},
				{dirX, dirY},
				{dirX, dirY},
				{dirX, dirY},
			}
		} else {
			texCoords = [4][2]float32{
				{halfW, halfH},
				{halfW, halfH},
				{halfW, halfH},
				{halfW, halfH},
			}
		}

		cornerPositions := [6][3]float32{
			{x0, y1, 0},
			{x1, y1, 0},
			{x0, y0, 0},
			{x0, y0, 0},
			{x1, y1, 0},
			{x1, y0, 0},
		}

		localPositions := [6][2]float32{
			{-1, -1},
			{1, -1},
			{-1, 1},
			{-1, 1},
			{1, -1},
			{1, 1},
		}

		texCoordIndices := [6]int{0, 1, 2, 2, 1, 3}
		halfSize := [2]float32{halfW, halfH}

		for i := range 6 {
			vertices = append(vertices, PrimitiveVertex{
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

// ExtractClipRectsFromPrimitives extracts clip rects from primitives.
// Returns one clip rect per vertex (6 per primitive, matching ConvertPrimitivesToVertices output).
func ExtractClipRectsFromPrimitives(primitives []Primitive) []*[4]int {
	if len(primitives) == 0 {
		return nil
	}

	clipRects := make([]*[4]int, 0, len(primitives)*6)
	for _, prim := range primitives {
		// Each primitive generates 6 vertices
		for range 6 {
			clipRects = append(clipRects, prim.ClipRect)
		}
	}
	return clipRects
}
