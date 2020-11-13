package main

import (
	"context"
	"fmt"
	"image"
	_ "image/png"
	"math"
	"net/http"
	"os"

	"github.com/PieterD/warp/pkg/dom"
	"github.com/PieterD/warp/pkg/dom/gl"
	"github.com/PieterD/warp/pkg/driver/wasmjs"
	"github.com/PieterD/warp/pkg/mdl"
	"github.com/go-gl/mathgl/mgl32"
)

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
		canvas.SetInnerSize(w, h)
		glx.Viewport(0, 0, w, h)

		_, rot := math.Modf(millis / 2000.0)
		if err := render(rot); err != nil {
			return fmt.Errorf("calling render: %w", err)
		}
		return nil
	})

	return nil
}

func loadTexture(fileName string) (image.Image, error) {
	resp, err := http.DefaultClient.Get(fileName)
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

func loadModel(fileName string) (*mdl.Model, error) {
	resp, err := http.DefaultClient.Get(fileName)
	if err != nil {
		return nil, fmt.Errorf("getting texture: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("unsuccessful status code while getting texture: %d", resp.StatusCode)
	}
	model, err := mdl.FromObj(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("creating model from obj file: %w", err)
	}
	return model, nil
}

func buildRenderer(glx *gl.Context) (renderFunc func(rot float64) error, err error) {
	textureImage, err := loadTexture("/texture.png")
	if err != nil {
		return nil, fmt.Errorf("getting texture: %w", err)
	}
	//heartModel, err := loadModel("/models/12190_Heart_v1_L3.obj")
	heartModel, err := loadModel("/models/square.obj")
	if err != nil {
		return nil, fmt.Errorf("getting model: %w", err)
	}

	programConfig := gl.ProgramConfig{
		VertexCode: `#version 100
precision highp float; // mediump

attribute vec3 Coordinates;
attribute vec2 TexCoord;
uniform mat4 Transform;
varying vec2 texCoord;

void main(void) {
	texCoord = TexCoord;
	gl_Position = Transform * vec4(Coordinates, 1.0);
}
`,
		FragmentCode: `#version 100
precision highp float; // mediump

varying vec2 texCoord;
uniform sampler2D Texture;

void main(void) {
	vec4 texColor = texture2D(Texture, texCoord);
	gl_FragColor = mix(vec4(1.0,0.0,0.0,1.0), texColor, 0.5);
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

	texAttr, err := program.Attribute("TexCoord")
	if err != nil {
		return nil, fmt.Errorf("fetching TexCoord attribute: %w", err)
	}

	heartVertices, heartIndices, err := heartModel.Interleaved()
	if err != nil {
		return nil, fmt.Errorf("generating interleaved arrays: %w", err)
	}
	verticesToRender := len(heartIndices)
	heartVertexBuffer := glx.Buffer()
	heartVertexBuffer.VertexData(heartVertices)
	heartElementBuffer := glx.Buffer()
	heartElementBuffer.IndexData(heartIndices)

	stride := (heartModel.VertexItems + heartModel.TextureItems + heartModel.VertexItems) * 4
	vao, err := glx.VertexArray(
		gl.VertexArrayAttribute{
			Buffer: heartVertexBuffer,
			Attr:   coordAttr,
			Layout: gl.VertexArrayAttributeLayout{
				ByteOffset: 0,
				ByteStride: stride,
			},
		},
		gl.VertexArrayAttribute{
			Buffer: heartVertexBuffer,
			Attr:   texAttr,
			Layout: gl.VertexArrayAttributeLayout{
				ByteOffset: 4 * heartModel.VertexItems,
				ByteStride: stride,
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("creating vertex array object: %w", err)
	}

	texture := glx.Texture(gl.Texture2DConfig{}, textureImage)
	glx.BindTextureUnits(texture)

	return func(rot float64) error {
		glx.Clear()
		err := glx.Draw(gl.DrawConfig{
			Use: program,
			Uniforms: func(us *gl.UniformSetter) {
				angle := 2 * math.Pi * rot
				deg2rad := float32(math.Pi) / 180.0
				fov := 70*deg2rad
				transform := mgl32.Ident4()
				transform = transform.Mul4(mgl32.Perspective(fov, 4.0/3.0, 0.1, 100.0))
				//transform = transform.Mul4(mgl32.LookAtV(
				//	mgl32.Vec3{0, 0, 0},
				//	mgl32.Vec3{0, 10, 25},
				//	mgl32.Vec3{0, 1, 0},
				//))
				transform = transform.Mul4(mgl32.Translate3D(0.0, 0.0, -5.0))
				transform = transform.Mul4(mgl32.HomogRotate3DY(float32(angle)))
				//transform = transform.Mul4(mgl32.HomogRotate3DX(float32(-math.Pi / 2.0)))
				//transform = transform.Mul4(mgl32.Scale3D(1/20.0, 1/20.0, 1/20.0))
				us.Mat4(uniformTransform, transform)
				us.Int(uniformSampler, 0)
			},
			VAO:          vao,
			ElementArray: heartElementBuffer,
			DrawMode:     gl.Triangles,
			Vertices: gl.VertexRange{
				FirstOffset: 0,
				VertexCount: verticesToRender,
			},
			Options: gl.DrawOptions{
				DepthTest: true,
			},
		})
		if err != nil {
			return fmt.Errorf("drawing: %w", err)
		}
		return nil
	}, nil
}
