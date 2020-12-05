package gl

import (
	"fmt"
	"sort"

	"github.com/PieterD/warp/pkg/driver"
)

type VertexArrayAttribute struct {
	Name   string
	Type   Type
	Buffer *Buffer
	Layout VertexArrayAttributeLayout
}

type VertexArrayAttributeLayout struct {
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

func newVertexArray(glx *Context, attrs ...VertexArrayAttribute) (*VertexArray, error) {
	glVAO := glx.constants.CreateVertexArray()
	glx.constants.BindVertexArray(glVAO)
	defer glx.constants.BindVertexArray(glx.factory.Null())

	attrMap := make(map[string]*vertexAttr)
	for attrIndex, attr := range attrs {
		glAttrIndex := glx.factory.Number(float64(attrIndex))
		attrType := attr.Type
		bufferType, bufferItemsPerVertex, err := attrType.asAttribute()
		if err != nil {
			return nil, fmt.Errorf("converting attribute type %s to attribute: %w", attrType, err)
		}
		glBufferType := glx.typeConverter.ToJs(bufferType)
		glItemsPerVertex := glx.factory.Number(float64(bufferItemsPerVertex))
		glNormalized := glx.factory.Boolean(false)
		glByteStride := glx.factory.Number(float64(attr.Layout.ByteStride))
		glByteOffset := glx.factory.Number(float64(attr.Layout.ByteOffset))
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

func (vao *VertexArray) Attributes() (attrs []AttributeDescription) {
	//indexCounts := make(map[int]struct{})
	for _, attr := range vao.attrs {
		if !attr.enabled {
			continue
		}
		//TODO: move this to Verify
		//if _, ok := indexCounts[attr.index]; ok {
		//	return nil, fmt.Errorf("index %d set by multiple attributes", attr.index)
		//}
		//indexCounts[attr.index] = struct{}{}
		attrs = append(attrs, AttributeDescription{
			Name:  attr.name,
			Type:  attr.typ,
			Index: attr.index,
		})
	}
	sort.Slice(attrs, func(i, j int) bool {
		return attrs[i].Index < attrs[j].Index
	})
	return attrs
}
