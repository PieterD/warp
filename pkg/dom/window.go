package dom

import (
	"context"
	"fmt"

	"github.com/PieterD/warp/pkg/driver"
)

type Window struct {
	factory driver.Factory
	obj     driver.Object
}

func (w *Window) Document() *Document {
	dValue := w.obj.Get("document")
	dObj, ok := dValue.ToObject()
	if !ok {
		return nil
	}
	return &Document{
		factory: w.factory,
		obj:     dObj,
	}
}

func (w *Window) Animate(ctx context.Context, f func(ctx context.Context, millis float64) error) {
	fRequestAnimationFrame := driver.Bind(w.obj, "requestAnimationFrame")
	var cb driver.Function
	cb = w.factory.Function(func(this driver.Object, args ...driver.Value) driver.Value {
		if len(args) != 1 {
			panic(fmt.Errorf("expecteed 1 argument, got: %d", len(args)))
		}
		millis, ok := args[0].ToFloat64()
		if !ok {
			panic(fmt.Errorf("expected first argument to be a number: %T", args[0]))
		}
		if err := f(ctx, millis); err != nil {
			fmt.Printf("[ERROR] animation callback: %v\n", err)
			return nil
		}
		fRequestAnimationFrame(cb)
		return nil
	})
	fRequestAnimationFrame(cb)
}
