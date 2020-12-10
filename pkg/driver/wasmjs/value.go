package wasmjs

import (
	"fmt"
	"syscall/js"

	"github.com/PieterD/warp/pkg/driver"
)

type jsUndefined struct {
	jsEmpty
	v js.Value
}

func (j jsUndefined) jsValue() js.Value {
	return j.v
}

func (j jsUndefined) IsUndefined() bool {
	return true
}

var _ vValue = jsUndefined{}

type jsNull struct {
	jsEmpty
	v js.Value
}

func (j jsNull) jsValue() js.Value {
	return j.v
}

func (j jsNull) IsNull() bool {
	return true
}

var _ vValue = jsNull{}

type jsBoolean struct {
	jsEmpty
	v js.Value
}

func (j jsBoolean) jsValue() js.Value {
	return j.v
}

func (j jsBoolean) ToBoolean() (bool, bool) {
	switch jsType := j.v.Type(); jsType {
	case js.TypeBoolean:
		return j.v.Bool(), true
	default:
		return false, false
	}
}

var _ vValue = jsBoolean{}

type jsNumber struct {
	jsEmpty
	v js.Value
}

func (j jsNumber) jsValue() js.Value {
	return j.v
}

func (j jsNumber) ToFloat64() (float64, bool) {
	switch jsType := j.v.Type(); jsType {
	case js.TypeNumber:
		return j.v.Float(), true
	default:
		return 0, false
	}
}

var _ vValue = jsNumber{}

type jsString struct {
	jsEmpty
	v js.Value
}

func (j jsString) jsValue() js.Value {
	return j.v
}

func (j jsString) ToString() (string, bool) {
	switch jsType := j.v.Type(); jsType {
	case js.TypeString:
		return j.v.String(), true
	default:
		return "", false
	}
}

var _ vValue = jsString{}

type jsObject struct {
	jsEmpty
	v js.Value
}

func (j jsObject) jsValue() js.Value {
	return j.v
}

func (j jsObject) ToObject() (driver.Object, bool) {
	switch jsType := j.v.Type(); jsType {
	case js.TypeObject:
		return j, true
	default:
		return nil, false
	}
}

var _ vValue = jsObject{}

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

func (j jsFunction) jsValue() js.Value {
	return j.v
}

func (j jsFunction) ToFunction() (driver.Function, bool) {
	switch jsType := j.v.Type(); jsType {
	case js.TypeFunction:
		return j, true
	default:
		return nil, false
	}
}

var _ vValue = jsFunction{}

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
	vThis, ok := this.(vValue)
	if !ok {
		panic(fmt.Errorf("unknown this type: %T", this))
	}
	jsArgs := []interface{}{vThis.jsValue()}
	for _, arg := range args {
		jsArgs = append(jsArgs, value2js(arg))
	}
	jsReturn := j.v.Call("call", jsArgs...)
	return js2value(jsReturn)
}

type jsBuffer struct {
	jsEmpty
	factory driver.Factory
	v       js.Value
	obj     driver.Object
}

func (j jsBuffer) jsValue() js.Value {
	return j.v
}

func newBuffer(factory jsFactory, size int) jsBuffer {
	fUint8Array, ok := factory.Global().Get("Uint8Array").ToFunction()
	if !ok {
		panic(fmt.Errorf("Uint8Array constructor is missing"))
	}
	obj := fUint8Array.New(factory.Number(float64(size)))
	vObj, ok := obj.(vValue)
	if !ok {
		panic(fmt.Errorf("buffer object was somehow not an object"))
	}
	return jsBuffer{
		factory: factory,
		v:       vObj.jsValue(),
		obj:     obj,
	}
}

var _ vValue = jsBuffer{}

func (j jsBuffer) Size() int {
	length, ok := j.obj.Get("length").ToFloat64()
	if !ok {
		panic(fmt.Errorf("buffer length was not a number"))
	}
	return int(length)
}

func (j jsBuffer) Put(data []byte) int {
	return js.CopyBytesToJS(j.v, data)
}

func (j jsBuffer) Get(data []byte) int {
	return js.CopyBytesToGo(data, j.v)
}

func (j jsBuffer) AsUint8Array() driver.Object {
	con, ok := j.factory.Global().Get("Uint8Array").ToFunction()
	if !ok {
		panic(fmt.Errorf("Uint8Array was not a function"))
	}
	return con.New(j.obj.Get("buffer"))
}

func (j jsBuffer) AsUint16Array() driver.Object {
	con, ok := j.factory.Global().Get("Uint16Array").ToFunction()
	if !ok {
		panic(fmt.Errorf("Uint16Array was not a function"))
	}
	return con.New(j.obj.Get("buffer"))
}

func (j jsBuffer) AsFloat32Array() driver.Object {
	con, ok := j.factory.Global().Get("Float32Array").ToFunction()
	if !ok {
		panic(fmt.Errorf("Float32Array was not a function"))
	}
	return con.New(j.obj.Get("buffer"))
}

type jsEmpty struct{}

func (j jsEmpty) ToBoolean() (value, ok bool) {
	return false, false
}

func (j jsEmpty) IsUndefined() (ok bool) {
	return false
}

func (j jsEmpty) IsNull() (ok bool) {
	return false
}

func (j jsEmpty) ToFloat64() (value float64, ok bool) {
	return 0, false
}

func (j jsEmpty) ToString() (value string, ok bool) {
	return "", false
}

func (j jsEmpty) ToObject() (optionalValue driver.Object, ok bool) {
	return nil, ok
}

func (j jsEmpty) ToFunction() (optionalValue driver.Function, ok bool) {
	return nil, ok
}

var _ driver.Value = jsEmpty{}
