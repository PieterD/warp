package glutil

import (
	"fmt"
	"image"

	"github.com/PieterD/warp/pkg/gl"
)

type SpriteMapConfig struct {
	SpriteSize int // Square
	GridSize   int // Square
}

type SpriteAtlas struct {
	glx        *gl.Context
	spriteSize int
	gridSize   int
	used       int
	texture    gl.TextureObject
}

func NewSpriteAtlas(glx *gl.Context, cfg SpriteMapConfig) *SpriteAtlas {
	atlas := &SpriteAtlas{
		glx:        glx,
		spriteSize: cfg.SpriteSize,
		gridSize:   cfg.GridSize,
		used:       0,
		texture:    glx.CreateTexture(),
	}
	return atlas
}

func (atlas *SpriteAtlas) Destroy() {
	atlas.texture.Destroy()
}

func (atlas *SpriteAtlas) Bind(textureUnit int) {
	glx := atlas.glx
	glx.Targets().ActiveTextureUnit(textureUnit)
	glx.Targets().Texture2D().Bind(atlas.texture)
}

func (atlas *SpriteAtlas) Unbind() {
	glx := atlas.glx
	glx.Targets().Texture2D().Unbind()
}

func (atlas *SpriteAtlas) Allocate() {
	glx := atlas.glx
	glx.Targets().Texture2D().Allocate(atlas.gridSize*atlas.spriteSize, atlas.gridSize*atlas.spriteSize, 0)
	glx.Targets().Texture2D().Settings(gl.Texture2DConfig{
		Minify:  gl.Nearest,
		Magnify: gl.Nearest,
		WrapS:   gl.ClampToEdge,
		WrapT:   gl.ClampToEdge,
	})
}

func (atlas *SpriteAtlas) Add(img image.Image) (textureId [2]float32, err error) {
	glx := atlas.glx
	if size := img.Bounds().Size(); size.X != atlas.spriteSize || size.Y != atlas.spriteSize {
		return [2]float32{}, fmt.Errorf("image size %v does not match sprite map's sprite size %d", size, atlas.spriteSize)
	}
	index := atlas.used
	if index >= atlas.gridSize*atlas.gridSize {
		return [2]float32{}, fmt.Errorf("sprite map full: it only fits %d images", atlas.gridSize*atlas.gridSize)
	}
	col := index % atlas.gridSize
	row := index / atlas.gridSize

	glx.Targets().Texture2D().SubImage(col*atlas.spriteSize, row*atlas.spriteSize, 0, img)
	atlas.used++
	return [2]float32{float32(col), float32(row)}, nil
}

func (atlas *SpriteAtlas) Scale() float32 {
	return 1.0 / float32(atlas.gridSize)
}

func (atlas *SpriteAtlas) GenerateMipmaps() {
	glx := atlas.glx
	glx.Targets().Texture2D().GenerateMipmap()
}
