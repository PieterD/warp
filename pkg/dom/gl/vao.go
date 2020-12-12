package gl

import (
	"fmt"

	"github.com/PieterD/warp/pkg/driver"
)

type VertexArrayAttributeConfig struct {
	Name       string
	Type       Type
	Buffer     *Buffer
	ByteOffset int
	ByteStride int
}

type VertexArray struct {
	glx      *Context
	glObject driver.Value
	attrs    map[string]*vertexAttr
}

type vertexAttr struct {
	enabled bool
	name    string
	typ     Type
	index   int
}

func newVertexArray(glx *Context, attrs ...VertexArrayAttributeConfig) (*VertexArray, error) {
	glVAO := glx.constants.CreateVertexArray()
	glx.constants.BindVertexArray(glVAO)
	defer glx.constants.BindVertexArray(glx.factory.Null())

	attrMap := make(map[string]*vertexAttr)
	for attrIndex, attr := range attrs {
		fmt.Printf("VAO ATTR: %s %s %d %d\n", attr.Name, attr.Type, attr.ByteOffset, attr.ByteStride)
		glAttrIndex := glx.factory.Number(float64(attrIndex))
		attrType := attr.Type
		bufferType, bufferItemsPerVertex, err := attrType.asAttribute()
		if err != nil {
			return nil, fmt.Errorf("converting attribute type %s to attribute: %w", attrType, err)
		}
		glBufferType := glx.typeConverter.ToJs(bufferType)
		glItemsPerVertex := glx.factory.Number(float64(bufferItemsPerVertex))
		glNormalized := glx.factory.Boolean(false)
		glByteStride := glx.factory.Number(float64(attr.ByteStride))
		glByteOffset := glx.factory.Number(float64(attr.ByteOffset))
		glx.constants.BindBuffer(glx.constants.ARRAY_BUFFER, attr.Buffer.glObject)
		glx.constants.VertexAttribPointer(
			glAttrIndex,
			glItemsPerVertex,
			glBufferType,
			glNormalized,
			glByteStride,
			glByteOffset,
		)
		glx.constants.EnableVertexAttribArray(glAttrIndex)
		attrMap[attr.Name] = &vertexAttr{
			enabled: true,
			name:    attr.Name,
			typ:     attr.Type,
			index:   attrIndex,
		}
	}
	glx.constants.BindBuffer(glx.constants.ARRAY_BUFFER, glx.factory.Null())

	return &VertexArray{
		glx:      glx,
		glObject: glVAO,
		attrs:    attrMap,
	}, nil
}

func (vao *VertexArray) Enable(attrNames ...string) error {
	glx := vao.glx
	enabledMap := make(map[string]struct{})
	for _, attrName := range attrNames {
		enabledMap[attrName] = struct{}{}
	}
	for attrName, va := range vao.attrs {
		_, shouldBeEnabled := enabledMap[attrName]
		if va.enabled == shouldBeEnabled {
			continue
		}
		if shouldBeEnabled {
			glx.constants.EnableVertexAttribArray(glx.factory.Number(float64(va.index)))
		} else {
			glx.constants.DisableVertexAttribArray(glx.factory.Number(float64(va.index)))
		}
		va.enabled = shouldBeEnabled
	}
	return nil
}
