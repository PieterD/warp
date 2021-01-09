package ecs_test

import (
	"fmt"
	"testing"

	"github.com/PieterD/warp/pkg/ecs"
)

type EntityID struct {
	ID uint64
}

func (e *EntityID) Less(_ bool, than ecs.Value) bool {
	eThan, ok := than.(*EntityID)
	if !ok {
		panic(fmt.Errorf("expected %T, got %T", e, than))
	}
	return e.ID < eThan.ID
}

var _ ecs.Value = &EntityID{}

type Mark int

const (
	Naught Mark = iota + 1
	Cross
)

type CellPosition int

const (
	TopLeft CellPosition = iota + 1
	Top
	TopRight
	Left
	Middle
	Right
	BottomLeft
	Bottom
	BottomRight
)

type CheckState struct {
	Mark Mark
}

func (v *CheckState) Less(highResolution bool, than ecs.Value) bool {
	if !highResolution {
		return false
	}
	vThan, ok := than.(*CheckState)
	if !ok {
		panic(fmt.Errorf("expected %T, got %T", v, than))
	}
	return v.Mark < vThan.Mark
}

var _ ecs.Value = &CheckState{}

func (v *CheckState) IsNaught() bool {
	switch v.Mark {
	case Naught:
		return true
	case Cross:
		return false
	default:
		panic(fmt.Errorf("invalid mark: %d", v.Mark))
	}
}

func (v *CheckState) IsCross() bool {
	switch v.Mark {
	case Naught:
		return false
	case Cross:
		return true
	default:
		panic(fmt.Errorf("invalid mark: %d", v.Mark))
	}
}

type Location struct {
	Pos CellPosition
}

func (v *Location) Less(highResolution bool, than ecs.Value) bool {
	if !highResolution {
		return false
	}
	vThan, ok := than.(*Location)
	if !ok {
		panic(fmt.Errorf("expected %T, got %T", v, than))
	}
	return v.Pos < vThan.Pos
}

type GameState struct {
	db *ecs.DB
}

func NewGameState() *GameState {
	db := ecs.New()
	gs := &GameState{
		db: db,
	}
	var mostRecentEntityID uint64
	nextId := func() *EntityID {
		mostRecentEntityID++
		return &EntityID{mostRecentEntityID}
	}
	db.Assign(nextId(), &Location{TopLeft})
	db.Assign(nextId(), &Location{Top})
	db.Assign(nextId(), &Location{TopRight})

	db.Assign(nextId(), &Location{Left})
	db.Assign(nextId(), &Location{Middle})
	db.Assign(nextId(), &Location{Right})

	db.Assign(nextId(), &Location{BottomLeft})
	db.Assign(nextId(), &Location{Bottom})
	db.Assign(nextId(), &Location{BottomRight})

	return gs
}

func (gs *GameState) Get(pos CellPosition) (m Mark, ok bool) {
	db := gs.db
	id := &EntityID{}
	if !db.FirstPrimary(id, &Location{pos}) {
		return 0, false
	}
	cs := &CheckState{}
	if !db.FirstValue(id, cs) {
		return 0, false
	}
	return cs.Mark, true
}

func (gs *GameState) Set(pos CellPosition, m Mark) {
	db := gs.db
	id := &EntityID{}
	if !db.FirstPrimary(id, &Location{pos}) {
		panic(fmt.Errorf("no primary for location %v", pos))
	}
	if m == 0 {
		db.Unassign(id, &CheckState{})
	}
	cs := &CheckState{m}
	db.Assign(id, cs)
}

func TestTictactoe(t *testing.T) {
	gs := NewGameState()
	_, ok := gs.Get(Middle)
	if ok {
		t.Fatalf("getting returned something")
	}
	gs.Set(Middle, Cross)
	m, ok := gs.Get(Middle)
	if !ok {
		t.Fatalf("expected something")
	}
	if m != Cross {
		t.Fatalf("expected Cross, got %v", m)
	}
}
