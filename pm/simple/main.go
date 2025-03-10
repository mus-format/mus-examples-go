package main

import (
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

// This example demonstrates how the pm package works.
func main() {
	SerializeTwoPtrs()
	SerializeThreePtrs()
}

func SerializeTwoPtrs() {
	var (
		twoPtrsSer = MakeTwoPtrsSer()
		v          = NewTwoPtrs("the same pointer in two fields")
	)

	// Marshal TwoPtrs.
	bs := make([]byte, twoPtrsSer.Size(v))
	twoPtrsSer.Marshal(v, bs)

	// After unmarshal, v.ptr1 and av.ptr2 fields will contain the same
	// pointer.
	av, _, err := twoPtrsSer.Unmarshal(bs)
	assert.EqualError(err, nil)
	assert.Equal[*string](av.ptr1, av.ptr2)
}

func SerializeThreePtrs() {
	var (
		// ThreePtrs structure serializer uses TwoPtrs serializer.
		threePtrsSer = MakeThreePtrsSer()
		v            = NewThreePtrs("the same pointer in three fields")
	)

	// Marshal ThreePtrs.
	bs := make([]byte, threePtrsSer.Size(v))
	threePtrsSer.Marshal(v, bs)

	// After unmarshal, all fields will contain the same pointer.
	av, _, err := threePtrsSer.Unmarshal(bs)
	assert.EqualError(err, nil)
	assert.Equal[*string](av.ptr1, av.ptr2)
	assert.Equal[*string](av.ptr1, av.ptr3)
}

func NewTwoPtrs(str string) TwoPtrs {
	ptr := &str
	return TwoPtrs{
		ptr1: ptr,
		ptr2: ptr,
	}
}

func NewThreePtrs(str string) ThreePtrs {
	ptr := &str
	return ThreePtrs{
		TwoPtrs: TwoPtrs{
			ptr1: ptr,
			ptr2: ptr,
		},
		ptr3: ptr,
	}
}
