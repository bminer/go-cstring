// This package provides a buffer that can be converted to a Go string or used
// in C functions that require a pointer to a C string. The advantage of this
// library is that the string is stored as a null-terminated C string in a Go
// slice, and memory management is handled by Go.
package cstring

// #include <string.h>
import "C"
import (
	"errors"
	"unsafe"

	"golang.org/x/exp/constraints"
)

// ByteConstraint is a type constraint that restricts the type to an 8-bit
// number. Since C.char can be signed or unsigned, this library permits either.
type ByteConstraint interface{ ~int8 | ~uint8 }

// CString is a buffer that can be converted to a Go string or used in C
// functions that require a pointer to a C string. Go will handle the memory
// allocation and freeing.
type CString[T ByteConstraint] []T

// Make creates a new CString with the given length in bytes. The provided
// length includes the null terminator. For example, the string "hello\0" has
// length 6. Panics if the length is less than 1.
func Make[T ByteConstraint, I constraints.Integer](n I) CString[T] {
	if n < 1 {
		panic("length must be at least 1")
	}
	return make(CString[T], n)
}

// New creates a new null-terminated CString from the given Go string. Panics
// if the specified Go string contains a null character.
func New[T ByteConstraint](s string) CString[T] {
	cStr, err := NewWithCheck[T](s)
	if err != nil {
		panic(err)
	}
	return cStr
}

// NewWithCheck creates a new null-terminated CString from the given Go string.
// Returns an error if the specified Go string contains a null character.
func NewWithCheck[T ByteConstraint](s string) (CString[T], error) {
	n := len(s) + 1 // to include the null terminator
	cStr := make(CString[T], n)
	for i := 0; i < n-1; i++ {
		if s[i] == 0 {
			return nil, errors.New("string contains null character")
		}
		cStr[i] = T(s[i])
	}
	cStr[n-1] = 0 // null-terminate the string
	return cStr, nil
}

// String returns the Go string representation of the CString by calling
// C.GoString.
func (s CString[T]) String() string {
	return C.GoString((*C.char)(unsafe.Pointer(&s[0])))
}

// Bytes returns a slice containing the Go string representation of the CString,
// but under the hood the underlying slice data is not copied.
func (s CString[T]) Bytes() []byte {
	if len(s) < 1 {
		return nil
	}
	dataPtr := unsafe.Pointer(unsafe.SliceData(s))
	return unsafe.Slice((*byte)(dataPtr), len(s)-1)
}

// Pointer returns the pointer to the first element of the CString. This
// function does not perform any conversions because the string is already
// stored internally as a null-terminated C string, so it is very fast.
func (s CString[T]) Pointer() *T {
	return &s[0]
}

// The following functions are used only for testing, as _test.go files cannot
// use cgo.

// strlen returns the length of the string pointed to by s.
func strlen(ptr unsafe.Pointer) int {
	cs := (*C.char)(ptr)
	return int(C.strlen(cs))
}

// cStringEquals returns true if the C string pointed to by h is equal
// to the Go string s.
func cStringEquals(ptr unsafe.Pointer, s string) bool {
	cs := (*C.char)(ptr)
	for i := 0; i < len(s); i++ {
		if *cs != C.char(s[i]) {
			return false
		}
		// Increment the pointer to the next byte
		cs = (*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cs)) + 1))
	}
	return *cs == 0
}
