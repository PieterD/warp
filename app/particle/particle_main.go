package main

import (
	"context"
	"fmt"
	"os"

	"github.com/PieterD/warp/pkg/dom"
	"github.com/PieterD/warp/pkg/driver/wasmjs"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := run(ctx)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "running particle application: %v", err)
	}
	<-make(chan struct{})
}

func run(ctx context.Context) error {
	factory := wasmjs.Open()
	global := dom.Open(factory)
	win := global.Window()
	doc := win.Document()
	body := doc.Body()
	canvasElem := doc.CreateElem("canvas", func(canvasElem *dom.Elem) {
		canvasElem.SetPropString("width", "500")
		canvasElem.SetPropString("height", "500")
	})
	body.AppendChildren(canvasElem)
	canvas := dom.AsCanvas(canvasElem)
	glx := canvas.GetContextWebgl()
	defer glx.Destroy()
	glx.Viewport(0, 0, 500, 500)

	win.Animate(ctx, func(ctx context.Context, millis float64) error {
		panic("not implemented")
	})

	return nil
}
