package dom

import (
	"fmt"
	"strings"

	"github.com/PieterD/warp/driver"
	"github.com/PieterD/warp/driver/driverutil"
)

type Elem struct {
	factory driver.Factory
	obj     driver.Object
}

func (e *Elem) Tag() string {
	v, ok := e.obj.Get("tagName").IsString()
	if !ok {
		return "!UNKNOWN-TAG-TYPE!"
	}
	return strings.ToLower(v)
}

func (e *Elem) SetText(text string) {
	e.obj.Set("innerText", e.factory.String(text))
}

func (e *Elem) Children() (children []*Elem) {
	childrenObject := e.obj.Get("children").IsObject()
	if childrenObject == nil {
		return
	}
	for _, childValue := range driverutil.IndexableToSlice(e.factory, childrenObject) {
		childObject := childValue.IsObject()
		if childObject == nil {
			return
		}
		children = append(children, &Elem{
			factory: e.factory,
			obj:     childObject,
		})
	}
	return children
}

func (e *Elem) AppendChildren(children ...*Elem) {
	fAppendChild := driverutil.Bind(e.obj, "appendChild")
	for _, child := range children {
		fAppendChild(child.obj)
	}
}

func (e *Elem) ClearChildren() {
	fRemoveChild := driverutil.Bind(e.obj, "removeChild")
	for {
		firstChildObject := e.obj.Get("firstChild").IsObject()
		if firstChildObject == nil {
			break
		}
		fRemoveChild(firstChildObject)
	}
}

func (e *Elem) EventHandler(eventName string, f func(this *Elem, event *Event)) (deregister func()) {
	fAddEventListener := driverutil.Bind(e.obj, "addEventListener")
	cbFunction := e.factory.Function(func(this driver.Object, args ...driver.Value) driver.Value {
		if len(args) != 1 {
			panic(fmt.Errorf("expected 1 argument, got: %d", len(args)))
		}
		eventObject := args[0].IsObject()
		if eventObject == nil {
			panic(fmt.Errorf("first argument is not an object: %T", args[0]))
		}
		evt := &Event{
			factory: e.factory,
			obj:     eventObject,
		}
		f(e, evt)
		return nil
	})
	dEventName := e.factory.String(eventName)
	fAddEventListener(dEventName, cbFunction)
	return func() {
		fRemoveEventListener := driverutil.Bind(e.obj, "removeEventListener")
		fRemoveEventListener(dEventName, cbFunction)
	}
}
