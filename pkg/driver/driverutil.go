package driver

func Bind(o Object, methodName string) func(args ...Value) Value {
	got := o.Get(methodName)
	function, ok := got.ToFunction()
	if !ok {
		return nil
	}
	return func(args ...Value) Value {
		return function.Call(o, args...)
	}
}

func IndexableToSlice(factory Factory, o Object) []Value {
	numChildren, ok := o.Get("length").ToFloat64()
	if !ok {
		return nil
	}
	var values []Value
	fIndex := Bind(o, "item")
	for i := 0; i < int(numChildren); i++ {
		dValue := fIndex(factory.Number(float64(i)))
		values = append(values, dValue)
	}
	return values
}
