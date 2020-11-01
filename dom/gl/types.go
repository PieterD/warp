package gl

import (
	"fmt"

	"github.com/PieterD/warp/driver"
)

//go:generate stringer -type=Type
type Type int

const (
	Float Type = iota
	Vec2
	Vec3
	Vec4
	Mat2
	Mat3
	Mat4
	Byte
	UnsignedByte
	Short
	UnsignedShort
	Int
	UnsignedInt
)

/*
GL_FLOAT_MAT2x3
GL_FLOAT_MAT2x4
GL_FLOAT_MAT3x2
GL_FLOAT_MAT3x4
GL_FLOAT_MAT4x2
GL_FLOAT_MAT4x3
GL_INT_VEC2
GL_INT_VEC3
GL_INT_VEC4
GL_UNSIGNED_INT_VEC2
GL_UNSIGNED_INT_VEC3
GL_UNSIGNED_INT_VEC4
*/

type typeConverter struct {
	jsConstants map[Type]driver.Value
	reverse     map[int]Type
}

func newTypeConverter(constants glConstants) *typeConverter {
	jsConstants := map[Type]driver.Value{
		Float:         constants.FLOAT,
		Vec2:          constants.FLOAT_VEC2,
		Vec3:          constants.FLOAT_VEC3,
		Vec4:          constants.FLOAT_VEC4,
		Mat2:          constants.FLOAT_MAT2,
		Mat3:          constants.FLOAT_MAT3,
		Mat4:          constants.FLOAT_MAT4,
		Byte:          constants.BYTE,
		UnsignedByte:  constants.UNSIGNED_BYTE,
		Short:         constants.SHORT,
		UnsignedShort: constants.UNSIGNED_SHORT,
		Int:           constants.INT,
		UnsignedInt:   constants.UNSIGNED_INT,
	}
	reverse := make(map[int]Type)
	for typ, v := range jsConstants {
		fNum, ok := v.IsNumber()
		if !ok {
			panic(fmt.Errorf("js constant for %s not a Number: %T", typ, v))
		}
		iNum := int(fNum)
		if _, ok := reverse[iNum]; ok {
			panic(fmt.Errorf("js constant for %s is double booked: %d", typ, iNum))
		}
		reverse[iNum] = typ
	}
	return &typeConverter{
		jsConstants: jsConstants,
		reverse:     reverse,
	}
}

func (tc *typeConverter) FromJs(value driver.Value) (Type, error) {
	fValue, ok := value.IsNumber()
	if !ok {
		return 0, fmt.Errorf("expected number value: %T", value)
	}
	iValue := int(fValue)
	typ, ok := tc.reverse[iValue]
	if !ok {
		return 0, fmt.Errorf("value %d not found in reverse", iValue)
	}
	return typ, nil
}

func (tc *typeConverter) ToJs(typ Type) driver.Value {
	v, ok := tc.jsConstants[typ]
	if !ok {
		panic(fmt.Errorf("invalid type: %d", typ))
	}
	return v
}
