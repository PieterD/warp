package ecs

import (
	"fmt"
	"github.com/google/btree"
)

type record struct {
	id      ID
	value   Lesser
	singlet bool
}

func (rec *record) Less(than btree.Item) bool {
	thanRec, ok := than.(*record)
	if !ok {
		panic(fmt.Errorf("casting than to *record, wrong type: %T", than))
	}
	if rec.singlet != thanRec.singlet {
		panic(fmt.Errorf("rec.singlet=%t, thanRec.singlet=%t", rec.singlet, thanRec.singlet))
	}
	if rec.id < thanRec.id {
		return true
	}
	if rec.id > thanRec.id {
		return false
	}
	// rec.id == thanRec.id
	if rec.singlet {
		return false
	}
	return rec.value.Less(thanRec.value)
}

var _ btree.Item = &record{}

type indexRecord struct {
	value Lesser
	id    ID
}

func (rec *indexRecord) Less(than btree.Item) bool {
	thanRec, ok := than.(*indexRecord)
	if !ok {
		panic(fmt.Errorf("casting than to *record, wrong type: %T", than))
	}
	if rec.value.Less(thanRec.value) {
		return true
	}
	if thanRec.value.Less(rec.value) {
		return false
	}
	// rec.value == thanRec.value

	if rec.id < thanRec.id {
		return true
	}
	if rec.id > thanRec.id {
		return false
	}
	// rec.id == thanRec.id

	return false
}

var _ btree.Item = &indexRecord{}
