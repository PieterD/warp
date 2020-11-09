package gl

import (
	"fmt"
	"image"

	"github.com/PieterD/warp/pkg/driver"
	"github.com/PieterD/warp/pkg/driver/driverutil"
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
	fGetContext := driverutil.Bind(canvasObject, "getContext")
	ctxObject := fGetContext(factory.String("webgl2")).IsObject()
	if ctxObject == nil {
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

func (glx *Context) VertexArray(attrs ...VertexArrayAttribute) (*VertexArray, error) {
	return newVertexArray(glx, attrs...)
}

func (glx *Context) Use(p *Program) {
	glx.constants.UseProgram(p.glObject)
}

func (glx *Context) Draw(cfg DrawConfig) error {
	return doDraw(glx, cfg)
}

func (glx *Context) Texture(cfg Texture2DConfig, img image.Image) *Texture2D {
	return newTexture2D(glx, cfg, img)
}

func (glx *Context) BindTextureUnits(textures ...*Texture2D) {
	maxUnits := glx.Parameters().MaxCombinedTextureImageUnits()
	if len(textures) >= maxUnits {
		panic(fmt.Errorf("only %d texture units allowed, got: %d", maxUnits, len(textures)))
	}
	fTexture0, ok := glx.constants.TEXTURE0.IsNumber()
	if !ok {
		panic(fmt.Errorf("expected TEXTURE0 to be a number: %T", glx.constants.TEXTURE0))
	}
	t0 := int(fTexture0)
	for textureUnit := 0; textureUnit < maxUnits; textureUnit++ {
		jsTextureUnit := glx.factory.Number(float64(t0 + textureUnit))
		glx.constants.ActiveTexture(jsTextureUnit)
		glObject := glx.factory.Null()
		if textureUnit < len(textures) && textures[textureUnit] != nil {
			glObject = textures[textureUnit].glObject
		} else {
			// Do we want to disable all non-selected units, or do we leave them alone?
			// For now, leave them alone.
			break
		}
		glx.constants.BindTexture(glx.constants.TEXTURE_2D, glObject)
	}
	glx.constants.ActiveTexture(glx.constants.TEXTURE0)
}

func (glx *Context) Parameters() *ParameterSet {
	return newParameterSet(glx)
}