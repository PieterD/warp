package main

import (
	"fmt"

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
	buffers := []*gl.Buffer{
		glx.Buffer(),
		glx.Buffer(),
	}
	var vaos []*gl.VertexArray
	totalStride := 3*4 + 3*4 + 3*4 + 4
	for _, buffer := range buffers {
		buffer.VertexData(make([]float32, particleLimit))
		vao, err := glx.VertexArray(
			gl.VertexArrayAttribute{
				Name:   "Location",
				Type:   gl.Vec3,
				Buffer: buffer,
				Layout: gl.VertexArrayAttributeLayout{
					ByteOffset: 0,
					ByteStride: totalStride,
				},
			},
			gl.VertexArrayAttribute{
				Name:   "Momentum",
				Type:   gl.Vec3,
				Buffer: buffer,
				Layout: gl.VertexArrayAttributeLayout{
					ByteOffset: 3 * 4,
					ByteStride: totalStride,
				},
			},
			gl.VertexArrayAttribute{
				Name:   "Color",
				Type:   gl.Vec3,
				Buffer: buffer,
				Layout: gl.VertexArrayAttributeLayout{
					ByteOffset: 3*4 + 3*4,
					ByteStride: totalStride,
				},
			},
			gl.VertexArrayAttribute{
				Name:   "Lifetime",
				Type:   gl.Float,
				Buffer: buffer,
				Layout: gl.VertexArrayAttributeLayout{
					ByteOffset: 3*4 + 3*4 + 3*4,
					ByteStride: totalStride,
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("building VAO: %w", err)
		}
		vaos = append(vaos, vao)
	}

	return p, nil
}
