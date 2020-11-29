package ecs

import "fmt"

type ID uint64

type EntityIDSet map[ID]struct{}

func NewEntityIDSet(ids ...ID) EntityIDSet {
	s := make(map[ID]struct{})
	for _, id := range ids {
		s[id] = struct{}{}
	}
	return s
}

type IDCursor struct {
}

func (c *IDCursor) Fetch() bool {
	panic(fmt.Errorf("not implemented"))
}

func (c *IDCursor) Get() ID {
	panic(fmt.Errorf("not implemented"))
}

func (c *IDCursor) Skip(to ID) {
	panic(fmt.Errorf("not implemented"))
}

func (c *IDCursor) Err() error {
	panic(fmt.Errorf("not implemented"))
}
