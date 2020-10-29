package driverutil

import "github.com/PieterD/warp/driver"

func Bind(o driver.Object, methodName string) func(args ...driver.Value) driver.Value {
	got := o.Get(methodName)
	function := got.IsFunction()
	if function == nil {
		return nil
	}
	return func(args ...driver.Value) driver.Value {
		return function.Call(o, args...)
	}
}
