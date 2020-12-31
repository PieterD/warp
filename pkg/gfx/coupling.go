package gfx

import (
	"fmt"
	"io"
)

type (
	//TODO: transforming []byte data from one DataCoupling to another
	// This will greatly simplify our implementation of std140
	DataCoupling struct {
		attributes    []vertexAttribute
		attrByName    map[string]int
		attrByIndex   map[int]int
		attrsByBuffer map[string][]int
	}
	ActiveCoupling struct {
		dc      *DataCoupling
		enabled map[string]struct{}
	}
	vertexAttribute struct {
		name    string
		typ     Type
		buffer  string
		padding int
		index   int
		offset  int
		stride  int
	}
	DataCouplingConfig struct {
		Vertices         []VertexConfig
		WithElements     bool
		ElementArrayType Type
	}
	VertexConfig struct {
		Name    string
		Type    Type
		Buffer  string
		Padding int
	}
)

func MustNewDataCoupling(config DataCouplingConfig) *DataCoupling {
	dc, err := NewDataCoupling(config)
	if err != nil {
		panic(fmt.Errorf("creating new DataCoupling: %w", err))
	}
	return dc
}

func NewDataCoupling(config DataCouplingConfig) (*DataCoupling, error) {
	vertexConfigs := config.Vertices
	if len(vertexConfigs) == 0 {
		return nil, fmt.Errorf("no vertex configurations provided")
	}
	dc := &DataCoupling{
		attributes:    make([]vertexAttribute, len(vertexConfigs)),
		attrByName:    make(map[string]int),
		attrByIndex:   make(map[int]int),
		attrsByBuffer: make(map[string][]int),
	}
	index := 0
	for i, vertexConfig := range vertexConfigs {
		if vertexConfig.Name == "" {
			return nil, fmt.Errorf("vertex config at index %d has an empty name", i)
		}
		if vertexConfig.Buffer == "" {
			return nil, fmt.Errorf("attribute %s has an empty buffer name", vertexConfig.Name)
		}
		if vertexConfig.Padding < 0 {
			return nil, fmt.Errorf("attribute %s has negative padding", vertexConfig.Name)
		}
		switch vertexConfig.Type {
		case Float, Vec2, Vec3, Vec4, Byte, UnsignedByte, Short, UnsignedShort, Int, UnsignedInt:
		default:
			return nil, fmt.Errorf("unsupported attribute type: %v", vertexConfig.Type)
		}
		if _, ok := dc.attrByName[vertexConfig.Name]; ok {
			return nil, fmt.Errorf("attribute name appears multiple times: %s", vertexConfig.Name)
		}
		dc.attributes[i] = vertexAttribute{
			name:    vertexConfig.Name,
			typ:     vertexConfig.Type,
			buffer:  vertexConfig.Buffer,
			padding: vertexConfig.Padding,
			index:   index,
			offset:  0,
			stride:  0,
		}
		dc.attrByName[vertexConfig.Name] = i
		dc.attrByIndex[index] = i
		dc.attrsByBuffer[vertexConfig.Buffer] = append(dc.attrsByBuffer[vertexConfig.Buffer], i)

		//TODO: I believe larger datatypes require larger index jumps.
		index++
	}
	for _, attrs := range dc.attrsByBuffer {
		stride := 0
		for _, iAttr := range attrs {
			attr := &dc.attributes[iAttr]
			attr.offset = stride
			stride += attr.typ.glSize() + attr.padding
		}
		for _, iAttr := range attrs {
			attr := &dc.attributes[iAttr]
			attr.stride = stride
		}
	}
	return dc, nil
}

func (dc *DataCoupling) Allocate(vertices int) (map[string]*Buffer, error) {
	panic("not implemented")
}

func (dc *DataCoupling) Upload(byteOffset int, data map[string][]byte) error {
	panic("not implemented")
}

func (dc *DataCoupling) Active(enabledAttrNames ...string) ActiveCoupling {
	enabled := make(map[string]struct{})
	for _, attrName := range enabledAttrNames {
		enabled[attrName] = struct{}{}
	}
	return ActiveCoupling{
		dc:      dc,
		enabled: enabled,
	}
}

// DataCouplingConfigFromStruct will return a DataCouplingConfig corresponding to
// the memory layout of rawStruct.
// Together with Map, this can be used to create generators.
func DataCouplingConfigFromStruct(rawStruct interface{}) (DataCouplingConfig, error) {
	panic("not implemented")
}

func (dc *DataCoupling) Translate(to *DataCoupling, source io.Reader) io.Reader {
	panic("not implemented")
}

func (adc ActiveCoupling) BufferNames() (bufferNames []string) {
	bufSet := make(map[string]struct{})
	for _, attr := range adc.dc.attributes {
		if _, ok := adc.enabled[attr.name]; !ok {
			continue
		}
		if _, ok := bufSet[attr.buffer]; ok {
			continue
		}
		bufSet[attr.buffer] = struct{}{}
		bufferNames = append(bufferNames, attr.name)
	}
	return bufferNames
}

func (adc ActiveCoupling) AttrNum() int {
	return len(adc.enabled)
}
