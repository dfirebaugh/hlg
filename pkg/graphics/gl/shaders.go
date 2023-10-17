package gl

import _ "embed"

//go:embed shaders/basic.vert
var BasicVert string

//go:embed shaders/basic.frag
var BasicFrag string

//go:embed shaders/texture.vert
var TextureVert string

//go:embed shaders/texture.frag
var TextureFrag string

//go:embed shaders/shape.vert
var ShapeVert string

//go:embed shaders/shape.frag
var ShapeFrag string

//go:embed shaders/mesh.vert
var MeshVert string

//go:embed shaders/mesh.frag
var MeshFrag string

//go:embed shaders/model.vert
var ModelVert string

//go:embed shaders/model.frag
var ModelFrag string

//go:embed shaders/nontextured_model.vert
var NonTexturedModelVert string

//go:embed shaders/nontextured_model.frag
var NonTexturedModelFrag string
