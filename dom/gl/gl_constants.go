package gl

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/PieterD/warp/driver"
	"github.com/PieterD/warp/driver/driverutil"
)

type glConstants struct {
	/* Unsorted. */

	VERTEX_SHADER           driver.Value
	FRAGMENT_SHADER         driver.Value
	COMPILE_STATUS          driver.Value
	LINK_STATUS             driver.Value
	ARRAY_BUFFER            driver.Value
	ELEMENT_ARRAY_BUFFER    driver.Value
	STATIC_DRAW             driver.Value
	COLOR_BUFFER_BIT        driver.Value
	DEPTH_TEST              driver.Value
	CreateShader            func(args ...driver.Value) driver.Value
	ShaderSource            func(args ...driver.Value) driver.Value
	CompileShader           func(args ...driver.Value) driver.Value
	GetShaderParameter      func(args ...driver.Value) driver.Value
	GetShaderInfoLog        func(args ...driver.Value) driver.Value
	CreateProgram           func(args ...driver.Value) driver.Value
	AttachShader            func(args ...driver.Value) driver.Value
	LinkProgram             func(args ...driver.Value) driver.Value
	GetProgramParameter     func(args ...driver.Value) driver.Value
	GetProgramInfoLog       func(args ...driver.Value) driver.Value
	UseProgram              func(args ...driver.Value) driver.Value
	GetAttribLocation       func(args ...driver.Value) driver.Value
	GetUniformLocation      func(args ...driver.Value) driver.Value
	GetActiveAttrib         func(args ...driver.Value) driver.Value
	GetActiveUniform        func(args ...driver.Value) driver.Value
	CreateVertexArray       func(args ...driver.Value) driver.Value
	BindVertexArray         func(args ...driver.Value) driver.Value
	Uniform1f               func(args ...driver.Value) driver.Value
	UniformMatrix4fv        func(args ...driver.Value) driver.Value
	VertexAttribPointer     func(args ...driver.Value) driver.Value
	EnableVertexAttribArray func(args ...driver.Value) driver.Value
	ClearColor              func(args ...driver.Value) driver.Value
	Enable                  func(args ...driver.Value) driver.Value
	Clear                   func(args ...driver.Value) driver.Value
	Viewport                func(args ...driver.Value) driver.Value
	DrawArrays              func(args ...driver.Value) driver.Value
	DrawElements            func(args ...driver.Value) driver.Value

	/* Drawn modes. */

	TRIANGLES driver.Value

	/* Data types. */

	FLOAT          driver.Value
	FLOAT_VEC2     driver.Value
	FLOAT_VEC3     driver.Value
	FLOAT_VEC4     driver.Value
	FLOAT_MAT2     driver.Value
	FLOAT_MAT3     driver.Value
	FLOAT_MAT4     driver.Value
	BYTE           driver.Value
	UNSIGNED_BYTE  driver.Value
	SHORT          driver.Value
	UNSIGNED_SHORT driver.Value
	INT            driver.Value
	UNSIGNED_INT   driver.Value

	/* Buffer stuff */

	CreateBuffer func(args ...driver.Value) driver.Value
	BindBuffer   func(args ...driver.Value) driver.Value
	BufferData   func(args ...driver.Value) driver.Value

	/* Texture stuff */

	RGBA               driver.Value
	TEXTURE_2D         driver.Value
	TEXTURE_MIN_FILTER driver.Value
	TEXTURE_MAG_FILTER driver.Value
	NEAREST            driver.Value
	TEXTURE_WRAP_S     driver.Value
	TEXTURE_WRAP_T     driver.Value
	CLAMP_TO_EDGE      driver.Value
	CreateTexture      func(args ...driver.Value) driver.Value
	BindTexture        func(args ...driver.Value) driver.Value
	TexParameteri      func(args ...driver.Value) driver.Value
	TexImage2D         func(args ...driver.Value) driver.Value
}

func newGlConstants(obj driver.Object) (c glConstants) {
	var driverValue driver.Value
	var driverFunc func(args ...driver.Value) driver.Value
	typeDriverValue := reflect.TypeOf(&driverValue).Elem()
	typeDriverFunc := reflect.TypeOf(&driverFunc).Elem()
	v := reflect.ValueOf(&c).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		fieldValue := v.Field(i)
		field := t.Field(i)
		value := obj.Get(field.Name)
		switch fieldValue.Type() {
		case typeDriverValue:
			v, _ := value.IsNumber()
			fmt.Printf("loading constant: %s = %d\n", field.Name, int(v))
			fieldValue.Set(reflect.ValueOf(value))
		case typeDriverFunc:
			functionName := strings.ToLower(field.Name[:1]) + field.Name[1:]
			fmt.Printf("loading function: %s\n", functionName)
			function := driverutil.Bind(obj, functionName)
			if function == nil {
				panic(fmt.Errorf("function %s is apparently not a function", functionName))
			}
			fieldValue.Set(reflect.ValueOf(function))
		default:
			panic(fmt.Errorf("unhandled type: %v", fieldValue.Type()))
		}
	}
	return c
}
