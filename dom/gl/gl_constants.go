package gl

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/PieterD/warp/driver/driverutil"

	"github.com/PieterD/warp/driver"
)

type glConstants struct {
	VERTEX_SHADER        driver.Value
	FRAGMENT_SHADER      driver.Value
	COMPILE_STATUS       driver.Value
	LINK_STATUS          driver.Value
	ARRAY_BUFFER         driver.Value
	ELEMENT_ARRAY_BUFFER driver.Value
	STATIC_DRAW          driver.Value
	COLOR_BUFFER_BIT     driver.Value
	DEPTH_TEST           driver.Value
	TRIANGLES            driver.Value
	FLOAT                driver.Value
	FLOAT_VEC2           driver.Value
	FLOAT_VEC3           driver.Value
	FLOAT_VEC4           driver.Value
	FLOAT_MAT2           driver.Value
	FLOAT_MAT3           driver.Value
	FLOAT_MAT4           driver.Value
	BYTE                 driver.Value
	UNSIGNED_BYTE        driver.Value
	SHORT                driver.Value
	UNSIGNED_SHORT       driver.Value
	INT                  driver.Value
	UNSIGNED_INT         driver.Value
}

type glFunctions struct {
	CreateBuffer            func(args ...driver.Value) driver.Value
	BindBuffer              func(args ...driver.Value) driver.Value
	BufferData              func(args ...driver.Value) driver.Value
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
	Uniform1f               func(args ...driver.Value) driver.Value
	VertexAttribPointer     func(args ...driver.Value) driver.Value
	EnableVertexAttribArray func(args ...driver.Value) driver.Value
	ClearColor              func(args ...driver.Value) driver.Value
	Enable                  func(args ...driver.Value) driver.Value
	Clear                   func(args ...driver.Value) driver.Value
	Viewport                func(args ...driver.Value) driver.Value
	DrawArrays              func(args ...driver.Value) driver.Value
	DrawElements            func(args ...driver.Value) driver.Value
}

func newGlConstants(obj driver.Object) (c glConstants) {
	v := reflect.ValueOf(&c).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		fieldValue := v.Field(i)
		field := t.Field(i)
		value := obj.Get(field.Name)
		{
			v, _ := value.IsNumber()
			fmt.Printf("loading constant: %s = %d\n", field.Name, int(v))
		}
		fieldValue.Set(reflect.ValueOf(value))
	}
	return c
}

func newGlFunctions(obj driver.Object) (f glFunctions) {
	v := reflect.ValueOf(&f).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		fieldValue := v.Field(i)
		field := t.Field(i)
		functionName := strings.ToLower(field.Name[:1]) + field.Name[1:]
		fmt.Printf("loading function: %s\n", functionName)
		function := driverutil.Bind(obj, functionName)
		if function == nil {
			panic(fmt.Errorf("function %s is apparently not a function", functionName))
		}
		fieldValue.Set(reflect.ValueOf(function))
	}
	return f
}
