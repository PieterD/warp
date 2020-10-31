package gl

import (
	"fmt"

	"github.com/PieterD/warp/driver"
)

type DrawConfig struct {
	Attributes        []DrawAttribute
	DrawMode          DrawMode
	FirstVertexOffset int
	VertexCount       int
}

type DrawAttribute struct {
	ArrayBuffer    *Buffer
	Attr           *Attribute
	ItemsPerVertex int
}

type DrawMode int

const (
	Triangles DrawMode = iota
)

func doDraw(glx *Context, cfg DrawConfig) error {
	for _, da := range cfg.Attributes {
		glAttrIndex := glx.factory.Number(float64(da.Attr.index))
		glItemsPerVertex := glx.factory.Number(float64(da.ItemsPerVertex))
		glBufferType := glx.constants.FLOAT
		glNormalized := glx.factory.Boolean(false)
		glByteStride := glx.factory.Number(0)
		glByteOffset := glx.factory.Number(0)
		glx.functions.BindBuffer(glx.constants.ARRAY_BUFFER, da.ArrayBuffer.glObject)
		glx.functions.VertexAttribPointer(
			glAttrIndex,
			glItemsPerVertex,
			glBufferType,
			glNormalized,
			glByteStride,
			glByteOffset,
		)
		glx.functions.EnableVertexAttribArray(glx.factory.Number(float64(da.Attr.index)))
	}

	var glDrawMode driver.Value
	switch cfg.DrawMode {
	case Triangles:
		glDrawMode = glx.constants.TRIANGLES
	default:
		return fmt.Errorf("unsupported draw mode: %v", cfg.DrawMode)
	}
	glx.functions.DrawArrays(
		glDrawMode,
		glx.factory.Number(float64(cfg.FirstVertexOffset)),
		glx.factory.Number(float64(cfg.VertexCount)),
	)
	return nil
}
