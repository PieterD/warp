package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"os"
	"time"

	"github.com/PieterD/warp/pkg/dom"
	"github.com/PieterD/warp/pkg/dom/gl/raw"
	"github.com/PieterD/warp/pkg/driver/wasmjs"
)

//Testing:
//
//Use framebuffer objects
//Generate images of set coordinates from opengl code:

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "running warplay: %v", err)
	}
	<-make(chan struct{})
}

func run(ctx context.Context) error {
	factory := wasmjs.Open()
	global := dom.Open(factory)
	doc := global.Window().Document()
	canvasElem := doc.CreateElem("canvas", func(canvasElem *dom.Elem) {})
	doc.Body().AppendChildren()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	glx := raw.NewContext(canvasElem)
	defer glx.Destroy()

	testResults := Run(ctx, glx)
	for _, testResult := range testResults {
		fmt.Printf("%s: %s (%v)\n", testResult.Name, testResult.Description, testResult.Error)
		text := fmt.Sprintf("%s: %s", testResult.Name, testResult.Description)
		if testResult.Error != nil {
			text = fmt.Sprintf("%s: %s (%v)", testResult.Name, testResult.Description, testResult.Error)
		}
		doc.Body().AppendChildren(
			doc.CreateElem("p", func(pElem *dom.Elem) {
				pElem.AppendClasses("testResult")
				pElem.AppendChildren(
					doc.CreateElem("label", func(labelElem *dom.Elem) {
						labelElem.SetText(text)
					}),
					doc.CreateElem("div", func(divElem *dom.Elem) {
						divElem.AppendClasses("clearfix")
						divElem.AppendChildren(
							doc.CreateElem("img", func(imgElem *dom.Elem) {
								img := dom.AsImage(imgElem)
								img.SetSrc(toDataURI(testResult.Image))
							}),
						)
					}),
				)
			}),
		)
	}

	return nil
}

func toDataURI(img image.Image) string {
	buf := &bytes.Buffer{}
	buf.WriteString("data:image/png;base64,")
	base64Encoder := base64.NewEncoder(base64.StdEncoding, buf)
	pngEncoder := png.Encoder{
		CompressionLevel: png.NoCompression,
	}
	if err := pngEncoder.Encode(base64Encoder, img); err != nil {
		panic(fmt.Errorf("encoding png image: %w", err))
	}
	if err := base64Encoder.Close(); err != nil {
		panic(fmt.Errorf("closing base64 encoder: %w", err))
	}
	return buf.String()
}
