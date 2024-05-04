package main

import (
	assert "github.com/ymz-ncnk/assert"
)

func init() {
	assert.On = true
}

func main() {
	str := "the same pointer in two fields"
	ptr := &str
	twoPtr := TwoPtr{
		ptr1: ptr,
		ptr2: ptr,
	}
	bs := make([]byte, SizeTowPtrMUS(twoPtr))
	MarshalTwoPtrMUS(twoPtr, bs)
	twoPtr, _, err := UnmarshalTwoPtrMUS(bs)
	assert.EqualError(err, nil)
	assert.Equal[*string](twoPtr.ptr1, twoPtr.ptr2)
}
