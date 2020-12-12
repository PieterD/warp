package gl

import (
	"github.com/PieterD/warp/pkg/dom/glunsafe"
	"github.com/PieterD/warp/pkg/driver"
)

type Buffer struct {
	glx         *Context
	glObject    driver.Value
	currentType driver.Value
}

func newBuffer(glx *Context) *Buffer {
	bufferObject := glx.constants.CreateBuffer()
	return &Buffer{
		glx:      glx,
		glObject: bufferObject,
	}
}

func (b *Buffer) VertexData(data []float32) {
	byteData := glunsafe.FastFloat32ToByte(data)
	jsBuffer := b.glx.factory.Buffer(len(byteData))
	jsBuffer.Put(byteData)
	vertexArray := jsBuffer.AsFloat32Array()
	glx := b.glx

	glx.constants.BindBuffer(glx.constants.ARRAY_BUFFER, b.glObject)
	glx.constants.BufferData(glx.constants.ARRAY_BUFFER, vertexArray, glx.constants.STATIC_DRAW)
	glx.constants.BindBuffer(glx.constants.ARRAY_BUFFER, glx.factory.Null())
}

func (b *Buffer) IndexData(data []uint16) {
	byteData := glunsafe.FastUint16ToByte(data)
	jsBuffer := b.glx.factory.Buffer(len(byteData))
	jsBuffer.Put(byteData)
	indexArray := jsBuffer.AsUint16Array()
	glx := b.glx

	glx.constants.BindBuffer(glx.constants.ELEMENT_ARRAY_BUFFER, b.glObject)
	glx.constants.BufferData(glx.constants.ELEMENT_ARRAY_BUFFER, indexArray, glx.constants.STATIC_DRAW)
	glx.constants.BindBuffer(glx.constants.ELEMENT_ARRAY_BUFFER, glx.factory.Null())
}

func (b *Buffer) UniformData(byteData []byte) {
	jsBuffer := b.glx.factory.Buffer(len(byteData))
	jsBuffer.Put(byteData)
	indexArray := jsBuffer.AsUint16Array()
	glx := b.glx

	glx.constants.BindBuffer(glx.constants.UNIFORM_BUFFER, b.glObject)
	glx.constants.BufferData(glx.constants.UNIFORM_BUFFER, indexArray, glx.constants.STATIC_DRAW)
	glx.constants.BindBuffer(glx.constants.UNIFORM_BUFFER, glx.factory.Null())
}
