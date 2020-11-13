package gl

import (
	"fmt"

	"github.com/PieterD/warp/pkg/driver"
)

type DrawConfig struct {
	Use          *Program
	VAO          *VertexArray
	ElementArray *Buffer // Optional
	DrawMode     DrawMode
	Vertices     VertexRange
	Options      DrawOptions
}

type DrawOptions struct {
	DepthTest     bool
	DepthReadOnly bool // DepthMask
	//TODO: DepthFunc
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

	if cfg.Options.DepthTest {
		glx.constants.Enable(glx.constants.DEPTH_TEST)
		defer glx.constants.Disable(glx.constants.DEPTH_TEST)
		glx.constants.DepthMask(glx.factory.Boolean(!cfg.Options.DepthReadOnly))
		glx.constants.DepthFunc(glx.constants.LESS)
	}

	glx.constants.UseProgram(cfg.Use.glObject)
	defer glx.constants.UseProgram(glx.factory.Null())
	glx.constants.BindVertexArray(cfg.VAO.glObject)
	defer glx.constants.BindVertexArray(glx.factory.Null())
	if cfg.ElementArray == nil {
		glx.constants.DrawArrays(
			glDrawMode,
			glx.factory.Number(float64(cfg.Vertices.FirstOffset)),
			glx.factory.Number(float64(cfg.Vertices.VertexCount)),
		)
		return nil
	}
	glx.constants.BindBuffer(glx.constants.ELEMENT_ARRAY_BUFFER, cfg.ElementArray.glObject)
	defer glx.constants.BindBuffer(glx.constants.ELEMENT_ARRAY_BUFFER, glx.factory.Null())
	glx.constants.DrawElements(
		glDrawMode,
		glx.factory.Number(float64(cfg.Vertices.VertexCount)),
		glx.constants.UNSIGNED_SHORT,
		glx.factory.Number(float64(cfg.Vertices.FirstOffset*2)),
	)
	return nil
}
