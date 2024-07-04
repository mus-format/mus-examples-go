package main

import (
	"errors"
	"fmt"

	com "github.com/mus-format/common-go"
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/ord"
	"github.com/ymz-ncnk/assert"
)

// Shows how you can check the length and elements of a slice during
// Unmarshalling.
func ValidateSlice() {
	var (
		// Marshaller for slice elements.
		m mus.MarshallerFn[string] = func(t string, bs []byte) (n int) {
			return ord.MarshalString(t, nil, bs)
		}
		// Unmarshaller for slice elements.
		u mus.UnmarshallerFn[string] = func(bs []byte) (t string, n int, err error) {
			return ord.UnmarshalString(nil, bs)
		}
		// Sizer for slice elements.
		s mus.SizerFn[string] = func(t string) (size int) {
			return ord.SizeString(t, nil)
		}
		// Skipper for slice element.
		sk mus.SkipperFn = func(bs []byte) (n int, err error) {
			return ord.SkipString(nil, bs)
		}

		sl   = []string{"hello", "world", "🙂"}
		size = ord.SizeSlice[string](sl, nil, s) // == 18, where
		// 1 byte 	- length of the slice
		// 6 bytes 	- "hello" element
		// 6 bytes 	- "world" element
		// 5 bytes 	- "🙂" element.
		bs = make([]byte, size)
	)
	ord.MarshalSlice[string](sl, nil, m, bs)

	// Defines a slice length validator.
	var (
		ErrTooLongSlice                      = errors.New("too long slice")
		maxLength       com.ValidatorFn[int] = func(length int) (err error) {
			if length > 2 {
				err = ErrTooLongSlice
			}
			return
		}
	)

	// Decodes a slice, checking its length. Skips all bytes of an invalid slice
	// due to sk != nil.
	sl, n, err := ord.UnmarshalValidSlice[string](nil, maxLength, u, nil, sk, bs)
	assert.EqualDeep(sl, []string(nil))
	assert.Equal(n, 18) // All slice bytes was used by Unmarshal funcation.
	assert.EqualError(err, ErrTooLongSlice)

	// Decodes a slice, checking its length. Returns a validation error
	// immediately due to a sk == nil.
	sl, n, err = ord.UnmarshalValidSlice[string](nil, maxLength, u, nil, nil, bs)
	assert.EqualDeep(sl, []string(nil))
	assert.Equal(n, 1) // Only one byte (the length of the slice) was used by
	// Unmarshal function.
	assert.EqualError(err, ErrTooLongSlice)

	// Defines a slice elements validator.
	var (
		NewInvalidElemError = func(e string) error {
			return fmt.Errorf("invalid \"%v\" elem", e)
		}
		elemValidator com.ValidatorFn[string] = func(str string) (err error) {
			if str == "world" {
				err = NewInvalidElemError(str)
			}
			return
		}
	)

	// Decodes a slice, checking its elements. Skips all bytes after invalid
	// element due to sk != nil.
	sl, n, err = ord.UnmarshalValidSlice[string](nil, nil, u, elemValidator, sk,
		bs)
	assert.EqualDeep(sl, []string{"hello", "", ""})
	assert.Equal(n, 18) // All slice bytes was used by Unmarshal funcation.
	assert.EqualError(err, NewInvalidElemError("world"))

	// Decodes a slice, checking its elements. Returns a validation error
	// immediately due to a sk == nil.
	sl, n, err = ord.UnmarshalValidSlice[string](nil, nil, u, elemValidator, nil,
		bs)
	assert.EqualDeep(sl, []string{"hello", "", ""})
	assert.Equal(n, 13) // Bytes of the "world" element was used by Unmarshal
	// function too.
	assert.EqualError(err, NewInvalidElemError("world"))
}
