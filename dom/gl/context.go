package gl

import (
	"github.com/PieterD/warp/driver"
	"github.com/PieterD/warp/driver/driverutil"
)

type Context struct {
	factory   driver.Factory
	obj       driver.Object
	constants glConstants
	functions glFunctions
}

type Canvas interface {
	Driver() (factory driver.Factory, obj driver.Object)
}

func NewContext(canvas Canvas) *Context {
	factory, canvasObject := canvas.Driver()
	fGetContext := driverutil.Bind(canvasObject, "getContext")
	ctxObject := fGetContext(factory.String("webgl")).IsObject()
	if ctxObject == nil {
		return nil
	}

	return &Context{
		factory:   factory,
		obj:       ctxObject,
		constants: newGlConstants(ctxObject),
		functions: newGlFunctions(ctxObject),
	}
}

func (glx *Context) Program(cfg ProgramConfig) (*Program, error) {
	return newProgram(glx, cfg)
}

func (glx *Context) Buffer() (*Buffer, error) {
	return newBuffer(glx)
}
