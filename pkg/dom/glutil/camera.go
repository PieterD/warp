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

func (c *Camera) Location() mgl32.Vec3 {
	return c.Rotation.Inverse().Rotate(mgl32.Vec3{0, 0, c.Distance})

	//v4 := mgl32.Ident4().
	//	Mul4(mgl32.Translate3D(c.Target[0], c.Target[1], c.Target[2])).
	//	Mul4x1(mgl32.Vec4{0, 0, c.Distance, 1.0},
	//)
	//return mgl32.Vec3{v4[0], v4[1], v4[2]}
}
