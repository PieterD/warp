package main

import (
	"fmt"
	"image"

	"github.com/PieterD/warp/pkg/gl/glunsafe"

	"github.com/PieterD/warp/pkg/gl"
)

func gltSpriteMap(glx *gl.Context, _ gl.FramebufferObject) error {
	sprite1Image, err := loadTexture("sprite-1.png")
	if err != nil {
		return fmt.Errorf("loading texture image: %w", err)
	}
	sprite2Image, err := loadTexture("sprite-2.png")
	if err != nil {
		return fmt.Errorf("loading texture image: %w", err)
	}
	atlas := NewSpriteMap(glx, SpriteMapConfig{
		SpriteSize: 32,
		GridSize:   2,
	})
	defer atlas.Destroy()

	atlas.Bind(0)
	atlas.Allocate()
	if err := atlas.Add(sprite1Image); err != nil {
		return fmt.Errorf("adding image 1: %w", err)
	}
	if err := atlas.Add(sprite2Image); err != nil {
		return fmt.Errorf("adding image 2: %w", err)
	}
	atlas.GenerateMipmaps()
	atlas.Unbind()

	program := glx.CreateProgram()
	defer program.Destroy()

	vShader := glx.CreateShader(gl.VertexShader)
	defer vShader.Destroy()
	vShader.Source(`#version 300 es
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
}`)

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
	instanceData := []float32{
		-0.5, -0.5, 0.0, 0.2, 1.0, 0.0,
		0.5, -0.5, 0.0, 0.3, 0.0, 0.0,
		0.5, 0.5, 0.0, 0.4, 1.0, 0.0,
		-0.5, 0.5, 0.0, 0.5, 0.0, 0.0,
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

	instanceBuffer := glx.CreateBuffer()
	defer instanceBuffer.Destroy()
	glx.Targets().Array().BindBuffer(instanceBuffer)
	glx.Targets().Array().BufferData(glunsafe.Map(instanceData), gl.Static, gl.Draw)
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
	vao.VertexAttribPointer(2, gl.Vec3, false, 6*4, 0)
	vao.VertexAttribDivisor(2, 1)
	vao.EnableVertexAttribArray(2)
	vao.VertexAttribPointer(3, gl.Float, false, 6*4, 3*4)
	vao.VertexAttribDivisor(3, 1)
	vao.EnableVertexAttribArray(3)
	vao.VertexAttribPointer(4, gl.Vec2, false, 6*4, 4*4)
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

type SpriteMapConfig struct {
	SpriteSize int // Square
	GridSize   int // Square
}

type SpriteMap struct {
	glx        *gl.Context
	spriteSize int
	gridSize   int
	used       int
	texture    gl.TextureObject
}

func NewSpriteMap(glx *gl.Context, cfg SpriteMapConfig) *SpriteMap {
	m := &SpriteMap{
		glx:        glx,
		spriteSize: cfg.SpriteSize,
		gridSize:   cfg.GridSize,
		used:       0,
		texture:    glx.CreateTexture(),
	}
	return m
}

func (m *SpriteMap) Destroy() {
	m.texture.Destroy()
}

func (m *SpriteMap) Bind(textureUnit int) {
	glx := m.glx
	glx.Targets().ActiveTextureUnit(textureUnit)
	glx.Targets().Texture2D().Bind(m.texture)
}

func (m *SpriteMap) Unbind() {
	glx := m.glx
	glx.Targets().Texture2D().Unbind()
}

func (m *SpriteMap) Allocate() {
	glx := m.glx
	glx.Targets().Texture2D().Allocate(m.gridSize*m.spriteSize, m.gridSize*m.spriteSize, 0)
	glx.Targets().Texture2D().Settings(gl.Texture2DConfig{
		Minify:  gl.Nearest,
		Magnify: gl.Nearest,
		WrapS:   gl.ClampToEdge,
		WrapT:   gl.ClampToEdge,
	})
}

func (m *SpriteMap) Add(imgs ...image.Image) error {
	glx := m.glx
	for _, img := range imgs {
		if size := img.Bounds().Size(); size.X != m.spriteSize || size.Y != m.spriteSize {
			return fmt.Errorf("image size %v does not match sprite map's sprite size %d", size, m.spriteSize)
		}
		index := m.used
		if index >= m.gridSize*m.gridSize {
			return fmt.Errorf("sprite map full: it only fits %d images", m.gridSize*m.gridSize)
		}
		col := index % m.gridSize
		row := index / m.gridSize

		glx.Targets().Texture2D().SubImage(col*m.spriteSize, row*m.spriteSize, 0, img)
		m.used++
	}
	return nil
}

func (m *SpriteMap) Scale() float32 {
	return 1.0 / float32(m.gridSize)
}

func (m *SpriteMap) GenerateMipmaps() {
	glx := m.glx
	glx.Targets().Texture2D().GenerateMipmap()
}
