package dom

import "github.com/PieterD/warp/pkg/driver"

type Event struct {
	factory driver.Factory
	obj     driver.Object
}

func (e *Event) Type() string {
	typ, ok := e.obj.Get("type").ToString()
	if !ok {
		return "!UNKNOWN-TYPE!"
	}
	return typ
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
