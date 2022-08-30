package main

import (
	"unsafe"
)

// SliceHeader is equivalent to reflect.SliceHeader, but represents the pointer
// to the underlying array as unsafe.Pointer rather than uintptr, allowing
// SliceHeaders to be directly converted to slice objects.
type SliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

// StringHeader is equivalent to reflect.StringHeader, but represents the
// pointer to the underlying array as unsafe.Pointer rather than uintptr,
// allowing StringHeaders to be directly converted to strings.
type StringHeader struct {
	Data unsafe.Pointer
	Len  int
}

// ImmutableBytesFromString is equivalent to []byte(s), except that it uses the
// same memory backing s instead of making a heap-allocated copy. This is only
// valid if the returned slice is never mutated.
func ImmutableBytesFromString(s string) (bs []byte) {
	shdr := (*StringHeader)(unsafe.Pointer(&s))
	bshdr := (*SliceHeader)(unsafe.Pointer(&bs))
	bshdr.Data = shdr.Data
	bshdr.Len = shdr.Len
	bshdr.Cap = shdr.Len
	return
}

// StringFromImmutableBytes is equivalent to string(bs), except that it uses
// the same memory backing bs instead of making a heap-allocated copy. This is
// only valid if bs is never mutated after StringFromImmutableBytes returns.
func StringFromImmutableBytes(bs []byte) string {
	// This is cheaper than messing with StringHeader and SliceHeader, which as
	// of this writing produces many dead stores of zeroes. Compare
	// strings.Builder.String().
	return *(*string)(unsafe.Pointer(&bs))
}
