package main

import (
	"errors"
	"fmt"

	com "github.com/mus-format/common-go"
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/varint"
)

// Shows how you can check the length, keys and values of a map during
// unmarshalling.
func ValidateMap() {
	var (
		m1  mus.MarshallerFn[int]   = varint.MarshalInt   // Marshaler for map keys.
		u1  mus.UnmarshallerFn[int] = varint.UnmarshalInt // Unmarshaler for map keys.
		s1  mus.SizerFn[int]        = varint.SizeInt      // Sizer for map keys.
		sk1 mus.SkipperFn           = varint.SkipInt      // Skipper for map keys.

		m2  mus.MarshallerFn[string]   = ord.MarshalString   // Marshaler for map values.
		u2  mus.UnmarshallerFn[string] = ord.UnmarshalString // Unmarshaler for map values.
		s2  mus.SizerFn[string]        = ord.SizeString      // Sizer for map values.
		sk2 mus.SkipperFn              = ord.SkipString      // Skipper for map values.

		mp   = map[int]string{1: "hello", 2: "world", 3: "🙂"}
		size = ord.SizeMap[int, string](mp, s1, s2) // == 21, where
		// 1 byte 	- length of the map
		// 1 byte 	- 1 key
		// 6 bytes 	- "hello" value
		// 1 byte 	- 2 key
		// 6 bytes 	- "world" value
		// 1 byte 	- 3 key
		// 5 bytes  - "🙂" value
		bs = make([]byte, size)
	)
	ord.MarshalMap[int, string](mp, m1, m2, bs)

	// Defines a map length validator.
	var (
		ErrTooLongMap                      = errors.New("too long map")
		maxLength     com.ValidatorFn[int] = func(length int) (err error) {
			if length > 2 {
				err = ErrTooLongMap
			}
			return
		}
	)

	// Decodes a map, checking its length. Skips all bytes of an invalid map
	// due to sk != nil && sk2 != nil.
	mp, n, err := ord.UnmarshalValidMap[int, string](maxLength, u1, u2, nil, nil,
		sk1,
		sk2,
		bs)
	_ = mp
	_ = n
	_ = err
	// ------------------------------------------------------------------------ //
	// Assertions are commented out bacause the order of the map elements is
	// unpredictable. Expected order is [1: "hello", 2: "world", 3: "🙂"].
	// ------------------------------------------------------------------------ //
	// assert.EqualDeep(mp, map[int]string(nil))
	// assert.Equal(n, 21) // All map bytes was used by Unmarshal function.
	// assert.EqualError(err, ErrTooLongMap)

	// Decodes a map, checking its length. Returns a validation error
	// immediately due to a sk1 == nil || sk2 == nil.
	mp, n, err = ord.UnmarshalValidMap[int, string](maxLength, u1, u2, nil, nil,
		nil,
		nil,
		bs)
	_ = mp
	_ = n
	_ = err
	// assert.EqualDeep(mp, map[int]string(nil))
	// assert.Equal(n, 1) // Only one byte (the length of the map) was used by
	// // Unmarshal function.
	// assert.EqualError(err, ErrTooLongMap)

	// Defines a map keys validator.
	var (
		NewInvalidKeyError = func(key int) error {
			return fmt.Errorf("invalid %v key", key)
		}
		keyValidator com.ValidatorFn[int] = func(key int) (err error) {
			if key == 2 {
				err = NewInvalidKeyError(key)
			}
			return
		}
	)

	// Decodes a map, checking its keys. Skips all bytes after an invalid key
	// due to sk != nil && sk2 != nil.
	mp, n, err = ord.UnmarshalValidMap[int, string](nil, u1, u2, keyValidator,
		nil,
		sk1,
		sk2,
		bs)
	_ = mp
	_ = n
	_ = err
	// assert.EqualDeep(mp, map[int]string{1: "hello"})
	// assert.Equal(n, 21) // All map bytes was used by Unmarshal function.
	// assert.EqualError(err, NewInvalidKeyError(2))

	// Decodes a map, checking its keys. Returns a validation error
	// immediately due to sk1 == nil || sk2 == nil.
	mp, n, err = ord.UnmarshalValidMap[int, string](nil, u1, u2, keyValidator, nil,
		nil,
		nil,
		bs)
	_ = mp
	_ = n
	_ = err
	// assert.EqualDeep(mp, map[int]string{1: "hello"})
	// assert.Equal(n, 9) // Bytes of the 2 key was used by Unmarshal function too.
	// assert.EqualError(err, NewInvalidKeyError(2))

	// Defines a map values validator.
	var (
		NewInvalidValueError = func(value string) error {
			return fmt.Errorf("invalid \"%v\" value", value)
		}
		valValidator com.ValidatorFn[string] = func(value string) (err error) {
			if value == "world" {
				err = NewInvalidValueError(value)
			}
			return
		}
	)

	// Decodes a map, checking its values. Skips all bytes after an invalid value
	// due to sk != nil && sk2 != nil.
	mp, n, err = ord.UnmarshalValidMap[int, string](nil, u1, u2, nil,
		valValidator,
		sk1,
		sk2,
		bs)
	_ = mp
	_ = n
	_ = err
	// assert.EqualDeep(mp, map[int]string{1: "hello"})
	// assert.Equal(n, 21) // All map bytes was used by Unmarshal function.
	// assert.EqualError(err, NewInvalidValueError("world"))

	// Decodes a map, checking its values. Returns a validation error
	// immediately due to sk1 == nil || sk2 == nil.
	mp, n, err = ord.UnmarshalValidMap[int, string](nil, u1, u2, nil,
		valValidator,
		nil,
		nil,
		bs)
	_ = mp
	_ = n
	_ = err
	// assert.EqualDeep(mp, map[int]string{1: "hello"})
	// assert.Equal(n, 15) // Bytes of the "world" value was used by Unmarshal
	// // function too.
	// assert.EqualError(err, NewInvalidValueError("world"))

}
