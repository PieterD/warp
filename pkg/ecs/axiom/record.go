package axiom

import (
	"fmt"

	"github.com/google/btree"
)

type primaryRecord struct {
	primary Value
}

func (r primaryRecord) Less(than btree.Item) bool {
	rThan, ok := than.(primaryRecord)
	if !ok {
		panic(fmt.Errorf("invalid than: expected %T, got %T", rThan, than))
	}
	return r.primary.Less(true, rThan.primary)
}

var _ btree.Item = primaryRecord{}

type forwardRecord struct {
	primary   Value
	secondary Value
}

func (r forwardRecord) Less(than btree.Item) bool {
	rThan, ok := than.(forwardRecord)
	if !ok {
		panic(fmt.Errorf("invalid than: expected %T, got %T", rThan, than))
	}
	if r.primary.Less(true, rThan.primary) {
		return true
	}
	if rThan.primary.Less(true, r.primary) {
		return false
	}
	// r.primary == rThan.primary

	if r.secondary != nil && rThan.secondary != nil {
		return r.secondary.Less(false, rThan.secondary)
	}

	if r.secondary == nil && rThan.secondary == nil {
		return false
	}

	if r.secondary == nil {
		return true
	}
	return false
}

var _ btree.Item = forwardRecord{}

type reverseRecord struct {
	secondary Value
	primary   Value
}

func (r reverseRecord) Less(than btree.Item) bool {
	rThan, ok := than.(reverseRecord)
	if !ok {
		panic(fmt.Errorf("invalid than: expected %T, got %T", rThan, than))
	}
	if r.secondary.Less(true, rThan.secondary) {
		return true
	}
	if rThan.secondary.Less(true, r.secondary) {
		return false
	}
	// r.primary == rThan.primary

	if r.primary != nil && rThan.primary != nil {
		return r.primary.Less(false, rThan.primary)
	}

	if r.primary == nil && rThan.primary == nil {
		return false
	}

	if r.primary == nil {
		return true
	}
	return false

}

var _ btree.Item = reverseRecord{}
