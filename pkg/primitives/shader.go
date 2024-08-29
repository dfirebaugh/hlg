package primitives

import (
	_ "embed"

	"github.com/dfirebaugh/hlg"
)

//go:embed line.wgsl
var lineShaderCode string

//go:embed circle_outline.wgsl
var circleOutlineShaderCode string

//go:embed circle_fill.wgsl
var circleFillShaderCode string

//go:embed rect_outline.wgsl
var rectOutlineShaderCode string

//go:embed rect_fill.wgsl
var rectFillShaderCode string

var (
	LineShader          int
	CircleOutlineShader int
	CircleFillShader    int
	RectOutlineShader   int
	RectFillShader      int
)

func init() {
	LineShader = hlg.CompileShader(lineShaderCode)
	CircleOutlineShader = hlg.CompileShader(circleOutlineShaderCode)
	CircleFillShader = hlg.CompileShader(circleFillShaderCode)
	RectOutlineShader = hlg.CompileShader(rectOutlineShaderCode)
	RectFillShader = hlg.CompileShader(rectFillShaderCode)
}
