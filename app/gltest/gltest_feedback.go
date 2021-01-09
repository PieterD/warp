package main

import (
	"fmt"
	"reflect"

	"github.com/PieterD/warp/pkg/gl"
	"github.com/PieterD/warp/pkg/gl/glunsafe"
)

func gltFeedback(glx *gl.Context, _ gl.FramebufferObject) error {
	program := glx.CreateProgram()
	defer program.Destroy()

	vShader := glx.CreateShader(gl.VertexShader)
	defer vShader.Destroy()
	vShader.Source(`#version 300 es
precision mediump float;
layout (location = 0) in vec2 Position;
layout (location = 1) in vec3 Color;
layout (location = 2) in float Size;
out float FeedbackSize;
out vec4 fragColor;

void main(void) {
	gl_Position = vec4(Position, 0.0, 1.0);
	gl_PointSize = Size;
	fragColor = vec4(Color, 1.0);
	FeedbackSize = Size;
}`)

	fShader := glx.CreateShader(gl.FragmentShader)
	defer fShader.Destroy()
	fShader.Source(`#version 300 es
precision mediump float;
in vec4 fragColor;
out vec4 FragColor;

void main(void) {
	FragColor = fragColor;
}`)
	vShader.Compile()
	fShader.Compile()
	program.Attach(vShader)
	program.Attach(fShader)
	program.TransformFeedbackVaryings(true, "FeedbackSize")
	program.Link()
	if !program.LinkSuccess() {
		glx.Log("vert shader log: %s", vShader.InfoLog())
		glx.Log("frag shader log: %s", fShader.InfoLog())
		glx.Log("prog linker log: %s", program.InfoLog())
		return fmt.Errorf("complex error linking program, see log")
	}

	tfBuffer := glx.CreateBuffer()
	defer tfBuffer.Destroy()
	glx.Targets().TransformFeedback().Bind(tfBuffer)
	glx.Targets().TransformFeedback().Alloc(3*4, gl.Static, gl.Draw)
	glx.Targets().TransformFeedback().Unbind()

	feedback := glx.CreateFeedback()
	defer feedback.Destroy()

	vertices := []float32{
		-0.5, -0.5, 1.0, 0.0, 0.0, 10.0,
		0.5, -0.5, 0.0, 1.0, 0.0, 15.0,
		0.0, 0.5, 0.0, 0.0, 1.0, 20.0,
	}
	vBuffer := glx.CreateBuffer()
	defer vBuffer.Destroy()
	glx.Targets().Array().BindBuffer(vBuffer)
	glx.Targets().Array().BufferData(glunsafe.Map(vertices), gl.Static, gl.Draw)
	glx.Targets().Array().UnbindBuffer()

	vao := glx.CreateVertexArray()
	defer vao.Destroy()
	glx.BindVertexArray(vao)
	glx.Targets().Array().BindBuffer(vBuffer)
	stride := 6 * 4
	offset := 0
	vao.VertexAttribPointer(0, gl.Vec2, false, stride, 0)
	vao.EnableVertexAttribArray(0)
	offset += 2 * 4
	vao.VertexAttribPointer(1, gl.Vec3, false, stride, offset)
	vao.EnableVertexAttribArray(1)
	offset += 3 * 4
	vao.VertexAttribPointer(2, gl.Float, false, stride, offset)
	vao.EnableVertexAttribArray(2)
	offset += 1 * 4
	glx.UnbindVertexArray()

	glx.ClearColor(0.75, 0.8, 0.85, 1.0)
	glx.Clear()
	glx.UseProgram(program)
	glx.Targets().TransformFeedback().BindBase(0, tfBuffer)
	feedback.Begin(gl.Points)
	glx.BindVertexArray(vao)
	glx.DrawArrays(gl.Points, 0, 3)
	glx.UnbindVertexArray()
	feedback.End()
	glx.Targets().TransformFeedback().UnbindBase(0)
	glx.UnuseProgram()

	glx.Targets().TransformFeedback().Bind(tfBuffer)
	tfFloats := make([]float32, 3)
	tfBytes := glunsafe.Map(tfFloats)
	n := glx.Targets().TransformFeedback().Contents(tfBytes)
	glx.Targets().TransformFeedback().Unbind()

	if n != 4*3 {
		return fmt.Errorf("transform feedback contents call returned %d bytes read, expected 12", n)
	}

	if got, want := tfFloats, []float32{10.0, 15.0, 20.0}; !reflect.DeepEqual(got, want) {
		return fmt.Errorf("invalid data read from transform feedback: %v", got)
	}

	return nil
}
