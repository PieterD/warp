package gl

import (
	"fmt"

	"github.com/PieterD/warp/pkg/driver"
)

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

func newVertexArray(glx *Context, adc ActiveCoupling, buffers map[string]*Buffer) (*VertexArray, error) {
	// Verify that all enabled attributes really exist.
	for attrName := range adc.Enabled {
		if _, ok := adc.DC.attrByName[attrName]; !ok {
			return nil, fmt.Errorf("unknown enabled attribute %s in active coupling", attrName)
		}
	}

	glVAO := glx.constants.CreateVertexArray()
	glx.constants.BindVertexArray(glVAO)
	defer glx.constants.BindVertexArray(glx.factory.Null())

	attrMap := make(map[string]*vertexAttr)
	for attrIndex, attr := range adc.DC.attributes {
		//fmt.Printf("VAO ATTR: %s %s %d %d %s\n", attr.name, attr.typ, attr.offset, attr.stride, attr.buffer)
		buffer, ok := buffers[attr.buffer]
		if !ok {
			return nil, fmt.Errorf("missing buffer with name %s", attr.buffer)
		}
		glAttrIndex := glx.factory.Number(float64(attrIndex))
		attrType := attr.typ
		bufferType, bufferItemsPerVertex, err := attrType.asAttribute()
		if err != nil {
			return nil, fmt.Errorf("converting attribute type %s to attribute: %w", attrType, err)
		}
		glBufferType := glx.typeConverter.ToJs(bufferType)
		glItemsPerVertex := glx.factory.Number(float64(bufferItemsPerVertex))
		glNormalized := glx.factory.Boolean(false)
		glByteStride := glx.factory.Number(float64(attr.stride))
		glByteOffset := glx.factory.Number(float64(attr.offset))
		glx.constants.BindBuffer(glx.constants.ARRAY_BUFFER, buffer.glObject)
		glx.constants.VertexAttribPointer(
			glAttrIndex,
			glItemsPerVertex,
			glBufferType,
			glNormalized,
			glByteStride,
			glByteOffset,
		)
		glx.constants.EnableVertexAttribArray(glAttrIndex)
		attrMap[attr.name] = &vertexAttr{
			enabled: true,
			name:    attr.name,
			typ:     attr.typ,
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
