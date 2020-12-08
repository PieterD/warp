package driver

type (
	Value interface {
		IsUndefined() (ok bool)
		IsNull() (ok bool)
		
		ToBoolean() (value, ok bool)
		ToFloat64() (value float64, ok bool)
		ToString() (value string, ok bool)
		ToObject() (optionalValue Object, ok bool)
		ToFunction() (optionalValue Function, ok bool)
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
