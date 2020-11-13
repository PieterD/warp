package gl

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/PieterD/warp/pkg/driver"
)

type ProgramConfig struct {
	VertexCode   string
	FragmentCode string
}

type Program struct {
	glx      *Context
	glObject driver.Value
}

func newProgram(glx *Context, cfg ProgramConfig) (*Program, error) {
	vertShaderObject, err := compileShader(glx, glx.constants.VERTEX_SHADER, cfg.VertexCode)
	if err != nil {
		return nil, fmt.Errorf("compiling vertex shader: %w", err)
	}
	fragShaderObject, err := compileShader(glx, glx.constants.FRAGMENT_SHADER, cfg.FragmentCode)
	if err != nil {
		return nil, fmt.Errorf("compiling fragment shader: %w", err)
	}
	programObject := glx.constants.CreateProgram()
	glx.constants.AttachShader(programObject, vertShaderObject)
	glx.constants.AttachShader(programObject, fragShaderObject)
	glx.constants.LinkProgram(programObject)
	linkStatus, ok := glx.constants.GetProgramParameter(programObject, glx.constants.LINK_STATUS).IsBoolean()
	if !ok {
		return nil, fmt.Errorf("LINK_STATUS program parameter did not return boolean")
	}
	if !linkStatus {
		info, ok := glx.constants.GetProgramInfoLog(programObject).IsString()
		if !ok {
			return nil, fmt.Errorf("programInfoLog did not return string")
		}
		return nil, fmt.Errorf("program linking error: %s", info)
	}
	return &Program{
		glx:      glx,
		glObject: programObject,
	}, nil
}

func compileShader(glx *Context, shaderType driver.Value, code string) (driver.Value, error) {
	shaderObject := glx.constants.CreateShader(shaderType)
	glx.constants.ShaderSource(shaderObject, glx.factory.String(code))
	glx.constants.CompileShader(shaderObject)
	csValue := glx.constants.GetShaderParameter(shaderObject, glx.constants.COMPILE_STATUS)
	compileStatus, ok := csValue.IsBoolean()
	if !ok {
		return nil, fmt.Errorf("COMPILE_STATUS shader parameteer did not return boolean: %T", csValue)
	}
	if !compileStatus {
		info, ok := glx.constants.GetShaderInfoLog(shaderObject).IsString()
		if !ok {
			return nil, fmt.Errorf("shaderInfoLog did not return string")
		}
		return nil, fmt.Errorf("shader compile error: %s", info)
	}
	return shaderObject, nil
}

type Attribute struct {
	p     *Program
	name  string
	index int
	typ   Type
	siz   float64
}

func (p *Program) Attribute(name string) (*Attribute, error) {
	glx := p.glx
	glIndex := glx.constants.GetAttribLocation(p.glObject, glx.factory.String(name))
	fIndex, ok := glIndex.IsNumber()
	if !ok {
		return nil, fmt.Errorf("returned attribute index is somehow not a number: %T", glIndex)
	}
	if fIndex == -1 {
		return nil, fmt.Errorf("attribute does not exist")
	}
	attribInfoValue := glx.constants.GetActiveAttrib(p.glObject, glIndex)
	if attribInfoValue.IsNull() {
		return nil, fmt.Errorf("no attribute info found")
	}
	attribInfo := attribInfoValue.IsObject()
	if attribInfo.IsNull() {
		return nil, fmt.Errorf("attribute info is not an object: %T", attribInfo)
	}
	attribName, ok := attribInfo.Get("name").IsString()
	if !ok {
		return nil, fmt.Errorf("expected attribute name")
	}
	if attribName != name {
		return nil, fmt.Errorf("attribute info name does not match request name: %s", attribName)
	}
	attribType, err := glx.typeConverter.FromJs(attribInfo.Get("type"))
	if err != nil {
		return nil, fmt.Errorf("fetching attr info type: %w", err)
	}
	attribSize, ok := attribInfo.Get("size").IsNumber()
	if !ok {
		return nil, fmt.Errorf("expected attribute size")
	}

	return &Attribute{
		p:     p,
		name:  name,
		index: int(fIndex),
		typ:   attribType,
		siz:   attribSize,
	}, nil
}

func (a *Attribute) Type() (typ Type, size int) {
	return a.typ, int(a.siz)
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
	buf.Put(fastFloat32ToByte(m[:]))
	glx.constants.UniformMatrix4fv(u.location, glx.factory.Boolean(false), buf.AsFloat32Array())
}
