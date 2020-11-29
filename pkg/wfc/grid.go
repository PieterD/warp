package wfc

import (
	"bytes"
	"fmt"
	"sort"
)

type (
	Grid3 struct {
		slides       []Slide3
		collapsed    map[Position]Slide3
		dead         map[Position]struct{}
		coefficients map[Position]map[int]struct{} //TODO: bitset
		//transform func(pos Position) (newPos Position, allowed bool)
	}
)

func NewGrid3(deck *Deck3) (*Grid3, error) {
	if deck == nil {
		return nil, fmt.Errorf("deck is nil")
	}
	if len(deck.slides) == 0 {
		return nil, fmt.Errorf("no slides in deck")
	}
	slides := make([]Slide3, 0, len(deck.slides))
	for slide := range deck.slides {
		slides = append(slides, slide)
	}
	sort.Slice(slides, func(i, j int) bool {
		return bytes.Compare(slides[i][:], slides[j][:]) < 0
	})
	return &Grid3{
		slides:       slides,
		collapsed:    make(map[Position]Slide3),
		dead:         make(map[Position]struct{}),
		coefficients: make(map[Position]map[int]struct{}),
	}, nil
}

func (g *Grid3) Paint(pos Position, color byte) error {
	err := Slide3{}.Visit(func(offset Position, _ byte) error {
		gridPos := pos.Add(offset)
		slidePos := offset.Invert()
		if err := g.filter(gridPos, slidePos, color); err != nil {
			return fmt.Errorf("filtering gridPos %v slidePos %v color %v: %w", gridPos, slidePos, color, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("visiting 3x3 centered on 0,0: %w", err)
	}
	return nil
}

func (g *Grid3) filter(gridPos, slidePos Position, color byte) error {
	if _, ok := g.dead[gridPos]; ok {
		return nil
	}

	if slide, ok := g.collapsed[gridPos]; ok {
		slideColor, err := slide.At(slidePos)
		if err != nil {
			return fmt.Errorf("fetching color from collapsed slide: %w", err)
		}
		if slideColor != color {
			return fmt.Errorf("color mismatch")
		}
		return nil
	}

	coefficients, ok := g.coefficients[gridPos]
	if !ok {
		// empty coefficients count as all coefficients when not present in the other maps.
		coefficients = make(map[int]struct{}, len(g.slides))
		for i := range g.slides {
			coefficients[i] = struct{}{}
		}
		g.coefficients[gridPos] = coefficients
	}

	for i := range coefficients {
		slide := g.slides[i]
		slideColor, err := slide.At(slidePos)
		if err != nil {
			return fmt.Errorf("slide %v fetching color at %v", gridPos, slidePos)
		}
		if slideColor == color {
			continue
		}
		delete(coefficients, i)
	}
	if len(coefficients) == 1 {
		for idx := range coefficients {
			g.collapsed[gridPos] = g.slides[idx]
		}
		delete(g.coefficients, gridPos)
	}
	if len(coefficients) == 0 {
		g.dead[gridPos] = struct{}{}
	}
	return nil
}
