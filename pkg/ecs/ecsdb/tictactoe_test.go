package ecsdb_test

import (
	"fmt"
	"github.com/PieterD/warp/pkg/ecs/ecsdb"
	"testing"
)

type Mark int

const (
	Naught Mark = iota + 1
	Cross
)

type CellID int

const (
	TopLeft CellID = iota + 1
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

func (v *CheckState) Less(than interface{}) bool {
	return false
}

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
	CellId CellID
}

func (v *Location) Less(than interface{}) bool {
	return false
}

type GameState struct {
	db *ecsdb.DB
}

func NewGameState() *GameState {
	db, err := ecsdb.New(&Location{}, &CheckState{})
	if err != nil {
		panic(fmt.Errorf("creating db: %v", err))
	}
	gs := &GameState{
		db: db,
	}
	db.SetComponent(db.Sequence("entity"), &Location{TopLeft})
	db.SetComponent(db.Sequence("entity"), &Location{Top})
	db.SetComponent(db.Sequence("entity"), &Location{TopRight})
	db.SetComponent(db.Sequence("entity"), &Location{Left})
	db.SetComponent(db.Sequence("entity"), &Location{Middle})
	db.SetComponent(db.Sequence("entity"), &Location{Right})
	db.SetComponent(db.Sequence("entity"), &Location{BottomLeft})
	db.SetComponent(db.Sequence("entity"), &Location{Bottom})
	db.SetComponent(db.Sequence("entity"), &Location{BottomRight})

	return gs
}

func (gs *GameState) Check(coord CellID, mark Mark) error {
	db := gs.db
	location := &Location{}
	db.FindEntities(func(id ecsdb.ID) bool {
		if location.CellId == coord {
			
		}
		return true
	}, location)
	db.GetComponent()
	switch mark {
	case Naught:
	case Cross:
	}
}

func TestTictactoe(t *testing.T) {
	db, err := ecsdb.New(&Location{}, &CheckState{})
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
}
