package glunsafe

import (
	"fmt"
	"reflect"
)

// Map returns a byte slice mapped to the given rawPtr.
// It will be exactly the size of the type.
// changes to rawPtr will be reflected there immediately, and vice versa.
//TODO: expand to understand slices.
func Map(rawPtr interface{}) ([]byte, error) {
	ptrValue := reflect.ValueOf(rawPtr)
	if ptrValue.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("rawPtr is not a pointer type: %T", rawPtr)
	}
	addr := ptrValue.Pointer()
	elemSize := int(ptrValue.Type().Elem().Size())
	return AddrToByteSlice(addr, elemSize), nil
}
