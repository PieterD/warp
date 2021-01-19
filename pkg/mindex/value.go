package mindex

import (
	"fmt"
	"reflect"
)

type Singleton struct{}

func (s Singleton) Less(value Value) bool {
	return false
}

var _ Value = Singleton{}

type Value interface {
	Less(value Value) bool
}

type Relation interface {
	Value
	Relate(primary Value) (relationPrimary, relationValue Value)
}

func MustRelation(v Value) Relation {
	rel, ok := v.(Relation)
	if !ok {
		panic(fmt.Errorf("value %T does not implement Relation", v))
	}
	return rel
}

// A Copier is an optional Value interface.
// Every instance of Value is copied when they are passed into DB.Assign.
// When retrieving values, some destination value must have the value from the btree copied to it.
// Ordinarily, shallow copies are made using reflection.
// If a Value also implements Copier, it can implement proper copy methods.
type Copier interface {
	Value
	Copy() Value
	CopyTo(dst Value)
}

func CopyValue(v Value) Value {
	if copier, ok := v.(Copier); ok {
		return copier.Copy()
	}
	//TODO: Freeable. Automatic pools might get rid of a lot of new values.
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

func CopyValueToValue(dst, src Value) {
	if copier, ok := src.(Copier); ok {
		copier.CopyTo(dst)
		return
	}
	reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(src).Elem())
}
