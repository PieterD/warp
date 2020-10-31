package gl

import (
	"fmt"

	"github.com/PieterD/warp/driver"
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
	programObject := glx.functions.CreateProgram()
	glx.functions.AttachShader(programObject, vertShaderObject)
	glx.functions.AttachShader(programObject, fragShaderObject)
	glx.functions.LinkProgram(programObject)
	linkStatus, ok := glx.functions.GetProgramParameter(programObject, glx.constants.LINK_STATUS).IsBoolean()
	if !ok {
		return nil, fmt.Errorf("LINK_STATUS program parameter did not return boolean")
	}
	if !linkStatus {
		info, ok := glx.functions.GetProgramInfoLog(programObject).IsString()
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
	shaderObject := glx.functions.CreateShader(shaderType)
	glx.functions.ShaderSource(shaderObject, glx.factory.String(code))
	glx.functions.CompileShader(shaderObject)
	csValue := glx.functions.GetShaderParameter(shaderObject, glx.constants.COMPILE_STATUS)
	compileStatus, ok := csValue.IsBoolean()
	if !ok {
		return nil, fmt.Errorf("COMPILE_STATUS shader parameteer did not return boolean: %T", csValue)
	}
	if !compileStatus {
		info, ok := glx.functions.GetShaderInfoLog(shaderObject).IsString()
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
	typ   float64
	siz   float64
}

func (p *Program) Attribute(name string) (*Attribute, error) {
	glx := p.glx
	glIndex := glx.functions.GetAttribLocation(p.glObject, glx.factory.String(name))
	fIndex, ok := glIndex.IsNumber()
	if !ok {
		return nil, fmt.Errorf("returned attribute index is somehow not a number: %T", glIndex)
	}
	if fIndex == -1 {
		return nil, fmt.Errorf("attribute does not exist")
	}
	attribInfoValue := glx.functions.GetActiveAttrib(p.glObject, glIndex)
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
	attribType, ok := attribInfo.Get("type").IsNumber()
	if !ok {
		return nil, fmt.Errorf("expected attribute type")
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

func (a *Attribute) Type() (glType, size int) {
	return int(a.typ), int(a.siz)
}

type Uniform struct {
	p        *Program
	name     string
	location driver.Value
}

func (p *Program) Uniform(name string) *Uniform {
	glx := p.glx
	location := glx.functions.GetUniformLocation(p.glObject, glx.factory.String(name))
	if location.IsNull() {
		return nil
	}
	return &Uniform{
		p:        p,
		name:     name,
		location: location,
	}
}

func (u *Uniform) SetFloat32(v float32) {
	glx := u.p.glx
	glx.functions.Uniform1f(u.location, glx.factory.Number(float64(v)))
}
