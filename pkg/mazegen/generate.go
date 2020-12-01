package mazegen

import (
	"math/rand"
)

/*
https://weblog.jamisbuck.org/2011/1/24/maze-generation-hunt-and-kill-algorithm.html

- pick a starting cell, make it the maze
- perform a random walk, until the current cell has no unvisited neighbors.
- scan the grid for an unvisited cell adjacent to a visited cell.
- carve a passage between the two, and call it the new starting location.
- repeat from step 2.
*/

func Generate(radius int, seed int64) *Field {
	prng := rand.New(rand.NewSource(seed))
	field := NewField(Box(Coord{X: -radius, Y: -radius}, Coord{X: radius, Y: radius}))
	field.Burrow(Coord{X: 0, Y: 0})
	pickFunc := func(_ Coord, bits Direction) Direction {
		choices := bits.Split()
		pick := prng.Intn(len(choices))
		return choices[pick]
	}
	pick := Coord{X: 0, Y: 0}
	RandomWalk(field, pick, pickFunc)
	for {
		choices := Neighbors(field)
		if len(choices) == 0 {
			break
		}
		pick = choices[prng.Intn(len(choices))]
		RandomWalk(field, pick, pickFunc)
	}
	return field
}

func RandomWalk(field *Field, start Coord, pick func(c Coord, bits Direction) Direction) {
	field.Burrow(start)
	current := start
	for {
		var validBits Direction
		for _, wallDir := range field.Walls(current).Split() {
			destination := current.Step(wallDir)
			if !field.Allowed(destination) {
				continue
			}
			if field.Exists(destination) {
				continue
			}
			validBits |= wallDir
		}
		if validBits == 0 {
			return
		}
		direction := pick(current, validBits)
		field.Carve(current, direction)
		current = current.Step(direction)
	}
}

func Neighbors(f *Field) []Coord {
	var choices []Coord
	for c, passageBits := range f.cells {
		wallBits := passageBits.Invert()
		for _, wallDir := range wallBits.Split() {
			destination := c.Step(wallDir)
			if !f.Allowed(destination) {
				continue
			}
			if f.Exists(destination) {
				continue
			}
			choices = append(choices, c)
		}
	}
	return choices
}
