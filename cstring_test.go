package cstring

import (
	"testing"
	"unsafe"
)

func TestMakeCString(t *testing.T) {
	length := 10
	cStr := Make[byte](length)
	if len(cStr) != length {
		t.Errorf("expected length %d, got %d", length, len(cStr))
	}

	// Explicitly check zero-length C string
	length = 1
	cStr = Make[byte](length)
	if len(cStr) != length {
		t.Errorf("expected length %d, got %d", length, len(cStr))
	}
}

func TestMakeCStringZeroBytes(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for length less than 1")
		}
	}()
	Make[byte](0)
}

func TestNewCString(t *testing.T) {
	goStr := "hello\n\tworld 😊!"
	cStr := New[byte](goStr)
	if cStr.String() != goStr {
		t.Errorf("expected %s, got %s", goStr, cStr.String())
	}
	if len(cStr) != len(goStr)+1 {
		t.Errorf("expected length %d, got %d", len(goStr)+1, len(cStr))
	}
}

func TestNewCStringWithNullCharacter(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for string with null character")
		}
	}()
	New[byte]("hello\x00world")
}

func TestCStringEqualityLength(t *testing.T) {
	goStr := "hello\n\tworld 😊!"
	cStr := New[byte](goStr)
	ptr := unsafe.Pointer(cStr.Pointer())
	if !cStringEquals(ptr, goStr) {
		t.Errorf("expected string to be equal to %s", goStr)
	}
	if exp, got := len(goStr), strlen(ptr); exp != got {
		t.Errorf("expected length %d, got %d", exp, got)
	}
}
