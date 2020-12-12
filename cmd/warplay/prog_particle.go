package main

import (
	"github.com/PieterD/warp/pkg/dom/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type ParticleProgram struct {
	Uniforms struct {
		Model          mgl32.Mat4
		View           mgl32.Mat4
		Projection     mgl32.Mat4
		LightLocation  mgl32.Vec3
		CameraLocation mgl32.Vec3
	}
	glx *gl.Context
	vao *gl.VertexArray
}

func NewParticleProgram(glx *gl.Context, particleLimit int) (*ParticleProgram, error) {
	p := &ParticleProgram{
		glx: glx,
	}
	return p, nil
}
