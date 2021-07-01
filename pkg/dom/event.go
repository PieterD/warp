package dom

import (
	"fmt"
	"github.com/PieterD/warp/pkg/driver"
)

type Driver interface {
	Driver() (factory driver.Factory, obj driver.Object)
}

type Event struct {
	factory driver.Factory
	obj     driver.Object
}

func (e *Event) Driver() (factory driver.Factory, obj driver.Object) {
	return e.factory, e.obj
}

func (e *Event) Type() string {
	typ, ok := e.obj.Get("type").ToString()
	if !ok {
		return "!UNKNOWN-TYPE!"
	}
	return typ
}

type DeviceMotionEvent struct {
	X float64
	Y float64
	Z float64
}

func (e *Event) AsDeviceMotionEvent() (DeviceMotionEvent, bool) {
	o := e.obj
	o, ok := o.Get("accelerationIncludingGravity").ToObject()
	if !ok {
		return DeviceMotionEvent{}, false
	}
	x, ok := o.Get("x").ToFloat64()
	if !ok {
		return DeviceMotionEvent{}, false
	}
	y, ok := o.Get("y").ToFloat64()
	if !ok {
		return DeviceMotionEvent{}, false
	}
	z, ok := o.Get("z").ToFloat64()
	if !ok {
		return DeviceMotionEvent{}, false
	}
	return DeviceMotionEvent{
		X: x,
		Y: y,
		Z: z,
	}, true
}

type DeviceOrientationEvent struct {
	Alpha float64
	Beta  float64
	Gamma float64
}

func (e *Event) AsDeviceOrientationEvent() (DeviceOrientationEvent, bool) {
	o := e.obj
	alpha, ok := o.Get("alpha").ToFloat64()
	if !ok {
		return DeviceOrientationEvent{}, false
	}
	beta, ok := o.Get("beta").ToFloat64()
	if !ok {
		return DeviceOrientationEvent{}, false
	}
	gamma, ok := o.Get("gamma").ToFloat64()
	if !ok {
		return DeviceOrientationEvent{}, false
	}
	return DeviceOrientationEvent{
		Alpha: alpha,
		Beta:  beta,
		Gamma: gamma,
	}, true
}

type MouseEvent struct {
	OffsetX int
	OffsetY int
}

func (e *Event) AsMouse() (MouseEvent, bool) {
	o := e.obj
	var me MouseEvent

	offsetX, ok := o.Get("offsetX").ToFloat64()
	if !ok {
		return MouseEvent{}, false
	}
	me.OffsetX = int(offsetX)

	offsetY, ok := o.Get("offsetY").ToFloat64()
	if !ok {
		return MouseEvent{}, false
	}
	me.OffsetY = int(offsetY)

	return me, true
}

type KeyboardEvent struct {
	Key        string
	Code       string
	Repeat     bool
	ShiftKey   bool
	MetaKey    bool
	ControlKey bool
	AltKey     bool
}

func (e *Event) AsKeyboard() (KeyboardEvent, bool) {
	o := e.obj
	var ke KeyboardEvent
	var ok bool
	ke.Key, ok = o.Get("key").ToString()
	if !ok {
		return KeyboardEvent{}, false
	}
	ke.Code, ok = o.Get("code").ToString()
	if !ok {
		return KeyboardEvent{}, false
	}
	ke.Repeat, ok = o.Get("repeat").ToBoolean()
	if !ok {
		return KeyboardEvent{}, false
	}
	ke.ShiftKey, ok = o.Get("shiftKey").ToBoolean()
	if !ok {
		return KeyboardEvent{}, false
	}
	ke.MetaKey, ok = o.Get("metaKey").ToBoolean()
	if !ok {
		return KeyboardEvent{}, false
	}
	ke.ControlKey, ok = o.Get("ctrlKey").ToBoolean()
	if !ok {
		return KeyboardEvent{}, false
	}
	ke.AltKey, ok = o.Get("altKey").ToBoolean()
	if !ok {
		return KeyboardEvent{}, false
	}
	return ke, true
}

func AddEventListener(elem Driver, eventName string, f func(this driver.Value, event *Event)) (deregister func(), err error) {
	factory, object := elem.Driver()
	addEventListener := driver.Bind(object, "addEventListener")
	if addEventListener == nil {
		return nil, fmt.Errorf("no addEventListener present")
	}

	jsEventName := factory.String(eventName)
	jsCallback := factory.Function(func(this driver.Object, args ...driver.Value) driver.Value {
		if len(args) != 1 {
			panic(fmt.Errorf("expected 1 argument, got: %d", len(args)))
		}
		eventObject, ok := args[0].ToObject()
		if !ok {
			panic(fmt.Errorf("first argument is not an object: %T", args[0]))
		}
		evt := &Event{
			factory: factory,
			obj:     eventObject,
		}
		f(object, evt)
		return nil
	})

	addEventListener(jsEventName, jsCallback)
	return func() {
		fRemoveEventListener := driver.Bind(object, "removeEventListener")
		fRemoveEventListener(jsEventName, jsCallback)
	}, nil
}
