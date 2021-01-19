package main

import (
	"context"
	"fmt"

	"github.com/PieterD/warp/pkg/gl"
	"github.com/PieterD/warp/pkg/gl/glunsafe"
	"github.com/PieterD/warp/pkg/gl/glutil"
)

func gltSpriteAtlas(_ context.Context, glx *gl.Context, _ gl.FramebufferObject) error {
	var (
		vSource = `#version 300 es
precision mediump float;

layout (location = 0) in vec3 Position;
layout (location = 1) in vec2 TexCoord;

layout (location = 2) in vec3 Translation;
layout (location = 3) in float Scale;
layout (location = 4) in vec2 GridCoord;
uniform float SpriteMapScale;
out vec2 texCoord;

void main(void) {
	gl_Position = vec4(Position*Scale + Translation, 1.0);
	texCoord = SpriteMapScale*(TexCoord + GridCoord);
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
		instanceData = []float32{
			-0.5, -0.5, 0.0, 0.2,
			0.5, -0.5, 0.0, 0.3,
			0.5, 0.5, 0.0, 0.4,
			-0.5, 0.5, 0.0, 0.5,
		}
	)

	sprite1Image, err := loadTexture("sprite_0.png")
	if err != nil {
		return fmt.Errorf("loading texture image 0: %w", err)
	}
	sprite2Image, err := loadTexture("sprite_1.png")
	if err != nil {
		return fmt.Errorf("loading texture image 1: %w", err)
	}
	atlas := glutil.NewSpriteAtlas(glx, glutil.SpriteMapConfig{
		SpriteSize: 32,
		GridSize:   2,
	})
	defer atlas.Destroy()

	atlas.Bind(0)
	atlas.Allocate()
	sprite1Coord, err := atlas.Add(sprite1Image)
	if err != nil {
		return fmt.Errorf("adding image 1: %w", err)
	}
	sprite2Coord, err := atlas.Add(sprite2Image)
	if err != nil {
		return fmt.Errorf("adding image 2: %w", err)
	}
	atlas.GenerateMipmaps()
	atlas.Unbind()

	texIndexData := [][2]float32{
		sprite1Coord,
		sprite2Coord,
		sprite1Coord,
		sprite2Coord,
	}

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
	scaleUniform, err := program.Uniform("SpriteMapScale")
	if err != nil {
		return fmt.Errorf("getting SpriteMapScale uniform: %w", err)
	}

	vBuffer := glx.CreateBuffer()
	defer vBuffer.Destroy()
	glx.Targets().Array().BindBuffer(vBuffer)
	glx.Targets().Array().BufferData(glunsafe.Map(vertices), gl.Static, gl.Draw)

	iBuffer := glx.CreateBuffer()
	defer iBuffer.Destroy()
	glx.Targets().ElementArray().BindBuffer(iBuffer)
	glx.Targets().ElementArray().BufferData(glunsafe.Map(indices), gl.Static, gl.Draw)

	instanceBuffer := glx.CreateBuffer()
	defer instanceBuffer.Destroy()
	glx.Targets().Array().BindBuffer(instanceBuffer)
	glx.Targets().Array().BufferData(glunsafe.Map(instanceData), gl.Static, gl.Draw)

	texIndexBuffer := glx.CreateBuffer()
	defer texIndexBuffer.Destroy()
	glx.Targets().Array().BindBuffer(texIndexBuffer)
	glx.Targets().Array().BufferData(glunsafe.Map(texIndexData), gl.Static, gl.Draw)

	glx.Targets().Array().UnbindBuffer()

	vao := glx.CreateVertexArray()
	defer vao.Destroy()
	glx.BindVertexArray(vao)
	glx.Targets().Array().BindBuffer(vBuffer)
	vao.VertexAttribPointer(0, gl.Vec3, false, 5*4, 0)
	vao.EnableVertexAttribArray(0)
	vao.VertexAttribPointer(1, gl.Vec2, false, 5*4, 3*4)
	vao.EnableVertexAttribArray(1)
	glx.Targets().Array().BindBuffer(instanceBuffer)
	vao.VertexAttribPointer(2, gl.Vec3, false, 4*4, 0)
	vao.VertexAttribDivisor(2, 1)
	vao.EnableVertexAttribArray(2)
	vao.VertexAttribPointer(3, gl.Float, false, 4*4, 3*4)
	vao.VertexAttribDivisor(3, 1)
	vao.EnableVertexAttribArray(3)
	glx.Targets().Array().BindBuffer(texIndexBuffer)
	vao.VertexAttribPointer(4, gl.Vec2, false, 2*4, 0)
	vao.VertexAttribDivisor(4, 1)
	vao.EnableVertexAttribArray(4)
	glx.Targets().Array().UnbindBuffer()
	glx.UnbindVertexArray()

	glx.Features().Blend(true)
	glx.Features().BlendFunc(gl.BlendFuncConfig{
		Source:      gl.SrcAlpha,
		Destination: gl.OneMinusSrcAlpha,
	})
	glx.ClearColor(0.75, 0.8, 0.85, 1.0)
	glx.Clear()
	glx.UseProgram(program)
	defer glx.UnuseProgram()
	textureUniform.Sampler(0)
	scaleUniform.Float(atlas.Scale())
	atlas.Bind(0)
	defer atlas.Unbind()
	glx.BindVertexArray(vao)
	defer glx.UnbindVertexArray()
	glx.Targets().ElementArray().BindBuffer(iBuffer)
	defer glx.Targets().ElementArray().UnbindBuffer()
	glx.DrawElementsInstanced(gl.Triangles, 6, gl.UnsignedShort, 0, 4)

	return nil
}
