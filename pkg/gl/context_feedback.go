package gl

import "github.com/PieterD/warp/pkg/driver"

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
