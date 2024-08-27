package shader

import (
	_ "embed"

	"github.com/dfirebaugh/hlg/graphics"
)

var (
	//go:embed shapes.wgsl
	shapeShaderCode string

	//go:embed texture.wgsl
	textureShaderCode string

	//go:embed line.wgsl
	lineShaderCode string

	//go:embed sdf_circle_outline.wgsl
	circleOutlineShaderCode string

	//go:embed sdf_circle_fill.wgsl
	circleFillShaderCode string

	//go:embed sdf_rect_outline.wgsl
	rectOutlineShaderCode string

	//go:embed sdf_rect_fill.wgsl
	rectFillShaderCode string

	ShapeShader         graphics.ShaderHandle
	TextureShader       graphics.ShaderHandle
	LineShader          graphics.ShaderHandle
	CircleOutlineShader graphics.ShaderHandle
	CircleFillShader    graphics.ShaderHandle
	RectOutlineShader   graphics.ShaderHandle
	RectFillShader      graphics.ShaderHandle
)

func CompileShaders(sm *ShaderManager) {
	ShapeShader = sm.CompileShader(shapeShaderCode)
	TextureShader = sm.CompileShader(textureShaderCode)

	// LineShader = sm.CompileShader(lineShaderCode)
	// CircleOutlineShader = sm.CompileShader(circleOutlineShaderCode)
	// CircleFillShader = sm.CompileShader(circleFillShaderCode)
	// RectOutlineShader = sm.CompileShader(rectOutlineShaderCode)
	// RectFillShader = sm.CompileShader(rectFillShaderCode)
}
