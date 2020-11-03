package gl

import (
	"fmt"

	"github.com/PieterD/warp/driver"
)

type DrawConfig struct {
	Attributes []DrawAttribute
	DrawMode   DrawMode
	Vertices   VertexRange
}

type VertexRange struct {
	FirstOffset int
	VertexCount int
}

type DrawAttribute struct {
	ArrayBuffer *Buffer
	Attr        *Attribute
	Layout      AttributeLayout
}

type AttributeLayout struct {
	ByteOffset int
	ByteStride int
}

type DrawMode int

const (
	Triangles DrawMode = iota
)

func doDraw(glx *Context, cfg DrawConfig) error {
	for _, da := range cfg.Attributes {
		glAttrIndex := glx.factory.Number(float64(da.Attr.index))
		attrType, attrSize := da.Attr.Type()
		if attrSize != 1 {
			return fmt.Errorf("unable to handle attrSize: %d", attrSize)
		}
		bufferType, bufferItemsPerVertex, err := attrType.asAttribute()
		if err != nil {
			return fmt.Errorf("converting attribute type %s to attribute: %w", attrType, err)
		}
		glBufferType := glx.typeConverter.ToJs(bufferType)
		glItemsPerVertex := glx.factory.Number(float64(bufferItemsPerVertex))
		glNormalized := glx.factory.Boolean(false)
		glByteStride := glx.factory.Number(float64(da.Layout.ByteStride))
		glByteOffset := glx.factory.Number(float64(da.Layout.ByteOffset))
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
		glx.factory.Number(float64(cfg.Vertices.FirstOffset)),
		glx.factory.Number(float64(cfg.Vertices.VertexCount)),
	)
	return nil
}
