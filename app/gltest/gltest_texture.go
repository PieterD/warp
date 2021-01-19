package main

import (
	"context"
	"fmt"

	"github.com/PieterD/warp/pkg/gl"
	"github.com/PieterD/warp/pkg/gl/glunsafe"
)

func gltTexture(_ context.Context, glx *gl.Context, _ gl.FramebufferObject) error {
	var (
		vSource = `#version 300 es
precision mediump float;

layout (location = 0) in vec3 Position;
layout (location = 1) in vec2 TexCoord;
out vec2 texCoord;

void main(void) {
	gl_Position = vec4(Position, 1.0);
	texCoord = TexCoord;
}`
		fSource = `#version 300 es
precision mediump float;
in vec2 texCoord;
uniform sampler2D Texture;
out vec4 FragColor;

void main(void) {
	FragColor = texture(Texture, texCoord);
}`
		vertices = []float32{
			-0.5, -0.5, 0.0, 0.0, 0.0,
			0.5, -0.5, 0.0, 1.0, 0.0,
			0.5, 0.5, 0.0, 1.0, 1.0,
			-0.5, 0.5, 0.0, 0.0, 1.0,
		}
		indices = []uint16{
			0, 1, 2,
			2, 3, 0,
		}
	)

	textureImage, err := loadTexture("texture.png")
	if err != nil {
		return fmt.Errorf("loading texture image: %w", err)
	}
	textureObject := glx.CreateTexture()
	defer textureObject.Destroy()
	glx.Targets().ActiveTextureUnit(0)
	glx.Targets().Texture2D().Bind(textureObject)
	textureSize := textureImage.Bounds().Size()
	glx.Targets().Texture2D().Settings(gl.Texture2DConfig{
		Minify:  gl.Linear,
		Magnify: gl.Linear,
		WrapS:   gl.ClampToEdge,
		WrapT:   gl.ClampToEdge,
	})
	glx.Targets().Texture2D().Allocate(textureSize.X, textureSize.Y, 0)
	glx.Targets().Texture2D().SubImage(0, 0, 0, textureImage)
	glx.Targets().Texture2D().GenerateMipmap()
	glx.Targets().Texture2D().Unbind()

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

	textureUniform, err := program.Uniform("Texture")
	if err != nil {
		return fmt.Errorf("getting Texture uniform: %w", err)
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
