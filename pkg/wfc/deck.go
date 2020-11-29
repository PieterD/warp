package wfc

import "fmt"

type (
	Deck3 struct {
		slides map[Slide3]int
	}
)

type DecomposeConfig struct {
	WrapH  bool
	WrapV  bool
	Flip   bool
	Rot90  bool
	Rot180 bool
}

func NewDeck3() *Deck3 {
	return &Deck3{
		slides: make(map[Slide3]int),
	}
}

func (deck *Deck3) Add(pixels []byte, width, height int, cfg DecomposeConfig) error {
	if width < 3 {
		return fmt.Errorf("width %d is too small for Deck3", width)
	}
	if height < 3 {
		return fmt.Errorf("height %d is too small for Deck3", height)
	}
	if len(pixels) != width*height {
		return fmt.Errorf("%d pixels in pixel data, expected %d", len(pixels), width*height)
	}
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			var slide Slide3
			pixelIndex := 0
			for slideY := -1; slideY <= 1; slideY++ {
				for slideX := -1; slideX <= 1; slideX++ {
					tX := x + slideX
					tY := y + slideY
					slide[pixelIndex] = pixels[tY*width+tX]
					pixelIndex++
				}
			}
			deck.slides[slide]++
		}
	}
	return nil
}
