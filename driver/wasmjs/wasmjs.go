package wasmjs

import (
	"syscall/js"

	"github.com/PieterD/warp/driver"
)

func Open() (global driver.Object, factory driver.Factory) {
	global = jsObject{
		v: js.Global(),
	}
	factory = jsFactory{}
	return global, factory
}
