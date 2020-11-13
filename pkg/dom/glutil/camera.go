package glutil

import "github.com/go-gl/mathgl/mgl32"

type Camera struct {
	Distance float32
	Target   mgl32.Vec3
	Rotation mgl32.Quat
}

func NewCamera(distance float32) *Camera {
	return &Camera{
		Distance: distance,
		Target:   mgl32.Vec3{},
		Rotation: mgl32.QuatIdent(),
	}
}

func (c *Camera) ViewMatrix() mgl32.Mat4 {
	cameraMatrix := mgl32.Ident4().
		Mul4(mgl32.Translate3D(0, 0, -c.Distance)).
		Mul4(c.Rotation.Mat4()).
		Mul4(mgl32.Translate3D(c.Target[0], c.Target[1], c.Target[2]))
	return cameraMatrix
}
