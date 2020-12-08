package gl

import "fmt"

type ParameterSet struct {
	glx *Context
}

func newParameterSet(glx *Context) *ParameterSet {
	return &ParameterSet{
		glx: glx,
	}
}

func (ps ParameterSet) MaxCombinedTextureImageUnits() int {
	glx := ps.glx
	paramValue := glx.constants.GetParameter(glx.constants.MAX_COMBINED_TEXTURE_IMAGE_UNITS)
	fMaxTextureUnits, ok := paramValue.ToFloat64()
	if !ok {
		panic(fmt.Errorf("parameter MAX_COMBINED_TEXTURE_IMAGE_UNITS should return number: %T", paramValue))
	}
	return int(fMaxTextureUnits)
}
