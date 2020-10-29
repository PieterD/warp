package main

import (
	"fmt"
	"os"

	"github.com/PieterD/warp/dom"
	"github.com/PieterD/warp/driver/wasmjs"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "running warplay: %v", err)
	}
	<-make(chan struct{})
}

func run() error {
	dGlobal, factory := wasmjs.Open()
	global := dom.Open(dGlobal, factory)
	doc := global.Window().Document()
	doc.Body().AppendChildren(
		doc.CreateElem("label", func(newElem *dom.Elem) {
			newElem.SetText("hello!")
		}),
	)
	return nil
}
