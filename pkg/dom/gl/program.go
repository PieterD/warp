package gl

import (
	"fmt"

	"github.com/PieterD/warp/pkg/dom/glunsafe"

	"github.com/PieterD/warp/pkg/dom/glutil"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/PieterD/warp/pkg/driver"
)

type ProgramConfig struct {
	HighPrecision bool
	Uniforms      interface{}
	Attributes    ActiveCoupling
	//Feedback      ActiveCoupling
	VertexCode   string
	FragmentCode string
}

type Program struct {
	glx           *Context
	glObject      driver.Value
	rawUniforms   interface{}
	uniBuffer     *Buffer
	uniBlockIndex driver.Value
}

var baseHeader = `#version 300 es
#define WARP_GL_ENABLED 1
`

var headerHighPrecision = baseHeader + `precision highp float;
`

var headerMediumPrecision = baseHeader + `precision mediump float;
`

func newProgram(glx *Context, cfg ProgramConfig) (*Program, error) {
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

	//for _, attr := range cfg.Feedback {
	//	vertHdr += fmt.Sprintf("out %s %s;\n", attr.Type.glString(), attr.Name)
	//}

	vertShaderObject, err := compileShader(glx, glx.constants.VERTEX_SHADER, hdr+vertHdr+cfg.VertexCode)
	if err != nil {
		return nil, fmt.Errorf("compiling vertex shader: %w", err)
	}
	fragShaderObject, err := compileShader(glx, glx.constants.FRAGMENT_SHADER, hdr+cfg.FragmentCode)
	if err != nil {
		return nil, fmt.Errorf("compiling fragment shader: %w", err)
	}
	programObject := glx.constants.CreateProgram()
	glx.constants.AttachShader(programObject, vertShaderObject)
	glx.constants.AttachShader(programObject, fragShaderObject)
	//if len(cfg.Feedback) > 0 {
	//	var arrayValues []driver.Value
	//	for _, attr := range cfg.Feedback {
	//		arrayValues = append(arrayValues, glx.factory.String(attr.Name))
	//	}
	//	feedbackNames := glx.factory.Array(arrayValues...)
	//	glx.constants.TransformFeedbackVaryings(programObject, feedbackNames, glx.constants.SEPARATE_ATTRIBS)
	//}
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
	}
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

type Uniform struct {
	p        *Program
	name     string
	location driver.Value
}

func (p *Program) Uniform(name string) *Uniform {
	glx := p.glx
	location := glx.constants.GetUniformLocation(p.glObject, glx.factory.String(name))
	if location.IsNull() {
		return nil
	}
	return &Uniform{
		p:        p,
		name:     name,
		location: location,
	}
}

func (p *Program) Update(f func(us *UniformSetter)) {
	glx := p.glx
	glx.constants.UseProgram(p.glObject)
	defer glx.constants.UseProgram(glx.factory.Null())
	f(&UniformSetter{glx: p.glx})
}

type UniformSetter struct {
	glx *Context
}

func (us *UniformSetter) Int(u *Uniform, v int) {
	glx := us.glx
	glx.constants.Uniform1i(u.location, glx.factory.Number(float64(v)))
}

func (us *UniformSetter) Float32(u *Uniform, v float32) {
	glx := us.glx
	glx.constants.Uniform1f(u.location, glx.factory.Number(float64(v)))
}

func (us *UniformSetter) Vec3(u *Uniform, v mgl32.Vec3) {
	glx := us.glx
	glx.constants.Uniform3f(u.location,
		glx.factory.Number(float64(v[0])),
		glx.factory.Number(float64(v[1])),
		glx.factory.Number(float64(v[2])),
	)
}

func (us *UniformSetter) Mat4(u *Uniform, m mgl32.Mat4) {
	glx := us.glx
	buf := glx.factory.Buffer(4 * 4 * 4)
	buf.Put(glunsafe.FastFloat32ToByte(m[:]))
	glx.constants.UniformMatrix4fv(u.location, glx.factory.Boolean(false), buf.AsFloat32Array())
}
