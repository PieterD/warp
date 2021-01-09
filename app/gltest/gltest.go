package main

import (
	"context"
	"fmt"
	"image"
	"os"
	"time"

	"github.com/PieterD/warp/pkg/dom"
	"github.com/PieterD/warp/pkg/driver/wasmjs"
	"github.com/PieterD/warp/pkg/gl"
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
			Description: "Sprite atlas, alpha blending, instanced rendering",
			URL:         `https://github.com/PieterD/warp/blob/master/app/gltest/gltest_spriteatlas.go`,
			TF:          gltSpriteAtlas,
		},
		&Test{
			Description: "Instanced rendering",
			URL:         `https://github.com/PieterD/warp/blob/master/app/gltest/gltest_instanced.go`,
			TF:          gltInstancedQuads,
		},
		&Test{
			Description: "Feedback transform",
			URL:         `https://github.com/PieterD/warp/blob/master/app/gltest/gltest_feedback.go`,
			TF:          gltFeedback,
		},
		&Test{
			Description: "Texture",
			URL:         `https://github.com/PieterD/warp/blob/master/app/gltest/gltest_texture.go`,
			TF:          gltTexture,
		},
		&Test{
			Description: "Uniform block",
			URL:         `https://github.com/PieterD/warp/blob/master/app/gltest/gltest_uniform.go`,
			TF:          gltUniformBlock,
		},
		&Test{
			Description: "Points with different sizes and colors",
			URL:         `https://github.com/PieterD/warp/blob/master/app/gltest/gltest_point.go`,
			TF:          gltPoint,
		},
		&Test{
			Description: "Simple triangle",
			URL:         `https://github.com/PieterD/warp/blob/master/app/gltest/gltest_triangle.go`,
			TF:          gltTriangle,
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
	glx := gl.NewContext(canvasElem)
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
					doc.CreateElem("a", func(anchorElem *dom.Elem) {
						anchorElem.SetPropString("href", test.URL)
						anchorElem.SetText("[github]")
					}),
					doc.CreateElem("label", func(labelElem *dom.Elem) {
						labelElem.SetText(" ")
					}),
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

type TestableFunc func(glx *gl.Context, fbo gl.FramebufferObject) error

type Test struct {
	Description string
	URL         string
	TF          TestableFunc
	Image       image.Image
	Error       error
}

func newFramebuffer(glx *gl.Context, width, height int) (gl.FramebufferObject, func()) {
	rboColor := glx.CreateRenderbuffer()
	//defer rboColor.Destroy()
	glx.Targets().RenderBuffer().Bind(rboColor)
	glx.Targets().RenderBuffer().Storage(gl.RenderbufferConfig{
		Type:   gl.ColorBuffer,
		Width:  width,
		Height: height,
	})
	rboDepthStencil := glx.CreateRenderbuffer()
	//defer rboDepthStencil.Destroy()
	glx.Targets().RenderBuffer().Bind(rboDepthStencil)
	glx.Targets().RenderBuffer().Storage(gl.RenderbufferConfig{
		Type:   gl.DepthStencilBuffer,
		Width:  width,
		Height: height,
	})
	glx.Targets().RenderBuffer().Unbind()
	fbo := glx.CreateFramebuffer()
	//defer fbo.Destroy()
	glx.Targets().Framebuffer().Bind(fbo)
	glx.Targets().RenderBuffer().Bind(rboColor)
	glx.Targets().Framebuffer().AttachRenderbuffer(gl.ColorBuffer, rboColor)
	glx.Targets().RenderBuffer().Bind(rboDepthStencil)
	glx.Targets().Framebuffer().AttachRenderbuffer(gl.DepthStencilBuffer, rboDepthStencil)
	glx.Targets().RenderBuffer().Unbind()
	glx.Targets().Framebuffer().Unbind()
	return fbo, func() {
		fbo.Destroy()
		rboDepthStencil.Destroy()
		rboColor.Destroy()
	}
}
