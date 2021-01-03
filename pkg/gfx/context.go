package gfx

import (
	"fmt"
	"image"

	"github.com/PieterD/warp/pkg/driver"
)

type Context struct {
	factory       driver.Factory
	obj           driver.Object
	constants     glConstants
	typeConverter *typeConverter
}

type Canvas interface {
	Driver() (factory driver.Factory, obj driver.Object)
}

func NewContext(canvas Canvas) *Context {
	factory, canvasObject := canvas.Driver()
	fGetContext := driver.Bind(canvasObject, "getContext")
	ctxObject, ok := fGetContext(factory.String("webgl2")).ToObject()
	if !ok {
		return nil
	}

	constants := newGlConstants(ctxObject, false)
	typeConverter := newTypeConverter(constants)
	glx := &Context{
		factory:       factory,
		obj:           ctxObject,
		constants:     constants,
		typeConverter: typeConverter,
	}
	return glx
}

func (glx *Context) Program(cfg ProgramConfig) (*Program, error) {
	return newProgram(glx, cfg)
}

func (glx *Context) Buffer() *Buffer {
	return newBuffer(glx)
}

func (glx *Context) VertexArray(adc ActiveCoupling, buffers map[string]*Buffer) (*VertexArray, error) {
	return newVertexArray(glx, adc, buffers)
}

func (glx *Context) Use(p *Program) {
	glx.constants.UseProgram(p.glObject)
}

func (glx *Context) Clear() {
	glx.constants.Clear(glx.constants.COLOR_BUFFER_BIT)
	glx.constants.Clear(glx.constants.DEPTH_BUFFER_BIT)
}

func (glx *Context) Draw(cfg DrawConfig) error {
	return doDraw(glx, cfg)
}

func (glx *Context) Texture(cfg Texture2DConfig, img image.Image) *Texture2D {
	return newTexture2D(glx, cfg, img)
}

func (glx *Context) bindTextureUnit(unit int, texture *Texture2D) {
	maxUnits := glx.Parameters().MaxCombinedTextureImageUnits()
	if unit >= maxUnits {
		panic(fmt.Errorf("trying to bind to %d, but only %d texture units are allowed %d", unit, maxUnits))
	}
	fTexture0, ok := glx.constants.TEXTURE0.ToFloat64()
	if !ok {
		panic(fmt.Errorf("expected TEXTURE0 to be a number: %T", glx.constants.TEXTURE0))
	}
	t0 := int(fTexture0)
	jsTextureUnit := glx.factory.Number(float64(t0 + unit))
	glx.constants.ActiveTexture(jsTextureUnit)
	glx.constants.BindTexture(glx.constants.TEXTURE_2D, texture.glObject)
}

func (glx *Context) Parameters() *ParameterSet {
	return newParameterSet(glx)
}

func (glx *Context) Viewport(x, y, w, h int) {
	glx.constants.Viewport(
		glx.factory.Number(float64(x)),
		glx.factory.Number(float64(y)),
		glx.factory.Number(float64(w)),
		glx.factory.Number(float64(h)),
	)
}