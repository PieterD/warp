package mindex

type dbMini interface {
	Assign(primary Value, values ...Value)
	Remove(primary Value, values ...Value)
	First(primaryPtr Value, valuePtr Value) bool
	Next(primaryPtr Value, valuePtr Value) bool
}

func exact(db dbMini, primaryPtr Value, valuePtr Value) bool {
	primaryOrig := primaryPtr
	primaryCopy := CopyValue(primaryPtr)
	valueOrig := valuePtr
	valueCopy := CopyValue(valuePtr)
	if !db.First(primaryCopy, valueCopy) {
		return false
	}
	if primaryOrig.Less(primaryCopy) {
		return false
	}
	if valueOrig.Less(valueCopy) {
		return false
	}
	return true
}

func firstValue(db dbMini, primaryPtr Value, valuePtr Value) bool {
	primaryOrig := primaryPtr
	primaryCopy := CopyValue(primaryPtr)
	valueCopy := CopyValue(valuePtr)
	found := false
	traverseFrom(db, primaryCopy, valueCopy, func() bool {
		if primaryOrig.Less(primaryCopy) {
			// primary key advanced, bail
			return false
		}
		found = true
		return false
	})
	if !found {
		return false
	}
	CopyValueToValue(primaryPtr, primaryCopy)
	CopyValueToValue(valuePtr, valueCopy)
	return true
}

func nextValue(db dbMini, primaryPtr Value, valuePtr Value) bool {
	primaryOrig := primaryPtr
	primaryCopy := CopyValue(primaryPtr)
	valueOrig := valuePtr
	valueCopy := CopyValue(valuePtr)
	found := false
	traverseFrom(db, primaryCopy, valueCopy, func() bool {
		if primaryOrig.Less(primaryCopy) {
			// primary key advanced, bail
			return false
		}
		if valueOrig.Less(valuePtr) {
			// new value not equal to original
			found = true
			return false
		}
		return true
	})
	if !found {
		return false
	}
	CopyValueToValue(primaryPtr, primaryCopy)
	CopyValueToValue(valuePtr, valueCopy)
	return true
}

func traverseFrom(db dbMini, primaryPtr Value, valuePtr Value, f func() bool) {
	if !db.First(primaryPtr, valuePtr) {
		return
	}
	if !f() {
		return
	}
	for {
		if !db.Next(primaryPtr, valuePtr) {
			return
		}
		if !f() {
			return
		}
	}
}
