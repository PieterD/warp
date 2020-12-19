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
	inputBuffer := glx.Buffer()
	inputBuffer.VertexData([]float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	outputBuffer := glx.Buffer()
	outputBuffer.VertexData([]float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	dc, err := gl.NewDataCoupling(gl.DataCouplingConfig{
		Vertices: []gl.VertexConfig{
			{
				Name:   "Input",
				Type:   gl.Float,
				Buffer: "input",
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("creating data coupling: %w", err)
	}
	adc := dc.Active("Input")
	vao, err := glx.VertexArray(adc, map[string]*gl.Buffer{
		"input": inputBuffer,
	})
	if err != nil {
		return nil, fmt.Errorf("creating vertex array object: %w", err)
	}

	_ = vao
	return p, nil
}
