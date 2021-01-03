package ecsdb

import (
	"fmt"
	"github.com/google/btree"
	"reflect"
	"sort"
)

const (
	btreeDegree = 100
)

type ID uint64

type DB struct {
	sequences    map[string]ID
	typesReverse map[reflect.Type]int
	types        []dbType // indexed by typeID
}

type Lesser interface {
	Less(than interface{}) bool
}

type dbType struct {
	rType reflect.Type
	data  *btree.BTree
	index *btree.BTree
}

func New(singletTypes, multiTypes []Lesser) (*DB, error) {
	if len(multiTypes)+len(singletTypes) == 0 {
		return nil, fmt.Errorf("must supply at least one type for storage")
	}
	db := &DB{
		sequences:    make(map[string]ID),
		typesReverse: make(map[reflect.Type]int),
		types:        make([]dbType, len(multiTypes)),
	}
	create := func(typeIndex int, storageType Lesser, singlet bool) error {
		rType := reflect.TypeOf(storageType)
		if err := validateType(rType); err != nil {
			return fmt.Errorf("invalid type %T: %w", storageType, err)
		}
		if _, ok := db.typesReverse[rType]; ok {
			return fmt.Errorf("duplicate type %T", storageType)
		}
		db.typesReverse[rType] = typeIndex
		db.types[typeIndex] = dbType{
			rType: rType,
			data:  btree.New(btreeDegree),
		}
		return nil
	}
	for typeIndex, storageType := range singletTypes {
		if err := create(typeIndex, storageType, true); err != nil {
			return nil, fmt.Errorf("creating singlet table %T: %w", storageType, err)
		}
	}
	for typeIndex, storageType := range multiTypes {
		if err := create(typeIndex, storageType, true); err != nil {
			return nil, fmt.Errorf("creating multikey table %T: %w", storageType, err)
		}
	}
	return db, nil
}

func validateType(rType reflect.Type) error {
	if rType.Kind() != reflect.Ptr {
		return fmt.Errorf("type %T is not a pointer", rType)
	}
	return nil
}

func (db *DB) Sequence(name string) ID {
	db.sequences[name]++
	return db.sequences[name]
}

func (db *DB) typeLookup(data Lesser) *dbType {
	rType := reflect.TypeOf(data)
	index, ok := db.typesReverse[rType]
	if !ok {
		panic(fmt.Errorf("invalid data type: please provide %T at construction time", data))
	}
	return &db.types[index]
}

func (db *DB) SetComponent(id ID, componentData Lesser) {
	typ := db.typeLookup(componentData)
	_ = typ.data.ReplaceOrInsert(&record{
		id:    id,
		value: componentData,
	})
	_ = typ.index.ReplaceOrInsert(&indexRecord{
		value: componentData,
		id:    id,
	})
}

func (db *DB) GetComponent(id ID, searchKey Lesser) (componentData Lesser) {
	typ := db.typeLookup(searchKey)
	found := typ.data.Get(&record{
		id:    id,
		value: searchKey,
	})
	if found == nil {
		return nil
	}
	foundRec, ok := found.(*record)
	if !ok {
		panic(fmt.Errorf("expected *record type, got: %T", found))
	}
	if foundRec.value == nil {
		panic(fmt.Errorf("found record has no data"))
	}
	return foundRec.value
}

func (db *DB) DelComponent(id ID, searchKey Lesser) bool {
	panic("not implemented")
}

func (db *DB) FindEntities(keys ...Lesser) []ID {
	if len(keys) == 0 {
		panic(fmt.Errorf("FindEntities requires at least 1 key"))
	}
	matches := make(map[ID]struct{})
	for i, key := range keys {
		ids := db.findEntitiesByComponent(key)
		if i == 0 {
			for _, id := range ids {
				matches[id] = struct{}{}
			}
			continue
		}
		newMatches := make(map[ID]struct{})
		for _, id := range ids {
			if _, ok := matches[id]; !ok {
				continue
			}
			newMatches[id] = struct{}{}
		}
		matches = newMatches
	}
	var ids []ID
	for id := range matches {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
	return ids
}

func (db *DB) findEntitiesByComponent(key Lesser) []ID {
	typ := db.typeLookup(key)
	var ids []ID
	typ.index.AscendGreaterOrEqual(&indexRecord{value: key, id: 0}, func(item btree.Item) bool {
		ir, ok := item.(*indexRecord)
		if !ok {
			panic(fmt.Errorf("expected *indexRecord type, got: %T", item))
		}
		if key.Less(ir.value) {
			// Value has advanced beyond key equality.
			return false
		}
		ids = append(ids, ir.id)
		return true
	})
	//sort.Slice(ids, func(i, j int) bool {
	//	return ids[i] < ids[j]
	//})
	return ids
}
