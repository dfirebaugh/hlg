package gl

import (
	"github.com/dfirebaugh/ggez/pkg/math/geom"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	cameraPos           = mgl32.Vec3{0.0, 0.0, 3.0}
	cameraFront         = mgl32.Vec3{0.0, 0.0, -1.0}
	cameraUp            = mgl32.Vec3{0.0, 1.0, 0.0}
	zoom        float32 = 45.0
)

type Camera struct {
	Position    mgl32.Vec3
	Target      mgl32.Vec3
	Up          mgl32.Vec3
	Fov         float32
	AspectRatio float32
	Near        float32
	Far         float32
}

func NewCamera() *Camera {
	return &Camera{
		Position:    mgl32.Vec3{0, 0, 3},
		Target:      mgl32.Vec3{0, 0, -1},
		Up:          mgl32.Vec3{0, 1, 0},
		Fov:         60.0,
		AspectRatio: 16.0 / 9.0,
		Near:        0.1,
		Far:         100.0,
	}
}

func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.Position, c.Target, c.Up)
}

func (c *Camera) GetProjectionMatrix() mgl32.Mat4 {
	return mgl32.Perspective(c.Fov, c.AspectRatio, c.Near, c.Far)
}

func (c *Camera) SetCameraPosition(position geom.Vector3D, target geom.Vector3D, up geom.Vector3D) {
	c.Position = mgl32.Vec3{position[0], position[1], position[2]}
	c.Target = mgl32.Vec3{target[0], target[1], target[2]}
	c.Up = mgl32.Vec3{up[0], up[1], up[2]}
}
