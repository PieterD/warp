package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/PieterD/warp/dom"
	"github.com/PieterD/warp/dom/gl"
	"github.com/PieterD/warp/driver/wasmjs"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "running warplay: %v", err)
	}
	<-make(chan struct{})
}

func run() error {
	ctx, _ := context.WithCancel(context.Background())
	factory := wasmjs.Open()
	global := dom.Open(factory)
	doc := global.Window().Document()
	canvasElem := doc.CreateElem("canvas", nil)
	doc.Body().AppendChildren(
		canvasElem,
	)

	canvas := dom.AsCanvas(canvasElem)
	glx := canvas.GetContextWebgl()
	render, err := buildRenderer(glx)
	if err != nil {
		panic(fmt.Errorf("building renderer: %w", err))
	}
	global.Window().Animate(ctx, func(ctx context.Context, millis float64) error {
		select {
		case <-ctx.Done():
			return fmt.Errorf("animate call for renderSquares: %w", ctx.Err())
		default:
		}
		w, h := canvas.OuterSize()
		if err := render(w, h); err != nil {
			return fmt.Errorf("calling render: %w", err)
		}
		return nil
	})

	return nil
}

func buildRenderer(glx *gl.Context) (renderFunc func(w, h int) error, err error) {
	vertices := []float32{
		0.75, 0.75, 0.0, 1.0,
		0.75, -0.75, 0.0, 1.0,
		-0.75, -0.75, 0.0, 1.0,
	}

	programConfig := gl.ProgramConfig{
		VertexCode: `
attribute vec4 coordinates;

uniform mat4 transform;

void main(void) {
	gl_Position = transform * coordinates;
}
`,
		FragmentCode: `
precision mediump float; // highp

uniform float height;

void main(void) {
	float lerpValue = gl_FragCoord.y / height;
	gl_FragColor = mix(vec4(0.25, 0.25, 0.25, 1.0), vec4(1.0, 1.0, 1.0, 1.0), lerpValue);
}
`,
	}

	program, err := glx.Program(programConfig)
	if err != nil {
		return nil, fmt.Errorf("compiling shader: %w", err)
	}
	uniformHeight := program.Uniform("height")
	if uniformHeight == nil {
		return nil, fmt.Errorf("height uniform not found")
	}
	uniformTransform := program.Uniform("transform")
	if uniformTransform == nil {
		return nil, fmt.Errorf("transform uniform not found")
	}
	coordAttr, err := program.Attribute("coordinates")
	if err != nil {
		return nil, fmt.Errorf("fetching coordinate attribute: %w", err)
	}
	vertexBuffer, err := glx.Buffer()
	if err != nil {
		return nil, fmt.Errorf("creating vertex buffer: %w", err)
	}
	vertexBuffer.VertexData(vertices)
	vao, err := glx.VertexArray(gl.VertexArrayConfig{
		Attributes: []gl.VertexArrayAttribute{
			{
				Buffer: vertexBuffer,
				Attr:   coordAttr,
				Layout: gl.VertexArrayAttributeLayout{},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("creating vertex array object: %w", err)
	}
	transform := mgl32.Ident4().Mul4(mgl32.HomogRotate3DZ(1.0))

	return func(w, h int) error {
		err := glx.Draw(gl.DrawConfig{
			Use: program,
			Uniforms: func(us *gl.UniformSetter) {
				us.Float32(uniformHeight, float32(h))
				us.Mat4(uniformTransform, transform)
			},
			VAO:      vao,
			DrawMode: gl.Triangles,
			Vertices: gl.VertexRange{
				FirstOffset: 0,
				VertexCount: 3,
			},
		})
		if err != nil {
			return fmt.Errorf("drawing: %w", err)
		}
		return nil
	}, nil
}
