package main

import (
	"fmt"
	"github.com/PieterD/warp/pkg/mdl"
	"image"
	"net/http"
)

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
