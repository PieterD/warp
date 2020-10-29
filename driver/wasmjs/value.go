package wasmjs

import (
	"fmt"
	"syscall/js"

	"github.com/PieterD/warp/driver"
)

type jsUndefined struct {
	jsEmpty
	v js.Value
}

func (jsUndefined) IsUndefined() bool {
	return true
}

type jsNull struct {
	jsEmpty
	v js.Value
}

func (jsNull) IsNull() bool {
	return true
}

type jsBoolean struct {
	jsEmpty
	v js.Value
}

func (j jsBoolean) IsBoolean() (bool, bool) {
	switch jsType := j.v.Type(); jsType {
	case js.TypeBoolean:
		return j.v.Bool(), true
	default:
		return false, false
	}
}

type jsNumber struct {
	jsEmpty
	v js.Value
}

func (j jsNumber) IsNumber() (float64, bool) {
	switch jsType := j.v.Type(); jsType {
	case js.TypeNumber:
		return j.v.Float(), true
	default:
		return 0, false
	}
}

type jsString struct {
	jsEmpty
	v js.Value
}

func (j jsString) IsString() (string, bool) {
	switch jsType := j.v.Type(); jsType {
	case js.TypeString:
		return j.v.String(), true
	default:
		return "", false
	}
}

type jsObject struct {
	jsEmpty
	v js.Value
}

func (j jsObject) IsObject() driver.Object {
	switch jsType := j.v.Type(); jsType {
	case js.TypeObject:
		return j
	default:
		return nil
	}
}

func (j jsObject) Get(key string) driver.Value {
	return js2value(j.v.Get(key))
}

func (j jsObject) Set(key string, value driver.Value) {
	j.v.Set(key, value2js(value))
}

type jsFunction struct {
	jsEmpty
	v js.Value
}

func (j jsFunction) IsFunction() driver.Function {
	switch jsType := j.v.Type(); jsType {
	case js.TypeFunction:
		return j
	default:
		return nil
	}
}

func (j jsFunction) New(args ...driver.Value) driver.Object {
	var jsArgs []interface{}
	for _, arg := range args {
		jsArgs = append(jsArgs, value2js(arg))
	}
	jsObj := j.v.New(jsArgs...)
	return jsObject{
		v: jsObj,
	}
}

func (j jsFunction) Call(this driver.Object, args ...driver.Value) driver.Value {
	ourThis, ok := this.(jsObject)
	if !ok {
		panic(fmt.Errorf("unknown this type: %T", this))
	}
	jsArgs := []interface{}{ourThis.v}
	for _, arg := range args {
		jsArgs = append(jsArgs, value2js(arg))
	}
	jsReturn := j.v.Call("call", jsArgs...)
	return js2value(jsReturn)
}

type jsEmpty struct{}

func (j jsEmpty) IsBoolean() (value, ok bool) {
	return false, false
}

func (j jsEmpty) IsUndefined() (ok bool) {
	return false
}

func (j jsEmpty) IsNull() (ok bool) {
	return false
}

func (j jsEmpty) IsNumber() (value float64, ok bool) {
	return 0, false
}

func (j jsEmpty) IsString() (value string, ok bool) {
	return "", false
}

func (j jsEmpty) IsObject() (value driver.Object) {
	return nil
}

func (j jsEmpty) IsFunction() (value driver.Function) {
	return nil
}

var _ driver.Value = jsEmpty{}
