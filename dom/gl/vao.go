package gl

import "github.com/PieterD/warp/driver"

type VertexArrayConfig struct {
}

type VertexArray struct {
	glx      *Context
	glObject driver.Value
}

func newVertexArray(glx *Context, cfg VertexArrayConfig) (*VertexArray, error) {
	glVAO := glx.functions.CreateVertexArray()
	return &VertexArray{
		glx:      glx,
		glObject: glVAO,
	}, nil
}
