package mindex

import (
	"fmt"

	"github.com/google/btree"
)

type DB struct {
	tableKeys map[tableKey]int
	tables    []*dbTable
}

func New() *DB {
	return &DB{
		tableKeys: make(map[tableKey]int),
	}
}

func (db *DB) Assign(primary Value, vs ...Value) {
	for _, v := range vs {
		table := db.table(primary, v)
		table.assign(primary, v)
	}
}

func (db *DB) Remove(primary Value, vs ...Value) {
	for _, v := range vs {
		table := db.table(primary, v)
		table.remove(primary, v)
	}
}

func (db *DB) First(primaryPtr Value, valuePtr Value) bool {
	table := db.table(primaryPtr, valuePtr)
	tup := newTuple(primaryPtr, valuePtr)
	found := false
	table.index.AscendGreaterOrEqual(tup, func(foundItem btree.Item) bool {
		foundTuple, ok := foundItem.(tupleItem)
		if !ok {
			panic(fmt.Errorf("got %T, want %T", foundItem, foundTuple))
		}
		found = true
		CopyValueToValue(primaryPtr, foundTuple.primary)
		CopyValueToValue(valuePtr, foundTuple.value)
		return false
	})
	return found
}

func (db *DB) Next(primaryPtr Value, valuePtr Value) bool {
	table := db.table(primaryPtr, valuePtr)
	tup := newTuple(primaryPtr, valuePtr)
	found := false
	table.index.AscendGreaterOrEqual(tup, func(foundItem btree.Item) bool {
		foundTuple, ok := foundItem.(tupleItem)
		if !ok {
			panic(fmt.Errorf("got %T, want %T", foundItem, foundTuple))
		}
		if tup.Less(foundItem) {
			// tup < foundItem, thus foundItem != tup
			found = true
			CopyValueToValue(primaryPtr, foundTuple.primary)
			CopyValueToValue(valuePtr, foundTuple.value)
			return false
		}
		return true
	})
	return found
}

func (db *DB) FirstValue(primaryPtr Value, valuePtr Value) bool {
	return firstValue(db, primaryPtr, valuePtr)
}

func (db *DB) NextValue(primaryPtr Value, valuePtr Value) bool {
	return nextValue(db, primaryPtr, valuePtr)
}

func (db *DB) TraverseFrom(primaryPtr Value, valuePtr Value, f func() bool) {
	traverseFrom(db, primaryPtr, valuePtr, f)
}

//func (db *DB) TraverseFixedPrimary(primaryPtr Value, valuePtr Value, f func() bool) {
//	panic("not implemented")
//}
//
//func (db *DB) ReverseFixedValue(primaryPtr Value, valuePtr Value, f func() bool) {
//	relPrimary, relValue := MustRelation(valuePtr).Relate(primaryPtr)
//	_ = relPrimary
//	_ = relValue
//	panic("not implemented")
//}
