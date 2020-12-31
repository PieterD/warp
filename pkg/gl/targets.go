package gl

import (
	"fmt"
	"github.com/PieterD/warp/pkg/driver"
)

type ModificationUsage int

const (
	// The data store contents are modified by the application, and used as the source for GL drawing and image specification commands.
	Draw ModificationUsage = iota + 1
	// The data store contents are modified by reading data from the GL, and used to return that data when queried by the application.
	Read
	// The data store contents are modified by reading data from the GL, and used as the source for GL drawing and image specification commands.
	Copy
)

type AccessUsage int

const (
	// The data store contents will be modified once and used at most a few times.
	Stream AccessUsage = iota + 1
	// The data store contents will be modified once and used many times.
	Static
	// The data store contents will be modified repeatedly and used many times.
	Dynamic
)

type ArrayTarget struct {
	glx     *Context
	which   driver.Value
	indexed bool
}

func (target ArrayTarget) BindBuffer(buffer BufferObject) {
	glx := target.glx
	glx.constants.BindBuffer(target.which, buffer.value)
}

func (target ArrayTarget) UnbindBuffer() {
	glx := target.glx
	glx.constants.BindBuffer(target.which, glx.factory.Null())
}

func (target ArrayTarget) BufferData(data []byte, accessUsage AccessUsage, modificationUsage ModificationUsage) {
	glx := target.glx
	bufferData(glx, target.which, data, accessUsage, modificationUsage)
}

func bufferData(glx *Context, target driver.Value, data []byte, accessUsage AccessUsage, modificationUsage ModificationUsage) {
	jsBuffer := glx.factory.Buffer(len(data))
	jsBuffer.Put(data)
	jsByteArray := jsBuffer.AsUint8Array()
	glUsage := combineUsage(glx, accessUsage, modificationUsage)
	glx.constants.BufferData(target, jsByteArray, glUsage)
}

func combineUsage(glx *Context, accessUsage AccessUsage, modificationUsage ModificationUsage) driver.Value {
	a := accessUsage
	m := modificationUsage
	if a == Stream {
		if m == Draw {
			return glx.constants.STREAM_DRAW
		}
		if m == Read {
			return glx.constants.STREAM_READ
		}
		if m == Copy {
			return glx.constants.STREAM_COPY
		}
	}
	if a == Static {
		if m == Draw {
			return glx.constants.STATIC_DRAW
		}
		if m == Read {
			return glx.constants.STATIC_READ
		}
		if m == Copy {
			return glx.constants.STATIC_COPY
		}
	}
	if a == Dynamic {
		if m == Draw {
			return glx.constants.DYNAMIC_DRAW
		}
		if m == Read {
			return glx.constants.DYNAMIC_READ
		}
		if m == Copy {
			return glx.constants.DYNAMIC_COPY
		}
	}
	panic(fmt.Errorf("invalid usage: %v %v", accessUsage, modificationUsage))
}

func (target ArrayTarget) IsIndexed() bool {
	return target.indexed
}

func (target ArrayTarget) BindBufferBase(index int, buffer BufferObject) {
	panic("not implemented")
}

type Targets struct {
	glx *Context
}

func (targets Targets) Array() ArrayTarget {
	glx := targets.glx
	return ArrayTarget{
		glx:     glx,
		which:   glx.constants.ARRAY_BUFFER,
		indexed: false,
	}
}

func (targets Targets) ElementArray() ArrayTarget {
	glx := targets.glx
	return ArrayTarget{
		glx:     glx,
		which:   glx.constants.ELEMENT_ARRAY_BUFFER,
		indexed: false,
	}
}

type RenderbufferTarget struct {
	glx *Context
}

func (targets Targets) RenderBuffer() RenderbufferTarget {
	return RenderbufferTarget{
		glx: targets.glx,
	}
}

func (target RenderbufferTarget) Bind(rbo RenderbufferObject) {
	glx := target.glx
	glx.constants.BindRenderbuffer(
		glx.constants.RENDERBUFFER,
		rbo.value,
	)
}

func (target RenderbufferTarget) Unbind() {
	glx := target.glx
	glx.constants.BindRenderbuffer(
		glx.constants.RENDERBUFFER,
		glx.factory.Null(),
	)
}

type RenderbufferConfig struct {
	Type   RenderbufferType
	Width  int
	Height int
	//Samples int
}

//go:generate stringer -type=RenderbufferType
type RenderbufferType int

const (
	ColorBuffer RenderbufferType = iota + 1
	DepthStencilBuffer
)

//TODO: multisampling. To read multisampled pixels, we first have to resolve the
// samples by blitting the multisample framebuffer into a regular one.
// See BlitFramebuffer.
func (target RenderbufferTarget) Storage(cfg RenderbufferConfig) {
	glx := target.glx
	var glType driver.Value
	switch cfg.Type {
	case ColorBuffer:
		glType = glx.constants.RGBA8
	case DepthStencilBuffer:
		glType = glx.constants.DEPTH24_STENCIL8
	default:
		panic(fmt.Errorf("invalid renderbuffer type: %v", cfg.Type))
	}
	glx.constants.RenderbufferStorageMultisample(
		glx.constants.RENDERBUFFER,
		glx.factory.Number(float64(0)),
		glType,
		glx.factory.Number(float64(cfg.Width)),
		glx.factory.Number(float64(cfg.Height)),
	)
}

func (targets Targets) Framebuffer() FramebufferTarget {
	return FramebufferTarget{
		glx: targets.glx,
	}
}

func (targets Targets) Texture2D() Texture2DTarget {
	return Texture2DTarget{
		glx: targets.glx,
	}
}

func (targets Targets) ActiveTextureUnit(unit int) {
	glx := targets.glx
	fTexture0, ok := glx.constants.TEXTURE0.ToFloat64()
	if !ok {
		panic(fmt.Errorf("expected TEXTURE0 to be a number: %T", glx.constants.TEXTURE0))
	}
	t0 := int(fTexture0)
	jsTextureUnit := glx.factory.Number(float64(t0 + unit))
	glx.constants.ActiveTexture(jsTextureUnit)
}

type UniformTarget struct {
	glx *Context
}

func (targets Targets) Uniform() UniformTarget {
	return UniformTarget{
		glx: targets.glx,
	}
}

func (target UniformTarget) Bind(buffer BufferObject) {
	glx := target.glx
	glx.constants.BindBuffer(
		glx.constants.UNIFORM_BUFFER,
		buffer.value,
	)
}

func (target UniformTarget) Unbind() {
	glx := target.glx
	glx.constants.BindBuffer(
		glx.constants.UNIFORM_BUFFER,
		glx.factory.Null(),
	)
}

// also Binds
func (target UniformTarget) BindBase(index int, buffer BufferObject) {
	glx := target.glx
	glx.constants.BindBufferBase(
		glx.constants.UNIFORM_BUFFER,
		glx.factory.Number(float64(index)),
		buffer.value,
	)
}

// also Unbinds
func (target UniformTarget) UnbindBase(index int) {
	glx := target.glx
	glx.constants.BindBufferBase(
		glx.constants.UNIFORM_BUFFER,
		glx.factory.Number(float64(index)),
		glx.factory.Null(),
	)
}

func (target UniformTarget) BufferData(data []byte, accessUsage AccessUsage, modificationUsage ModificationUsage) {
	glx := target.glx
	bufferData(glx, glx.constants.UNIFORM_BUFFER, data, accessUsage, modificationUsage)
}

type TransformFeedbackTarget struct {
	glx *Context
}

func (targets Targets) TransformFeedback() TransformFeedbackTarget {
	glx := targets.glx
	return TransformFeedbackTarget{
		glx: glx,
	}
}

func (target TransformFeedbackTarget) Bind(buffer BufferObject) {
	glx := target.glx
	glx.constants.BindBuffer(glx.constants.TRANSFORM_FEEDBACK_BUFFER, buffer.value)
}

func (target TransformFeedbackTarget) Unbind() {
	glx := target.glx
	glx.constants.BindBuffer(glx.constants.TRANSFORM_FEEDBACK_BUFFER, glx.factory.Null())
}

func (target TransformFeedbackTarget) BindBase(index int, buffer BufferObject) {
	glx := target.glx
	glx.constants.BindBufferBase(
		glx.constants.TRANSFORM_FEEDBACK_BUFFER,
		glx.factory.Number(float64(index)),
		buffer.value,
	)
}

func (target TransformFeedbackTarget) UnbindBase(index int) {
	glx := target.glx
	glx.constants.BindBufferBase(
		glx.constants.TRANSFORM_FEEDBACK_BUFFER,
		glx.factory.Number(float64(index)),
		glx.factory.Null(),
	)
}

func (target TransformFeedbackTarget) Alloc(size int, accessUsage AccessUsage, modificationUsage ModificationUsage) {
	glx := target.glx
	bufferData(glx, glx.constants.TRANSFORM_FEEDBACK_BUFFER, make([]byte, size), accessUsage, modificationUsage)
}

func (target TransformFeedbackTarget) Contents(data []byte) int {
	glx := target.glx
	if len(data) == 0 {
		return 0
	}
	jsBuffer := glx.factory.Buffer(len(data))
	jsArray := jsBuffer.AsUint8Array()
	glx.constants.GetBufferSubData(
		glx.constants.TRANSFORM_FEEDBACK_BUFFER,
		glx.factory.Number(float64(0)),
		jsArray,
	)
	return jsBuffer.Get(data)
}
