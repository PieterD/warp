package ecs

type (
	EntityID    uint64
	EntityIDSet map[EntityID]struct{}
)

func NewEntityIDSet(ids ...EntityID) EntityIDSet {
	s := make(map[EntityID]struct{})
	for _, id := range ids {
		s[id] = struct{}{}
	}
	return s
}
