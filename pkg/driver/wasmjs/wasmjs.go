package wasmjs

import (
	"github.com/PieterD/warp/pkg/driver"
)

func Open() (factory driver.Factory) {
	factory = jsFactory{}
	return factory
}
