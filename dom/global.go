package dom

import "github.com/PieterD/warp/driver"

type Global struct {
	factory driver.Factory
	obj     driver.Object
}

func Open(factory driver.Factory) *Global {
	return &Global{
		factory: factory,
		obj:     factory.Global(),
	}
}

func (g *Global) Window() *Window {
	dValue := g.obj.Get("window")
	dObj := dValue.IsObject()
	if dObj == nil {
		return nil
	}
	return &Window{
		factory: g.factory,
		obj:     dObj,
	}
}
