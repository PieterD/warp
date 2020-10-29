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
	Factory interface {
		Undefined() Value
		Null() Value
		Number(f float64) Value
		String(s string) Value
		Function(f func(this Object, args ...Value) Value) Function
	}
)
