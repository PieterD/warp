package ecs_test

import (
	"fmt"
	"github.com/PieterD/warp/pkg/ecs"
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
	vThan, ok := than.(*CheckState)
	if !ok {
		panic(fmt.Errorf("expected %T, got %T", v, than))
	}
	return v.Mark < vThan.Mark
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
	vThan, ok := than.(*Location)
	if !ok {
		panic(fmt.Errorf("expected %T, got %T", v, than))
	}
	return v.CellId < vThan.CellId
}

type GameState struct {
	db *ecs.DB
}

func NewGameState() *GameState {
	db, err := ecs.New([]ecs.Lesser{&Location{}, &CheckState{}}, nil)
	if err != nil {
		panic(fmt.Errorf("creating db: %v", err))
	}
	gs := &GameState{
		db: db,
	}
	db.SetComponent(db.Seq("entity"), &Location{TopLeft})
	db.SetComponent(db.Seq("entity"), &Location{Top})
	db.SetComponent(db.Seq("entity"), &Location{TopRight})
	db.SetComponent(db.Seq("entity"), &Location{Left})
	db.SetComponent(db.Seq("entity"), &Location{Middle})
	db.SetComponent(db.Seq("entity"), &Location{Right})
	db.SetComponent(db.Seq("entity"), &Location{BottomLeft})
	db.SetComponent(db.Seq("entity"), &Location{Bottom})
	db.SetComponent(db.Seq("entity"), &Location{BottomRight})

	return gs
}

func TestTictactoe(t *testing.T) {
	gs := NewGameState()
	_ = gs
	fmt.Printf("hello\n")
}
