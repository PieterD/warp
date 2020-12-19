package raw

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
	jsBuffer := glx.factory.Buffer(len(data))
	jsBuffer.Put(data)
	jsByteArray := jsBuffer.AsUint8Array()
	glUsage := combineUsage(glx, accessUsage, modificationUsage)
	glx.constants.BufferData(target.which, jsByteArray, glUsage)
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

type FramebufferTarget struct {
	glx *Context
}

func (targets Targets) FrameBuffer() FramebufferTarget {
	return FramebufferTarget{
		glx: targets.glx,
	}
}

func (target FramebufferTarget) Bind(fbo FramebufferObject) {
	glx := target.glx
	glx.constants.BindFramebuffer(
		glx.constants.FRAMEBUFFER,
		fbo.value,
	)
}

func (target FramebufferTarget) Unbind() {
	glx := target.glx
	glx.constants.BindFramebuffer(
		glx.constants.FRAMEBUFFER,
		glx.factory.Null(),
	)
}

func (target FramebufferTarget) AttachRenderbuffer(attachmentType RenderbufferType, rbo RenderbufferObject) {
	glx := target.glx
	var glType driver.Value
	switch attachmentType {
	case ColorBuffer:
		glType = glx.constants.COLOR_ATTACHMENT0
	case DepthStencilBuffer:
		glType = glx.constants.DEPTH_STENCIL_ATTACHMENT
	default:
		panic(fmt.Errorf("invalid renderbuffer type: %v", attachmentType))
	}
	glx.constants.FramebufferRenderbuffer(
		glx.constants.FRAMEBUFFER,
		glType,
		glx.constants.RENDERBUFFER,
		rbo.value,
	)
}

func (target FramebufferTarget) IsComplete() bool {
	glx := target.glx
	fbsJs := glx.constants.CheckFramebufferStatus()
	fbsFloat, ok := fbsJs.ToFloat64()
	if !ok {
		panic(fmt.Errorf("CheckFramebufferStatus return value was not a number: %T", fbsJs))
	}
	completeFloat, ok := glx.constants.FRAMEBUFFER_COMPLETE.ToFloat64()
	if !ok {
		panic(fmt.Errorf("FRAMEBUFFER_COMPLETE was not a number: %T", glx.constants.FRAMEBUFFER_COMPLETE))
	}
	return fbsFloat == completeFloat
}

func (target FramebufferTarget) ReadPixels(x, y, w, h int) []byte {
	glx := target.glx
	pixelDataSize := w * h * 4
	jsBuffer := glx.factory.Buffer(pixelDataSize)
	jsArray := jsBuffer.AsUint8Array()
	glx.constants.ReadPixels(
		glx.factory.Number(float64(x)),
		glx.factory.Number(float64(y)),
		glx.factory.Number(float64(w)),
		glx.factory.Number(float64(h)),
		glx.constants.RGBA,
		glx.constants.UNSIGNED_BYTE,
		jsArray,
		glx.factory.Number(0),
	)
	data := make([]byte, pixelDataSize)
	num := jsBuffer.Get(data)
	if num != pixelDataSize {
		panic(fmt.Errorf("expected jsBuffer.Get to return %d, got %d", pixelDataSize, num))
	}
	return data
}
