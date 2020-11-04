package gl

import (
	"fmt"

	"github.com/PieterD/warp/driver"
)

type DrawConfig struct {
	Use          *Program
	Uniforms     func(us *UniformSetter) // Optional
	VAO          *VertexArray
	ElementArray *Buffer // Optional
	DrawMode     DrawMode
	Vertices     VertexRange
}

type VertexRange struct {
	FirstOffset int
	VertexCount int
}

type DrawMode int

const (
	Triangles DrawMode = iota
)

func doDraw(glx *Context, cfg DrawConfig) error {
	var glDrawMode driver.Value
	switch cfg.DrawMode {
	case Triangles:
		glDrawMode = glx.constants.TRIANGLES
	default:
		return fmt.Errorf("unsupported draw mode: %v", cfg.DrawMode)
	}

	glx.functions.UseProgram(cfg.Use.glObject)
	defer glx.functions.UseProgram(glx.factory.Null())
	if cfg.Uniforms != nil {
		cfg.Uniforms(&UniformSetter{glx: glx})
	}
	glx.functions.BindVertexArray(cfg.VAO.glObject)
	defer glx.functions.BindVertexArray(glx.factory.Null())
	if cfg.ElementArray == nil {
		glx.functions.DrawArrays(
			glDrawMode,
			glx.factory.Number(float64(cfg.Vertices.FirstOffset)),
			glx.factory.Number(float64(cfg.Vertices.VertexCount)),
		)
		return nil
	}
	glx.functions.BindBuffer(glx.constants.ELEMENT_ARRAY_BUFFER, cfg.ElementArray.glObject)
	defer glx.functions.BindBuffer(glx.constants.ELEMENT_ARRAY_BUFFER, glx.factory.Null())
	glx.functions.DrawElements(
		glDrawMode,
		glx.factory.Number(float64(cfg.Vertices.VertexCount)),
		glx.constants.UNSIGNED_SHORT,
		glx.factory.Number(float64(cfg.Vertices.FirstOffset*2)),
	)
	return nil
}
