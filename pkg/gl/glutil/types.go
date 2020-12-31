package glutil

import (
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

type value interface {
	IsStruct() (tStruct, bool)
	IsSlice() (tSlice, bool)
	IsInt8() (tInt8, bool)
	IsUint8() (tUint8, bool)
	IsInt16() (tInt16, bool)
	IsUint16() (tUint16, bool)
	IsInt32() (tInt32, bool)
	IsUint32() (tUint32, bool)
	IsFloat32() (tFloat32, bool)
	IsVec3() (tVec3, bool)
	IsVec4() (tVec4, bool)
	IsMat3() (tMat3, bool)
	IsMat4() (tMat4, bool)
	Size() (size int)
	Marshal(w io.Writer) error
}

type tEmpty struct{}

func (t tEmpty) IsInt8() (tInt8, bool) {
	return tInt8{}, false
}

func (t tEmpty) IsUint8() (tUint8, bool) {
	return tUint8{}, false
}

func (t tEmpty) IsInt16() (tInt16, bool) {
	return tInt16{}, false
}

func (t tEmpty) IsUint16() (tUint16, bool) {
	return tUint16{}, false
}

func (t tEmpty) IsInt32() (tInt32, bool) {
	return tInt32{}, false
}

func (t tEmpty) IsUint32() (tUint32, bool) {
	return tUint32{}, false
}

func (t tEmpty) IsStruct() (tStruct, bool) {
	return tStruct{}, false
}

func (t tEmpty) IsSlice() (tSlice, bool) {
	return tSlice{}, false
}

func (t tEmpty) IsFloat32() (tFloat32, bool) {
	return tFloat32{}, false
}

func (t tEmpty) IsVec3() (tVec3, bool) {
	return tVec3{}, false
}

func (t tEmpty) IsVec4() (tVec4, bool) {
	return tVec4{}, false
}

func (t tEmpty) IsMat3() (tMat3, bool) {
	return tMat3{}, false
}

func (t tEmpty) IsMat4() (tMat4, bool) {
	return tMat4{}, false
}

func (t tEmpty) Marshal(_ io.Writer) error {
	return fmt.Errorf("empty type does not implement Marshal")
}

type tStruct struct {
	tEmpty
	name   string
	fields []tStructElement
}

func (t tStruct) Size() (size int) {
	for _, field := range t.fields {
		size += field.value.Size()
	}
	return size
}

func (t tStruct) IsStruct() (tStruct, bool) {
	return t, true
}

var _ value = tStruct{}

type tStructElement struct {
	name  string
	value value
}

type tSlice struct {
	tEmpty
	hasStructElem bool
	value    []tSliceElement
}

func (t tSlice) Size() (size int) {
	if t.hasStructElem {
		panic(fmt.Errorf("not yet implemented"))
	}
	for _, elem := range t.value {
		size += elem.value.Size()
		if mod := size % 4 * 4; mod > 0 {
			size += 4*4 - mod
		}
	}
	return size
}

func (t tSlice) IsSlice() (tSlice, bool) {
	return t, true
}

var _ value = tSlice{}

type tSliceElement struct {
	index int
	value value
}

type tInt8 struct {
	tEmpty
	p *int8
}

func (t tInt8) Size() (size int) {
	return 1
}

func (t tInt8) IsInt8() (tInt8, bool) {
	return t, true
}

var _ value = tInt8{}

type tUint8 struct {
	tEmpty
	p *uint8
}

func (t tUint8) Size() (size int) {
	return 1
}

func (t tUint8) IsUint8() (tUint8, bool) {
	return t, true
}

var _ value = tUint8{}

type tInt16 struct {
	tEmpty
	p *int16
}

func (t tInt16) Size() (size int) {
	return 2
}

func (t tInt16) IsInt16() (tInt16, bool) {
	return t, true
}

var _ value = tInt16{}

type tUint16 struct {
	tEmpty
	p *uint16
}

func (t tUint16) Size() (size int) {
	return 2
}

func (t tUint16) IsUint16() (tUint16, bool) {
	return t, true
}

var _ value = tUint16{}

type tInt32 struct {
	tEmpty
	p *int32
}

func (t tInt32) Size() (size int) {
	return 4
}

func (t tInt32) IsInt32() (tInt32, bool) {
	return t, true
}

var _ value = tInt32{}

type tUint32 struct {
	tEmpty
	p *uint32
}

func (t tUint32) Size() (size int) {
	return 4
}

func (t tUint32) IsUint32() (tUint32, bool) {
	return t, true
}

var _ value = tUint32{}

type tFloat32 struct {
	tEmpty
	p *float32
}

func (t tFloat32) Size() (size int) {
	return 4
}

func (t tFloat32) IsFloat32() (tFloat32, bool) {
	return t, true
}

var _ value = tFloat32{}

type tVec3 struct {
	tEmpty
	p *[3]float32
}

func (t tVec3) Size() (size int) {
	return 4 * 4
}

func (t tVec3) IsVec3() (tVec3, bool) {
	return t, true
}

var _ value = tVec3{}

type tVec4 struct {
	tEmpty
	p *[4]float32
}

func (t tVec4) Size() (size int) {
	return 4 * 4
}

func (t tVec4) IsVec4() (tVec4, bool) {
	return t, true
}

var _ value = tVec4{}

type tMat3 struct {
	tEmpty
	p *[9]float32
}

func (t tMat3) Size() (size int) {
	return (3*3 + /*padding*/ 3 /*padding*/) * 4
}

func (t tMat3) IsMat3() (tMat3, bool) {
	return t, true
}

var _ value = tMat3{}

type tMat4 struct {
	tEmpty
	p *[16]float32
}

func (t tMat4) Size() (size int) {
	return 4 * 4 * 4
}

func (t tMat4) IsMat4() (tMat4, bool) {
	return t, true
}

var _ value = tMat4{}

func build(raw interface{}) (value, error) {
	rawValue := reflect.ValueOf(raw)
	if kind := rawValue.Kind(); kind != reflect.Ptr {
		return nil, fmt.Errorf("data value does not have pointer kind: %s", kind)
	}
	if rawValue.IsNil() {
		return nil, fmt.Errorf("data value pointer is nil")
	}
	rawValue = rawValue.Elem()
	if kind := rawValue.Kind(); kind != reflect.Struct {
		return nil, fmt.Errorf("data value does not have struct kind: %s", kind)
	}
	v, err := rBuild(rawValue)
	if err != nil {
		return nil, fmt.Errorf("building %s: %w", rawValue.Type(), err)
	}
	return v, nil
}

func rBuild(rawValue reflect.Value) (value, error) {
	ptr := unsafe.Pointer(rawValue.Addr().Pointer())
	rawType := rawValue.Type()
	switch rawType.Kind() {
	default:
		return nil, fmt.Errorf("unhandled kind: %s", rawType)
	case reflect.Int8:
		return tInt8{
			p: (*int8)(ptr),
		}, nil
	case reflect.Uint8:
		return tUint8{
			p: (*uint8)(ptr),
		}, nil
	case reflect.Int16:
		return tInt16{
			p: (*int16)(ptr),
		}, nil
	case reflect.Uint16:
		return tUint16{
			p: (*uint16)(ptr),
		}, nil
	case reflect.Int32:
		return tInt32{
			p: (*int32)(ptr),
		}, nil
	case reflect.Uint32:
		return tUint32{
			p: (*uint32)(ptr),
		}, nil
	case reflect.Float32:
		return tFloat32{
			p: (*float32)(ptr),
		}, nil
	case reflect.Struct:
		var fields []tStructElement
		for fieldNum := 0; fieldNum < rawType.NumField(); fieldNum++ {
			field := rawType.Field(fieldNum)
			fieldValue := rawValue.Field(fieldNum)
			value, err := rBuild(fieldValue)
			if err != nil {
				return nil, fmt.Errorf("building %s: %w", field.Name, err)
			}
			fields = append(fields, tStructElement{
				name:  field.Name,
				value: value,
			})
		}
		return tStruct{
			name:   rawType.Name(),
			fields: fields,
		}, nil
	case reflect.Array:
		switch rawType.Elem().Kind() {
		default:
			return nil, fmt.Errorf("unhandled array element kind: %s", rawType.Elem().Kind())
		case reflect.Float32:
			switch rawType.Len() {
			default:
				return nil, fmt.Errorf("unhandled float32 array length: %d", rawType.Len())
			case 3:
				return tVec3{
					p: (*[3]float32)(ptr),
				}, nil
			case 4:
				return tVec4{
					p: (*[4]float32)(ptr),
				}, nil
			case 9:
				return tMat3{
					p: (*[9]float32)(ptr),
				}, nil
			case 16:
				return tMat4{
					p: (*[16]float32)(ptr),
				}, nil
			}
		}
	}
}
