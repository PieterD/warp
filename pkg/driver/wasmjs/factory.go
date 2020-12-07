package wasmjs

import (
	"fmt"
	"syscall/js"

	"github.com/PieterD/warp/pkg/driver"
)

type vValue interface {
	jsValue() js.Value
	driver.Value
}

type jsFactory struct{}

func (j jsFactory) Array(values ...driver.Value) driver.Object {
	var jsValues []interface{}
	for _, value := range values {
		vv, ok := value.(vValue)
		if !ok {
			panic(fmt.Errorf("non-wasmjs value: %T %v", value, value))
		}
		jsValues = append(jsValues, vv.jsValue())
	}
	jsArrayObject := js.Global().Get("Array").New(jsValues...)
	return jsObject{
		v: jsArrayObject,
	}
}

func (j jsFactory) Boolean(t bool) driver.Value {
	return jsBoolean{
		v: js.ValueOf(t),
	}
}

func (j jsFactory) Buffer(size int) driver.Buffer {
	return newBuffer(j, size)
}

func (j jsFactory) Global() driver.Object {
	return jsObject{
		v: js.Global(),
	}
}

func (j jsFactory) Undefined() driver.Value {
	return jsUndefined{
		v: js.Undefined(),
	}
}

func (j jsFactory) Null() driver.Value {
	return jsNull{
		v: js.Null(),
	}
}

func (j jsFactory) Number(f float64) driver.Value {
	return jsNumber{
		v: js.ValueOf(f),
	}
}

func (j jsFactory) String(s string) driver.Value {
	return jsString{
		v: js.ValueOf(s),
	}
}

func (j jsFactory) Function(f func(this driver.Object, args ...driver.Value) driver.Value) driver.Function {
	return jsFunction{
		v: js.ValueOf(js.FuncOf(func(jsThis js.Value, jsArgs []js.Value) interface{} {
			var vArgs []driver.Value
			for _, arg := range jsArgs {
				vArgs = append(vArgs, js2value(arg))
			}
			rv := f(jsObject{v: jsThis}, vArgs...)
			return value2js(rv)
		})),
	}
}

var _ driver.Factory = jsFactory{}
