package dom

import (
	"github.com/PieterD/warp/pkg/driver"
)

type Document struct {
	factory driver.Factory
	obj     driver.Object
}

func (doc *Document) Body() *Elem {
	dValue := doc.obj.Get("body")
	dObj, ok := dValue.ToObject()
	if !ok {
		return nil
	}
	return &Elem{
		factory: doc.factory,
		obj:     dObj,
	}
}

func (doc *Document) CreateElem(tag string, constructor func(newElem *Elem)) *Elem {
	fCreateElement := driver.Bind(doc.obj, "createElement")
	elementValue := fCreateElement(doc.factory.String(tag))
	elementObject, ok := elementValue.ToObject()
	if !ok {
		return nil
	}
	elem := &Elem{
		factory: doc.factory,
		obj:     elementObject,
	}
	if constructor != nil {
		constructor(elem)
	}
	return elem
}
