package main

import (
	"fmt"

	"github.com/PieterD/warp/pkg/gl"
	"github.com/PieterD/warp/pkg/gl/glunsafe"
)

func gltTexture(glx *gl.Context, _ gl.FramebufferObject) error {
	textureImage, err := loadTexture("texture.png")
	if err != nil {
		return fmt.Errorf("loading texture image: %w", err)
	}
	textureObject := glx.CreateTexture()
	defer textureObject.Destroy()
	glx.Targets().ActiveTextureUnit(0)
	glx.Targets().Texture2D().Bind(textureObject)
	glx.Targets().Texture2D().Image(textureImage, gl.Texture2DConfig{})
	glx.Targets().Texture2D().Unbind()

	program := glx.CreateProgram()
	defer program.Destroy()

	vShader := glx.CreateShader(gl.VertexShader)
	defer vShader.Destroy()
	vShader.Source(`#version 300 es
precision mediump float;

layout (location = 0) in vec3 Position;
layout (location = 1) in vec2 TexCoord;
out vec2 texCoord;

void main(void) {
	gl_Position = vec4(Position, 1.0);
	texCoord = TexCoord;
}`)
	if err := vShader.Compile(); err != nil {
		return fmt.Errorf("compiling vertex shader: %w", err)
	}

	fShader := glx.CreateShader(gl.FragmentShader)
	defer fShader.Destroy()
	fShader.Source(`#version 300 es
precision mediump float;
in vec2 texCoord;
uniform sampler2D Texture;
out vec4 FragColor;

void main(void) {
	FragColor = texture(Texture, texCoord);
}`)
	if err := fShader.Compile(); err != nil {
		return fmt.Errorf("compiling fragment shader: %w", err)
	}

	program.Attach(vShader)
	program.Attach(fShader)
	if err := program.Link(); err != nil {
		return fmt.Errorf("linking program: %w", err)
	}
	textureUniform, err := program.Uniform("Texture")
	if err != nil {
		return fmt.Errorf("getting Texture uniform: %w", err)
	}

	vertices := []float32{
		-0.5, -0.5, 0.0, 0.0, 0.0,
		0.5, -0.5, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.0, 1.0, 1.0,
		-0.5, 0.5, 0.0, 0.0, 1.0,
	}
	indices := []uint16{
		0, 1, 2,
		2, 3, 0,
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
	vao.VertexAttribPointer(0, gl.Vec3, false, 5*4, 0)
	vao.VertexAttribPointer(1, gl.Vec2, false, 5*4, 3*4)
	vao.EnableVertexAttribArray(0)
	vao.EnableVertexAttribArray(1)
	glx.Targets().Array().UnbindBuffer()
	glx.UnbindVertexArray()

	glx.ClearColor(0.75, 0.8, 0.85, 1.0)
	glx.Clear()
	glx.UseProgram(program)
	defer glx.UnuseProgram()
	textureUniform.Sampler(0)
	glx.Targets().ActiveTextureUnit(0)
	glx.Targets().Texture2D().Bind(textureObject)
	defer glx.Targets().Texture2D().Unbind()
	glx.BindVertexArray(vao)
	defer glx.UnbindVertexArray()
	glx.Targets().ElementArray().BindBuffer(iBuffer)
	defer glx.Targets().ElementArray().UnbindBuffer()
	glx.DrawElements(gl.Triangles, 6, gl.UnsignedShort, 0)

	return nil
}
