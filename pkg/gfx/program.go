package gfx

import (
	"fmt"

	"github.com/PieterD/warp/pkg/gl/glutil"

	"github.com/PieterD/warp/pkg/driver"
)

type SampleMode int

const (
	Sample2D SampleMode = iota
)

type ProgramConfig struct {
	HighPrecision bool
	Uniforms      interface{}
	Attributes    ActiveCoupling
	Textures      []ProgramSamplerConfig
	Feedback      ActiveCoupling
	VertexCode    string
	FragmentCode  string
}

type ProgramSamplerConfig struct {
	Name string
	Mode SampleMode
}

type Program struct {
	glx           *Context
	glObject      driver.Value
	rawUniforms   interface{}
	uniBuffer     *Buffer
	uniBlockIndex driver.Value
	textures      []programTexture
	feedback      *feedback
}

type programTexture struct {
	name string
	mode SampleMode
}

var baseHeader = `#version 300 es
#define WARP_GL_ENABLED 1
`

var headerHighPrecision = baseHeader + `precision highp float;
`

var headerMediumPrecision = baseHeader + `precision mediump float;
`

func newProgram(glx *Context, cfg ProgramConfig) (*Program, error) {
	programObject := glx.constants.CreateProgram()

	hdr := headerMediumPrecision
	if cfg.HighPrecision {
		hdr = headerHighPrecision
	}
	rawUniform := cfg.Uniforms
	var uniBuffer *Buffer
	if rawUniform != nil {
		uniformDef, err := glutil.Std140Uniform(rawUniform)
		if err != nil {
			return nil, fmt.Errorf("creating uniform definition: %w", err)
		}
		hdr += uniformDef
		uniBuffer = glx.Buffer()
	}
	for _, textureCfg := range cfg.Textures {
		modeName := ""
		switch textureCfg.Mode {
		case Sample2D:
			modeName = "sampler2D"
		default:
			return nil, fmt.Errorf("unknown texture sampler mode %v", textureCfg.Mode)
		}
		hdr += fmt.Sprintf("uniform %s %s;\n", modeName, textureCfg.Name)
	}

	// Verify that all enabled attributes really exist.
	for attrName := range cfg.Attributes.enabled {
		if _, ok := cfg.Attributes.dc.attrByName[attrName]; !ok {
			return nil, fmt.Errorf("unknown enabled attribute %s in active coupling", attrName)
		}
	}
	vertHdr := ""
	for _, attr := range cfg.Attributes.dc.attributes {
		if _, ok := cfg.Attributes.enabled[attr.name]; !ok {
			continue
		}
		vertHdr += fmt.Sprintf("layout(location = %d) in %s %s;\n", attr.index, attr.typ.glString(), attr.name)
	}

	var feedback *feedback
	if cfg.Feedback.AttrNum() > 0 {
		feedback = newFeedback(glx, cfg.Feedback)
		interleaved := false
		switch {
		case len(cfg.Feedback.dc.attrsByBuffer) == 0:
			return nil, fmt.Errorf("feedback buffer with enabled but without attrsByBuffer")
		case len(cfg.Feedback.dc.attrsByBuffer) == 1:
			interleaved = true
		case len(cfg.Feedback.dc.attrsByBuffer) > 1:
			interleaved = false
			for bufferName, attrIndices := range cfg.Feedback.dc.attrsByBuffer {
				if len(attrIndices) != 1 {
					return nil, fmt.Errorf("multiple buffers, but buffer %s has more than one (%d) attributes", bufferName, len(attrIndices))
				}
			}
		}
		var arrayValues []driver.Value
		for _, attr := range cfg.Feedback.dc.attributes {
			if _, ok := cfg.Feedback.enabled[attr.name]; !ok {
				continue
			}
			vertHdr += fmt.Sprintf("out %s %s;\n", attr.typ.glString(), attr.name)
			arrayValues = append(arrayValues, glx.factory.String(attr.name))
		}
		feedbackNames := glx.factory.Array(arrayValues...)
		glBufferMode := glx.constants.SEPARATE_ATTRIBS
		if interleaved {
			glBufferMode = glx.constants.INTERLEAVED_ATTRIBS
		}
		glx.constants.TransformFeedbackVaryings(programObject, feedbackNames, glBufferMode)
	}

	vertShaderObject, err := compileShader(glx, glx.constants.VERTEX_SHADER, hdr+vertHdr+cfg.VertexCode)
	if err != nil {
		return nil, fmt.Errorf("compiling vertex shader: %w", err)
	}
	fragShaderObject, err := compileShader(glx, glx.constants.FRAGMENT_SHADER, hdr+cfg.FragmentCode)
	if err != nil {
		return nil, fmt.Errorf("compiling fragment shader: %w", err)
	}
	glx.constants.AttachShader(programObject, vertShaderObject)
	glx.constants.AttachShader(programObject, fragShaderObject)
	glx.constants.LinkProgram(programObject)
	linkStatus, ok := glx.constants.GetProgramParameter(programObject, glx.constants.LINK_STATUS).ToBoolean()
	if !ok {
		return nil, fmt.Errorf("LINK_STATUS program parameter did not return boolean")
	}
	if !linkStatus {
		info, ok := glx.constants.GetProgramInfoLog(programObject).ToString()
		if !ok {
			return nil, fmt.Errorf("programInfoLog did not return string")
		}
		return nil, fmt.Errorf("program linking error: %s", info)
	}
	uniBlockIndex := glx.constants.GetUniformBlockIndex(
		programObject,
		glx.factory.String("Uniform"),
	)
	p := &Program{
		glx:           glx,
		glObject:      programObject,
		rawUniforms:   rawUniform,
		uniBuffer:     uniBuffer,
		uniBlockIndex: uniBlockIndex,
		feedback:      feedback,
	}
	glx.constants.UseProgram(p.glObject)
	defer glx.constants.UseProgram(glx.factory.Null())
	if uniBuffer != nil {
		if err := p.UpdateUniforms(); err != nil {
			return nil, fmt.Errorf("updating uniforms: %w", err)
		}
		bufferBaseIndex := 0
		glx.constants.BindBufferBase(
			glx.constants.UNIFORM_BUFFER,
			glx.factory.Number(float64(bufferBaseIndex)),
			uniBuffer.glObject,
		)
		glx.constants.UniformBlockBinding(
			programObject,
			uniBlockIndex,
			glx.factory.Number(float64(bufferBaseIndex)),
		)
		glx.constants.BindBuffer(glx.constants.UNIFORM_BUFFER, glx.factory.Null())
	}
	for textureIndex, textureCfg := range cfg.Textures {
		uniformSampler := p.uniform(textureCfg.Name)
		if uniformSampler == nil {
			return nil, fmt.Errorf("sampler uniform not found")
		}
		glx.constants.Uniform1i(uniformSampler.location, glx.factory.Number(float64(textureIndex)))
		//UniformSetter{glx: glx}.Int(uniformSampler, textureIndex)
		p.textures = append(p.textures, programTexture{
			name: textureCfg.Name,
			mode: textureCfg.Mode,
		})
	}
	return p, nil
}

func (p *Program) UpdateUniforms() error {
	glx := p.glx
	data, err := glutil.Std140Data(p.rawUniforms)
	if err != nil {
		return fmt.Errorf("converting uniform to std140 data: %w", err)
	}
	jsBuffer := glx.factory.Buffer(len(data))
	jsBuffer.Put(data)
	jsBytes := jsBuffer.AsUint8Array()
	glx.constants.BindBuffer(glx.constants.UNIFORM_BUFFER, p.uniBuffer.glObject)
	glx.constants.BufferData(glx.constants.UNIFORM_BUFFER, jsBytes, glx.constants.STATIC_DRAW)
	glx.constants.BindBuffer(glx.constants.UNIFORM_BUFFER, glx.factory.Null())
	return nil
}

func compileShader(glx *Context, shaderType driver.Value, code string) (driver.Value, error) {
	fmt.Printf("compileShader: %s\n", code)
	shaderObject := glx.constants.CreateShader(shaderType)
	glx.constants.ShaderSource(shaderObject, glx.factory.String(code))
	glx.constants.CompileShader(shaderObject)
	csValue := glx.constants.GetShaderParameter(shaderObject, glx.constants.COMPILE_STATUS)
	compileStatus, ok := csValue.ToBoolean()
	if !ok {
		return nil, fmt.Errorf("COMPILE_STATUS shader parameteer did not return boolean: %T", csValue)
	}
	if !compileStatus {
		info, ok := glx.constants.GetShaderInfoLog(shaderObject).ToString()
		if !ok {
			return nil, fmt.Errorf("shaderInfoLog did not return string")
		}
		return nil, fmt.Errorf("shader compile error: %s", info)
	}
	return shaderObject, nil
}

type uniform struct {
	p        *Program
	name     string
	location driver.Value
}

func (p *Program) uniform(name string) *uniform {
	glx := p.glx
	location := glx.constants.GetUniformLocation(p.glObject, glx.factory.String(name))
	if location.IsNull() {
		return nil
	}
	return &uniform{
		p:        p,
		name:     name,
		location: location,
	}
}

//func (us UniformSetter) Int(u *uniform, v int) {
//	glx := us.glx
//	glx.constants.Uniform1i(u.location, glx.factory.Number(float64(v)))
//}
//
//func (us UniformSetter) Float32(u *uniform, v float32) {
//	glx := us.glx
//	glx.constants.Uniform1f(u.location, glx.factory.Number(float64(v)))
//}
//
//func (us UniformSetter) Vec3(u *uniform, v mgl32.Vec3) {
//	glx := us.glx
//	glx.constants.Uniform3f(u.location,
//		glx.factory.Number(float64(v[0])),
//		glx.factory.Number(float64(v[1])),
//		glx.factory.Number(float64(v[2])),
//	)
//}
//
//func (us UniformSetter) Mat4(u *uniform, m mgl32.Mat4) {
//	glx := us.glx
//	buf := glx.factory.Buffer(4 * 4 * 4)
//	buf.Put(glunsafe.FastFloat32ToByte(m[:]))
//	glx.constants.UniformMatrix4fv(u.location, glx.factory.Boolean(false), buf.AsFloat32Array())
//}
