package raw

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/PieterD/warp/pkg/driver"
)

type glConstants struct {
	/* Unsorted. */

	VERTEX_SHADER            driver.Value
	FRAGMENT_SHADER          driver.Value
	COMPILE_STATUS           driver.Value
	LINK_STATUS              driver.Value
	DEPTH_TEST               driver.Value
	CreateShader             func(args ...driver.Value) driver.Value
	DeleteShader             func(args ...driver.Value) driver.Value
	ShaderSource             func(args ...driver.Value) driver.Value
	CompileShader            func(args ...driver.Value) driver.Value
	GetShaderParameter       func(args ...driver.Value) driver.Value
	GetShaderInfoLog         func(args ...driver.Value) driver.Value
	CreateProgram            func(args ...driver.Value) driver.Value
	DeleteProgram            func(args ...driver.Value) driver.Value
	AttachShader             func(args ...driver.Value) driver.Value
	LinkProgram              func(args ...driver.Value) driver.Value
	GetProgramParameter      func(args ...driver.Value) driver.Value
	GetProgramInfoLog        func(args ...driver.Value) driver.Value
	UseProgram               func(args ...driver.Value) driver.Value
	GetAttribLocation        func(args ...driver.Value) driver.Value
	GetUniformLocation       func(args ...driver.Value) driver.Value
	GetUniformBlockIndex     func(args ...driver.Value) driver.Value
	UniformBlockBinding      func(args ...driver.Value) driver.Value
	GetActiveAttrib          func(args ...driver.Value) driver.Value
	GetActiveUniform         func(args ...driver.Value) driver.Value
	CreateVertexArray        func(args ...driver.Value) driver.Value
	DeleteVertexArray        func(args ...driver.Value) driver.Value
	BindVertexArray          func(args ...driver.Value) driver.Value
	Uniform1i                func(args ...driver.Value) driver.Value
	Uniform1f                func(args ...driver.Value) driver.Value
	Uniform2f                func(args ...driver.Value) driver.Value
	Uniform3f                func(args ...driver.Value) driver.Value
	Uniform4f                func(args ...driver.Value) driver.Value
	UniformMatrix4fv         func(args ...driver.Value) driver.Value
	VertexAttribPointer      func(args ...driver.Value) driver.Value
	EnableVertexAttribArray  func(args ...driver.Value) driver.Value
	DisableVertexAttribArray func(args ...driver.Value) driver.Value
	ClearColor               func(args ...driver.Value) driver.Value
	Enable                   func(args ...driver.Value) driver.Value
	Disable                  func(args ...driver.Value) driver.Value
	Viewport                 func(args ...driver.Value) driver.Value

	/* Depth. */

	ALWAYS    driver.Value
	NEVER     driver.Value
	LESS      driver.Value
	LEQUAL    driver.Value
	GREATER   driver.Value
	GEQUAL    driver.Value
	NOTEQUAL  driver.Value
	DepthMask func(args ...driver.Value) driver.Value
	DepthFunc func(args ...driver.Value) driver.Value

	/* Parameters. */

	MAX_COMBINED_TEXTURE_IMAGE_UNITS driver.Value
	MAX_TEXTURE_SIZE  driver.Value
	GetParameter                     func(args ...driver.Value) driver.Value

	/* Clearing. */

	COLOR_BUFFER_BIT driver.Value
	DEPTH_BUFFER_BIT driver.Value
	Clear            func(args ...driver.Value) driver.Value

	/* Drawing. */

	POINTS       driver.Value
	LINES        driver.Value
	TRIANGLES    driver.Value
	DrawArrays   func(args ...driver.Value) driver.Value
	DrawElements func(args ...driver.Value) driver.Value

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

	STREAM_DRAW          driver.Value
	STREAM_READ          driver.Value
	STREAM_COPY          driver.Value
	STATIC_DRAW          driver.Value
	STATIC_READ          driver.Value
	STATIC_COPY          driver.Value
	DYNAMIC_DRAW         driver.Value
	DYNAMIC_READ         driver.Value
	DYNAMIC_COPY         driver.Value
	ARRAY_BUFFER         driver.Value
	ELEMENT_ARRAY_BUFFER driver.Value
	UNIFORM_BUFFER       driver.Value
	CreateBuffer         func(args ...driver.Value) driver.Value
	DeleteBuffer         func(args ...driver.Value) driver.Value
	BindBuffer           func(args ...driver.Value) driver.Value
	BufferData           func(args ...driver.Value) driver.Value
	BindBufferBase       func(args ...driver.Value) driver.Value
	BindBufferRange      func(args ...driver.Value) driver.Value

	/* Transform feedback */

	TRANSFORM_FEEDBACK        driver.Value
	INTERLEAVED_ATTRIBS       driver.Value
	SEPARATE_ATTRIBS          driver.Value
	CreateTransformFeedback   func(args ...driver.Value) driver.Value
	BindTransformFeedback     func(args ...driver.Value) driver.Value
	TransformFeedbackVaryings func(args ...driver.Value) driver.Value
	BeginTransformFeedback    func(args ...driver.Value) driver.Value
	EndTransformFeedback      func(args ...driver.Value) driver.Value

	/* Internal formats */

	R32F             driver.Value
	RG32F            driver.Value
	RGB32F           driver.Value
	RGBA32F          driver.Value
	DEPTH24_STENCIL8 driver.Value
	RGBA8            driver.Value

	/* Texture stuff */

	RGBA               driver.Value
	TEXTURE_2D         driver.Value
	TEXTURE_MIN_FILTER driver.Value
	TEXTURE_MAG_FILTER driver.Value
	NEAREST            driver.Value
	LINEAR             driver.Value
	TEXTURE_WRAP_S     driver.Value
	TEXTURE_WRAP_T     driver.Value
	CLAMP_TO_EDGE      driver.Value
	TEXTURE0           driver.Value
	ActiveTexture      func(args ...driver.Value) driver.Value
	CreateTexture      func(args ...driver.Value) driver.Value
	DeleteTexture      func(args ...driver.Value) driver.Value
	BindTexture        func(args ...driver.Value) driver.Value
	TexParameteri      func(args ...driver.Value) driver.Value
	TexImage2D         func(args ...driver.Value) driver.Value
	GenerateMipmap     func(args ...driver.Value) driver.Value

	/* Renderbuffer stuff */

	RENDERBUFFER                   driver.Value
	CreateRenderbuffer             func(args ...driver.Value) driver.Value
	DeleteRenderbuffer             func(args ...driver.Value) driver.Value
	BindRenderbuffer               func(args ...driver.Value) driver.Value
	RenderbufferStorage            func(args ...driver.Value) driver.Value
	RenderbufferStorageMultisample func(args ...driver.Value) driver.Value

	/* Framebuffer stuff */

	COLOR_ATTACHMENT0        driver.Value
	DEPTH_STENCIL_ATTACHMENT driver.Value
	FRAMEBUFFER              driver.Value
	FRAMEBUFFER_COMPLETE     driver.Value
	CreateFramebuffer        func(args ...driver.Value) driver.Value
	DeleteFramebuffer        func(args ...driver.Value) driver.Value
	BindFramebuffer          func(args ...driver.Value) driver.Value
	FramebufferRenderbuffer  func(args ...driver.Value) driver.Value
	CheckFramebufferStatus   func(args ...driver.Value) driver.Value
	ReadPixels               func(args ...driver.Value) driver.Value
}

func newGlConstants(obj driver.Object, trace bool) (c glConstants) {
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
			_, ok := value.ToFloat64()
			if !ok {
				panic(fmt.Errorf("constant %s is apparently not a value", field.Name))
			}
			//fmt.Printf("loading constant: %s = %d\n", field.Name, int(v))
			fieldValue.Set(reflect.ValueOf(value))
		case typeDriverFunc:
			functionName := strings.ToLower(field.Name[:1]) + field.Name[1:]
			//fmt.Printf("loading function: %s\n", functionName)
			function := driver.Bind(obj, functionName)
			if function == nil {
				panic(fmt.Errorf("function %s is apparently not a function", functionName))
			}
			if trace {
				function = wrapTrace(functionName, function)
			}
			fieldValue.Set(reflect.ValueOf(function))
		default:
			panic(fmt.Errorf("unhandled type: %v", fieldValue.Type()))
		}
	}
	return c
}

func wrapTrace(functionName string, f func(args ...driver.Value) driver.Value) func(args ...driver.Value) driver.Value {
	return func(args ...driver.Value) driver.Value {
		fmt.Printf("[TRACE] %s(", functionName)
		for i, arg := range args {
			isLastArg := i == len(args)-1
			if arg.IsNull() {
				fmt.Printf("null")
			} else if arg.IsUndefined() {
				fmt.Printf("undefined")
			} else if v, ok := arg.ToBoolean(); ok {
				fmt.Printf("%t", v)
			} else if v, ok := arg.ToFloat64(); ok {
				fmt.Printf("%f", v)
			} else if v, ok := arg.ToString(); ok {
				fmt.Printf("%q", v)
			} else if _, ok := arg.ToObject(); ok {
				fmt.Printf("object")
			} else if _, ok := arg.ToFunction(); ok {
				fmt.Printf("function")
			} else {
				fmt.Printf("UNKNOWN")
			}
			if !isLastArg {
				fmt.Printf(", ")
			}
		}
		fmt.Printf(")\n")
		rv := f(args...)
		return rv
	}

}
