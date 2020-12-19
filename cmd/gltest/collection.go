package main

import (
	"context"
	"fmt"
	"image"

	"github.com/PieterD/warp/pkg/dom/gl/raw"
)

type TestableFunc func(glx *raw.Context, fbo raw.FramebufferObject) error

type testCollection struct {
	testsByName   map[string]test
	resultsByName map[string]result
}

var globalTestCollection = &testCollection{
	testsByName:   make(map[string]test),
	resultsByName: make(map[string]result),
}

func Register(name string, description string, tf TestableFunc) {
	globalTestCollection.register(name, description, tf)
}

func (tc *testCollection) register(name string, description string, f TestableFunc) {
	_, ok := tc.testsByName[name]
	if ok {
		panic(fmt.Errorf("test with name %s already registered", name))
	}
	tc.testsByName[name] = test{
		name:        name,
		description: description,
		f:           f,
	}
}

type Test struct {
	Name        string
	Description string
	Image       image.Image
	Error       error
}

func Run(ctx context.Context, glx *raw.Context) (tests []Test) {
	globalTestCollection.runAll(ctx, glx)
	return globalTestCollection.results()
}

func (tc *testCollection) runAll(ctx context.Context, glx *raw.Context) {
	const (
		fbWidth  = 128
		fbHeight = 128
	)
	rboColor := glx.CreateRenderbuffer()
	defer rboColor.Destroy()
	glx.Targets().RenderBuffer().Bind(rboColor)
	glx.Targets().RenderBuffer().Storage(raw.RenderbufferConfig{
		Type:   raw.ColorBuffer,
		Width:  fbWidth,
		Height: fbHeight,
	})
	rboDepthStencil := glx.CreateRenderbuffer()
	defer rboDepthStencil.Destroy()
	glx.Targets().RenderBuffer().Bind(rboDepthStencil)
	glx.Targets().RenderBuffer().Storage(raw.RenderbufferConfig{
		Type:   raw.DepthStencilBuffer,
		Width:  fbWidth,
		Height: fbHeight,
	})
	glx.Targets().RenderBuffer().Unbind()
	fbo := glx.CreateFramebuffer()
	defer fbo.Destroy()
	glx.Targets().FrameBuffer().Bind(fbo)
	glx.Targets().RenderBuffer().Bind(rboColor)
	glx.Targets().FrameBuffer().AttachRenderbuffer(raw.ColorBuffer, rboColor)
	glx.Targets().RenderBuffer().Bind(rboDepthStencil)
	glx.Targets().FrameBuffer().AttachRenderbuffer(raw.DepthStencilBuffer, rboDepthStencil)
	glx.Targets().RenderBuffer().Unbind()
	glx.Targets().FrameBuffer().Unbind()
	glx.Viewport(0, 0, fbWidth, fbHeight)
	for testName, test := range tc.testsByName {
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
			glx.Targets().FrameBuffer().Bind(fbo)
			glx.ClearColor(0, 0, 0, 1)
			glx.Clear()

			if err := test.f(glx, fbo); err != nil {
				return nil, fmt.Errorf("running test function: %w", err)
			}

			glx.Targets().FrameBuffer().Bind(fbo)
			pixels := glx.Targets().FrameBuffer().ReadPixels(0, 0, fbWidth, fbHeight)
			glx.Targets().FrameBuffer().Unbind()

			img = pixelsToImage(fbWidth, fbHeight, pixels)
			return img, nil
		}()
		//TODO: convert rendered context to image
		tc.resultsByName[testName] = result{
			img: img,
			err: err,
		}
	}
}

func (tc *testCollection) results() (tests []Test) {
	for testName, test := range globalTestCollection.testsByName {
		result, ok := globalTestCollection.resultsByName[testName]
		if !ok {
			panic(fmt.Errorf("expected a result with test %s", testName))
		}
		tests = append(tests, Test{
			Name:        test.name,
			Description: test.description,
			Image:       result.img,
			Error:       result.err,
		})
	}
	return tests
}

type test struct {
	name        string
	description string
	f           TestableFunc
}

type result struct {
	img image.Image
	err error
}
