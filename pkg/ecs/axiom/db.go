package axiom

import (
	"fmt"
	"reflect"

	"github.com/google/btree"
)

const btreeMaxSize = 100

type DB struct {
	tablesReverse map[reflect.Type]int
	tables        []*dbTable // indexed by typeID
}

type dbTable struct {
	typeName     string
	rType        reflect.Type
	forwardIndex *btree.BTree
	reverseIndex *btree.BTree
}

type Value interface {
	Less(highResolution bool, than Value) bool
}

func New() *DB {
	db := &DB{
		tablesReverse: make(map[reflect.Type]int),
	}
	return db
}

func (db *DB) tableLookup(v Value) *dbTable {
	rType := reflect.TypeOf(v)
	if rType.Kind() != reflect.Ptr {
		panic(fmt.Errorf("invalid value type %T: non-pointer types are not supported", v))
	}
	rType = rType.Elem()
	if rType.Kind() == reflect.Ptr {
		panic(fmt.Errorf("invalid value type %T: pointers of pointers are not supported", v))
	}
	index, ok := db.tablesReverse[rType]
	if !ok {
		dbt := &dbTable{
			typeName:     rType.String(),
			forwardIndex: btree.New(btreeMaxSize),
			reverseIndex: btree.New(btreeMaxSize),
		}
		index = len(db.tables)
		db.tablesReverse[rType] = index
		db.tables = append(db.tables, dbt)
	}
	return db.tables[index]
}

func (db *DB) Assign(primary Value, vs ...Value) {
	for _, v := range vs {
		table := db.tableLookup(v)
		table.forwardIndex.ReplaceOrInsert(forwardRecord{
			primary:   primary,
			secondary: v,
		})
		table.reverseIndex.ReplaceOrInsert(reverseRecord{
			primary:   primary,
			secondary: v,
		})
	}
}

// Next finds (for primary's record) the next value after valuePtr, and assigns it to valuePtr.
func (db *DB) NextValue(primary Value, valuePtr Value) bool {
	valueTable := db.tableLookup(valuePtr)
	found := false
	valueTable.forwardIndex.AscendGreaterOrEqual(forwardRecord{
		primary:   primary,
		secondary: valuePtr,
	}, func(i btree.Item) bool {
		rec, ok := i.(forwardRecord)
		if !ok {
			panic(fmt.Errorf("got %T, want %T", i, rec))
		}
		if primary.Less(true, rec.primary) {
			// next primary key, bail.
			return false
		}
		if valuePtr.Less(true, rec.secondary) {
			copyValueToValue(valuePtr, rec.secondary)
			found = true
			return false
		}
		return true
	})
	return found
}

// Search finds the next record after primaryPtr which matches all the given Values,
// and assigns its primary key to primaryPtr, and record data to vs.
func (db *DB) Search(primaryPtr Value, vs ...Value) bool {
	panic("not implemented")
	/*
		TODO:
		- use reverse index
		- IDs are now in ascending order
		- lockstep
		- copyValueToValue
	*/
}

func copyValue(v Value) Value {
	rValue := reflect.ValueOf(v)
	rValue = rValue.Elem()
	newValue := reflect.New(rValue.Type())
	iface := newValue.Interface()
	iValue, ok := iface.(Value)
	if !ok {
		panic(fmt.Errorf("new value %t of %T did not implement Value", iValue, v))
	}
	newValue.Elem().Set(rValue)
	return iValue
}

func copyValueToValue(dst, src Value) {
	reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(src).Elem())
}
