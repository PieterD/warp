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
		//w, h := canvas.OuterSize()
		_, rot := math.Modf(millis / 2000.0)
		if err := render(rot); err != nil {
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

func buildRenderer(glx *gl.Context) (renderFunc func(rot float64) error, err error) {
	textureImage, err := loadTexture("texture.png")
	if err != nil {
		return nil, fmt.Errorf("getting texture: %w", err)
	}

	vertices := []float32{
		0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.0, 1.0,
		-0.5, 0.5, 0.0, 1.0,
		0.5, -0.5, 0.0, 1.0,
	}
	texCoords := []float32{
		1.0, 1.0,
		0.0, 0.0,
		0.0, 1.0,
		1.0, 0.0,
	}
	color := []float32{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 0.0, 1.0,
		1.0, 1.0, 1.0,
	}
	elements := []uint16{
		0, 1, 2,
		3, 1, 0,
	}

	programConfig := gl.ProgramConfig{
		VertexCode: `#version 100
precision highp float; // mediump

attribute vec4 Coordinates;
attribute vec3 Color;
attribute vec2 TexCoord;
uniform mat4 Transform;
varying vec4 color;
varying vec2 texCoord;

void main(void) {
	color = vec4(Color, 1.0);
	texCoord = TexCoord;
	gl_Position = Transform * Coordinates;
}
`,
		FragmentCode: `#version 100
precision highp float; // mediump

varying vec4 color;
varying vec2 texCoord;
uniform sampler2D Texture;

void main(void) {
	vec4 texColor = texture2D(Texture, texCoord);
	gl_FragColor = mix(color, texColor, 0.5);
}
`,
	}

	program, err := glx.Program(programConfig)
	if err != nil {
		return nil, fmt.Errorf("compiling shader: %w", err)
	}
	uniformTransform := program.Uniform("Transform")
	if uniformTransform == nil {
		return nil, fmt.Errorf("transform uniform not found")
	}
	uniformSampler := program.Uniform("Texture")
	if uniformSampler == nil {
		return nil, fmt.Errorf("sampler uniform not found")
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

	texAttr, err := program.Attribute("TexCoord")
	if err != nil {
		return nil, fmt.Errorf("fetching TexCoord attribute: %w", err)
	}
	texBuffer := glx.Buffer()
	texBuffer.VertexData(texCoords)

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
			{
				Buffer: texBuffer,
				Attr:   texAttr,
				Layout: gl.VertexArrayAttributeLayout{},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("creating vertex array object: %w", err)
	}
	elementBuffer := glx.Buffer()
	elementBuffer.IndexData(elements)

	texture := glx.Texture(gl.Texture2DConfig{}, textureImage)
	glx.BindTextureUnits(texture)

	return func(rot float64) error {
		err := glx.Draw(gl.DrawConfig{
			Use: program,
			Uniforms: func(us *gl.UniformSetter) {
				angle := 2 * math.Pi * rot
				transform := mgl32.HomogRotate3DZ(float32(angle))
				us.Mat4(uniformTransform, transform)
				us.Int(uniformSampler, 0)
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
