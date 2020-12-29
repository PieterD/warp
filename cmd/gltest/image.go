package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"net/http"
)

func toDataURI(img image.Image) string {
	buf := &bytes.Buffer{}
	buf.WriteString("data:image/png;base64,")
	base64Encoder := base64.NewEncoder(base64.StdEncoding, buf)
	pngEncoder := png.Encoder{
		CompressionLevel: png.NoCompression,
	}
	if err := pngEncoder.Encode(base64Encoder, img); err != nil {
		panic(fmt.Errorf("encoding png image: %w", err))
	}
	if err := base64Encoder.Close(); err != nil {
		panic(fmt.Errorf("closing base64 encoder: %w", err))
	}
	return buf.String()
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

func pixelsToImage(width, height int, pixels []byte) image.Image {
	flippedPixels := flipPixels(width, height, pixels)
	img := &image.NRGBA{
		Pix:    flippedPixels,
		Stride: width * 4,
		Rect: image.Rectangle{
			Min: image.Point{
				X: 0,
				Y: 0,
			},
			Max: image.Point{
				X: width,
				Y: height,
			},
		},
	}
	return img
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
