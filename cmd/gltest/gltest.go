package main

import (
	"context"
	"fmt"
	"image"
	"os"
	"time"

	"github.com/PieterD/warp/pkg/dom"
	"github.com/PieterD/warp/pkg/dom/gl/raw"
	"github.com/PieterD/warp/pkg/driver/wasmjs"
)

const (
	fbWidth  = 128
	fbHeight = 128
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := run(ctx,
		&Test{
			Description: "Render a triangle",
			TF:          gltTriangle,
		},
		&Test{
			Description: "Use a uniform block",
			TF:          gltUniformBlock,
		},
		&Test{
			Description: "Render some points with different sizes and colors",
			TF:          gltPoint,
		},
		&Test{
			Description: "Render a texture",
			TF:          gltTexture,
		},
	)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "running warplay: %v", err)
	}
	<-make(chan struct{})
}

func run(ctx context.Context, tests ...*Test) error {
	factory := wasmjs.Open()
	global := dom.Open(factory)
	doc := global.Window().Document()
	canvasElem := doc.CreateElem("canvas", func(canvasElem *dom.Elem) {})
	doc.Body().AppendChildren()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	glx := raw.NewContext(canvasElem)
	defer glx.Destroy()

	glx.Viewport(0, 0, fbWidth, fbHeight)
	fbo, cleanup := newFramebuffer(glx, fbWidth, fbHeight)
	defer cleanup()

	// Run tests
	for _, test := range tests {
		img, err := func() (img image.Image, err error) {
			defer func() {
				p := recover()
				if p == nil {
					return
				}
				if pErr, ok := p.(error); ok {
					img = nil
					err = fmt.Errorf("recovered error from panic: %w", pErr)
					return
				}
				panic(p)
			}()
			glx.Targets().Framebuffer().Bind(fbo)
			glx.ClearColor(0, 0, 0, 1)
			glx.Clear()

			if err := test.TF(glx, fbo); err != nil {
				return nil, fmt.Errorf("running test function: %w", err)
			}

			glx.Targets().Framebuffer().Bind(fbo)
			pixels := glx.Targets().Framebuffer().ReadPixels(0, 0, fbWidth, fbHeight)
			glx.Targets().Framebuffer().Unbind()

			img = pixelsToImage(fbWidth, fbHeight, pixels)
			return img, nil
		}()
		//TODO: convert rendered context to image
		test.Image = img
		test.Error = err
	}

	// Render test results
	for _, test := range tests {
		text := fmt.Sprintf("%s", test.Description)
		if test.Error != nil {
			text = fmt.Sprintf("%s (error: %v)", test.Description, test.Error)
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
						if test.Image == nil {
							return
						}
						divElem.AppendChildren(
							doc.CreateElem("img", func(imgElem *dom.Elem) {
								img := dom.AsImage(imgElem)
								img.SetSrc(toDataURI(test.Image))
							}),
						)
					}),
				)
			}),
		)
	}

	return nil
}

type TestableFunc func(glx *raw.Context, fbo raw.FramebufferObject) error

type Test struct {
	Description string
	TF          TestableFunc
	Image       image.Image
	Error       error
}

func newFramebuffer(glx *raw.Context, width, height int) (raw.FramebufferObject, func()) {
	rboColor := glx.CreateRenderbuffer()
	//defer rboColor.Destroy()
	glx.Targets().RenderBuffer().Bind(rboColor)
	glx.Targets().RenderBuffer().Storage(raw.RenderbufferConfig{
		Type:   raw.ColorBuffer,
		Width:  width,
		Height: height,
	})
	rboDepthStencil := glx.CreateRenderbuffer()
	//defer rboDepthStencil.Destroy()
	glx.Targets().RenderBuffer().Bind(rboDepthStencil)
	glx.Targets().RenderBuffer().Storage(raw.RenderbufferConfig{
		Type:   raw.DepthStencilBuffer,
		Width:  width,
		Height: height,
	})
	glx.Targets().RenderBuffer().Unbind()
	fbo := glx.CreateFramebuffer()
	//defer fbo.Destroy()
	glx.Targets().Framebuffer().Bind(fbo)
	glx.Targets().RenderBuffer().Bind(rboColor)
	glx.Targets().Framebuffer().AttachRenderbuffer(raw.ColorBuffer, rboColor)
	glx.Targets().RenderBuffer().Bind(rboDepthStencil)
	glx.Targets().Framebuffer().AttachRenderbuffer(raw.DepthStencilBuffer, rboDepthStencil)
	glx.Targets().RenderBuffer().Unbind()
	glx.Targets().Framebuffer().Unbind()
	return fbo, func() {
		fbo.Destroy()
		rboDepthStencil.Destroy()
		rboColor.Destroy()
	}
}
