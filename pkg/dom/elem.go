package dom

import (
	"fmt"
	"sort"
	"strings"

	"github.com/PieterD/warp/pkg/driver"
)

type Elem struct {
	factory driver.Factory
	obj     driver.Object
}

func (e *Elem) Driver() (factory driver.Factory, obj driver.Object) {
	return e.factory, e.obj
}

func (e *Elem) Tag() string {
	v, ok := e.obj.Get("tagName").ToString()
	if !ok {
		return "!UNKNOWN-TAG-TYPE!"
	}
	return strings.ToLower(v)
}

func (e *Elem) SetText(text string) {
	e.obj.Set("innerText", e.factory.String(text))
}

func (e *Elem) SetPropString(propName string, propValue string) {
	e.obj.Set(propName, e.factory.String(propValue))
}

func (e *Elem) Classes() []string {
	allClasses, ok := e.obj.Get("className").ToString()
	if !ok {
		return nil
	}
	return strings.Split(allClasses, " ")
}

func (e *Elem) AppendClasses(classNames ...string) {
	set := make(map[string]struct{})
	for _, name := range e.Classes() {
		set[name] = struct{}{}
	}
	for _, name := range classNames {
		set[name] = struct{}{}
	}
	list := make([]string, len(set))
	for name := range set {
		list = append(list, name)
	}
	sort.Strings(list)
	e.obj.Set("className", e.factory.String(strings.Join(list, " ")))
}

func (e *Elem) ClearClasses() {
	e.obj.Set("className", e.factory.String(""))
}

func (e *Elem) Children() (children []*Elem) {
	childrenObject, ok := e.obj.Get("children").ToObject()
	if !ok {
		return
	}
	for _, childValue := range driver.IndexableToSlice(e.factory, childrenObject) {
		childObject, ok := childValue.ToObject()
		if !ok {
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
	fAppendChild := driver.Bind(e.obj, "appendChild")
	for _, child := range children {
		fAppendChild(child.obj)
	}
}

func (e *Elem) ClearChildren() {
	fRemoveChild := driver.Bind(e.obj, "removeChild")
	for {
		firstChildObject, ok := e.obj.Get("firstChild").ToObject()
		if !ok {
			break
		}
		fRemoveChild(firstChildObject)
	}
}

func (e *Elem) EventHandler(eventName string, f func(this *Elem, event *Event)) (deregister func()) {
	fAddEventListener := driver.Bind(e.obj, "addEventListener")
	cbFunction := e.factory.Function(func(this driver.Object, args ...driver.Value) driver.Value {
		if len(args) != 1 {
			panic(fmt.Errorf("expected 1 argument, got: %d", len(args)))
		}
		eventObject, ok := args[0].ToObject()
		if !ok {
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
		fRemoveEventListener := driver.Bind(e.obj, "removeEventListener")
		fRemoveEventListener(dEventName, cbFunction)
	}
}
