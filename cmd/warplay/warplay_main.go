package main

import (
	"context"
	"fmt"
	"image"
	_ "image/png"
	"math"
	"net/http"
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
		_, rot := math.Modf(millis / 2000.0)
		if err := render(w, h, rot); err != nil {
			return fmt.Errorf("calling render: %w", err)
		}
		return nil
	})

	return nil
}

func loadTexture(fileName string) (image.Image, error) {
	resp, err := http.DefaultClient.Get(fmt.Sprintf("/%s", fileName))
	if err != nil {
		return nil, fmt.Errorf("getting texture: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("unsuccessful status code while getting texture: %d", resp.StatusCode)
	}
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}
	return img, nil
}

func imageToBytes(img image.Image) []byte {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	buffer := make([]byte, 0, 4*width*height)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Min.Y; x++ {
			color := img.At(x, y)
			r, g, b, a := color.RGBA()
			buffer = append(buffer, byte(r/0xff), byte(g/0xff), byte(b/0xff), byte(a/0xff))
		}
	}
	return buffer
}

func buildRenderer(glx *gl.Context) (renderFunc func(w, h int, rot float64) error, err error) {
	//textureImage, err := loadTexture("texture.png")
	//if err != nil {
	//	return nil, fmt.Errorf("getting texture: %w", err)
	//}
	//textureBytes := imageToBytes(textureImage)

	vertices := []float32{
		0.75, 0.75, 0.0, 1.0,
		0.75, -0.75, 0.0, 1.0,
		-0.75, -0.75, 0.0, 1.0,
		-0.75, 0.75, 0.0, 1.0,
	}
	//texCoords := []float32{
	//	0.0, 1.0,
	//	1.0, 1.0,
	//	1.0, 0.0,
	//	0.0, 0.0,
	//}
	color := []float32{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 0.0, 1.0,
		1.0, 0.0, 1.0,
	}
	elements := []uint16{
		0, 1, 2,
		0, 2, 3,
	}

	programConfig := gl.ProgramConfig{
		VertexCode: `#version 100
attribute vec4 Coordinates;
attribute vec3 Color;
varying vec4 color;
uniform mat4 Transform;

void main(void) {
	color = vec4(Color, 1.0);
	gl_Position = Transform * Coordinates;
}
`,
		FragmentCode: `#version 100
precision mediump float; // highp
varying vec4 color;
uniform float Height;

void main(void) {
	float lerpValue = gl_FragCoord.y / Height;
	gl_FragColor = mix(color, vec4(1.0, 1.0, 1.0, 1.0), lerpValue);
}
`,
	}

	program, err := glx.Program(programConfig)
	if err != nil {
		return nil, fmt.Errorf("compiling shader: %w", err)
	}
	uniformHeight := program.Uniform("Height")
	if uniformHeight == nil {
		return nil, fmt.Errorf("height uniform not found")
	}
	uniformTransform := program.Uniform("Transform")
	if uniformTransform == nil {
		return nil, fmt.Errorf("transform uniform not found")
	}

	coordAttr, err := program.Attribute("Coordinates")
	if err != nil {
		return nil, fmt.Errorf("fetching coordinate attribute: %w", err)
	}
	vertexBuffer := glx.Buffer()
	vertexBuffer.VertexData(vertices)

	colorAttr, err := program.Attribute("Color")
	if err != nil {
		return nil, fmt.Errorf("fetching color attribute: %w", err)
	}
	colorBuffer := glx.Buffer()
	colorBuffer.VertexData(color)

	vao, err := glx.VertexArray(gl.VertexArrayConfig{
		Attributes: []gl.VertexArrayAttribute{
			{
				Buffer: vertexBuffer,
				Attr:   coordAttr,
				Layout: gl.VertexArrayAttributeLayout{},
			},
			{
				Buffer: colorBuffer,
				Attr:   colorAttr,
				Layout: gl.VertexArrayAttributeLayout{},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("creating vertex array object: %w", err)
	}
	elementBuffer := glx.Buffer()
	elementBuffer.IndexData(elements)

	return func(w, h int, rot float64) error {
		err := glx.Draw(gl.DrawConfig{
			Use: program,
			Uniforms: func(us *gl.UniformSetter) {
				angle := 2 * math.Pi * rot
				transform := mgl32.HomogRotate3DZ(float32(angle))
				us.Mat4(uniformTransform, transform)
				us.Float32(uniformHeight, float32(h))
			},
			VAO:          vao,
			ElementArray: elementBuffer,
			DrawMode:     gl.Triangles,
			Vertices: gl.VertexRange{
				FirstOffset: 0,
				VertexCount: 6,
			},
		})
		if err != nil {
			return fmt.Errorf("drawing: %w", err)
		}
		return nil
	}, nil
}
