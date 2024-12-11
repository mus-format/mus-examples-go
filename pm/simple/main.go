package main

import (
	assert "github.com/ymz-ncnk/assert"
)

func init() {
	assert.On = true
}

// This example demonstrates how the pm package works.
func main() {
	// The same pointer is used for TwoPtr.ptr1 and TwoPtr.ptr2 fields.
	var (
		str    = "the same pointer in two fields"
		ptr    = &str
		twoPtr = TwoPtr{
			ptr1: ptr,
			ptr2: ptr,
		}
	)

	// Marshal twoPtr.
	bs := make([]byte, SizeTowPtrMUS(twoPtr))
	MarshalTwoPtrMUS(twoPtr, bs)

	// After unmarshal, TwoPtr.ptr1 and TwoPtr.ptr2 fields will contain the same
	// pointer too.
	twoPtr, _, err := UnmarshalTwoPtrMUS(bs)
	assert.EqualError(err, nil)
	assert.Equal[*string](twoPtr.ptr1, twoPtr.ptr2)
}
