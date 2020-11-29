package wfc

type Position struct {
	X, Y int
}

func NewPosition(x, y int) Position {
	return Position{
		X: x,
		Y: y,
	}
}

func (pos Position) Add(pos2 Position) Position {
	return Position{
		X: pos.X + pos2.X,
		Y: pos.Y + pos2.Y,
	}
}

func (pos Position) byte(width int) int {
	return pos.X + pos.Y*width
}

func (pos Position) Invert() Position {
	return NewPosition(-pos.X, -pos.Y)
}
