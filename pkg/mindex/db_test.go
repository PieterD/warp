package mindex

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type playerID struct {
	id uint64
}

func (id *playerID) Less(value Value) bool {
	id2, ok := value.(*playerID)
	if !ok {
		panic(fmt.Errorf("got %T, wanted %T", value, id))
	}
	return id.id < id2.id
}

var _ Value = &playerID{}

type idToName struct {
	name string
}

func (r *idToName) Less(value Value) bool {
	return false
}

func (r *idToName) Relate(primary Value) (relationPrimary, relationValue Value) {
	//TODO: a type switch on primary values for multiple primaries
	return &nameToId{r.name}, primary
}

var _ Relation = &idToName{}

type nameToId struct {
	name string
}

func (name *nameToId) Less(value Value) bool {
	name2, ok := value.(*nameToId)
	if !ok {
		panic(fmt.Errorf("got %T, wanted %T", value, name2))
	}
	return name.name < name2.name
}

var _ Value = &nameToId{}

type idToAge struct {
	dateOfBirth time.Time
}

func (r *idToAge) Relate(primary Value) (relationPrimary, relationValue Value) {
	return &ageToId{r.dateOfBirth}, primary
}

func (r *idToAge) Less(value Value) bool {
	return false
}

var _ Relation = &idToAge{}

type ageToId struct {
	dateOfBirth time.Time
}

func (r *ageToId) Less(value Value) bool {
	age2, ok := value.(*ageToId)
	if !ok {
		panic(fmt.Errorf("got %T, wanted %T", value, age2))
	}
	return r.dateOfBirth.Before(age2.dateOfBirth)
}

var _ Value = &ageToId{}

func TestNew(t *testing.T) {
	db := New()

	t.Run("content", func(t *testing.T) {
		pId := &playerID{}
		pName := &idToName{}
		pAge := &idToAge{}

		pId.id = 1
		pName.name = "Carol"
		pAge.dateOfBirth = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		db.Assign(pId, pName, pAge)

		pId.id = 2
		pName.name = "Alice"
		pAge.dateOfBirth = time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)
		db.Assign(pId, pName, pAge)

		pId.id = 3
		pName.name = "Bob"
		pAge.dateOfBirth = time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)
		db.Assign(pId, pName, pAge)

		pId.id = 4
		pName.name = "Alice"
		pAge.dateOfBirth = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		db.Assign(pId, pName, pAge)
	})

	type idAndName struct {
		id   uint64
		name string
	}

	t.Run("name traversal", func(t *testing.T) {
		nameRec := &nameToId{}
		pid := &playerID{}
		var got []idAndName
		traverseFrom(db, nameRec, pid, func() bool {
			got = append(got, idAndName{
				id:   pid.id,
				name: nameRec.name,
			})
			return true
		})
		want := []idAndName{
			{id: 0x2, name: "Alice"},
			{id: 0x4, name: "Alice"},
			{id: 0x3, name: "Bob"},
			{id: 0x1, name: "Carol"},
		}
		if !reflect.DeepEqual(got, want) {
			t.Logf("got : %#v", got)
			t.Logf("want: %#v", want)
			t.Fatalf("mismatch")
		}
	})

	t.Run("age traversal", func(t *testing.T) {
		pid := &playerID{}
		ageRec := &ageToId{}
		var got []idAndName
		traverseFrom(db, ageRec, pid, func() bool {
			nameRec := &idToName{}
			if !db.FirstValue(pid, nameRec) {
				t.Errorf("player %d doesn't have a name", pid.id)
			}
			got = append(got, idAndName{
				id:   pid.id,
				name: nameRec.name,
			})
			return true
		})
		want := []idAndName{
			{id: 0x2, name: "Alice"},
			{id: 0x3, name: "Bob"},
			{id: 0x1, name: "Carol"},
			{id: 0x4, name: "Alice"},
		}
		if !reflect.DeepEqual(got, want) {
			t.Logf("got : %#v", got)
			t.Logf("want: %#v", want)
			t.Fatalf("mismatch")
		}
	})

	t.Run("removal", func(t *testing.T) {
		nameRec := &nameToId{}
		pid := &playerID{}
		var got []idAndName
		db.Remove(&playerID{0x4}, &idToName{"Alice"})
		traverseFrom(db, nameRec, pid, func() bool {
			got = append(got, idAndName{
				id:   pid.id,
				name: nameRec.name,
			})
			return true
		})
		want := []idAndName{
			{id: 0x2, name: "Alice"},
			{id: 0x3, name: "Bob"},
			{id: 0x1, name: "Carol"},
		}
		if !reflect.DeepEqual(got, want) {
			t.Logf("got : %#v", got)
			t.Logf("want: %#v", want)
			t.Fatalf("mismatch")
		}
	})
}
