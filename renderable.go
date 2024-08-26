package hlg

import (
	"github.com/dfirebaugh/hlg/graphics"
)

type Renderable graphics.ShaderRenderable

type Uniform struct {
	Binding uint32
	Size    uint64
}

func convertToGraphicsUniform(hu Uniform) graphics.Uniform {
	return graphics.Uniform{
		Binding: hu.Binding,
		Size:    hu.Size,
	}
}

func convertUniformsToGraphics(uniforms map[string]Uniform) map[string]graphics.Uniform {
	gus := make(map[string]graphics.Uniform)
	for k, u := range uniforms {
		gus[k] = convertToGraphicsUniform(u)
	}
	return gus
}

// CreateRenderable creates a `graphics.ShaderRenderable` with the provided shader handle, uniforms, and vertices.
func CreateRenderable(shaderHandle int, vertices []Vertex, uniforms map[string]Uniform, dataMap map[string][]byte) Renderable {
	graphicsVertices := convertVerticesToGraphics(vertices)
	graphicsUniforms := convertUniformsToGraphics(uniforms)

	return hlg.graphicsBackend.AddDynamicRenderable(graphicsVertices, shaderHandle, graphicsUniforms, dataMap)
}
