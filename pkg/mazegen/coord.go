package mazegen

import "fmt"

type Coord struct {
	X, Y int
}

func (coord Coord) Step(dir Direction) Coord {
	switch dir {
	case North:
		return Coord{X: coord.X + 0, Y: coord.Y - 1}
	case South:
		return Coord{X: coord.X + 0, Y: coord.Y + 1}
	case East:
		return Coord{X: coord.X + 1, Y: coord.Y + 0}
	case West:
		return Coord{X: coord.X - 1, Y: coord.Y + 0}
	default:
		panic(fmt.Errorf("expected North, East, South or West, got: %s", dir.String()))
	}
}
