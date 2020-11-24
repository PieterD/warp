package ecs

import (
	"fmt"
	"math/rand"
)

type (
	Location struct {
		X, Y int
	}
	Passable struct{}
	Wall     struct {
		Faces Direction
	}
	Direction byte
)

type Repository struct {
	entities  EntityIDSet
	locations map[EntityID]Location
	passables map[EntityID]Passable
	walls     map[EntityID]Wall

	rLocations map[Location]EntityIDSet
}

func NewRepository() *Repository {
	return &Repository{
		entities:  NewEntityIDSet(),
		locations: make(map[EntityID]Location),
		passables: make(map[EntityID]Passable),
		walls:     make(map[EntityID]Wall),

		rLocations: make(map[Location]EntityIDSet),
	}
}

func (repo *Repository) NewEntity() EntityID {
	for n := 0; n < 256; n++ {
		randomId := EntityID(rand.Int63())
		if _, ok := repo.entities[randomId]; ok {
			continue
		}
		repo.entities[randomId] = struct{}{}
		return randomId
	}
	panic(fmt.Errorf("took too long to find unused id"))
}

func (repo *Repository) EntityExists(id EntityID) bool {
	_, ok := repo.entities[id]
	return ok
}

func (repo *Repository) DelEntity(id EntityID) {
	repo.DelLocation(id)
	repo.DelPassable(id)
	repo.DelWall(id)
	delete(repo.entities, id)
}

func (repo *Repository) EntitiesByLocation(l Location) EntityIDSet {
	oSet := repo.rLocations[l]
	if oSet == nil {
		return nil
	}
	set := make(EntityIDSet)
	for loc := range oSet {
		set[loc] = struct{}{}
	}
	return set
}

func (repo *Repository) Location(id EntityID) (Location, bool) {
	l, ok := repo.locations[id]
	return l, ok
}

func (repo *Repository) PutLocation(id EntityID, l Location) {
	repo.entities[id] = struct{}{}
	repo.locations[id] = l
	s, ok := repo.rLocations[l]
	if !ok {
		s = NewEntityIDSet()
		repo.rLocations[l] = s
	}
	s[id] = struct{}{}
}

func (repo *Repository) DelLocation(id EntityID) {
	location, ok := repo.locations[id]
	if !ok {
		return
	}
	delete(repo.locations, id)
	entitySet, ok := repo.rLocations[location]
	if !ok {
		return
	}
	delete(entitySet, id)
}

func (repo *Repository) Passable(id EntityID) (Passable, bool) {
	p, ok := repo.passables[id]
	return p, ok
}

func (repo *Repository) PutPassable(id EntityID, p Passable) {
	repo.entities[id] = struct{}{}
	repo.passables[id] = p
}

func (repo *Repository) DelPassable(id EntityID) {
	delete(repo.passables, id)
}

func (repo *Repository) Wall(id EntityID) (Wall, bool) {
	w, ok := repo.walls[id]
	return w, ok
}

func (repo *Repository) PutWall(id EntityID, w Wall) {
	repo.entities[id] = struct{}{}
	repo.walls[id] = w
}

func (repo *Repository) DelWall(id EntityID) {
	delete(repo.walls, id)
}
