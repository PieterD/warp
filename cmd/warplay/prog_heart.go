package main

import (
	"fmt"
	"math"

	"github.com/PieterD/warp/pkg/gl/glutil"

	"github.com/PieterD/warp/pkg/gfx"
	"github.com/go-gl/mathgl/mgl32"
)

type HeartProgram struct {
	Uniforms struct {
		Model          mgl32.Mat4
		View           mgl32.Mat4
		Projection     mgl32.Mat4
		LightLocation  mgl32.Vec3
		CameraLocation mgl32.Vec3
	}
	glx              *gfx.Context
	verticesToRender int
	vertexBuffer     *gfx.Buffer
	elementBuffer    *gfx.Buffer
	vao              *gfx.VertexArray
	texture          *gfx.Texture2D
	program          *gfx.Program
}

func NewHeartProgram(glx *gfx.Context, modelPath string, texturePath string) (*HeartProgram, error) {
	model, err := loadModel(modelPath)
	if err != nil {
		return nil, fmt.Errorf("getting model: %w", err)
	}
	texture, err := loadTexture(texturePath)
	if err != nil {
		return nil, fmt.Errorf("getting texture: %w", err)
	}
	p := &HeartProgram{
		glx: glx,
	}
	heartVertices, heartIndices, err := model.Interleaved()
	if err != nil {
		return nil, fmt.Errorf("generating interleaved arrays: %w", err)
	}
	p.verticesToRender = len(heartIndices)
	p.vertexBuffer = glx.Buffer()
	p.vertexBuffer.VertexData(heartVertices)
	p.elementBuffer = glx.Buffer()
	p.elementBuffer.IndexData(heartIndices)

	dc, err := gfx.NewDataCoupling(gfx.DataCouplingConfig{
		Vertices: []gfx.VertexConfig{
			{
				Name:    "Coordinates",
				Type:    gfx.Vec3,
				Buffer:  "vertex",
				Padding: 0,
			},
			{
				Name:    "TexCoord",
				Type:    gfx.Vec2,
				Buffer:  "vertex",
				Padding: 4,
			},
			{
				Name:    "Normal",
				Type:    gfx.Vec3,
				Buffer:  "vertex",
				Padding: 0,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("creating data coupling: %w", err)
	}
	adc := dc.Active("Coordinates", "TexCoord", "Normal")

	p.vao, err = glx.VertexArray(adc, map[string]*gfx.Buffer{
		"vertex": p.vertexBuffer,
	})
	if err != nil {
		return nil, fmt.Errorf("creating vertex array object: %w", err)
	}

	programConfig := gfx.ProgramConfig{
		HighPrecision: true,
		Uniforms:      &p.Uniforms,
		Attributes:    adc,
		Textures: []gfx.ProgramSamplerConfig{{
			Name: "Texture",
			Mode: gfx.Sample2D,
		}},
		VertexCode: `
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

	p.program, err = glx.Program(programConfig)
	if err != nil {
		return nil, fmt.Errorf("compiling shader: %w", err)
	}

	p.texture = glx.Texture(gfx.Texture2DConfig{}, texture)

	return p, nil
}

func (p *HeartProgram) Draw(camera *glutil.Camera, rot float64) error {
	deg2rad := float32(math.Pi) / 180.0
	fov := 70 * deg2rad
	lightAngle := float32(rot * 2 * math.Pi)
	lightLocation := mgl32.Vec3{0, 0, 20}
	lightLocation = mgl32.Rotate3DY(lightAngle).Mul3x1(lightLocation)

	drawConfig := gfx.DrawConfig{
		Use:          p.program,
		VAO:          p.vao,
		Textures:     map[string]*gfx.Texture2D{"Texture": p.texture},
		ElementArray: p.elementBuffer,
		DrawMode:     gfx.Triangles,
		Vertices: gfx.VertexRange{
			FirstOffset: 0,
			VertexCount: p.verticesToRender,
		},
		Options: gfx.DrawOptions{
			DepthTest: true,
		},
	}

	p.Uniforms.Model = mgl32.Ident4().
		Mul4(mgl32.Translate3D(0, -10, 0)).
		Mul4(mgl32.HomogRotate3DX(-math.Pi / 2.0))
	p.Uniforms.View = camera.ViewMatrix()
	p.Uniforms.Projection = mgl32.Perspective(fov, 4.0/3.0, 0.1, 100.0)
	p.Uniforms.LightLocation = lightLocation
	p.Uniforms.CameraLocation = camera.Location()
	if err := p.program.UpdateUniforms(); err != nil {
		return fmt.Errorf("updaring uniforms: %w", err)
	}
	err := p.glx.Draw(drawConfig)
	if err != nil {
		return fmt.Errorf("drawing: %w", err)
	}

	p.Uniforms.Model = mgl32.Ident4().
		Mul4(mgl32.Translate3D(lightLocation[0], lightLocation[1], lightLocation[2])).
		Mul4(mgl32.Scale3D(0.2, 0.2, 0.2))
	if err := p.program.UpdateUniforms(); err != nil {
		return fmt.Errorf("updaring uniforms: %w", err)
	}
	err = p.glx.Draw(drawConfig)
	if err != nil {
		return fmt.Errorf("drawing: %w", err)
	}

	return nil
}
