package mindex

import (
	"fmt"

	"github.com/google/btree"
)

type tupleItem struct {
	primary  Value
	value    Value
	minValue bool
}

func newTuple(primary, value Value) tupleItem {
	return tupleItem{
		primary:  primary,
		value:    value,
		minValue: false,
	}
}

func tupleCopy(primary, value Value) tupleItem {
	return newTuple(CopyValue(primary), CopyValue(value))
}

func (t tupleItem) Less(than btree.Item) bool {
	t2, ok := than.(tupleItem)
	if !ok {
		panic(fmt.Errorf("invalid than: expected %T, got %T", t, than))
	}
	if t.primary == nil {
		panic(fmt.Errorf("nil primary"))
	}
	if t.primary.Less(t2.primary) {
		return true
	}
	if t2.primary.Less(t.primary) {
		return false
	}
	// t primary == t2 primary

	if t.minValue {
		return true
	}
	if t2.minValue {
		return false
	}

	if t.value != nil && t2.value != nil {
		return t.value.Less(t2.value)
	}
	if t.value == nil && t2.value != nil {
		return true
	}
	return false
}

var _ btree.Item = tupleItem{}
