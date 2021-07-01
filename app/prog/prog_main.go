package main

import (
	"context"
	"fmt"
	"github.com/PieterD/warp/pkg/dom"
	"github.com/PieterD/warp/pkg/driver"
	"github.com/PieterD/warp/pkg/driver/wasmjs"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := run(ctx)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "running progressive web application: %v", err)
	}
	<-make(chan struct{})
}

func run(ctx context.Context) error {
	factory := wasmjs.Open()
	global := dom.Open(factory)
	win := global.Window()
	doc := win.Document()
	body := doc.Body()
	labelElem := doc.CreateElem("label", func(labelElem *dom.Elem) {
		labelElem.SetText("No measurements")
	})
	body.AppendChildren(labelElem)

	//_, err := dom.AddEventListener(win, "deviceorientation", func(this driver.Value, event *dom.Event) {
	//	doe, ok := event.AsDeviceOrientationEvent()
	//	if !ok {
	//		panic("not a device orientation event")
	//	}
	//	labelElem.SetText(fmt.Sprintf("orientation: %f %f %f", doe.Alpha, doe.Beta, doe.Gamma))
	//})
	_, err := dom.AddEventListener(win, "devicemotion", func(this driver.Value, event *dom.Event) {
		doe, ok := event.AsDeviceMotionEvent()
		if !ok {
			panic("not a device orientation event")
		}
		labelElem.SetText(fmt.Sprintf("motion: %f %f %f", doe.X, doe.Y, doe.Z))
	})
	if err != nil {
		return fmt.Errorf("error adding event listener: %w", err)
	}

	if err := RegisterServiceWorker(win, "/sw.js"); err != nil {
		return fmt.Errorf("registering service worker: %w", err)
	}

	return nil
}

func RegisterServiceWorker(win *dom.Window, serviceWorkerCodePath string) error {
	factory, winObj := win.Driver()
	navObj, ok := winObj.Get("navigator").ToObject()
	if !ok {
		return fmt.Errorf("no navigator present in window")
	}
	swObj, ok := navObj.Get("serviceWorker").ToObject()
	if !ok {
		return fmt.Errorf("no serviceWorker present in navigator")
	}
	jsRegister := driver.Bind(swObj, "register")
	if jsRegister == nil {
		return fmt.Errorf("no register method present in serviceWorker")
	}
	registerReturn, ok := jsRegister(factory.String(serviceWorkerCodePath)).ToObject()
	if !ok {
		return fmt.Errorf("register did not return an object")
	}
	jsThen := driver.Bind(registerReturn, "then")
	jsThen(
		factory.Function(func(this driver.Object, args ...driver.Value) driver.Value {
			fmt.Println("Service worker registered")
			return nil
		}),
		factory.Function(func(this driver.Object, args ...driver.Value) driver.Value {
			fmt.Println("Error registering service worker")
			return nil
		}),
	)
	return nil
}
