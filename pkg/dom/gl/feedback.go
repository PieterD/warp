package gl

import (
	"fmt"

	"github.com/PieterD/warp/pkg/driver"
)

type feedback struct {
	glx         *Context
	glObject    driver.Value
	bufferNames []string
}

func newFeedback(glx *Context, coupling ActiveCoupling) *feedback {
	feedbackObject := glx.constants.CreateTransformFeedback()
	return &feedback{
		glx:         glx,
		glObject:    feedbackObject,
		bufferNames: coupling.BufferNames(),
	}
}

func (f *feedback) begin(m PrimitiveDrawMode, buffers map[string]*Buffer) error {
	glx := f.glx
	var jsMode driver.Value
	switch m {
	case Points:
		jsMode = glx.constants.POINTS
	case Lines:
		jsMode = glx.constants.LINES
	case Triangles:
		jsMode = glx.constants.TRIANGLES
	}
	for index, bufferName := range f.bufferNames {
		buffer, ok := buffers[bufferName]
		if !ok {
			return fmt.Errorf("missing buffer %s", bufferName)
		}
		glx.constants.BindBufferBase(
			glx.constants.TRANSFORM_FEEDBACK,
			glx.factory.Number(float64(index)),
			buffer.glObject,
		)
	}
	glx.constants.BeginTransformFeedback(jsMode)
	return nil
}

func (f *feedback) end() {
	glx := f.glx
	glx.constants.EndTransformFeedback()
	for index, _ := range f.bufferNames {
		glx.constants.BindBufferBase(
			glx.constants.TRANSFORM_FEEDBACK,
			glx.factory.Number(float64(index)),
			glx.factory.Null(),
		)
	}
}
