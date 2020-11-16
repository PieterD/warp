package main

import (
	"context"
	"fmt"
	"github.com/PieterD/warp/pkg/dom/glutil"
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
				fmt.Printf("%v\n", rs.camera.Location())
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
	render, err := buildRenderer(glx, rs)
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
		_, rot := math.Modf(millis / 4000.0)
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

func buildRenderer(glx *gl.Context, rs *rendererState) (renderFunc func(rot float64) error, err error) {
	textureImage, err := loadTexture("/texture.png")
	if err != nil {
		return nil, fmt.Errorf("getting texture: %w", err)
	}
	heartModel, err := loadModel("/models/12190_Heart_v1_L3.obj")
	//heartModel, err := loadModel("/models/cube.obj")
	if err != nil {
		return nil, fmt.Errorf("getting model: %w", err)
	}

	uniforms := &struct {
		Model          mgl32.Mat4
		View           mgl32.Mat4
		Projection     mgl32.Mat4
		LightLocation  mgl32.Vec3
		CameraLocation mgl32.Vec3
	}{}
	programConfig := gl.ProgramConfig{
		HighPrecision: true,
		Uniform:       uniforms,
		VertexCode: `
in vec3 Coordinates;
in vec2 TexCoord;
in vec3 Normal;
out vec2 texCoord;
out vec3 normal;
out vec3 fragPos;

void main(void) {
	texCoord = TexCoord;
    normal = mat3(transpose(inverse(Uniforms.Model))) * Normal;
    fragPos = vec3(Uniforms.Model * vec4(Coordinates, 1.0));
    gl_Position = Uniforms.Projection * Uniforms.View * vec4(fragPos, 1.0);
}
`,
		FragmentCode: `
in vec2 texCoord;
in vec3 normal;
in vec3 fragPos;
out vec4 FragColor;
uniform sampler2D Texture;

void main(void) {
	vec3 lightColor = vec3(1.0, 0.0, 0.0);
	float shininess = 32.0;
	float ambientStrength = 0.1;
	float diffuseStrength = 0.5;
    float specularStrength = 0.5;

	vec3 ambient = ambientStrength * lightColor;

	vec3 norm = normalize(normal);
	vec3 lightDir = normalize(Uniforms.LightLocation - fragPos);
	float diff = max(dot(norm, lightDir), 0.0);
	vec3 diffuse = diffuseStrength * diff * lightColor;

    vec3 viewDir = normalize(Uniforms.CameraLocation - fragPos);
    vec3 reflectDir = reflect(-lightDir, norm);  
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), shininess);
    vec3 specular = specularStrength * spec * lightColor;

	vec4 texColor = texture(Texture, texCoord);
	FragColor = vec4(ambient + diffuse + specular, 1.0) * texColor;
}
`}

	program, err := glx.Program(programConfig)
	if err != nil {
		return nil, fmt.Errorf("compiling shader: %w", err)
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

	program.Update(func(us *gl.UniformSetter) {
		us.Int(uniformSampler, 0)
	})

	return func(rot float64) error {
		glx.Clear()

		lightAngle := float32(rot * 2 * math.Pi)
		lightLocation := mgl32.Vec3{0, 0, 20}
		lightLocation = mgl32.Rotate3DY(lightAngle).Mul3x1(lightLocation)

		deg2rad := float32(math.Pi) / 180.0
		fov := 70 * deg2rad
		uniforms.Model = mgl32.Ident4().
			Mul4(mgl32.Translate3D(0, -10, 0)).
			Mul4(mgl32.HomogRotate3DX(-math.Pi / 2.0))
		uniforms.View = rs.camera.ViewMatrix()
		uniforms.Projection = mgl32.Perspective(fov, 4.0/3.0, 0.1, 100.0)
		uniforms.LightLocation = lightLocation
		uniforms.CameraLocation = rs.camera.Location()
		if err := program.UpdateUniforms(); err != nil {
			return fmt.Errorf("updaring uniforms: %w", err)
		}

		err := glx.Draw(gl.DrawConfig{
			Use:          program,
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

		uniforms.Model = mgl32.Ident4().
			Mul4(mgl32.Translate3D(lightLocation[0], lightLocation[1], lightLocation[2])).
			Mul4(mgl32.Scale3D(0.2, 0.2, 0.2))
		if err := program.UpdateUniforms(); err != nil {
			return fmt.Errorf("updaring uniforms: %w", err)
		}
		err = glx.Draw(gl.DrawConfig{
			Use:          program,
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
