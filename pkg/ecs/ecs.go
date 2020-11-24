package ecs

import (
	"fmt"
	"math/rand"
	"sync"
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
	lock       sync.RWMutex
	entities   EntityIDSet
	locations  map[EntityID]Location
	rLocations map[Location]EntityIDSet
	passables  map[EntityID]Passable
	walls      map[EntityID]Wall
}

func NewRepository() *Repository {
	return &Repository{
		entities:   NewEntityIDSet(),
		locations:  make(map[EntityID]Location),
		rLocations: make(map[Location]EntityIDSet),
		passables:  make(map[EntityID]Passable),
		walls:      make(map[EntityID]Wall),
	}
}

func (repo *Repository) NewEntity() EntityID {
	repo.lock.Lock()
	defer repo.lock.Unlock()

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
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	_, ok := repo.entities[id]
	return ok
}

func (repo *Repository) EntitiesByLocation(l Location) EntityIDSet {
	repo.lock.Lock()
	defer repo.lock.Unlock()

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
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	l, ok := repo.locations[id]
	return l, ok
}

func (repo *Repository) PutLocation(id EntityID, l Location) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

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
	repo.lock.Lock()
	defer repo.lock.Unlock()

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
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	p, ok := repo.passables[id]
	return p, ok
}

func (repo *Repository) PutPassable(id EntityID, p Passable) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	repo.entities[id] = struct{}{}
	repo.passables[id] = p
}

func (repo *Repository) DelPassable(id EntityID) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	delete(repo.passables, id)
}

func (repo *Repository) Wall(id EntityID) (Wall, bool) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	w, ok := repo.walls[id]
	return w, ok
}

func (repo *Repository) PutWall(id EntityID, w Wall) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	repo.entities[id] = struct{}{}
	repo.walls[id] = w
}

func (repo *Repository) DelWall(id EntityID) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	delete(repo.walls, id)
}
