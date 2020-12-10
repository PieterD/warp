package wasmjs

import (
	"fmt"
	"syscall/js"

	"github.com/PieterD/warp/pkg/driver"
)

func js2value(jsValue js.Value) (v driver.Value) {
	switch jsType := jsValue.Type(); jsType {
	case js.TypeUndefined:
		return jsUndefined{
			v: jsValue,
		}
	case js.TypeNull:
		return jsNull{
			v: jsValue,
		}
	case js.TypeBoolean:
		return jsBoolean{
			v: jsValue,
		}
	case js.TypeNumber:
		return jsNumber{
			v: jsValue,
		}
	case js.TypeString:
		return jsString{
			v: jsValue,
		}
	case js.TypeObject:
		return jsObject{
			v: jsValue,
		}
	case js.TypeFunction:
		return jsFunction{
			v: jsValue,
		}
	default:
		return nil
	}
}

func value2js(dv driver.Value) (jsValue js.Value) {
	if dv == nil {
		return js.Null()
	}
	vdv, ok := dv.(vValue)
	if !ok {
		panic(fmt.Errorf("value was not our type: %T", dv))
	}
	return vdv.jsValue()
}
