package gl

import (
	"fmt"

	"github.com/PieterD/warp/pkg/driver"
)

type FramebufferTarget struct {
	glx *Context
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
