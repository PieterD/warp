package mindex

import (
	"fmt"
	"reflect"

	"github.com/google/btree"
)

const btreeMaxSize = 100

type dbTable struct {
	db           *DB
	key          tableKey
	fromRelation *dbTable
	withRelation *dbTable
	index        *btree.BTree
}

func (db *DB) table(primary, value Value) *dbTable {
	tableKey := newTableKey(primary, value)
	index, ok := db.tableKeys[tableKey]
	if ok {
		return db.tables[index]
	}
	dbt := &dbTable{
		db:    db,
		key:   tableKey,
		index: btree.New(btreeMaxSize),
	}
	index = len(db.tables)
	db.tableKeys[tableKey] = index
	db.tables = append(db.tables, dbt)

	return dbt
}

func (db *DB) relation(primary Value, value Value) (relPrimary, relValue Value, relTable *dbTable, ok bool) {
	rel, ok := value.(Relation)
	if !ok {
		return nil, nil, nil, false
	}
	table := db.table(primary, value)
	relPrimary, relValue = rel.Relate(primary)
	relTable = db.table(relPrimary, relValue)
	if relTable.fromRelation == nil {
		relTable.fromRelation = table
	}
	if relTable.fromRelation != table {
		panic(fmt.Errorf("circular relationship (1) detected"))
	}
	if table.withRelation == nil {
		table.withRelation = relTable
	}
	if table.withRelation != relTable {
		panic(fmt.Errorf("circular relationship (2) detected"))
	}
	return relPrimary, relValue, relTable, true
}

func (table *dbTable) assign(primary, value Value) {
	db := table.db

	if exact(db, primary, value) {
		db.Remove(primary, value)
	}
	_ = table.index.ReplaceOrInsert(tupleCopy(primary, value))
	if relPrimary, relValue, relTable, ok := db.relation(primary, value); ok {
		relTable.assign(relPrimary, relValue)
	}
}

func (table *dbTable) remove(primary, value Value) {
	db := table.db
	table.index.Delete(newTuple(primary, value))
	if relPrimary, relValue, relTable, ok := db.relation(primary, value); ok {
		relTable.remove(relPrimary, relValue)
	}
}

type tableKey struct {
	primary reflect.Type
	value   reflect.Type
}

func newTableKey(primary, value Value) tableKey {
	return tableKey{
		primary: valueToType(primary),
		value:   valueToType(value),
	}
}

func valueToType(v Value) reflect.Type {
	reflectType := reflect.TypeOf(v)
	if reflectType.Kind() != reflect.Ptr {
		panic(fmt.Errorf("invalid value type %T: non-pointer types are not supported", v))
	}
	reflectType = reflectType.Elem()
	if reflectType.Kind() == reflect.Ptr {
		panic(fmt.Errorf("invalid value type %T: pointers of pointers are not supported", v))
	}
	return reflectType
}
