package gl

import "github.com/PieterD/warp/driver"

type Texture struct {
	glx      *Context
	glObject driver.Value
}

type TextureConfig struct {
}

func newTexture(glx *Context, cfg TextureConfig) (*Texture, error) {
	textureObject := glx.constants.CreateTexture()
	return &Texture{
		glx:      glx,
		glObject: textureObject,
	}, nil
}
