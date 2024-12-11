package main

import (
	"errors"

	com "github.com/mus-format/common-go"
	"github.com/mus-format/mus-go/ord"
	"github.com/ymz-ncnk/assert"
)

// Shows how you can check the length of a string during unmarshalling.
func ValidateString() {
	var (
		str  = "hello world 🙂"
		size = ord.SizeString(str, nil) // == 17, where
		// 1 byte		- length of the string
		// 16 bytes	- string content.
		bs = make([]byte, size)
	)
	ord.MarshalString(str, nil, bs)

	// Defines a string length validator.
	var (
		ErrTooLongString                      = errors.New("too long string")
		maxLength        com.ValidatorFn[int] = func(length int) error {
			if length > 5 {
				return ErrTooLongString
			}
			return nil
		}
	)

	// Decodes a string, checking its length. Skips all bytes of an invalid string
	// due to skip == true.
	str, n, err := ord.UnmarshalValidString(nil, maxLength, true, bs)
	assert.Equal(str, "")
	assert.Equal(n, 17) // All string bytes was used by Unmarshal funcation.
	assert.EqualError(err, ErrTooLongString)

	// Decodes a string, checking its length. Returns a validation error
	// immediately due to a skip == false.
	str, n, err = ord.UnmarshalValidString(nil, maxLength, false, bs)
	assert.Equal(str, "")
	assert.Equal(n, 1) // Only one byte (the length of the string) was used by
	// Unmarshal function.
	assert.EqualError(err, ErrTooLongString)
}
