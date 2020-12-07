package gl

import "github.com/PieterD/warp/pkg/driver"

// https://developer.mozilla.org/en-US/docs/Web/API/WebGLTransformFeedback

type Feedback struct {
	glx      *Context
	glObject driver.Value
}

type FeedbackAttribute struct {
	Name   string
	Type   Type
	Buffer *Buffer
	Layout FeedbackAttributeLayout
}

type FeedbackAttributeLayout struct {
	ByteOffset int
	ByteStride int
}

func newFeedback(glx *Context, attrs ...FeedbackAttribute) *Feedback {
	feedbackObject := glx.constants.CreateTransformFeedback()
	return &Feedback{
		glx:      glx,
		glObject: feedbackObject,
	}
}

func (f *Feedback) begin(m PrimitiveDrawMode) {
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
	glx.constants.BeginTransformFeedback(jsMode)
}

func (f *Feedback) end() {
	glx := f.glx
	glx.constants.EndTransformFeedback()
}
