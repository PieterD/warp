package dom

import "github.com/PieterD/warp/driver"

type Window struct {
	factory driver.Factory
	obj     driver.Object
}

func (w *Window) Document() *Document {
	dValue := w.obj.Get("document")
	dObj := dValue.IsObject()
	if dObj == nil {
		return nil
	}
	return &Document{
		factory: w.factory,
		obj:     dObj,
	}
}
