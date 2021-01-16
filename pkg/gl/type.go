package gl

import (
	"fmt"
	"strings"

	"github.com/PieterD/warp/pkg/driver"
)

//go:generate stringer -type=Type
type Type int

const (
	Float Type = iota + 1
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

func (t Type) glSize() int {
	switch t {
	case Byte, UnsignedByte:
		return 1
	case Short, UnsignedShort:
		return 2
	case Int, UnsignedInt, Float:
		return 4
	case Vec2:
		return 2 * 4
	case Vec3:
		return 3 * 4
	case Vec4:
		return 4 * 4
	case Mat2:
		return 2 * 2 * 4
	case Mat3:
		return 3 * 3 * 4
	case Mat4:
		return 4 * 4 * 4
	default:
		panic(fmt.Errorf("invalid Type: %v", t))
	}
}

func (t Type) GLSL() string {
	switch t {
	case Vec2, Vec3, Vec4, Mat2, Mat3, Mat4:
		return strings.ToLower(t.String())
	default:
		panic(fmt.Errorf("unimplemented %s.GLSL", t.String()))
	}
}

func (t Type) asAttribute() (bufferType Type, itemsPerVertex int, err error) {
	switch t {
	case Float:
		return Float, 1, nil
	case Vec2:
		return Float, 2, nil
	case Vec3:
		return Float, 3, nil
	case Vec4:
		return Float, 4, nil
	default:
		return 0, 0, fmt.Errorf("unable to decompose: %s", t)
	}
}

//TODO: replace this with Type.glType()
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
		fNum, ok := v.ToFloat64()
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
	fValue, ok := value.ToFloat64()
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
