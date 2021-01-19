package gl

import (
	"github.com/PieterD/warp/pkg/driver"
)

type QueryTarget struct {
	glx    *Context
	target driver.Value
}

func (target QueryTarget) Begin(qo QueryObject) {
	glx := target.glx
	glx.constants.BeginQuery(target.target, qo.value)
}

func (target QueryTarget) End() {
	glx := target.glx
	glx.constants.EndQuery(target.target)
}
