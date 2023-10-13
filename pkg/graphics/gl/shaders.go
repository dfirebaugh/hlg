package gl

import _ "embed"

//go:embed shaders/basic.vert
var BasicVert string

//go:embed shaders/basic.frag
var BasicFrag string

//go:embed shaders/shape.vert
var ShapeVert string

//go:embed shaders/shape.frag
var ShapeFrag string

//go:embed shaders/cube.vert
var CubeVert string

//go:embed shaders/cube.frag
var CubeFrag string

//go:embed shaders/mesh.vert
var MeshVert string

//go:embed shaders/mesh.frag
var MeshFrag string

//go:embed shaders/model.vert
var ModelVert string

//go:embed shaders/model.frag
var ModelFrag string
