package wasmjs

import (
	"github.com/PieterD/warp/driver"
)

func Open() (factory driver.Factory) {
	factory = jsFactory{}
	return factory
}
