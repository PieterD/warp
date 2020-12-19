package main

import (
	"fmt"
	"image"

	"github.com/PieterD/warp/pkg/dom/gl/raw"
	"github.com/PieterD/warp/pkg/dom/glunsafe"
)

func init() {
	Register("triangle", "Render a triangle", triangle)
}

func triangle(glx *raw.Context, fbo raw.FramebufferObject) error {
	program := glx.CreateProgram()
	defer program.Destroy()

	vShader := glx.CreateShader(raw.VertexShader)
	defer vShader.Destroy()
	vShader.Source(`#version 300 es
precision mediump float;
layout (location = 0) in vec3 Position;

void main(void) {
	gl_Position = vec4(Position, 1.0);
}`)
	if err := vShader.Compile(); err != nil {
		return fmt.Errorf("compiling vertex shader: %w", err)
	}

	fShader := glx.CreateShader(raw.FragmentShader)
	defer fShader.Destroy()
	fShader.Source(`#version 300 es
precision mediump float;
out vec4 FragColor;

void main(void) {
	FragColor = vec4(0.0, 0.5, 1.0, 1.0);
}`)
	if err := fShader.Compile(); err != nil {
		return fmt.Errorf("compiling fragment shader: %w", err)
	}

	program.Attach(vShader)
	program.Attach(fShader)
	if err := program.Link(); err != nil {
		return fmt.Errorf("linking program: %w", err)
	}

	vertices := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}
	indices := []uint16{
		0, 1, 2,
	}
	vData := glunsafe.Map(vertices)
	vBuffer := glx.CreateBuffer()
	defer vBuffer.Destroy()
	glx.Targets().Array().BindBuffer(vBuffer)
	glx.Targets().Array().BufferData(vData, raw.Static, raw.Draw)
	glx.Targets().Array().UnbindBuffer()

	iData := glunsafe.Map(indices)
	iBuffer := glx.CreateBuffer()
	defer iBuffer.Destroy()
	glx.Targets().ElementArray().BindBuffer(iBuffer)
	glx.Targets().ElementArray().BufferData(iData, raw.Static, raw.Draw)
	glx.Targets().ElementArray().UnbindBuffer()

	vao := glx.CreateVertexArray()
	defer vao.Destroy()
	glx.BindVertexArray(vao)
	glx.Targets().Array().BindBuffer(vBuffer)
	vao.VertexAttribPointer(0, raw.Vec3, false, 3*4, 0)
	vao.EnableVertexAttribArray(0)
	glx.Targets().Array().UnbindBuffer()
	glx.UnbindVertexArray()

	glx.ClearColor(0.75, 0.8, 0.85, 1.0)
	glx.Clear()
	glx.UseProgram(program)
	glx.BindVertexArray(vao)
	glx.Targets().ElementArray().BindBuffer(iBuffer)
	glx.DrawElements(raw.Triangles, 0, 3, raw.UnsignedShort)
	glx.Targets().ElementArray().UnbindBuffer()
	glx.UnuseProgram()

	return nil
}

func pixelsToImage(width, height int, pixels []byte) image.Image {
	flippedPixels := make([]byte, len(pixels))
	rowSize := width * 4
	for y := 0; y < height; y++ {
		fy := height - 1 - y
		copy(flippedPixels[fy*rowSize:(fy+1)*rowSize], pixels[y*rowSize:(y+1)*rowSize])
	}
	img := &image.NRGBA{
		Pix:    flippedPixels,
		Stride: width * 4,
		Rect: image.Rectangle{
			Min: image.Point{
				X: 0,
				Y: 0,
			},
			Max: image.Point{
				X: width,
				Y: height,
			},
		},
	}
	return img
}
