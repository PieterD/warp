package main

import (
	"fmt"

	"github.com/PieterD/warp/pkg/gl"
	"github.com/PieterD/warp/pkg/gl/glunsafe"
	"github.com/go-gl/mathgl/mgl32"
)

func gltUniformBlock(glx *gl.Context, _ gl.FramebufferObject) error {
	program := glx.CreateProgram()
	defer program.Destroy()

	vShader := glx.CreateShader(gl.VertexShader)
	defer vShader.Destroy()
	vShader.Source(`#version 300 es
precision mediump float;
layout (location = 0) in vec3 Position;

void main(void) {
	gl_Position = vec4(Position, 1.0);
}`)

	fShader := glx.CreateShader(gl.FragmentShader)
	defer fShader.Destroy()
	fShader.Source(`#version 300 es
precision mediump float;
out vec4 FragColor;

uniform Uniform {
	vec4 Color;
};

void main(void) {
	FragColor = Color;
}`)

	vShader.Compile()
	fShader.Compile()
	program.Attach(vShader)
	program.Attach(fShader)
	program.Link()
	if !program.LinkSuccess() {
		glx.Log("vert shader log: %s", vShader.InfoLog())
		glx.Log("frag shader log: %s", fShader.InfoLog())
		glx.Log("prog linker log: %s", program.InfoLog())
		return fmt.Errorf("complex error linking program, see log")
	}

	uniformBlockIndex := program.GetUniformBlockIndex("Uniform")
	const uniformBufferIndex = 5
	program.UniformBlockBinding(uniformBlockIndex, uniformBufferIndex)

	color := mgl32.Vec4{
		0.0, 1.0, 0.0, 1.0,
	}
	uniformBuffer := glx.CreateBuffer()
	defer uniformBuffer.Destroy()
	uniformData := glunsafe.Map(color[:])
	glx.Targets().Uniform().Bind(uniformBuffer)
	glx.Targets().Uniform().BufferData(uniformData, gl.Static, gl.Draw)
	glx.Targets().Uniform().Unbind()

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
	glx.Targets().Array().BufferData(vData, gl.Static, gl.Draw)
	glx.Targets().Array().UnbindBuffer()

	iData := glunsafe.Map(indices)
	iBuffer := glx.CreateBuffer()
	defer iBuffer.Destroy()
	glx.Targets().ElementArray().BindBuffer(iBuffer)
	glx.Targets().ElementArray().BufferData(iData, gl.Static, gl.Draw)
	glx.Targets().ElementArray().UnbindBuffer()

	vao := glx.CreateVertexArray()
	defer vao.Destroy()
	glx.BindVertexArray(vao)
	glx.Targets().Array().BindBuffer(vBuffer)
	vao.VertexAttribPointer(0, gl.Vec3, false, 3*4, 0)
	vao.EnableVertexAttribArray(0)
	glx.Targets().Array().UnbindBuffer()
	glx.UnbindVertexArray()

	glx.ClearColor(0.75, 0.8, 0.85, 1.0)
	glx.Clear()
	glx.UseProgram(program)
	glx.BindVertexArray(vao)
	glx.Targets().Uniform().BindBase(uniformBufferIndex, uniformBuffer)
	glx.Targets().ElementArray().BindBuffer(iBuffer)
	glx.DrawElements(gl.Triangles, 3, gl.UnsignedShort, 0)
	glx.Targets().ElementArray().UnbindBuffer()
	glx.Targets().Uniform().UnbindBase(uniformBufferIndex)
	glx.UnuseProgram()

	return nil
}
