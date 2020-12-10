package main

import (
	"context"
	"fmt"
	_ "image/png"
	"math"
	"os"

	"github.com/PieterD/warp/pkg/dom/glutil"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/PieterD/warp/pkg/dom"
	"github.com/PieterD/warp/pkg/driver/wasmjs"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "running warplay: %v", err)
	}
	<-make(chan struct{})
}

type rendererState struct {
	currentlyRotating bool
	startVec          mgl32.Vec3
	startCamera       mgl32.Quat
	currentVec        mgl32.Vec3

	camera *glutil.Camera
}

func run(ctx context.Context) error {
	factory := wasmjs.Open()
	global := dom.Open(factory)
	doc := global.Window().Document()
	rs := &rendererState{
		camera: glutil.NewCamera(25.0),
	}
	mouseVec := func(canvasElem *dom.Elem, event *dom.Event) mgl32.Vec3 {
		me, ok := event.AsMouse()
		if !ok {
			panic(fmt.Errorf("expected mouse event"))
		}
		x := me.OffsetX
		y := me.OffsetY
		w, h := dom.AsCanvas(canvasElem).OuterSize()
		v := mgl32.Vec3{
			float32(x)/float32(w) - 0.5,
			float32(h-y)/float32(h) - 0.5,
			0.5,
		}.Normalize()
		return v
	}
	canvasElem := doc.CreateElem("canvas", func(canvasElem *dom.Elem) {
		canvasElem.EventHandler("mousedown", func(this *dom.Elem, event *dom.Event) {
			v := mouseVec(canvasElem, event)
			rs.currentlyRotating = true
			rs.startVec = v
			rs.startCamera = rs.camera.Rotation
			rs.currentVec = v
		})
		canvasElem.EventHandler("mousemove", func(this *dom.Elem, event *dom.Event) {
			v := mouseVec(canvasElem, event)
			if rs.currentlyRotating {
				rs.currentVec = v
				rs.camera.Rotation = mgl32.QuatBetweenVectors(rs.startVec, rs.currentVec).Mul(rs.startCamera)
			}
		})
		canvasElem.EventHandler("mouseup", func(this *dom.Elem, event *dom.Event) {
			v := mouseVec(canvasElem, event)
			if rs.currentlyRotating {
				rs.currentVec = v
				rs.camera.Rotation = mgl32.QuatBetweenVectors(rs.startVec, rs.currentVec).Mul(rs.startCamera)
				rs.currentlyRotating = false
			}
		})
		canvasElem.EventHandler("mouseout", func(this *dom.Elem, event *dom.Event) {
			v := mouseVec(canvasElem, event)
			if rs.currentlyRotating {
				rs.currentVec = v
				rs.camera.Rotation = mgl32.QuatBetweenVectors(rs.startVec, rs.currentVec).Mul(rs.startCamera)
				rs.currentlyRotating = false
			}
		})
	})
	doc.Body().AppendChildren(
		canvasElem,
	)

	canvas := dom.AsCanvas(canvasElem)
	glx := canvas.GetContextWebgl()

	heartProgram, err := NewHeartProgram(glx, "/models/12190_Heart_v1_L3.obj", "/texture.png")
	if err != nil {
		doc.Body().ClearChildren()
		doc.Body().AppendChildren(
			doc.CreateElem("label", func(labelElem *dom.Elem) {
				labelElem.SetText(fmt.Sprintf("error building heart program: %v", err))
			}),
		)
		return nil
	}
	global.Window().Animate(ctx, func(ctx context.Context, millis float64) error {
		select {
		case <-ctx.Done():
			return fmt.Errorf("animate call for renderSquares: %w", ctx.Err())
		default:
		}
		w, h := canvas.OuterSize()
		canvas.SetInnerSize(w, h)
		glx.Viewport(0, 0, w, h)
		glx.Clear()

		_, rot := math.Modf(millis / 4000.0)
		if err := heartProgram.Draw(rs.camera, rot); err != nil {
			return fmt.Errorf("drawing heart program: %w", err)
		}
		return nil
	})

	return nil
}
