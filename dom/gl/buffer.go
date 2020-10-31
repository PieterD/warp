package gl

import "github.com/PieterD/warp/driver"

type Buffer struct {
	glx         *Context
	glObject    driver.Value
	currentType driver.Value
}

func newBuffer(glx *Context) (*Buffer, error) {
	bufferObject := glx.functions.CreateBuffer()
	return &Buffer{
		glx:      glx,
		glObject: bufferObject,
	}, nil
}

func (b *Buffer) VertexData(data []float32) {
	byteData := fastFloat32ToByte(data)
	jsBuffer := b.glx.factory.Buffer(len(byteData))
	jsBuffer.Put(byteData)
	vertexArray := jsBuffer.AsFloat32Array()
	glx := b.glx

	glx.functions.BindBuffer(glx.constants.ARRAY_BUFFER, b.glObject)
	glx.functions.BufferData(glx.constants.ARRAY_BUFFER, vertexArray, glx.constants.STATIC_DRAW)
	glx.functions.BindBuffer(glx.constants.ARRAY_BUFFER, glx.factory.Null())
}

func (b *Buffer) IndexData(data []uint16) {
	byteData := fastUint16ToByte(data)
	jsBuffer := b.glx.factory.Buffer(len(byteData))
	jsBuffer.Put(byteData)
	indexArray := jsBuffer.AsUint16Array()
	glx := b.glx

	glx.functions.BindBuffer(glx.constants.ELEMENT_ARRAY_BUFFER, b.glObject)
	glx.functions.BufferData(glx.constants.ELEMENT_ARRAY_BUFFER, indexArray, glx.constants.STATIC_DRAW)
	glx.functions.BindBuffer(glx.constants.ELEMENT_ARRAY_BUFFER, glx.factory.Null())
}
