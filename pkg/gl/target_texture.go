package gl

import (
	"fmt"
	"image"
	"image/color"

	"github.com/PieterD/warp/pkg/driver"
)

type Texture2DTarget struct {
	glx *Context
}

func (target Texture2DTarget) Bind(texture TextureObject) {
	glx := target.glx
	glx.constants.BindTexture(
		glx.constants.TEXTURE_2D,
		texture.value,
	)
}

func (target Texture2DTarget) Unbind() {
	glx := target.glx
	glx.constants.BindTexture(
		glx.constants.TEXTURE_2D,
		glx.factory.Null(),
	)
}

func (to TextureObject) Destroy() {
	glx := to.glx
	glx.constants.DeleteTexture(to.value)
}

type TextureFilter int

const (
	Linear TextureFilter = iota + 1
	Nearest
)

func (f TextureFilter) glValue(glx *Context) driver.Value {
	switch f {
	case Linear:
		return glx.constants.LINEAR
	case Nearest:
		return glx.constants.NEAREST
	default:
		panic(fmt.Errorf("unknown texture filter value: %v", f))
	}
}

type WrapFunction int

const (
	Repeat WrapFunction = iota + 1
	ClampToEdge
	MirroredRepeat
)

func (f WrapFunction) glValue(glx *Context) driver.Value {
	switch f {
	case Repeat:
		return glx.constants.REPEAT
	case ClampToEdge:
		return glx.constants.CLAMP_TO_EDGE
	case MirroredRepeat:
		return glx.constants.MIRRORED_REPEAT
	default:
		panic(fmt.Errorf("unknown texture filter value: %v", f))
	}
}

type Texture2DConfig struct {
	Minify  TextureFilter
	Magnify TextureFilter
	WrapS   WrapFunction
	WrapT   WrapFunction
}

func (target Texture2DTarget) Settings(cfg Texture2DConfig) {
	glx := target.glx
	if cfg.Minify != 0 {
		glx.constants.TexParameteri(glx.constants.TEXTURE_2D, glx.constants.TEXTURE_MIN_FILTER, cfg.Minify.glValue(glx))
	}
	if cfg.Magnify != 0 {
		glx.constants.TexParameteri(glx.constants.TEXTURE_2D, glx.constants.TEXTURE_MAG_FILTER, cfg.Magnify.glValue(glx))
	}
	if cfg.WrapS != 0 {
		glx.constants.TexParameteri(glx.constants.TEXTURE_2D, glx.constants.TEXTURE_WRAP_S, cfg.WrapS.glValue(glx))
	}
	if cfg.WrapT != 0 {
		glx.constants.TexParameteri(glx.constants.TEXTURE_2D, glx.constants.TEXTURE_WRAP_T, cfg.WrapT.glValue(glx))
	}
}

func (target Texture2DTarget) GenerateMipmap() {
	glx := target.glx
	glx.constants.GenerateMipmap(glx.constants.TEXTURE_2D)
}

func (target Texture2DTarget) Allocate(width, height, level int) {
	glx := target.glx
	glx.constants.TexImage2D(
		glx.constants.TEXTURE_2D,
		glx.factory.Number(float64(0)),
		glx.constants.RGBA,
		glx.factory.Number(float64(width)),
		glx.factory.Number(float64(height)),
		glx.factory.Number(0), // border (must be 0)
		glx.constants.RGBA,
		glx.constants.UNSIGNED_BYTE,
		glx.factory.Null(), // pixels
	)
}

func (target Texture2DTarget) SubImage(x, y, level int, img image.Image) {
	glx := target.glx
	imageWidth, imageHeight, imageData := imageToBytes(img, img.Bounds())
	jsImageData := glx.factory.Buffer(len(imageData))
	jsImageData.Put(imageData)
	glx.constants.TexSubImage2D(
		glx.constants.TEXTURE_2D,
		glx.factory.Number(float64(0)),
		glx.factory.Number(float64(x)),
		glx.factory.Number(float64(y)),
		glx.factory.Number(float64(imageWidth)),
		glx.factory.Number(float64(imageHeight)),
		glx.constants.RGBA,
		glx.constants.UNSIGNED_BYTE,
		jsImageData.AsUint8Array(),
		glx.factory.Number(0), // offset
	)
}

func flipPixels(width, height int, pixels []byte) []byte {
	flippedPixels := make([]byte, len(pixels))
	rowSize := width * 4
	for y := 0; y < height; y++ {
		fy := height - 1 - y
		copy(flippedPixels[fy*rowSize:(fy+1)*rowSize], pixels[y*rowSize:(y+1)*rowSize])
	}
	return flippedPixels
}

func imageToBytes(img image.Image, bounds image.Rectangle) (width, height int, pixels []byte) {
	bounds = bounds.Intersect(img.Bounds())
	width = bounds.Dx()
	height = bounds.Dy()
	if width == 0 || height == 0 {
		return width, height, nil
	}
	bufferSize := 4 * width * height
	buffer := make([]byte, 0, bufferSize)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			iPixel := img.At(x, y)
			switch pixel := iPixel.(type) {
			case color.NRGBA:
				r := pixel.R
				g := pixel.G
				b := pixel.B
				a := pixel.A
				buffer = append(buffer, r, g, b, a)
			case color.RGBA:
				r := pixel.R
				g := pixel.G
				b := pixel.B
				a := pixel.A
				buffer = append(buffer, r, g, b, a)
			default:
				panic(fmt.Errorf("unknown pixel type: %T", iPixel))
			}
		}
	}
	if len(buffer) != bufferSize {
		panic(fmt.Errorf("expected buffer size %d, got: %d", bufferSize, len(buffer)))
	}
	buffer = flipPixels(width, height, buffer)
	return width, height, buffer
}
