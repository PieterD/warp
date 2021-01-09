package ecs

import (
	"fmt"
	"testing"
)

type EntityID struct {
	ID uint64
}

func (e EntityID) Less(highResolution bool, than Value) bool {
	eThan, ok := than.(*EntityID)
	if !ok {
		panic(fmt.Errorf("expected %T, got %T", eThan, than))
	}
	return e.ID < eThan.ID
}

type Name struct {
	Name string
}

func (n Name) Less(highResolution bool, than Value) bool {
	nThan, ok := than.(*Name)
	if !ok {
		panic(fmt.Errorf("expected %T, got %T", nThan, than))
	}
	return n.Name < nThan.Name
}

type Age struct {
	Age int
}

func (a Age) Less(highResolution bool, than Value) bool {
	if !highResolution {
		// Only one Age allowed per entity
		return false
	}
	aThan, ok := than.(*Age)
	if !ok {
		panic(fmt.Errorf("expected %T, got %T", aThan, than))
	}
	return a.Age < aThan.Age
}

func TestDB(t *testing.T) {
	db := New()
	db.Assign(&EntityID{1}, &Name{"Alicia"}, &Name{"Allison"})
	db.Assign(&EntityID{1}, &Age{19}, &Name{"Alice"})
	db.Assign(&EntityID{2}, &Name{"Ellie"}, &Name{"Alice"}, &Age{52})

	t.Run("find Values", func(t *testing.T) {
		id := &EntityID{1}
		name := &Name{}
		if ok := db.FirstValue(id, name); !ok {
			t.Fatalf("expected first value")
		}
		if name.Name != "Alice" {
			t.Errorf("expected Alice, got %s", name.Name)
		}
		if ok := db.NextValue(id, name); !ok {
			t.Fatalf("expected next value")
		}
		if name.Name != "Alicia" {
			t.Errorf("expected Alicia, got %s", name.Name)
		}
		if ok := db.NextValue(id, name); !ok {
			t.Fatalf("expected next value")
		}
		if name.Name != "Allison" {
			t.Errorf("expected Allison, got %s", name.Name)
		}
		if ok := db.NextValue(id, name); ok {
			t.Fatalf("expected no more values")
		}
	})

	t.Run("find Primaries", func(t *testing.T) {
		id := &EntityID{1}
		name := &Name{"Alice"}
		if ok := db.FirstPrimary(id, name); !ok {
			t.Fatalf("expected first primary")
		}
		if id.ID != 1 {
			t.Errorf("invalid id, got %d", id.ID)
		}
		if name.Name != "Alice" {
			t.Errorf("invalid name, got %s", name.Name)
		}
		if ok := db.NextPrimary(id, name); !ok {
			t.Fatalf("expected next primary")
		}
		if id.ID != 2 {
			t.Errorf("invalid id, got %d", id.ID)
		}
		if name.Name != "Alice" {
			t.Errorf("invalid name, got %s", name.Name)
		}
		if ok := db.NextPrimary(id, name); ok {
			t.Fatalf("expected no more primary")
		}
		age := &Age{}
		id.ID = 0
		if ok := db.NextPrimary(id, age); !ok {
			t.Fatalf("expected first primary")
		}
	})
}
