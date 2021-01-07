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

// FirstValue finds (for primary's record) the value equal to or greater than valuePtr, and sets it to valuePtr.
// Calling this method more than once on the same data will do nothing.
func (db *DB) FirstValue(primary Value, valuePtr Value) bool {
	return db.nextValue(true, primary, valuePtr)
}

// NextValue finds (for primary's record) the next value after valuePtr, and assigns it to valuePtr.
// Calling this method more than once will step through the assigned values in ascending order.
func (db *DB) NextValue(primary Value, valuePtr Value) bool {
	return db.nextValue(false, primary, valuePtr)
}

// FirstPrimary finds the first primary that's equal to or greater than primaryPtr (which must match value).
// Also assigns the record's data to valuePtr.
// Calling this method more than once on the same data will do nothing.
func (db *DB) FirstPrimary(primaryPtr Value, valuePtr Value) bool {
	return db.nextPrimary(true, primaryPtr, valuePtr)
}

// NextPrimary finds the next primary after primaryPtr which matches value.
// Also assigns the record's data to valuePtr.
// Calling this method more than once will step through the assigned primaries/values in ascending order.
func (db *DB) NextPrimary(primaryPtr Value, valuePtr Value) bool {
	return db.nextPrimary(false, primaryPtr, valuePtr)
}

func (db *DB) nextValue(first bool, primary Value, valuePtr Value) bool {
	valueTable := db.tableLookup(valuePtr)
	found := false
	fwdRec := forwardRecord{
		primary:   primary,
		secondary: valuePtr,
	}
	valueTable.forwardIndex.AscendGreaterOrEqual(fwdRec, func(i btree.Item) bool {
		rec, ok := i.(forwardRecord)
		if !ok {
			panic(fmt.Errorf("got %T, want %T", i, rec))
		}
		if primary.Less(true, rec.primary) {
			// next primary key, bail.
			return false
		}
		if !first && !valuePtr.Less(false, rec.secondary) {
			// skip the first record
			return true
		}
		copyValueToValue(valuePtr, rec.secondary)
		found = true
		return false
	})
	return found
}

func (db *DB) nextPrimary(first bool, primaryPtr Value, valuePtr Value) bool {
	valueTable := db.tableLookup(valuePtr)
	found := false
	fwdRec := reverseRecord{
		secondary: valuePtr,
		primary:   primaryPtr,
	}
	if first {
		fwdRec.primary = nil
	}
	valueTable.reverseIndex.AscendGreaterOrEqual(fwdRec, func(i btree.Item) bool {
		rec, ok := i.(reverseRecord)
		if !ok {
			panic(fmt.Errorf("got %T, want %T", i, rec))
		}
		if valuePtr.Less(false, rec.secondary) {
			// next value, bail.
			return false
		}
		if !first && !primaryPtr.Less(true, rec.primary) {
			// skip the first record
			return true
		}
		copyValueToValue(primaryPtr, rec.primary)
		copyValueToValue(valuePtr, rec.secondary)
		found = true
		return false
	})
	return found
}

func (db *DB) FirstSearch(primaryPtr Value, vs ...Value) bool {
	//if len(vs) == 0 {
	//	return false
	//}
	//type tuple struct {
	//	p Value
	//	v Value
	//}
	//tuples := make([]tuple, len(vs))
	//for i, v := range vs {
	//	p := copyValue(primaryPtr)
	//	tuples[i].p = p
	//	tuples[i].v = v
	//	if !db.FirstPrimary(tuples[i].p, tuples[i].v) {
	//		return false
	//	}
	//}
	//for {
	//
	//	sort.Slice(tuples, func(i, j int) bool {
	//		return tuples[i].p.Less(true, tuples[j].p)
	//	})
	//	allEqual := true
	//	for i := 1; i < len(tuples); i++ {
	//		if tuples[i-1].p.Less(true, tuples[i].p) {
	//			allEqual = false
	//			break
	//		}
	//	}
	//	highestPrimary := tuples[len(tuples)-1].p
	//	if allEqual {
	//		copyValueToValue(primaryPtr, highestPrimary)
	//		return true
	//	}
	//	for i := range tuples {
	//		copyValueToValue(tuples[i].p, highestPrimary)
	//		if !db.FirstPrimary(highestPrimary, tuples[i].v) {
	//			return false
	//		}
	//	}
	//}
	panic("not implemented")
}

// Search finds the next record after primaryPtr which matches all the given Values,
// and assigns its primary key to primaryPtr, and record data to vs.
func (db *DB) NextSearch(primaryPtr Value, vs ...Value) bool {
	panic("not implemented")
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
