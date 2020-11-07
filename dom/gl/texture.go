package gl

import (
	"fmt"
	"image"

	"github.com/PieterD/warp/driver"
)

type Texture2D struct {
	glx      *Context
	glObject driver.Value
}

type Texture2DConfig struct {
}

func newTexture2D(glx *Context, cfg Texture2DConfig, img image.Image) *Texture2D {
	textureObject := glx.constants.CreateTexture()
	glx.constants.BindTexture(glx.constants.TEXTURE_2D, textureObject)
	defer glx.constants.BindTexture(glx.constants.TEXTURE_2D, glx.factory.Null())

	imageWidth, imageHeight, imageData := imageToBytes(img, img.Bounds())
	jsImageData := glx.factory.Buffer(len(imageData))
	jsImageData.Put(imageData)

	texTarget := glx.constants.TEXTURE_2D
	texLevel := glx.factory.Number(0)
	texFormat := glx.constants.RGBA
	texWidth := glx.factory.Number(float64(imageWidth))
	texHeight := glx.factory.Number(float64(imageHeight))
	texBorder := glx.factory.Number(0)
	texType := glx.constants.UNSIGNED_BYTE
	texPixels := jsImageData.AsUint8Array()
	glx.constants.TexImage2D(texTarget, texLevel, texFormat, texWidth, texHeight, texBorder, texFormat, texType, texPixels)

	//TODO: config this
	glx.constants.TexParameteri(glx.constants.TEXTURE_2D, glx.constants.TEXTURE_MIN_FILTER, glx.constants.NEAREST)
	glx.constants.TexParameteri(glx.constants.TEXTURE_2D, glx.constants.TEXTURE_MAG_FILTER, glx.constants.NEAREST)
	glx.constants.TexParameteri(glx.constants.TEXTURE_2D, glx.constants.TEXTURE_WRAP_S, glx.constants.CLAMP_TO_EDGE)
	glx.constants.TexParameteri(glx.constants.TEXTURE_2D, glx.constants.TEXTURE_WRAP_T, glx.constants.CLAMP_TO_EDGE)

	return &Texture2D{
		glx:      glx,
		glObject: textureObject,
	}
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
			color := img.At(x, y)
			r, g, b, a := color.RGBA()
			buffer = append(buffer, byte(r/0xff), byte(g/0xff), byte(b/0xff), byte(a/0xff))
		}
	}
	if len(buffer) != bufferSize {
		panic(fmt.Errorf("expected buffer size %d, got: %d", bufferSize, len(buffer)))
	}
	return width, height, buffer
}
