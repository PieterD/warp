package dom

import "github.com/PieterD/warp/pkg/driver"

type Event struct {
	factory driver.Factory
	obj     driver.Object
}

func (e *Event) Type() string {
	typ, ok := e.obj.Get("type").IsString()
	if !ok {
		return "!UNKNOWN-TYPE!"
	}
	return typ
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
	ke.Key, ok = o.Get("key").IsString()
	if !ok {
		return KeyboardEvent{}, false
	}
	ke.Code, ok = o.Get("code").IsString()
	if !ok {
		return KeyboardEvent{}, false
	}
	ke.Repeat, ok = o.Get("repeat").IsBoolean()
	if !ok {
		return KeyboardEvent{}, false
	}
	ke.ShiftKey, ok = o.Get("shiftKey").IsBoolean()
	if !ok {
		return KeyboardEvent{}, false
	}
	ke.MetaKey, ok = o.Get("metaKey").IsBoolean()
	if !ok {
		return KeyboardEvent{}, false
	}
	ke.ControlKey, ok = o.Get("ctrlKey").IsBoolean()
	if !ok {
		return KeyboardEvent{}, false
	}
	ke.AltKey, ok = o.Get("altKey").IsBoolean()
	if !ok {
		return KeyboardEvent{}, false
	}
	return ke, true
}
