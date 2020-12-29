package main

import (
	"fmt"
	"github.com/PieterD/warp/pkg/dom/gl/raw"
	"github.com/PieterD/warp/pkg/dom/glunsafe"
)

func gltPoint(glx *raw.Context, _ raw.FramebufferObject) error {
	program := glx.CreateProgram()
	defer program.Destroy()

	vShader := glx.CreateShader(raw.VertexShader)
	defer vShader.Destroy()
	vShader.Source(`#version 300 es
precision mediump float;
layout (location = 0) in vec2 Position;
layout (location = 1) in vec3 Color;
layout (location = 2) in float Size;
out vec4 fragColor;

void main(void) {
	gl_Position = vec4(Position, 0.0, 1.0);
	gl_PointSize = Size;
	fragColor = vec4(Color, 1.0);
}`)
	if err := vShader.Compile(); err != nil {
		return fmt.Errorf("compiling vertex shader: %w", err)
	}

	fShader := glx.CreateShader(raw.FragmentShader)
	defer fShader.Destroy()
	fShader.Source(`#version 300 es
precision mediump float;
in vec4 fragColor;
out vec4 FragColor;

void main(void) {
	FragColor = fragColor;
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
		-0.5, -0.5, 1.0, 0.0, 0.0, 10.0,
		0.5, -0.5, 0.0, 1.0, 0.0, 15.0,
		0.0, 0.5, 0.0, 0.0, 1.0, 20.0,
	}
	vBuffer := glx.CreateBuffer()
	defer vBuffer.Destroy()
	glx.Targets().Array().BindBuffer(vBuffer)
	glx.Targets().Array().BufferData(glunsafe.Map(vertices), raw.Static, raw.Draw)
	glx.Targets().Array().UnbindBuffer()

	vao := glx.CreateVertexArray()
	defer vao.Destroy()
	glx.BindVertexArray(vao)
	glx.Targets().Array().BindBuffer(vBuffer)
	stride := 6 * 4
	offset := 0
	vao.VertexAttribPointer(0, raw.Vec2, false, stride, 0)
	vao.EnableVertexAttribArray(0)
	offset += 2 * 4
	vao.VertexAttribPointer(1, raw.Vec3, false, stride, offset)
	vao.EnableVertexAttribArray(1)
	offset += 3 * 4
	vao.VertexAttribPointer(2, raw.Float, false, stride, offset)
	vao.EnableVertexAttribArray(2)
	offset += 1 * 4
	glx.UnbindVertexArray()

	glx.ClearColor(0.75, 0.8, 0.85, 1.0)
	glx.Clear()
	glx.UseProgram(program)
	glx.BindVertexArray(vao)
	glx.DrawArrays(raw.Points, 0, 3)
	glx.Targets().ElementArray().UnbindBuffer()
	glx.UnuseProgram()

	return nil
}
