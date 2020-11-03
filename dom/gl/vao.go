package gl

import (
	"fmt"

	"github.com/PieterD/warp/driver"
)

type VertexArrayConfig struct {
	Attributes []VertexArrayAttribute
}

type VertexArrayAttribute struct {
	ArrayBuffer *Buffer
	Attr        *Attribute
	Layout      VertexArrayAttributeLayout
}

type VertexArrayAttributeLayout struct {
	ByteOffset int
	ByteStride int
}

type VertexArray struct {
	glx      *Context
	glObject driver.Value
}

func newVertexArray(glx *Context, cfg VertexArrayConfig) (*VertexArray, error) {
	glVAO := glx.functions.CreateVertexArray()
	glx.functions.BindVertexArray(glVAO)
	defer glx.functions.BindVertexArray(glx.factory.Null())

	for _, da := range cfg.Attributes {
		glAttrIndex := glx.factory.Number(float64(da.Attr.index))
		attrType, attrSize := da.Attr.Type()
		if attrSize != 1 {
			return nil, fmt.Errorf("unable to handle attrSize: %d", attrSize)
		}
		bufferType, bufferItemsPerVertex, err := attrType.asAttribute()
		if err != nil {
			return nil, fmt.Errorf("converting attribute type %s to attribute: %w", attrType, err)
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
	glx.functions.BindBuffer(glx.constants.ARRAY_BUFFER, glx.factory.Null())

	return &VertexArray{
		glx:      glx,
		glObject: glVAO,
	}, nil
}
