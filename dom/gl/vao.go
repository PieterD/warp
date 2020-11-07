package gl

import (
	"fmt"

	"github.com/PieterD/warp/driver"
)

type VertexArrayAttribute struct {
	Buffer *Buffer
	Attr   *Attribute
	Layout VertexArrayAttributeLayout
}

type VertexArrayAttributeLayout struct {
	ByteOffset int
	ByteStride int
}

type VertexArray struct {
	glx      *Context
	glObject driver.Value
}

func newVertexArray(glx *Context, attrs ...VertexArrayAttribute) (*VertexArray, error) {
	glVAO := glx.constants.CreateVertexArray()
	glx.constants.BindVertexArray(glVAO)
	defer glx.constants.BindVertexArray(glx.factory.Null())

	for _, da := range attrs {
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
		glx.constants.BindBuffer(glx.constants.ARRAY_BUFFER, da.Buffer.glObject)
		glx.constants.VertexAttribPointer(
			glAttrIndex,
			glItemsPerVertex,
			glBufferType,
			glNormalized,
			glByteStride,
			glByteOffset,
		)
		glx.constants.EnableVertexAttribArray(glx.factory.Number(float64(da.Attr.index)))
	}
	glx.constants.BindBuffer(glx.constants.ARRAY_BUFFER, glx.factory.Null())

	return &VertexArray{
		glx:      glx,
		glObject: glVAO,
	}, nil
}
