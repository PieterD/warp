package glutil

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

/*
Scalar bool, int, uint, float and double
	Both the size and alignment are the size of the scalar in basic machine types
Two-component vectors (e.g., ivec2)
	Both the size and alignment are twice the size of the underlying scalar type.
Three-component vectors (e.g., vec3) and Four-component vectors (e.g., vec4)
	Both the size and alignment are four times the size of the underlying scalar type.
An array of scalars or vectors
	The size of each element in the array will be the size of the element type,
	rounded up to a multiple of the size of a vec4. This is also the array’s alignment.
	The array’s size will be this rounded-up element’s size times the number of elements in the array.
A column-major matrix or an array of column-major matrices of size C columns and R rows
	Same layout as an array of N vectors each with R components, where N is the total number of columns present.
A row-major matrix or an array of row-major matrices with R rows and C columns
	Same layout as an array of N vectors each with C components, where N is the total number of rows present.
A single-structure definition, or an array of structures
	Structure alignment will be the alignment for the biggest structure member,
	according to the previous rules, rounded up to a multiple of the size of a vec4.
	Each structure will start on this alignment, and its size will be the space needed by its members,
	according to the previous rules, rounded up to a multiple of the structure alignment. (e.g., sizeof(GLfloat)).
*/

func Std140Uniform(raw interface{}) (string, error) {
	rawValue := reflect.ValueOf(raw)
	if kind := rawValue.Kind(); kind != reflect.Ptr {
		return "", fmt.Errorf("data value does not have pointer kind: %s", kind)
	}
	if rawValue.IsNil() {
		return "", fmt.Errorf("data value pointer is nil")
	}
	rawValue = rawValue.Elem()
	rawType := rawValue.Type()
	if kind := rawValue.Kind(); kind != reflect.Struct {
		return "", fmt.Errorf("data value does not have struct kind: %s", kind)
	}

	buf := &bytes.Buffer{}
	buf.WriteString("layout (std140) uniform Uniform {\n")
	for fieldNum := 0; fieldNum < rawType.NumField(); fieldNum++ {
		field := rawType.Field(fieldNum)
		fieldValue := rawValue.Field(fieldNum)
		typeOverride := field.Tag.Get("std140")
		if typeOverride == "-" {
			continue
		}
		if err := putUniformRow(buf, field.Name, typeOverride, fieldValue); err != nil {
			return "", fmt.Errorf("field %s: %w", field.Name, err)
		}
	}
	buf.WriteString("}Uniforms;\n")
	return buf.String(), nil
}

func putUniformRow(buf *bytes.Buffer, fieldName, typeOverride string, fieldValue reflect.Value) error {
	fieldType := fieldValue.Type()
	switch fieldType.Kind() {
	case reflect.Int:
		i := fieldValue.Int()
		if i > math.MaxInt32 {
			return fmt.Errorf("int over int32 value: %d", i)
		}
		return putUniformRow(buf, fieldName, typeOverride, reflect.ValueOf(int32(i)))
	case reflect.Uint32:
		ts := "uint"
		if typeOverride != "" {
			ts = typeOverride
		}
		buf.WriteString(fmt.Sprintf("	%s %s;\n", ts, fieldName))
		return nil
	case reflect.Int32:
		ts := "int"
		if typeOverride != "" {
			ts = typeOverride
		}
		buf.WriteString(fmt.Sprintf("	%s %s;\n", ts, fieldName))
		return nil
	case reflect.Float32:
		ts := "float"
		if typeOverride != "" {
			ts = typeOverride
		}
		buf.WriteString(fmt.Sprintf("	%s %s;\n", ts, fieldName))
		return nil
	case reflect.Array:
		elemType := fieldType.Elem()
		elemNum := fieldType.Len()
		if elemType.Kind() != reflect.Float32 {
			return fmt.Errorf("array of non-float type: %s", elemType)
		}
		switch elemNum {
		case 2:
			ts := "vec2"
			if typeOverride != "" {
				ts = typeOverride
			}
			buf.WriteString(fmt.Sprintf("	%s %s;\n", ts, fieldName))
			return nil
		case 3:
			ts := "vec3"
			if typeOverride != "" {
				ts = typeOverride
			}
			buf.WriteString(fmt.Sprintf("	%s %s;\n", ts, fieldName))
			return nil
		case 4:
			ts := "vec4"
			if typeOverride != "" {
				ts = typeOverride
			}
			buf.WriteString(fmt.Sprintf("	%s %s;\n", ts, fieldName))
			return nil
		case 16:
			ts := "mat4"
			if typeOverride != "" {
				ts = typeOverride
			}
			buf.WriteString(fmt.Sprintf("	%s %s;\n", ts, fieldName))
			return nil
		default:
			return fmt.Errorf("field %s: array of unhandled number of floats: %d", fieldName, elemNum)
		}
	default:
		return fmt.Errorf("field %s: unhandled type: %s", fieldName, fieldType)
	}
}

func Std140Data(raw interface{}) ([]byte, error) {
	rawValue := reflect.ValueOf(raw)
	if kind := rawValue.Kind(); kind != reflect.Ptr {
		return nil, fmt.Errorf("data value does not have pointer kind: %s", kind)
	}
	if rawValue.IsNil() {
		return nil, fmt.Errorf("data value pointer is nil")
	}
	rawValue = rawValue.Elem()
	rawType := rawValue.Type()
	if kind := rawValue.Kind(); kind != reflect.Struct {
		return nil, fmt.Errorf("data value does not have struct kind: %s", kind)
	}
	buf := &bytes.Buffer{}
	for fieldNum := 0; fieldNum < rawType.NumField(); fieldNum++ {
		field := rawType.Field(fieldNum)
		fieldValue := rawValue.Field(fieldNum)
		if err := putDataRow(buf, field.Name, fieldValue); err != nil {
			return nil, fmt.Errorf("field %s: %w", field.Name, err)
		}
	}
	return buf.Bytes(), nil
}

func putDataRow(buf io.Writer, fieldName string, fieldValue reflect.Value) error {
	fieldType := fieldValue.Type()
	switch fieldType.Kind() {
	case reflect.Uint32:
		var b [4]byte
		binary.LittleEndian.PutUint32(b[:], uint32(fieldValue.Uint()))
		if _, err := buf.Write(b[:]); err != nil {
			return fmt.Errorf("writing to buffer: %w", err)
		}
		return nil
	case reflect.Int32:
		ui := uint32(fieldValue.Int())
		if err := putDataRow(buf, fieldName, reflect.ValueOf(ui)); err != nil {
			return fmt.Errorf("recurse(float32): %w", err)
		}
		return nil
	case reflect.Int:
		i := fieldValue.Int()
		if i > math.MaxInt32 {
			return fmt.Errorf("int over int32 value: %d", i)
		}
		if err := putDataRow(buf, fieldName, reflect.ValueOf(int32(i))); err != nil {
			return fmt.Errorf("recurse(float32): %w", err)
		}
		return nil
	case reflect.Float32:
		bits := math.Float32bits(float32(fieldValue.Float()))
		if err := putDataRow(buf, fieldName, reflect.ValueOf(bits)); err != nil {
			return fmt.Errorf("recurse(float32): %w", err)
		}
		return nil
	case reflect.Array:
		elemType := fieldType.Elem()
		elemNum := fieldType.Len()
		if elemType.Kind() != reflect.Float32 {
			return fmt.Errorf("array of non-float type: %s", elemType)
		}
		var padding []byte
		switch elemNum {
		case 2:
			// vec2
		case 3:
			// vec3
			padding = make([]byte, 4)
		case 4:
			// vec4
		case 16:
			// mat4
		default:
			return fmt.Errorf("array of unhandled number of floats: %d", elemNum)
		}
		for i := 0; i < elemNum; i++ {
			if err := putDataRow(buf, fmt.Sprintf("%s[%d]", fieldName, i), fieldValue.Index(i)); err != nil {
				return fmt.Errorf("recurse: %w", err)
			}
		}
		if _, err := buf.Write(padding); err != nil {
			return fmt.Errorf("writing padding: %w", err)
		}
	default:
		return fmt.Errorf("unhandled type: %s", fieldType)
	}
	return nil
}
