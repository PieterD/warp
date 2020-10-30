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
	factory := wasmjs.Open()
	global := dom.Open(factory)
	doc := global.Window().Document()
	doc.Body().AppendChildren(
		doc.CreateElem("label", func(newElem *dom.Elem) {
			newElem.SetText("hello!")
			newElem.AppendClasses("first")
		}),
		doc.CreateElem("label", func(newElem *dom.Elem) {
			newElem.SetText("world!")
			newElem.EventHandler("click", func(this *dom.Elem, event *dom.Event) {
				fmt.Printf("click!\n")
			})
		}),
		doc.CreateElem("label", func(newElem *dom.Elem) {
			newElem.SetText("yes!")
		}),
	)
	return nil
}
