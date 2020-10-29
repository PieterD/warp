package dom

import (
	"github.com/PieterD/warp/driver"
	"github.com/PieterD/warp/driver/driverutil"
)

type Elem struct {
	factory driver.Factory
	obj     driver.Object
}

func (e *Elem) AppendChildren(children ...*Elem) {
	fAppendChild := driverutil.Bind(e.obj, "appendChild")
	for _, child := range children {
		fAppendChild(child.obj)
	}
}

func (e *Elem) SetText(text string) {
	e.obj.Set("innerText", e.factory.String(text))
}
