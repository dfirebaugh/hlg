package glapi

// Buffer represents a GL buffer object
// On OpenGL: uint32
// On WebGL: js.Value wrapper
type Buffer uint32

// Texture represents a GL texture object
type Texture uint32

// Program represents a GL shader program
type Program uint32

// Shader represents a GL shader object
type Shader uint32

// VertexArray represents a GL vertex array object
type VertexArray uint32

// UniformLocation represents a GL uniform location
type UniformLocation int32

// InvalidBuffer represents an invalid/null buffer
const InvalidBuffer Buffer = 0

// InvalidTexture represents an invalid/null texture
const InvalidTexture Texture = 0

// InvalidProgram represents an invalid/null program
const InvalidProgram Program = 0

// InvalidShader represents an invalid/null shader
const InvalidShader Shader = 0

// InvalidVertexArray represents an invalid/null VAO
const InvalidVertexArray VertexArray = 0

// InvalidUniformLocation represents an invalid uniform location
const InvalidUniformLocation UniformLocation = -1
