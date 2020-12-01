package mazegen

func Box(topLeft, bottomRight Coord) func(c Coord) bool {
	return func(c Coord) bool {
		return c.X >= topLeft.X && c.X <= bottomRight.X &&
			c.Y >= topLeft.Y && c.Y <= bottomRight.Y
	}
}

type Field struct {
	allow func(c Coord) bool
	cells map[Coord]Direction
}

func NewField(allow func(c Coord) bool) *Field {
	return &Field{
		allow: allow,
		cells: make(map[Coord]Direction),
	}
}

func (f *Field) Passages(c Coord) Direction {
	return f.cells[c]
}

func (f *Field) Walls(c Coord) Direction {
	return f.cells[c].Invert()
}

func (f *Field) Allowed(c Coord) bool {
	return f.allow(c)
}

func (f *Field) Exists(c Coord) bool {
	if _, ok := f.cells[c]; ok {
		return true
	}
	return false
}

func (f *Field) Burrow(c Coord) {
	if f.Exists(c) {
		return
	}
	f.cells[c] = Direction(0)
}

func (f *Field) Carve(c Coord, dir Direction) {
	f.cells[c] |= dir
	f.cells[c.Step(dir)] |= dir.Flip()
}
