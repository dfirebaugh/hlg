package hlg

import (
	"github.com/dfirebaugh/hlg/graphics"
)

type Renderable graphics.ShaderRenderable

type Uniform struct {
	Binding uint32
	Size    uint64
}

type VertexAttributeLayout struct {
	ShaderLocation uint32
	Offset         uint64
	Format         string
}

type VertexBufferLayout struct {
	ArrayStride uint64
	Attributes  []VertexAttributeLayout
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

func convertToPipelinesVertexBufferLayout(layout VertexBufferLayout) graphics.VertexBufferLayout {
	attributes := make([]graphics.VertexAttributeLayout, len(layout.Attributes))
	for i, attr := range layout.Attributes {
		attributes[i] = graphics.VertexAttributeLayout{
			ShaderLocation: attr.ShaderLocation,
			Offset:         attr.Offset,
			Format:         attr.Format,
		}
	}
	return graphics.VertexBufferLayout{
		ArrayStride: layout.ArrayStride,
		Attributes:  attributes,
	}
}

// CreateRenderable creates a `graphics.ShaderRenderable` with the provided shader handle, uniforms, and vertices.
func CreateRenderable(shaderHandle int, vertexData []byte, layout VertexBufferLayout, uniforms map[string]Uniform, dataMap map[string][]byte) Renderable {
	graphicsUniforms := convertUniformsToGraphics(uniforms)
	pipelinesLayout := convertToPipelinesVertexBufferLayout(layout)

	return hlg.graphicsBackend.AddDynamicRenderable(vertexData, pipelinesLayout, shaderHandle, graphicsUniforms, dataMap)
}
