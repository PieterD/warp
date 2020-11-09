package gl

import (
	"reflect"
	"unsafe"
)

func newSliceHeader(data unsafe.Pointer, size int) unsafe.Pointer {
	return unsafe.Pointer(&reflect.SliceHeader{
		Len:  size,
		Cap:  size,
		Data: uintptr(data),
	})
}

func fastUint16ToByte(b []uint16) []byte {
	return *(*[]byte)(newSliceHeader(unsafe.Pointer(&b[0]), len(b)*2))
}

func fastFloat32ToByte(b []float32) []byte {
	return *(*[]byte)(newSliceHeader(unsafe.Pointer(&b[0]), len(b)*4))
}
