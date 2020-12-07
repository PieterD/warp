package driver

type (
	Value interface {
		IsUndefined() (ok bool)
		IsNull() (ok bool)
		IsBoolean() (value, ok bool)
		IsNumber() (value float64, ok bool)
		IsString() (value string, ok bool)
		IsObject() (optionalValue Object)
		IsFunction() (optionalValue Function)
	}
	Object interface {
		Value
		Get(key string) Value
		Set(key string, value Value)
	}
	Function interface {
		Value
		New(args ...Value) Object
		Call(this Object, args ...Value) Value
	}
	Buffer interface {
		Size() int
		Put(data []byte) int
		Get(data []byte) int
		AsUint8Array() Object
		AsUint16Array() Object
		AsFloat32Array() Object
	}
	Factory interface {
		Global() Object
		Undefined() Value
		Null() Value
		Boolean(t bool) Value
		Number(f float64) Value
		String(s string) Value
		Function(f func(this Object, args ...Value) Value) Function
		Buffer(size int) Buffer
		Array(values ...Value) Object
	}
)
