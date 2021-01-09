package main

import (
	"fmt"

	"github.com/PieterD/warp/pkg/gl"
	"github.com/PieterD/warp/pkg/gl/glunsafe"
)

func gltInstancedQuads(glx *gl.Context, _ gl.FramebufferObject) error {
	var (
		vSource = `#version 300 es
precision mediump float;
layout (location = 0) in vec3 Position;

layout (location = 1) in vec3 Translation;
layout (location = 2) in float Scale;
layout (location = 3) in vec3 Color;
out vec4 color;

void main(void) {
	gl_Position = vec4(Position*Scale + Translation, 1.0);
	color = vec4(Color, 1.0);
}`
		fSource = `#version 300 es
precision mediump float;
in vec4 color;
out vec4 FragColor;

void main(void) {
	FragColor = color;
}`
		quadVertices = []float32{
			-0.5, -0.5, 0.0,
			0.5, -0.5, 0.0,
			0.5, 0.5, 0.0,
			-0.5, 0.5, 0.0,
		}
		quadIndices = []uint16{
			0, 1, 2,
			2, 3, 0,
		}
		instanceData = []float32{
			-0.5, -0.5, 0.0, 0.2, 1.0, 0.0, 0.0,
			0.5, -0.5, 0.0, 0.3, 0.0, 1.0, 0.0,
			0.5, 0.5, 0.0, 0.4, 0.0, 0.0, 1.0,
			-0.5, 0.5, 0.0, 0.5, 1.0, 1.0, 1.0,
		}
	)
	program := glx.CreateProgram()
	defer program.Destroy()

	vShader := glx.CreateShader(gl.VertexShader)
	defer vShader.Destroy()
	vShader.Source(vSource)

	fShader := glx.CreateShader(gl.FragmentShader)
	defer fShader.Destroy()
	fShader.Source(fSource)

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

	quadBuffer := glx.CreateBuffer()
	defer quadBuffer.Destroy()
	glx.Targets().Array().BindBuffer(quadBuffer)
	glx.Targets().Array().BufferData(glunsafe.Map(quadVertices), gl.Static, gl.Draw)
	glx.Targets().Array().UnbindBuffer()

	indexBuffer := glx.CreateBuffer()
	defer indexBuffer.Destroy()
	glx.Targets().ElementArray().BindBuffer(indexBuffer)
	glx.Targets().ElementArray().BufferData(glunsafe.Map(quadIndices), gl.Static, gl.Draw)
	glx.Targets().ElementArray().UnbindBuffer()

	instanceBuffer := glx.CreateBuffer()
	defer instanceBuffer.Destroy()
	glx.Targets().Array().BindBuffer(instanceBuffer)
	glx.Targets().Array().BufferData(glunsafe.Map(instanceData), gl.Static, gl.Draw)
	glx.Targets().Array().UnbindBuffer()

	vao := glx.CreateVertexArray()
	defer vao.Destroy()
	glx.BindVertexArray(vao)
	glx.Targets().Array().BindBuffer(quadBuffer)
	vao.VertexAttribPointer(0, gl.Vec3, false, 3*4, 0)
	vao.EnableVertexAttribArray(0)
	glx.Targets().Array().BindBuffer(instanceBuffer)
	vao.VertexAttribPointer(1, gl.Vec3, false, 7*4, 0)
	vao.VertexAttribDivisor(1, 1)
	vao.EnableVertexAttribArray(1)
	vao.VertexAttribPointer(2, gl.Float, false, 7*4, 3*4)
	vao.VertexAttribDivisor(2, 1)
	vao.EnableVertexAttribArray(2)
	vao.VertexAttribPointer(3, gl.Vec3, false, 7*4, 4*4)
	vao.VertexAttribDivisor(3, 1)
	vao.EnableVertexAttribArray(3)
	glx.Targets().Array().UnbindBuffer()
	glx.UnbindVertexArray()

	glx.ClearColor(0.75, 0.8, 0.85, 1.0)
	glx.Clear()
	glx.UseProgram(program)
	glx.BindVertexArray(vao)
	glx.Targets().ElementArray().BindBuffer(indexBuffer)
	glx.DrawElementsInstanced(gl.Triangles, 6, gl.UnsignedShort, 0, 4)
	glx.Targets().ElementArray().UnbindBuffer()
	glx.UnuseProgram()

	return nil
}
