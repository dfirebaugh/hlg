package assets

import _ "embed"

//go:embed models/the-utah-teapot/source/teapot.obj
var TeaPot string

//go:embed models/the-utah-teapot/source/default.png
var DefaultTextureImage []byte
