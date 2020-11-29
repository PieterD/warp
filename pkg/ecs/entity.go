package ecs

import (
	"fmt"
	"reflect"
)

type Entity struct {
	repo *Repository
	id   ID
}

func newEntity(repo *Repository, id ID) Entity {
	return Entity{
		repo: repo,
		id:   id,
	}
}

func (e Entity) ID() ID {
	return e.id
}

func (e Entity) SetComponent(rawComponent interface{}) error {
	repo := e.repo
	if !repo.HasEntity(e.id) {
		return errMissingEntity
	}
	val := reflect.ValueOf(rawComponent)
	typ := val.Type()
	switch kind := typ.Kind(); kind {
	case reflect.UnsafePointer, reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Invalid, reflect.Chan, reflect.Func, reflect.Map:
		return fmt.Errorf("disallowed component kind %v, type: %v", kind, typ)
	}
	components, ok := repo.components[typ]
	if !ok {
		components = make(map[ID]reflect.Value)
		repo.components[typ] = components
	}
	if _, ok := components[e.id]; !ok {
		components[e.id] = reflect.New(typ).Elem()
	}
	components[e.id].Set(val)
	return nil
}

func (e Entity) GetComponent(rawComponentPtr interface{}) error {
	repo := e.repo
	if !repo.HasEntity(e.id) {
		return errMissingEntity
	}
	val := reflect.ValueOf(rawComponentPtr)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr {
		return fmt.Errorf("only pointer types are allowed: %v", typ)
	}
	val = val.Elem()
	typ = typ.Elem()
	repoVal, ok := repo.components[typ][e.id]
	if !ok {
		return errMissingComponent
	}
	val.Set(repoVal)
	return nil
}

func (e Entity) DelComponent(rawComponent interface{}) error {
	repo := e.repo
	if !repo.HasEntity(e.id) {
		return errMissingEntity
	}
	typ := reflect.TypeOf(rawComponent)
	switch kind := typ.Kind(); kind {
	case reflect.UnsafePointer, reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Invalid, reflect.Chan, reflect.Func, reflect.Map:
		return fmt.Errorf("disallowed component kind %v, type: %v", kind, typ)
	}
	_, ok := repo.components[typ][e.id]
	if !ok {
		return errMissingComponent
	}
	delete(repo.components[typ], e.id)
	if len(repo.components[typ]) == 0 {
		delete(repo.components, typ)
	}
	return nil
}
