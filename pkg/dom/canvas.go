package dom

import (
	"fmt"

	"github.com/PieterD/warp/pkg/dom/gl"
)

type Canvas struct {
	elem *Elem
}

func AsCanvas(elem *Elem) *Canvas {
	if elem.Tag() != "canvas" {
		return nil
	}
	return &Canvas{
		elem: elem,
	}
}

func (c *Canvas) GetContextWebgl() *gl.Context {
	return gl.NewContext(c.elem)
}

func (c *Canvas) InnerSize() (width, height int) {
	fWidth, ok := c.elem.obj.Get("width").ToFloat64()
	if !ok {
		panic(fmt.Errorf("canvas width is not a number"))
	}
	fHeight, ok := c.elem.obj.Get("height").ToFloat64()
	if !ok {
		panic(fmt.Errorf("canvas height is not a number"))
	}
	return int(fWidth), int(fHeight)
}

func (c *Canvas) SetInnerSize(width, height int) {
	c.elem.obj.Set("width", c.elem.factory.Number(float64(width)))
	c.elem.obj.Set("height", c.elem.factory.Number(float64(height)))
}

func (c *Canvas) OuterSize() (width, height int) {
	fWidth, ok := c.elem.obj.Get("clientWidth").ToFloat64()
	if !ok {
		panic(fmt.Errorf("canvas clientWidth is not a number"))
	}
	fHeight, ok := c.elem.obj.Get("clientHeight").ToFloat64()
	if !ok {
		panic(fmt.Errorf("canvas clientHeight is not a number"))
	}
	return int(fWidth), int(fHeight)
}
