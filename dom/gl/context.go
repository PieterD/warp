package gl

import (
	"github.com/PieterD/warp/driver"
	"github.com/PieterD/warp/driver/driverutil"
)

type Context struct {
	factory       driver.Factory
	obj           driver.Object
	constants     glConstants
	functions     glFunctions
	typeConverter *typeConverter
}

type Canvas interface {
	Driver() (factory driver.Factory, obj driver.Object)
}

func NewContext(canvas Canvas) *Context {
	factory, canvasObject := canvas.Driver()
	fGetContext := driverutil.Bind(canvasObject, "getContext")
	ctxObject := fGetContext(factory.String("webgl2")).IsObject()
	if ctxObject == nil {
		return nil
	}

	constants := newGlConstants(ctxObject)
	functions := newGlFunctions(ctxObject)
	typeConverter := newTypeConverter(constants)
	glx := &Context{
		factory:       factory,
		obj:           ctxObject,
		constants:     constants,
		functions:     functions,
		typeConverter: typeConverter,
	}
	return glx
}

func (glx *Context) Program(cfg ProgramConfig) (*Program, error) {
	return newProgram(glx, cfg)
}

func (glx *Context) Buffer() (*Buffer, error) {
	return newBuffer(glx)
}

func (glx *Context) VertexArray(cfg VertexArrayConfig) (*VertexArray, error) {
	return newVertexArray(glx, cfg)
}

func (glx *Context) Use(p *Program) {
	glx.functions.UseProgram(p.glObject)
}

func (glx *Context) Draw(cfg DrawConfig) error {
	return doDraw(glx, cfg)
}
