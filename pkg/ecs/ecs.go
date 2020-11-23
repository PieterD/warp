package ecs

type (
	EntityID uint64
	System   struct {
		Entities  map[EntityID]struct{}
		Locations map[EntityID]Location
		Passable  map[EntityID]Passable
	}
	Location struct {
		X, Y int
	}
	Passable struct{}
	Walls    struct {
		Faces Direction
	}
	Direction byte
)
