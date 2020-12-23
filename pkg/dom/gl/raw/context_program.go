package raw

import (
	"fmt"

	"github.com/PieterD/warp/pkg/driver"
	"github.com/go-gl/mathgl/mgl32"
)

type ProgramObject struct {
	glx   *Context
	value driver.Value
}

func (program ProgramObject) Destroy() {
	glx := program.glx
	glx.constants.DeleteProgram(program.value)
}

func (program ProgramObject) Attach(shader ShaderObject) {
	glx := program.glx
	glx.constants.AttachShader(program.value, shader.value)
}

func (program ProgramObject) TransformFeedbackVaryings(interleaved bool, inputNames []string) {
	glx := program.glx
	glBufferMode := glx.constants.SEPARATE_ATTRIBS
	if interleaved {
		glBufferMode = glx.constants.INTERLEAVED_ATTRIBS
	}
	var arrayValues []driver.Value
	for _, name := range inputNames {
		arrayValues = append(arrayValues, glx.factory.String(name))
	}
	feedbackNames := glx.factory.Array(arrayValues...)
	glx.constants.TransformFeedbackVaryings(program.value, feedbackNames, glBufferMode)
}

func (program ProgramObject) Link() error {
	glx := program.glx
	glx.constants.LinkProgram(program.value)
	linkStatus, ok := glx.constants.GetProgramParameter(program.value, glx.constants.LINK_STATUS).ToBoolean()
	if !ok {
		return fmt.Errorf("LINK_STATUS program parameter did not return boolean")
	}
	if !linkStatus {
		info, ok := glx.constants.GetProgramInfoLog(program.value).ToString()
		if !ok {
			return fmt.Errorf("program linking error, but programInfoLog did not return string")
		}
		return fmt.Errorf("program linking error: %s", info)
	}
	return nil
}

func (program ProgramObject) GetUniformBlockIndex(blockName string) uint {
	glx := program.glx
	rv := glx.constants.GetUniformBlockIndex(program.value, glx.factory.String(blockName))
	f, ok := rv.ToFloat64()
	if !ok {
		panic(fmt.Errorf("unknown return type from GetUniformBlockIndex: %T %v", rv, rv))
	}
	return uint(f)
}

type Uniform struct {
	glx   *Context
	value driver.Value
}

func (program ProgramObject) Uniform(name string) (Uniform, error) {
	glx := program.glx
	value := glx.constants.GetUniformLocation(program.value, glx.factory.String(name))
	if value.IsNull() {
		return Uniform{}, fmt.Errorf("no uniform by that name")
	}
	return Uniform{
		glx:   glx,
		value: value,
	}, nil
}

func (u Uniform) Float(v float32) {
	glx := u.glx
	glx.constants.Uniform1f(
		u.value,
		glx.factory.Number(float64(v)),
	)
}

func (u Uniform) Vec2(v mgl32.Vec2) {
	glx := u.glx
	glx.constants.Uniform2f(
		u.value,
		glx.factory.Number(float64(v[0])),
		glx.factory.Number(float64(v[1])),
	)
}

func (u Uniform) Vec3(v mgl32.Vec3) {
	glx := u.glx
	glx.constants.Uniform3f(
		u.value,
		glx.factory.Number(float64(v[0])),
		glx.factory.Number(float64(v[1])),
		glx.factory.Number(float64(v[2])),
	)
}

func (u Uniform) Vec4(v mgl32.Vec4) {
	glx := u.glx
	glx.constants.Uniform4f(
		u.value,
		glx.factory.Number(float64(v[0])),
		glx.factory.Number(float64(v[1])),
		glx.factory.Number(float64(v[2])),
		glx.factory.Number(float64(v[3])),
	)
}

func (u Uniform) Int(v int) {
	glx := u.glx
	glx.constants.Uniform1i(u.value, glx.factory.Number(float64(v)))
}

func (u Uniform) Sampler(textureIndex int) {
	u.Int(textureIndex)
}

type ShaderObject struct {
	glx   *Context
	value driver.Value
}

//go:generate stringer -type=ShaderType
type ShaderType int

const (
	VertexShader ShaderType = iota + 1
	FragmentShader
)

func (shader ShaderObject) Destroy() {
	glx := shader.glx
	glx.constants.DeleteShader(shader.value)
}

func (shader ShaderObject) Source(source string) {
	glx := shader.glx
	glx.constants.ShaderSource(shader.value, glx.factory.String(source))
}

func (shader ShaderObject) Compile() error {
	glx := shader.glx
	glx.constants.CompileShader(shader.value)
	csValue := glx.constants.GetShaderParameter(shader.value, glx.constants.COMPILE_STATUS)
	compileStatus, ok := csValue.ToBoolean()
	if !ok {
		return fmt.Errorf("COMPILE_STATUS shader parameteer did not return boolean: %T", csValue)
	}
	if !compileStatus {
		info, ok := glx.constants.GetShaderInfoLog(shader.value).ToString()
		if !ok {
			return fmt.Errorf("shaderInfoLog did not return string")
		}
		return fmt.Errorf("shader compile error: %s", info)
	}
	return nil
}
