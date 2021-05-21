package bitstring

import (
	"bytes"
	"reflect"
	"unsafe"
)

// invariant: len(a) == len(b)
func u64cmp(a, b []uint64) bool {
	b = b[:len(a)] // remove BCE

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func bytesEq(a, b []uint64) bool {
	p := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&a)).Data)

	var aBytes []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&aBytes))
	hdr.Data = uintptr(p)
	hdr.Len = len(a) * 8
	hdr.Cap = len(a) * 8

	p = unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&b)).Data)

	var bBytes []byte
	hdr = (*reflect.SliceHeader)(unsafe.Pointer(&bBytes))
	hdr.Data = uintptr(p)
	hdr.Len = len(b) * 8
	hdr.Cap = len(b) * 8

	return bytes.Equal(aBytes, bBytes)
}
