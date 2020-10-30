package gl

import "github.com/PieterD/warp/driver"

type Buffer struct {
	glx      *Context
	glObject driver.Value
	jsObject driver.Object
}

func newBuffer(glx *Context) (*Buffer, error) {
	bufferObject := glx.functions.CreateBuffer()
	return &Buffer{
		glx:      glx,
		glObject: bufferObject,
	}, nil
}
