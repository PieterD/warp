package main

import (
	"context"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	_ "image/png"
	"math"
	"net/http"
	"os"

	"github.com/PieterD/warp/pkg/dom"
	"github.com/PieterD/warp/pkg/dom/gl"
	"github.com/PieterD/warp/pkg/driver/wasmjs"
	"github.com/PieterD/warp/pkg/mdl"
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

		//_, rot := math.Modf(millis / 2000.0)
		_, rot := math.Modf(millis / 10000.0)
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
	heartModel, err := loadModel("/models/cube.obj")
	if err != nil {
		return nil, fmt.Errorf("getting model: %w", err)
	}

	programConfig := gl.ProgramConfig{
		VertexCode: `#version 300 es
precision highp float; // mediump

in vec3 Coordinates;
in vec2 TexCoord;
in vec3 Normal;
uniform mat4 Model;
uniform mat4 View;
uniform mat4 Projection;
out vec2 texCoord;
out vec3 normal;
out vec3 fragPos;

void main(void) {
	texCoord = TexCoord;
	normal = Normal;
    //normal = mat3(transpose(inverse(Model))) * Normal;
    fragPos = vec3(Model * vec4(Coordinates, 1.0));
    gl_Position = Projection * View * vec4(fragPos, 1.0);
}
`,
		FragmentCode: `#version 300 es
precision highp float; // mediump

in vec2 texCoord;
in vec3 normal;
in vec3 fragPos;
uniform sampler2D Texture;
uniform vec3 LightPos;
out vec4 FragColor;

void main(void) {
	vec3 lightPos = vec3(5.0, 5.0, 5.0);
	float ambientStrength = 0.1;
	vec3 lightColor = vec3(1.0, 0.0, 0.0);
	vec3 ambient = ambientStrength * lightColor;

	vec3 norm = normalize(normal);
	vec3 lightDir = normalize(LightPos - fragPos);
	float diff = max(dot(norm, lightDir), 0.0);
	vec3 diffuse = diff * lightColor;

	vec4 texColor = texture(Texture, texCoord);
	vec4 objectColor = vec4(1.0, 0.5, 0.3, 1.0);
	FragColor = vec4(ambient + diffuse, 1.0) * objectColor;
}
`,
	}

	program, err := glx.Program(programConfig)
	if err != nil {
		return nil, fmt.Errorf("compiling shader: %w", err)
	}
	uniformModel := program.Uniform("Model")
	if uniformModel == nil {
		return nil, fmt.Errorf("model uniform not found")
	}
	uniformView := program.Uniform("View")
	if uniformView == nil {
		return nil, fmt.Errorf("view uniform not found")
	}
	uniformProjection := program.Uniform("Projection")
	if uniformProjection == nil {
		return nil, fmt.Errorf("projection uniform not found")
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

	normalAttr, err := program.Attribute("Normal")
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
		gl.VertexArrayAttribute{
			Buffer: heartVertexBuffer,
			Attr:   normalAttr,
			Layout: gl.VertexArrayAttributeLayout{
				ByteOffset: (heartModel.VertexItems + heartModel.TextureItems) * 4,
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
		program.Update(func(us *gl.UniformSetter) {
			angle := 2 * math.Pi * rot
			deg2rad := float32(math.Pi) / 180.0
			fov := 70 * deg2rad
			modelMatrix := mgl32.Ident4().
				Mul4(mgl32.HomogRotate3DY(float32(angle)))
			cameraLocation := mgl32.Vec3{0, 0, 5}
			cameraTarget := mgl32.Vec3{0, 0, 0}
			up := mgl32.Vec3{0, 1, 0}
			viewMatrix := mgl32.Ident4().
				Mul4(mgl32.LookAtV(
					cameraLocation,
					cameraTarget,
					up,
				))
			projectionMatrix := mgl32.Ident4().
				Mul4(mgl32.Perspective(fov, 4.0/3.0, 0.1, 100.0))
			us.Mat4(uniformModel, modelMatrix)
			us.Mat4(uniformView, viewMatrix)
			us.Mat4(uniformProjection, projectionMatrix)
			us.Int(uniformSampler, 0)
		})
		err := glx.Draw(gl.DrawConfig{
			Use: program,
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
