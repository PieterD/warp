package glunsafe

import (
	"reflect"
	"unsafe"
)

func AddrToByteSlice(addr uintptr, size int) []byte {
	slicePtr := unsafe.Pointer(&reflect.SliceHeader{
		Len:  size,
		Cap:  size,
		Data: addr,
	})
	return *(*[]byte)(slicePtr)
}

func FastUint16ToByte(b []uint16) []byte {
	return AddrToByteSlice(uintptr(unsafe.Pointer(&b[0])), len(b)*2)
}

func FastFloat32ToByte(b []float32) []byte {
	return AddrToByteSlice(uintptr(unsafe.Pointer(&b[0])), len(b)*4)
}
