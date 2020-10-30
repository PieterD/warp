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

func IndexableToSlice(factory driver.Factory, o driver.Object) []driver.Value {
	numChildren, ok := o.Get("length").IsNumber()
	if !ok {
		return nil
	}
	var values []driver.Value
	fIndex := Bind(o, "item")
	for i := 0; i < int(numChildren); i++ {
		dValue := fIndex(factory.Number(float64(i)))
		values = append(values, dValue)
	}
	return values
}
