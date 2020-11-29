package wfc

import "fmt"

type (
	Slide3 [9]byte
)

func (slide Slide3) At(p Position) (byte, error) {
	if p.X < -1 || p.X > 1 || p.Y < -1 || p.Y > 1 {
		return 0, fmt.Errorf("position out of bounds: %v", p)
	}
	return slide[p.X+p.Y*3], nil
}

func (slide Slide3) Visit(f func(pos Position, color byte) error) error {
	for y := -1; y <= 1; y++ {
		for x := -1; x <= 1; x++ {
			pos := NewPosition(x, y)
			color := slide[pos.byte(3)]
			if err := f(pos, color); err != nil {
				return err
			}
		}
	}
	return nil
}
