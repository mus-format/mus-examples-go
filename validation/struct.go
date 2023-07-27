package main

import (
	"errors"

	muscom "github.com/mus-format/mus-common-go"
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/varint"
	"github.com/ymz-ncnk/assert"
)

// Shows how you can check a struct field during unmarshalling.
func ValidateStruct() {
	// Defines a Foo.a field validator and Foo.a field Skipper.
	var (
		ErrTooBigAField                         = errors.New("too big 'a' field")
		avl             muscom.ValidatorFn[int] = func(a int) (err error) {
			if a > 10 {
				err = ErrTooBigAField
			}
			return
		}
		// Skips all subsequent Foo fields.
		ask mus.SkipperFn = func(bs []byte) (n int, err error) {
			n, err = ord.SkipBool(bs)
			if err != nil {
				return
			}
			var n1 int
			n1, err = ord.SkipString(bs[n:])
			n += n1
			return
		}
		foo  = Foo{a: 50, b: true, c: "🙂"}
		size = SizeFoo(foo) // == 7, where
		// 1 byte 	- 'a' field
		// 1 byte		- 'b' field
		// 5 bytes	- 'c' field
		bs = make([]byte, size)
	)
	MarshalFoo(foo, bs)

	// Decodes a struct, checking its 'a' field. Skips all bytes of an invalid
	// struct due to ask != nil.
	foo, n, err := UnmarshalValidFoo(avl, ask, bs)
	assert.EqualDeep(foo, Foo{})
	assert.Equal(n, 7) // All Foo bytes was used by Unmarshal funcation.
	assert.EqualError(err, ErrTooBigAField)

	// Decodes a struct, checking its 'a' field. Returns a validation error
	// immediately due to ask == nil.
	foo, n, err = UnmarshalValidFoo(avl, nil, bs)
	assert.EqualDeep(foo, Foo{})
	assert.Equal(n, 1) // Only one byte ('a' field) was used by Unmarshal
	// function.
	assert.EqualError(err, ErrTooBigAField)
}

type Foo struct {
	a int
	b bool
	c string
}

func MarshalFoo(v Foo, bs []byte) (n int) {
	n = varint.MarshalInt(v.a, bs)
	n += ord.MarshalBool(v.b, bs[n:])
	return n + ord.MarshalString(v.c, bs[n:])
}

func UnmarshalValidFoo(avl muscom.Validator[int], ask mus.Skipper, bs []byte) (
	v Foo, n int, err error) {
	a, n, err := varint.UnmarshalInt(bs)
	if err != nil {
		return
	}
	var (
		n1   int
		err1 error
	)
	if avl != nil {
		err = avl.Validate(a)
		if err != nil {
			if ask != nil { // If Skipper != nil, applies it, otherwise returns a
				// validation error immediately.
				n1, err1 = ask.SkipMUS(bs[n:])
				n += n1
				if err1 != nil {
					err = err1
				}
			}
			return
		}
	}
	v.a = a
	v.b, n1, err = ord.UnmarshalBool(bs[n:])
	n += n1
	if err != nil {
		return
	}
	v.c, n1, err = ord.UnmarshalString(bs[n:])
	n += n1
	return
}

func SizeFoo(v Foo) (size int) {
	size += varint.SizeInt(v.a)
	size += ord.SizeBool(v.b)
	return size + ord.SizeString(v.c)
}
