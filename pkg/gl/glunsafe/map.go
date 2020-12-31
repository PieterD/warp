package glunsafe

import (
	"fmt"
	"reflect"
)

// Map returns a byte slice mapped to the given rawPtr.
// It will be exactly the size of the type.
// rawPtr must be either a pointer or a slice.
// changes to rawPtr will be reflected there immediately, and vice versa.
func Map(rawPtr interface{}) []byte {
	ptrValue := reflect.ValueOf(rawPtr)
	switch ptrValue.Kind() {
	case reflect.Ptr:
		addr := ptrValue.Pointer()
		elemSize := int(ptrValue.Type().Elem().Size())
		return AddrToByteSlice(addr, elemSize)
	case reflect.Slice:
		addr := ptrValue.Pointer()
		totalSize := int(ptrValue.Type().Elem().Size()) * ptrValue.Len()
		return AddrToByteSlice(addr, totalSize)
	default:
		panic(fmt.Errorf("rawPtr is not of a usable type: %T", rawPtr))
	}
}
