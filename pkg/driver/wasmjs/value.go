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

type jsBuffer struct {
	jsEmpty
	factory driver.Factory
	v       js.Value
	obj     driver.Object
}

func newBuffer(factory jsFactory, size int) jsBuffer {
	fUint8Array := factory.Global().Get("Uint8Array").IsFunction()
	if fUint8Array == nil {
		panic(fmt.Errorf("Uint8Array constructor is missing"))
	}
	obj := fUint8Array.New(factory.Number(float64(size)))
	ours, ok := obj.(jsObject)
	if !ok {
		panic(fmt.Errorf("buffer object was somehow not an object"))
	}
	return jsBuffer{
		factory: factory,
		v:       ours.v,
		obj:     obj,
	}
}

func (j jsBuffer) Size() int {
	length, ok := j.obj.Get("length").IsNumber()
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
	con := j.factory.Global().Get("Uint8Array").IsFunction()
	if con == nil {
		panic(fmt.Errorf("Uint8Array was not a function"))
	}
	return con.New(j.obj.Get("buffer"))
}

func (j jsBuffer) AsUint16Array() driver.Object {
	con := j.factory.Global().Get("Uint16Array").IsFunction()
	if con == nil {
		panic(fmt.Errorf("Uint16Array was not a function"))
	}
	return con.New(j.obj.Get("buffer"))
}

func (j jsBuffer) AsFloat32Array() driver.Object {
	con := j.factory.Global().Get("Float32Array").IsFunction()
	if con == nil {
		panic(fmt.Errorf("Float32Array was not a function"))
	}
	return con.New(j.obj.Get("buffer"))
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
