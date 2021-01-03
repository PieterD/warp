package gl

import (
	"fmt"

	"github.com/PieterD/warp/pkg/driver"
)

type Canvas interface {
	Driver() (factory driver.Factory, obj driver.Object)
}

type Context struct {
	factory       driver.Factory
	obj           driver.Object
	constants     glConstants
	typeConverter *typeConverter
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

func (glx *Context) Destroy() {
	getExtension := driver.Bind(glx.obj, "getExtension")
	if getExtension == nil {
		panic(fmt.Errorf("missing getExtension"))
	}
	extension := getExtension(glx.factory.String("WEBGL_lose_context"))
	extensionObject, ok := extension.ToObject()
	if !ok {
		return
	}
	loseContext := driver.Bind(extensionObject, "loseContext")
	if loseContext == nil {
		panic(fmt.Errorf("missing loseContext"))
	}
	loseContext()
}

func (glx *Context) UseProgram(program ProgramObject) {
	glx.constants.UseProgram(program.value)
}

func (glx *Context) UnuseProgram() {
	glx.constants.UseProgram(glx.factory.Null())
}

func (glx *Context) CreateProgram() ProgramObject {
	programObject := glx.constants.CreateProgram()
	return ProgramObject{
		glx:   glx,
		value: programObject,
	}
}

func (glx *Context) CreateShader(shaderType ShaderType) ShaderObject {
	var glShaderType driver.Value
	switch shaderType {
	case VertexShader:
		glShaderType = glx.constants.VERTEX_SHADER
	case FragmentShader:
		glShaderType = glx.constants.FRAGMENT_SHADER
	default:
		panic(fmt.Errorf("unusable shaderType: %v", shaderType))
	}
	value := glx.constants.CreateShader(glShaderType)
	return ShaderObject{
		glx:   glx,
		value: value,
	}
}

func (glx *Context) Targets() Targets {
	return Targets{
		glx: glx,
	}
}

func (glx *Context) Parameters() Parameters {
	return Parameters{
		glx: glx,
	}
}

type BufferObject struct {
	glx   *Context
	value driver.Value
}

func (glx *Context) CreateBuffer() BufferObject {
	value := glx.constants.CreateBuffer()
	return BufferObject{
		glx:   glx,
		value: value,
	}
}

func (buffer BufferObject) Destroy() {
	glx := buffer.glx
	glx.constants.DeleteBuffer(buffer.value)
}

type VertexArrayObject struct {
	glx   *Context
	value driver.Value
}

func (glx *Context) CreateVertexArray() VertexArrayObject {
	value := glx.constants.CreateVertexArray()
	return VertexArrayObject{
		glx:   glx,
		value: value,
	}
}

func (vao VertexArrayObject) Destroy() {
	glx := vao.glx
	glx.constants.DeleteVertexArray(vao.value)
}

// A buffer must be bound to the ARRAY_BUFFER target.
func (vao VertexArrayObject) VertexAttribPointer(attrIndex int, attrType Type, normalized bool, byteStride int, byteOffset int) {
	glx := vao.glx
	bufferType, bufferItemsPerVertex, err := attrType.asAttribute()
	if err != nil {
		panic(fmt.Errorf("converting attribute type %s to attribute: %w", attrType, err))
	}
	glx.constants.VertexAttribPointer(
		glx.factory.Number(float64(attrIndex)),
		glx.factory.Number(float64(bufferItemsPerVertex)),
		glx.typeConverter.ToJs(bufferType),
		glx.factory.Boolean(normalized),
		glx.factory.Number(float64(byteStride)),
		glx.factory.Number(float64(byteOffset)),
	)
}

func (vao VertexArrayObject) EnableVertexAttribArray(attrIndex int) {
	glx := vao.glx
	glx.constants.EnableVertexAttribArray(glx.factory.Number(float64(attrIndex)))
}

func (vao VertexArrayObject) DisableVertexAttribArray(attrIndex int) {
	glx := vao.glx
	glx.constants.DisableVertexAttribArray(glx.factory.Number(float64(attrIndex)))
}

func (glx *Context) BindVertexArray(vao VertexArrayObject) {
	glx.constants.BindVertexArray(vao.value)
}

func (glx *Context) UnbindVertexArray() {
	glx.constants.BindVertexArray(glx.factory.Null())
}

//go:generate stringer -type=PrimitiveDrawMode
type PrimitiveDrawMode int

const (
	Points PrimitiveDrawMode = iota
	Lines
	Triangles
)

func (glx *Context) DrawArrays(mode PrimitiveDrawMode, vertexOffset, vertexCount int) {
	var glDrawMode driver.Value
	switch mode {
	case Points:
		glDrawMode = glx.constants.POINTS
	case Lines:
		glDrawMode = glx.constants.LINES
	case Triangles:
		glDrawMode = glx.constants.TRIANGLES
	default:
		panic(fmt.Errorf("unsupported draw mode: %v", mode))
	}
	glx.constants.DrawArrays(
		glDrawMode,
		glx.factory.Number(float64(vertexOffset)),
		glx.factory.Number(float64(vertexCount)),
	)
}

func (glx *Context) DrawElements(mode PrimitiveDrawMode, elementArrayByteOffset, vertexCount int, elementArrayType Type) {
	var glDrawMode driver.Value
	switch mode {
	case Points:
		glDrawMode = glx.constants.POINTS
	case Lines:
		glDrawMode = glx.constants.LINES
	case Triangles:
		glDrawMode = glx.constants.TRIANGLES
	default:
		panic(fmt.Errorf("unsupported draw mode: %v", mode))
	}
	glx.constants.DrawElements(
		glDrawMode,
		glx.factory.Number(float64(vertexCount)),
		glx.typeConverter.ToJs(elementArrayType),
		glx.factory.Number(float64(elementArrayByteOffset)),
	)
}

func (glx *Context) Viewport(x, y, w, h int) {
	glx.constants.Viewport(
		glx.factory.Number(float64(x)),
		glx.factory.Number(float64(y)),
		glx.factory.Number(float64(w)),
		glx.factory.Number(float64(h)),
	)
}

func (glx *Context) ClearColor(r, g, b, a float32) {
	glx.constants.ClearColor(
		glx.factory.Number(float64(r)),
		glx.factory.Number(float64(g)),
		glx.factory.Number(float64(b)),
		glx.factory.Number(float64(a)),
	)
}

//TODO: specify what to clear
func (glx *Context) Clear() {
	glx.constants.Clear(glx.constants.COLOR_BUFFER_BIT)
	glx.constants.Clear(glx.constants.DEPTH_BUFFER_BIT)
}

type RenderbufferObject struct {
	glx   *Context
	value driver.Value
}

func (glx *Context) CreateRenderbuffer() RenderbufferObject {
	value := glx.constants.CreateRenderbuffer()
	return RenderbufferObject{
		glx:   glx,
		value: value,
	}
}

func (rbo RenderbufferObject) Destroy() {
	glx := rbo.glx
	glx.constants.DeleteRenderbuffer(rbo.value)
}

type FramebufferObject struct {
	glx   *Context
	value driver.Value
}

func (glx *Context) CreateFramebuffer() FramebufferObject {
	value := glx.constants.CreateFramebuffer()
	return FramebufferObject{
		glx:   glx,
		value: value,
	}
}

func (fbo FramebufferObject) Destroy() {
	glx := fbo.glx
	glx.constants.DeleteFramebuffer(fbo.value)
}

type TextureObject struct {
	glx   *Context
	value driver.Value
}

func (glx *Context) CreateTexture() TextureObject {
	value := glx.constants.CreateTexture()
	return TextureObject{
		glx:   glx,
		value: value,
	}
}

type FeedbackObject struct {
	glx   *Context
	value driver.Value
}

func (glx *Context) CreateFeedback() FeedbackObject {
	feedbackObject := glx.constants.CreateTransformFeedback()
	return FeedbackObject{
		glx:   glx,
		value: feedbackObject,
	}
}

func (tfo FeedbackObject) Destroy() {
	glx := tfo.glx
	glx.constants.DeleteTransformFeedback(tfo.value)
}

func (tfo FeedbackObject) Begin(m PrimitiveDrawMode) {
	glx := tfo.glx
	var jsMode driver.Value
	switch m {
	case Points:
		jsMode = glx.constants.POINTS
	case Lines:
		jsMode = glx.constants.LINES
	case Triangles:
		jsMode = glx.constants.TRIANGLES
	}
	glx.constants.BeginTransformFeedback(jsMode)
}

func (tfo FeedbackObject) End() {
	glx := tfo.glx
	glx.constants.EndTransformFeedback()
}