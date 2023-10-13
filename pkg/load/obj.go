package load

import (
	"fmt"

	"github.com/dfirebaugh/ggez/pkg/math/geom"
	"github.com/udhos/gwob"
)

func LoadOBJModel(filePath string) (*geom.Model, error) {
	options := gwob.ObjParserOptions{}

	obj, err := gwob.NewObjFromFile(filePath, &options)
	if err != nil {
		return nil, fmt.Errorf("failed to parse obj: %w", err)
	}

	if len(obj.Coord) == 0 || len(obj.Indices) == 0 {
		return nil, fmt.Errorf("no valid coordinates or indices found in obj file")
	}

	vertices := make([]float32, len(obj.Coord))
	copy(vertices, obj.Coord)

	indices := make([]uint32, len(obj.Indices))
	for i, index := range obj.Indices {
		if index < 0 {
			return nil, fmt.Errorf("negative index found: %d", index)
		}
		indices[i] = uint32(index)
	}

	mesh := geom.NewMesh(vertices, indices)

	return &geom.Model{
		Meshes:      []*geom.Mesh{mesh},
		ScaleFactor: 1,
		Position:    geom.Vector3{X: 0, Y: 0, Z: 0},
		Rotation:    geom.Matrix4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1},
	}, nil
}
