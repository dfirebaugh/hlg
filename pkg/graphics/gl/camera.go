package gl

import "github.com/go-gl/mathgl/mgl32"

var (
	cameraPos           = mgl32.Vec3{0.0, 0.0, 3.0}
	cameraFront         = mgl32.Vec3{0.0, 0.0, -1.0}
	cameraUp            = mgl32.Vec3{0.0, 1.0, 0.0}
	zoom        float32 = 45.0 // 45 degrees FOV as a starting point
)
