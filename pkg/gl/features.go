package gl

import (
	"fmt"

	"github.com/PieterD/warp/pkg/driver"
)

type Features struct {
	glx *Context
}

func (fs Features) Blend(enable bool) {
	glx := fs.glx
	if enable {
		glx.constants.Enable(glx.constants.BLEND)
	} else {
		glx.constants.Disable(glx.constants.BLEND)
	}
}

type BlendFactor int

const (
	Zero BlendFactor = iota + 1
	One
	SrcColor
	OneMinusSrcColor
	DstColor
	OneMinusDstColor
	SrcAlpha
	OneMinusSrcAlpha
	DstAlpha
	OneMinusDstAlpha
)

func (f BlendFactor) glValue(glx *Context) driver.Value {
	switch f {
	case Zero:
		return glx.constants.ZERO
	case One:
		return glx.constants.ONE
	case SrcColor:
		return glx.constants.SRC_COLOR
	case OneMinusSrcColor:
		return glx.constants.ONE_MINUS_SRC_COLOR
	case DstColor:
		return glx.constants.DST_COLOR
	case OneMinusDstColor:
		return glx.constants.ONE_MINUS_DST_COLOR
	case SrcAlpha:
		return glx.constants.SRC_ALPHA
	case OneMinusSrcAlpha:
		return glx.constants.ONE_MINUS_SRC_ALPHA
	case DstAlpha:
		return glx.constants.DST_ALPHA
	case OneMinusDstAlpha:
		return glx.constants.ONE_MINUS_DST_ALPHA
	default:
		panic(fmt.Errorf("BlendFactor %v is not valid", f))
	}
}

type BlendEquation int

const (
	Add BlendEquation = iota + 1
	Subtract
	ReverseSubtract
	Min
	Max
)

func (f BlendEquation) glValue(glx *Context) driver.Value {
	switch f {
	case Add:
		return glx.constants.FUNC_ADD
	case Subtract:
		return glx.constants.FUNC_SUBTRACT
	case ReverseSubtract:
		return glx.constants.FUNC_REVERSE_SUBTRACT
	case Min:
		return glx.constants.MIN
	case Max:
		return glx.constants.MAX
	default:
		panic(fmt.Errorf("BlendFactor %v is not valid", f))
	}
}

type BlendFuncConfig struct {
	Source      BlendFactor
	Destination BlendFactor
	Equation    BlendEquation
}

func (fs Features) BlendFunc(cfg BlendFuncConfig) {
	glx := fs.glx
	if cfg.Source > 0 && cfg.Destination > 0 {
		glx.constants.BlendFunc(cfg.Source.glValue(glx), cfg.Destination.glValue(glx))
	}
	if cfg.Equation > 0 {
		glx.constants.BlendEquation(cfg.Equation.glValue(glx))
	}
}
